package miner

import (
	"log"
	"math/big"
	"runtime"

	"github.com/jbenet/go-base58"
	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-order-compute"
	"github.com/republicprotocol/go-rpc"
	"github.com/republicprotocol/go-swarm-network"
	"github.com/republicprotocol/go-xing"
)

// TODO: Do not make this values constant.
var (
	N        = int64(3)
	K        = int64(2)
	Prime, _ = big.NewInt(0).SetString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137859", 10)
)

type Miner struct {
	Computer *compute.ComputationMatrix
	Swarm    *swarm.Node
	Xing     *xing.Node
}

func NewMiner(config *Config) (*Miner, error) {
	miner := &Miner{
		Computer: compute.NewComputationMatrix(),
	}

	swarmOptions := swarm.Options{
		MultiAddress:            config.MultiAddress,
		BootstrapMultiAddresses: config.BootstrapMultiAddresses,
		Debug: swarm.DebugHigh,
	}
	swarmNode := swarm.NewNode(miner, swarmOptions)
	miner.Swarm = swarmNode

	xingOptions := xing.Options{
		Address: config.MultiAddress.Address(),
		Debug:   xing.DebugHigh,
	}
	xingNode := xing.NewNode(miner, xingOptions)
	miner.Xing = xingNode

	return miner, nil
}

// EstablishConnections to other peers in the swarm network by bootstrapping
// against a set of bootstrap network.Nodes.
func (miner *Miner) EstablishConnections() {
	miner.Swarm.Bootstrap()
}

// OnPingReceived implements the network.Delegate interface. It is used by the
// underlying network.Node whenever the Miner needs to handle a Ping RPC.
func (miner *Miner) OnPingReceived(peer identity.MultiAddress) {
}

// OnQueryCloserPeersReceived implements the network.Delegate interface. It is
// used by the underlying network.Node whenever the Miner needs to handle a
// QueryCloserPeers RPC.
func (miner *Miner) OnQueryCloserPeersReceived(peer identity.MultiAddress) {
}

// OnQueryCloserPeersOnFrontierReceived implements the network.Delegate
// interface. It is used by the underlying network.Node whenever the Miner
// needs to handle a QueryCloserPeersOnFrontier RPC.
func (miner *Miner) OnQueryCloserPeersOnFrontierReceived(peer identity.MultiAddress) {
}

func (miner *Miner) OnOrderFragmentReceived(from identity.MultiAddress, orderFragment *compute.OrderFragment) {
	miner.Computer.AddOrderFragment(orderFragment)
}

func (miner *Miner) OnResultFragmentReceived(from identity.MultiAddress, resultFragment *compute.ResultFragment) {
	miner.addResultFragments([]*compute.ResultFragment{resultFragment})
}

func (miner *Miner) OnOrderFragmentForwarding(to identity.Address, from identity.MultiAddress, orderFragment *compute.OrderFragment) {
}

func (miner *Miner) OnResultFragmentForwarding(to identity.Address, from identity.MultiAddress, resultFragment *compute.ResultFragment) {
}

func (miner *Miner) Mine(quit chan struct{}) {
	for {
		select {
		case <-quit:
			miner.Xing.Stop()
			miner.Swarm.Stop()
			return
		default:
			// FIXME: If this function call blocks forever then the quit signal
			// will never be received.
			miner.ComputeAll()
		}
	}
}

func (miner Miner) ComputeAll() {
	numberOfCPUs := runtime.NumCPU()
	computations := miner.Computer.WaitForComputations(numberOfCPUs)
	resultFragments := make([]*compute.ResultFragment, len(computations))

	do.CoForAll(computations, func(i int) {
		resultFragment, err := miner.Compute(computations[i])
		if err != nil {
			return
		}
		resultFragments[i] = resultFragment
	})

	go func() {
		resultFragmentsOk := make([]*compute.ResultFragment, 0, len(resultFragments))
		for _, resultFragment := range resultFragments {
			if resultFragment != nil {
				resultFragmentsOk = append(resultFragmentsOk, resultFragment)
			}
		}
		miner.addResultFragments(resultFragmentsOk)
	}()
}

// Compute the required computation on two OrderFragments and send the result
// to all Miners in the M Network.
// TODO: Send computed order fragments to the M Network instead of all peers.
func (miner Miner) Compute(computation *compute.Computation) (*compute.ResultFragment, error) {
	resultFragment, err := computation.Sub(Prime)
	if err != nil {
		return nil, err
	}
	go func() {
		for _, multiAddress := range miner.Swarm.DHT.MultiAddresses() {
			rpc.SendResultFragmentToTarget(multiAddress, multiAddress.Address(), miner.Swarm.MultiAddress(), resultFragment, miner.Swarm.Options.Timeout)
		}
	}()
	return resultFragment, nil
}

func (miner Miner) addResultFragments(resultFragments []*compute.ResultFragment) {
	results, _ := miner.Computer.AddResultFragments(resultFragments, K, Prime)
	for _, result := range results {
		if result.IsMatch(Prime) {
			log.Printf("match found for buy = %s, sell = %s\n", base58.Encode(result.BuyOrderID), base58.Encode(result.SellOrderID))
		}
	}
}
