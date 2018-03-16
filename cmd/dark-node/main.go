package main

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/republicprotocol/republic-go/contracts/connection"
	"github.com/republicprotocol/republic-go/contracts/dnr"
	"github.com/republicprotocol/republic-go/dark-node"
	"github.com/republicprotocol/republic-go/identity"
)

func main() {
	// Load configuration path from the command line
	configFilename := flag.String("config", "/home/.darknode/config.json", "Path to the JSON configuration file")
	flag.Parse()

	// Load the default configuration
	config, err := LoadConfig(*configFilename)
	if err != nil {
		log.Fatal(err)
	}

	// Create a dark node registrar.
	darkNodeRegistrar, err := CreateDarkNodeRegistrar(config.EthereumKey, config.EthereumRPC)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new node.node.
	node, err := node.NewDarkNode(*config, darkNodeRegistrar)
	if err != nil {
		log.Fatal(err)
	}

	go node.StartServices()
	go node.StartAPI()
	go node.StartUI()
	node.StartBackgroundWorkers()
	node.Bootstrap()
	node.WatchDarkOcean()
}

// LoadConfig returns a default Config object for the Falcon testnet.
func LoadConfig(filename string) (*node.Config, error) {
	var writeBack bool

	// Load configuration file
	config, err := node.LoadConfig(filename)
	if err != nil {
		return nil, err
	}

	// Generate our ethereum keypair
	if config.EthereumKey.PrivateKey == nil {
		writeBack = true
		config.EthereumKey = *keystore.NewKeyForDirectICAP(rand.Reader)
	}

	if config.KeyPair.PrivateKey == nil {
		writeBack = true
		// Get an address and keypair
		address, keyPair, err := identity.NewAddress()
		if err != nil {
			return nil, err
		}

		// Get our IP address
		out, err := exec.Command("curl", "https://ipinfo.io/ip").Output()
		out = []byte(strings.Trim(string(out), "\n "))
		if err != nil {
			return nil, err
		}
		if err != nil {
			return nil, err
		}

		// Generate our multiaddress
		multiAddress, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/%v/tcp/%v/republic/%v", string(out), config.Port, address.String()))
		if err != nil {
			return nil, err
		}

		config.KeyPair = keyPair
		config.NetworkOptions.MultiAddress = multiAddress
	}

	if writeBack {
		file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		if err := json.NewEncoder(file).Encode(config); err != nil {
			return nil, err
		}
	}

	return config, nil
}

func CreateDarkNodeRegistrar(ethereumKey keystore.Key, ethereumRPC string) (dnr.DarkNodeRegistrar, error) {
	auth := bind.NewKeyedTransactor(ethereumKey.PrivateKey)
	client, err := connection.FromURI(ethereumRPC, connection.ChainRopsten)
	if err != nil {
		return nil, err
	}
	return dnr.NewEthereumDarkNodeRegistrar(context.Background(), &client, auth, &bind.CallOpts{})
}