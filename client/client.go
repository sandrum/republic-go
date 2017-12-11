package main

import (
	"github.com/republicprotocol/republic/crypto"
	"github.com/republicprotocol/republic/rpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
)

const (
	address = "localhost:8080"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := rpc.NewNodeClient(conn)

	// Generating idendity for the client node
	secp, err := crypto.NewSECP256K1()
	if err != nil {
		log.Fatalf("failed to identify self: %v", err)
	}
	id := secp.PublicAddress()
	grpcID := &rpc.ID{Address: id}
	log.Println("Republic address:", id)

	// Ping the server
	log.Println("Ping: " + address)
	rID, err := c.Ping(context.Background(), grpcID)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Pong: %s \n", rID.Address)

	// Get all peers from the server
	log.Printf("Ask all peers from : %s \n", rID.Address)
	rMultiAddresses, err := c.Peers(context.Background(), grpcID)
	if err != nil {
		log.Fatalf("could not get peers: %v", err)
	}

	for _, j := range rMultiAddresses.Multis {
		log.Printf("Get Peer from server : %s \n", j)
	}

	// Get peers from a node that are closer to a target than the node itself
	log.Printf("Ask peers close to: %s \n", id)
	rMultiAddresses, err = c.CloserPeers(context.Background(), &rpc.Path{To: grpcID, From: grpcID})
	if err != nil {
		log.Fatalf("could not get peers: %v", err)
	}

	for _, j := range rMultiAddresses.Multis {
		log.Printf("Close peer : %s \n", j)
	}

}
