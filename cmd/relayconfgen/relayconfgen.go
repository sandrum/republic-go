package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/republicprotocol/republic-go/blockchain/ethereum"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/relay"
)

func main() {

	keyPair, err := identity.NewKeyPair()
	if err != nil {
		log.Fatalf("cannot create ecdsa key: %v", err)
	}
	auth := bind.NewKeyedTransactor(keyPair.PrivateKey)
	multiAddress, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/0.0.0.0/tcp/18515/republic/%s", keyPair.ID()))
	if err != nil {
		log.Fatalf("cannot create multiAddress: %v", err)
	}
	conf := relay.Config{
		KeyPair:                 keyPair,
		MultiAddress:            multiAddress,
		BootstrapMultiAddresses: identity.MultiAddresses{},
		Token:           "",
		EthereumAddress: auth.From.String(),
		Ethereum: ethereum.Config{
			Network:                 ethereum.NetworkRopsten,
			URI:                     "https://ropsten.infura.io",
			RepublicTokenAddress:    ethereum.RepublicTokenAddressOnRopsten.String(),
			DarknodeRegistryAddress: ethereum.DarknodeRegistryAddressOnRopsten.String(),
			HyperdriveAddress:       ethereum.HyperdriveAddressOnRopsten.String(),
			ArcAddress:              ethereum.ArcAddressOnRopsten.String(),
		},
	}

	file, err := os.Create("config.json")
	if err != nil {
		log.Fatalf("cannot create file: %v", err)
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(conf); err != nil {
		log.Fatalf("cannot write conf to file: %v", err)
	}
}