package consensus

import (
	"container/heap"
	"fmt"
	"github.com/jddixon/xlattice_go/rnglib"
	. "launchpad.net/gocheck"
	"os"
	"testing"
	"time"
)

// gocheck tie-in /////////////////////
func Test(t *testing.T) { TestingT(t) }

type XLSuite struct{}

var _ = Suite(&XLSuite{})

// end gocheck setup //////////////////

func (s *XLSuite) makeSimpleRNG() *rnglib.PRNG {
	t := time.Now().Unix()
	rng := rnglib.NewSimpleRNG(t)
	return rng
}

func (s *XLSuite) TestCmdQ(c *C) {
	q := pairQ{}
	heap.Init(&q)
	c.Assert(q.Len(), Equals, 0)

	pair0 := NumberedCmd{Seqn: 42, Cmd: "foo"}
	pair1 := NumberedCmd{Seqn: 1, Cmd: "bar"}
	pair2 := NumberedCmd{Seqn: 99, Cmd: "baz"}

	pp0 := cmdPlus{pair: &pair0}
	pp1 := cmdPlus{pair: &pair1}
	pp2 := cmdPlus{pair: &pair2}

	heap.Push(&q, &pp0)
	heap.Push(&q, &pp1)
	heap.Push(&q, &pp2)
	c.Assert(q.Len(), Equals, 3)

	out := heap.Pop(&q).(*cmdPlus)
	c.Assert(out.pair.Seqn, Equals, int64(1))
	c.Assert(out.pair.Cmd, Equals, "bar")

	out = heap.Pop(&q).(*cmdPlus)
	c.Assert(out.pair.Seqn, Equals, int64(42))
	c.Assert(out.pair.Cmd, Equals, "foo")

	out = heap.Pop(&q).(*cmdPlus)
	c.Assert(out.pair.Seqn, Equals, int64(99))
	c.Assert(out.pair.Cmd, Equals, "baz")

	c.Assert(q.Len(), Equals, 0)
	// XXX THIS PANICS - so if popping from a heap, always check
	// the length first.
	//zzz		:= heap.Pop(&q)
	// c.Assert(zzz, Equals, nil)
}
func (s *XLSuite) doTestCmdBufferI(c *C, p CmdBufferI, logging bool) {
	var pairMap = map[int64]string{
		1: "foo",
		2: "bar",
		3: "baz",
		4: "it's me!",
		5: "my chance will come soon",
		6: "it's my turn now",
		7: "wait for me",
	}
	// we send the messages somewhat out of order, with some duplicates
	order := [...]int{1, 2, 3, 6, 3, 2, 6, 5, 4, 1, 7}
	var out = make(chan NumberedCmd, len(order)+1) // must exceed len(order)
	var stopCh = make(chan bool, 1)
	var logFile string
	if logging {
		logFile = "tmp/logFile"
	}
	p.Init(out, stopCh, 0, 4, logFile) // 4 is bufSize
	if logging {
		_, err := os.Stat(logFile)
		c.Assert(err, Equals, nil)
	}
	c.Assert(p.Running(), Equals, false)

	fmt.Println("  starting p loop ...")
	// XXX Run() can return an error, which must be nil
	go p.Run()
	for !p.Running() {
		time.Sleep(time.Millisecond)
	}
	c.Assert(p.Running(), Equals, true)

	for n := 0; n < len(order); n++ {
		which := order[n]
		cmd := pairMap[int64(which)]
		pair := NumberedCmd{Seqn: int64(which), Cmd: cmd}
		// DEBUG
		fmt.Printf("sending %d : %s\n", order[n], cmd)
		// END
		p.InCh() <- pair
	}

	var results [7]NumberedCmd
	for n := 0; n < 7; n++ {
		results[n] = <-out
		c.Assert(results[n].Seqn, Equals, int64(n+1))
	}

	c.Assert(p.Running(), Equals, true)
	stopCh <- true
	time.Sleep(time.Microsecond)
	c.Assert(p.Running(), Equals, false)
}
func (s *XLSuite) TestCmdBuffer(c *C) {
	var buf CmdBuffer
	fmt.Println("running test without logging")
	s.doTestCmdBufferI(c, &buf, false)
	fmt.Println("running test -with- logging")
	s.doTestCmdBufferI(c, &buf, true)
}
