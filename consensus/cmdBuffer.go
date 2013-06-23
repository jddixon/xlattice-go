package consensus

// xlattice_go/consensus/cmdBuffer.go

import (
	"container/heap"
	"fmt" // DEBUG
	"sync"
)

// Mediates between a producer and a consumer, where the producer
// sends a series of <number, command> pairs which will generally
// be unordered and may contain pairs with the same sequence
// number.  The consumer expects a stream of  commands in ascending
// numerical order, with no duplicates and no gaps.  This code
// effects that requirement by discarding duplicates, buffering up
// pairs until gaps are filled, and releasing pairs from the
// internal sorted buffer in ascending order as soon as possible.

type CmdPair struct {
	Seqn int64
	Cmd  string
}

type pairPlus struct {
	pair  *CmdPair
	index int // used by heap logic
}

type pairQ []*pairPlus

func (q pairQ) Len() int { // not in heap interface
	return len(q)
}

// implementation of the heap interface /////////////////////////////
// These functions are invoke like heap.Push(&q, &whatever)

func (q pairQ) Less(i, j int) bool {
	return q[i].pair.Seqn < q[j].pair.Seqn
}

func (q pairQ) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
	q[i].index = i // remember this is post-swap
	q[j].index = j
}

func (q *pairQ) Push(x interface{}) {
	n := len(*q)
	thePair := x.(*pairPlus) // a cast
	thePair.index = n
	*q = append(*q, thePair)
}

/////////////////////////////////////////////////////////////////////
// XXX if the length is zero, we get a panic in the heap code. //////
/////////////////////////////////////////////////////////////////////
func (q *pairQ) Pop() interface{} {
	nowQ := *q
	n := len(nowQ)
	if n == 0 {
		return nil
	}
	lastPair := nowQ[n-1] // last element
	lastPair.index = -1   // doesn't matter
	*q = nowQ[0 : n-1]
	return lastPair
}

//
type CmdBuffer struct {
	InCh     chan CmdPair
	outCh    chan CmdPair
	stopCh   chan bool
	q        pairQ
	L        sync.Mutex
	padding  [64]byte
	lastSeqn int64
	running  bool
}

func (c *CmdBuffer) Init(out chan CmdPair, stopCh chan bool, lastSeqn int64) {
	c.q = pairQ{}
	c.InCh = make(chan CmdPair, 4) // buffered

	c.outCh = out // should also be buffered
	c.stopCh = stopCh
	c.lastSeqn = lastSeqn

	c.running = false

	fmt.Printf("exiting Init: running = %v, lastSeqn = %v\n",
		c.running, c.lastSeqn)
}
func (c *CmdBuffer) Running() bool {
	fmt.Printf("enter c.Running(): running %v, lastSeqn %v\n",
		c.running, c.lastSeqn)
	// c.running is volatile
	c.L.Lock()
	whether := c.running
	c.L.Unlock()
	fmt.Printf("Running() returning %v\n", whether)
	return whether
}
func (c *CmdBuffer) Run() {
	c.running = true
	fmt.Printf("buffer.running <== %v\n", c.running)
	for {
		c.L.Lock()
		whether := c.running
		c.L.Unlock()
		if !whether {
			break
		}
		fmt.Printf("waiting for select; running is %v\n", c.running)
		select {
		case junk := <-c.stopCh:
			fmt.Printf("RECEIVED STOP %v\n", junk)
			c.L.Lock()
			c.running = false
			c.L.Unlock()
			fmt.Println("c.running has been set to false")
		case inPair := <-c.InCh: // get the next command
			seqN := inPair.Seqn
			fmt.Printf("RECEIVED PAIR %v\n", seqN)
			if seqN <= c.lastSeqn { // already sent, so discard
				fmt.Printf("    ALREADY SEEN, DISCARDING\n")
				continue
			} else if seqN == c.lastSeqn+1 {
				c.outCh <- inPair
				c.lastSeqn += 1
				fmt.Printf("    SEQN %v MATCHED LAST + 1, SENDING\n", seqN)
				for c.q.Len() > 0 {
					first := c.q[0]
					if first.pair.Seqn <= c.lastSeqn {
						fmt.Printf("        Q: DISCARDING %v, DUPE\n",
							first.pair.Seqn)
						// a duplicate, so discard
						_ = heap.Pop(&c.q).(*pairPlus)
					} else if first.pair.Seqn == c.lastSeqn+1 {
						pp := heap.Pop(&c.q).(*pairPlus)
						c.outCh <- *pp.pair
						c.lastSeqn += 1
						fmt.Printf("        Q: SENT %v\n", c.lastSeqn)
					} else {
						fmt.Printf("        Q: LEAVING %v IN Q\n",
							first.pair.Seqn)
						break
					}
				}
			} else {
				// seqN > c.lastSeqn + 1, so buffer
				fmt.Printf("    HIGH SEQN %v, SO BUFFERING\n", seqN)
				pp := &pairPlus{pair: &inPair}
				heap.Push(&c.q, pp)
			}
		}
	}
}
