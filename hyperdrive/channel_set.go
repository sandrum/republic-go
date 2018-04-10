package hyper

import (
	"sync"

	"github.com/republicprotocol/republic-go/dispatch"
)

type ChannelSet struct {
	Proposal chan Proposal
	Prepare  chan Prepare
	Fault    chan Fault
	Commit   chan Commit
	Err      chan error
	Block    chan Block
}

func NewChannelSet(proposal chan Proposal, prepare chan Prepare, commit chan Commit, fault chan Fault, block chan Block, err chan error) ChannelSet {
	return ChannelSet{
		Proposal: proposal,
		Prepare:  prepare,
		Commit:   commit,
		Fault:    fault,
		Block:    block,
		Err:      err,
	}
}

func EmptyChannelSet() ChannelSet {
	return ChannelSet{
		Proposal: make(chan Proposal, 240),
		Prepare:  make(chan Prepare, 240),
		Fault:    make(chan Fault, 240),
		Commit:   make(chan Commit, 240),
		Err:      make(chan error, 240),
		Block:    make(chan Block, 240),
	}
}

func (c *ChannelSet) Close() {
	close(c.Proposal)
	close(c.Prepare)
	close(c.Fault)
	close(c.Commit)
	close(c.Err)
	close(c.Block)
}

func (c *ChannelSet) Split(cs []ChannelSet) {
	var wg sync.WaitGroup

	proposals := make([]chan Proposal, len(cs))
	prepares := make([]chan Prepare, len(cs))
	commits := make([]chan Commit, len(cs))
	faults := make([]chan Fault, len(cs))
	errs := make([]chan error, len(cs))
	blocks := make([]chan Block, len(cs))

	for i, chset := range cs {
		proposals[i] = chset.Proposal
		prepares[i] = chset.Prepare
		commits[i] = chset.Commit
		faults[i] = chset.Fault
		errs[i] = chset.Err
		blocks[i] = chset.Block
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		dispatch.Split(c.Proposal, proposals)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		dispatch.Split(c.Prepare, prepares)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		dispatch.Split(c.Commit, commits)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		dispatch.Split(c.Fault, faults)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		dispatch.Split(c.Block, blocks)
	}()

	wg.Wait()
}
