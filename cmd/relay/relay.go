package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/republicprotocol/republic-go/blockchain/ethereum/arc"
	"github.com/republicprotocol/republic-go/blockchain/ethereum/hd"
	"github.com/republicprotocol/republic-go/blockchain/swap"
	"github.com/republicprotocol/republic-go/stackint"

	"github.com/republicprotocol/republic-go/order"

	. "github.com/republicprotocol/republic-go/relay"

	abiBind "github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/republicprotocol/republic-go/blockchain/ethereum"
	"github.com/republicprotocol/republic-go/blockchain/ethereum/dnr"
	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/orderbook"
	"github.com/republicprotocol/republic-go/rpc/client"
	"github.com/republicprotocol/republic-go/rpc/dht"
	"github.com/republicprotocol/republic-go/rpc/relayer"
	"github.com/republicprotocol/republic-go/rpc/smpcer"
	"github.com/republicprotocol/republic-go/rpc/swarmer"
	"google.golang.org/grpc"
)

func main() {
	// keystore := flag.String("keystore", "", "Encrypted keystore file")
	// passphrase := flag.String("passphrase", "", "Passphrase for the encrypted keystore file")
	bind := flag.String("bind", "127.0.0.1", "Binding address for the gRPC and HTTP API")
	port := flag.Int("port", 18515, "Binding port for the HTTP API")
	// token := flag.String("token", "", "Bearer token for restricting access")
	configLocation := flag.String("config", "", "Relay configuration file location")
	maxConnections := flag.Int("maxConnections", 4, "Maximum number of connections to peers during synchronization")
	flag.Parse()

	// fmt.Println("Decrypting keystore...")
	// key, err := getKey(*keystore, *passphrase)
	// if err != nil {
	// 	fmt.Println(fmt.Errorf("cannot obtain key: %s", err))
	// 	return
	// }

	// keyPair, err := getKeyPair(key)
	// if err != nil {
	// 	fmt.Println(fmt.Errorf("cannot obtain keypair: %s", err))
	// 	return
	// }

	// multiAddr, err := getMultiaddress(keyPair, *port)
	// if err != nil {
	// 	fmt.Println(fmt.Errorf("cannot obtain multiaddress: %s", err))
	// 	return
	// }

	config, err := LoadConfig(*configLocation)
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	// Create gRPC server and TCP listener always using port 18514
	server := grpc.NewServer()
	listener, err := net.Listen("tcp", fmt.Sprintf("%v:%v", *bind, *port))
	if err != nil {
		log.Fatal(err)
	}

	// Create Relay
	// config := Config{
	// 	KeyPair:      keyPair,
	// 	MultiAddress: multiAddr,
	// 	Token:        *token,
	// }

	registrar, err := getRegistry(config)
	if err != nil {
		fmt.Println(fmt.Errorf("cannot obtain registrar: %s", err))
		return
	}

	hyperdrive, err := getHyperdrive(config)
	if err != nil {
		fmt.Println(fmt.Errorf("cannot obtain hyperdrive: %s", err))
		return
	}

	book := orderbook.NewOrderbook(100)
	crypter := crypto.NewWeakCrypter()
	dht := dht.NewDHT(config.MultiAddress.Address(), 100)
	connPool := client.NewConnPool(100)
	relayerClient := relayer.NewClient(&crypter, &dht, &connPool)
	smpcerClient := smpcer.NewClient(&crypter, config.MultiAddress, &connPool)
	swarmerClient := swarmer.NewClient(&crypter, config.MultiAddress, &dht, &connPool)
	relay := NewRelay(*config, registrar, &book, &relayerClient, &smpcerClient, &swarmerClient)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	for err := range swarmerClient.Bootstrap(ctx, config.BootstrapMultiAddresses, -1) {
		log.Printf("error while bootstrapping the relay: %v", err)
	}

	entries := make(chan orderbook.Entry)
	defer close(entries)
	go func() {
		defer book.Unsubscribe(entries)
		if err := book.Subscribe(entries); err != nil {
			log.Fatalf("cannot subscribe to orderbook: %v", err)
		}
	}()
	confirmedOrders := processOrderbookEntries(hyperdrive, entries)
	conn, err := ethereum.Connect(config.Ethereum)
	auth := abiBind.NewKeyedTransactor(config.KeyPair.PrivateKey)
	if err != nil {
		log.Fatalf("cannot fetch dark node registry: %s", err)
	}
	swaps := executeConfirmedOrders(context.Background(), conn, auth, hyperdrive, confirmedOrders)
	processAtomicSwaps(swaps)

	// Server gRPC and RESTful API
	fmt.Println(fmt.Sprintf("Relay API available at %s:%v", *bind, *port+1))
	dispatch.CoBegin(func() {
		if err := relay.ListenAndServe(*bind, fmt.Sprintf("%d", *port+1)); err != nil {
			log.Fatalf("error serving http: %v", err)
		}
	}, func() {
		relay.Register(server)
		if err := server.Serve(listener); err != nil {
			log.Fatalf("error serving grpc: %v", err)
		}
	}, func() {
		// if err := relay.Sync(context.Background(), *maxConnections); err != nil {
		// 	log.Fatalf("error syncing relay: %v", err)
		// }
		relay.Sync(context.Background(), *maxConnections)
	})

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	go func() {
		<-sig
		server.Stop()
	}()
}

func processOrderbookEntries(hyperdrive hd.HyperdriveContract, entryInCh <-chan orderbook.Entry) <-chan orderbook.Entry {
	unconfirmedOrders := make(chan orderbook.Entry, 100)
	confirmedEntries := make(chan orderbook.Entry)

	orderConfirmed := func(orderID []byte) bool {
		depth, err := hyperdrive.GetDepth(orderID)
		if err != nil {
			log.Fatalf("failed to get depth: %v", err)
		}
		return depth >= 5
	}

	go func() {
		defer close(confirmedEntries)
		for {
			select {
			case entry, ok := <-entryInCh:
				log.Println("we get order from the orderbook", entry.Order.ID.String(), entry.Status)
				if !ok {
					return
				}
				if !orderConfirmed(entry.Order.ID) {
					unconfirmedOrders <- entry
				} else {
					entry.Status = order.Confirmed
					confirmedEntries <- entry
				}
			}
		}
	}()

	go func() {
		defer close(unconfirmedOrders)
		for {
			select {
			case entry, ok := <-unconfirmedOrders:
				if !ok {
					return
				}
				if !orderConfirmed(entry.Order.ID) {
					unconfirmedOrders <- entry
					time.Sleep(time.Second)
				} else {
					entry.Status = order.Confirmed
					confirmedEntries <- entry
				}
			}
		}
	}()
	return confirmedEntries
}

func executeConfirmedOrders(ctx context.Context, conn ethereum.Conn, auth *abiBind.TransactOpts, hyperdrive hd.HyperdriveContract, entries <-chan orderbook.Entry) <-chan swap.Swap {
	swaps := make(chan swap.Swap)

	go func() {
		defer close(swaps)
		for {
			select {
			case entry, ok := <-entries:
				if !ok {
					return
				}
				orderID := [32]byte{}
				copy(orderID[:], entry.Order.ID)
				_, orderIDs, err := hyperdrive.GetOrderMatch(orderID)
				if err != nil {
					log.Fatalf("failed to get order match: %v", err)
					continue
				}
				if orderID == orderIDs[0] {
					swaps <- initSwap(ctx, conn, auth, entry, orderIDs[0], orderIDs[1])
				} else {
					swaps <- initSwap(ctx, conn, auth, entry, orderIDs[1], orderIDs[0])
				}
			}
		}
	}()

	return swaps
}

func getKey(filename, passphrase string) (*keystore.Key, error) {
	// Read data from the keystore file and generate the key
	encryptedKey, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("cannot read keystore file: %v", err)
	}

	key, err := keystore.DecryptKey(encryptedKey, passphrase)
	if err != nil {
		return nil, fmt.Errorf("cannot decrypt key with provided passphrase: %v", err)
	}

	return key, nil
}

func getKeyPair(key *keystore.Key) (identity.KeyPair, error) {
	id, err := identity.NewKeyPairFromPrivateKey(key.PrivateKey)
	if err != nil {
		return identity.KeyPair{}, fmt.Errorf("cannot generate id from key %v", err)
	}
	return id, nil
}

func getMultiaddress(id identity.KeyPair, port string) (identity.MultiAddress, error) {
	// Get our IP address
	ipInfoOut, err := exec.Command("curl", "https://ipinfo.io/ip").Output()
	if err != nil {
		return identity.MultiAddress{}, err
	}
	ipAddress := strings.Trim(string(ipInfoOut), "\n ")

	relayMultiaddress, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/%s/tcp/%s/republic/%s", ipAddress, port, id.Address().String()))
	if err != nil {
		return identity.MultiAddress{}, fmt.Errorf("cannot obtain trader multi address %v", err)
	}

	return relayMultiaddress, nil
}

func getRegistry(config *Config) (dnr.DarknodeRegistry, error) {
	conn, err := ethereum.Connect(config.Ethereum)
	auth := abiBind.NewKeyedTransactor(config.KeyPair.PrivateKey)
	if err != nil {
		fmt.Println(fmt.Errorf("cannot fetch dark node registry: %s", err))
		return dnr.DarknodeRegistry{}, err
	}
	auth.GasPrice = big.NewInt(6000000000)
	registrar, err := dnr.NewDarknodeRegistry(context.Background(), conn, auth, &abiBind.CallOpts{})
	if err != nil {
		fmt.Println(fmt.Errorf("cannot fetch dark node registry: %s", err))
		return dnr.DarknodeRegistry{}, err
	}
	return registrar, nil
}

func getHyperdrive(config *Config) (hd.HyperdriveContract, error) {
	conn, err := ethereum.Connect(config.Ethereum)
	auth := abiBind.NewKeyedTransactor(config.KeyPair.PrivateKey)
	if err != nil {
		fmt.Println(fmt.Errorf("cannot fetch hyperdrive: %s", err))
		return hd.HyperdriveContract{}, err
	}
	auth.GasPrice = big.NewInt(6000000000)
	hyperdrive, err := hd.NewHyperdriveContract(context.Background(), conn, auth, &abiBind.CallOpts{})
	if err != nil {
		fmt.Println(fmt.Errorf("cannot fetch hyperdrive: %s", err))
		return hd.HyperdriveContract{}, err
	}
	return hyperdrive, nil
}

func initSwap(ctx context.Context, conn ethereum.Conn, auth *abiBind.TransactOpts, entry orderbook.Entry, fstOrderID, sndOrderID [32]byte) swap.Swap {
	fstTokenAddress := getTokenAddress(entry.FstCode)
	sndTokenAddress := getTokenAddress(entry.SndCode)
	goesFirst := entry.Parity == order.ParitySell

	expiry := time.Now().Unix() + 24*60*60
	from := common.HexToAddress("0x1a459f0dF58cF0B9a4246Bd193a00125B45492Df").Bytes()
	to := common.HexToAddress("0xad560E16C7f474281A31c5F38F903382EaBAc107").Bytes()
	value := getFstValue(entry.Price, entry.MaxVolume) // Assuming max and min volume to be same
	if goesFirst {
		expiry = time.Now().Unix() + 48*60*60
		from = common.HexToAddress("0xad560E16C7f474281A31c5F38F903382EaBAc107").Bytes()
		to = common.HexToAddress("0x1a459f0dF58cF0B9a4246Bd193a00125B45492Df").Bytes()
		value = getSndValue(entry.Price, entry.MaxVolume) // Assuming max and min volume to be same
	}
	order := entry.Order.Bytes()
	fee := big.NewInt(2)

	fstArc, err := arc.NewArc(ctx, conn, auth, order, fstTokenAddress, fee)
	if err != nil {
		log.Fatalf("cannot create new arc: %s", err)
	}
	sndArc, err := arc.NewArc(ctx, conn, auth, []byte{}, sndTokenAddress, fee)
	if err != nil {
		log.Fatalf("cannot create new arc: %s", err)
	}
	fstItem := swap.SwapItem{
		OrderID:   fstOrderID,
		From:      from,
		To:        to,
		Value:     value,
		Expiry:    expiry,
		Arc:       fstArc,
		GoesFirst: goesFirst,
	}
	sndItem := swap.SwapItem{
		OrderID: sndOrderID,
		Arc:     sndArc,
	}
	return swap.NewSwap(fstItem, sndItem)
}

func processAtomicSwaps(swaps <-chan swap.Swap) {
	go func() {
		for {
			select {
			case swap, ok := <-swaps:
				if !ok {
					return
				}
				go func() {
					err := swap.Execute(context.Background())
					if err != nil {
						log.Fatalf("failed to execute the atomic swap: %v", err)
						return
					}
				}()
			}
		}
	}()
}

func getTokenAddress(currencyCode order.CurrencyCode) common.Address {
	switch currencyCode {
	case order.CurrencyCodeREN:
		return common.HexToAddress("0x65d54EDa5f032F2275Caa557E50c029cFbCCBB54")
	case order.CurrencyCodeETH:
		return common.HexToAddress("0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE")
	default:
		log.Fatalf("error unsupported currency code")
		return common.Address{}
	}
}

func getFstValue(price, volume stackint.Int1024) *big.Int {
	value := volume.Mul(&price)
	return value.ToBigInt()
}

func getSndValue(price, volume stackint.Int1024) *big.Int {
	one := stackint.One()
	price = one.Div(&price)
	value := volume.Mul(&price)
	return value.ToBigInt()
}
