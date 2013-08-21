package search

// xlattice_go/search/peer_map_test.go

import (
	"fmt"
	xn "github.com/jddixon/xlattice_go/node"
	xi "github.com/jddixon/xlattice_go/nodeID"
	"github.com/jddixon/xlattice_go/rnglib"
	. "launchpad.net/gocheck"
)

var _ = fmt.Print
var _ = rnglib.MakeSimpleRNG

const (
	SHA1_LEN  = 20
	VERBOSITY = 1
)

func (s *XLSuite) makeTopAndBottom(c *C) (topPeer, bottomPeer *xn.Peer) {
	t := make([]byte, SHA1_LEN)
	for i := 0; i < SHA1_LEN; i++ {
		t[i] = byte(0xf)
	}
	top, err := xi.NewNodeID(t)
	c.Assert(err, IsNil)
	c.Assert(top, Not(IsNil))

	topPeer, err = xn.NewNewPeer("top", top)
	c.Assert(err, IsNil)
	c.Assert(topPeer, Not(IsNil))

	bottom, err := xi.NewNodeID(make([]byte, SHA1_LEN))
	c.Assert(err, IsNil)
	c.Assert(bottom, Not(IsNil))

	bottomPeer, err = xn.NewNewPeer("bottom", bottom)
	c.Assert(err, IsNil)
	c.Assert(bottomPeer, Not(IsNil))

	return topPeer, bottomPeer
}
func (s *XLSuite) makeAPeer(c *C, name string, id ...int) (peer *xn.Peer) {
	t := make([]byte, SHA1_LEN)
	for i := 0; i < len(id); i++ {
		t[i] = byte(id[i])
	}
	nodeID, err := xi.NewNodeID(t)
	c.Assert(err, IsNil)
	c.Assert(nodeID, Not(IsNil))

	peer, err = xn.NewNewPeer(name, nodeID)
	c.Assert(err, IsNil)
	c.Assert(peer, Not(IsNil))
	return
}
func (s *XLSuite) TestPeerMapTools(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_PEER_MAP_TOOLS")
	}
	threePeer := s.makeAPeer(c, "threePeer", 1, 2, 3)
	nodeID := threePeer.GetNodeID()
	value := nodeID.Value()
	c.Assert(threePeer.GetName(), Equals, "threePeer")
	c.Assert(value[0], Equals, byte(1))
	c.Assert(value[1], Equals, byte(2))
	c.Assert(value[2], Equals, byte(3))
	for i := 3; i < SHA1_LEN; i++ {
		c.Assert(value[i], Equals, byte(0))
	}

}
func (s *XLSuite) TestTopBottomMap(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_PEER_MAP")
	}

	var pm PeerMap
	c.Assert(pm.nextCol, IsNil)

	topPeer, bottomPeer := s.makeTopAndBottom(c)
	err := pm.AddToPeerMap(topPeer)
	c.Assert(err, IsNil)
	c.Assert(pm.nextCol, Not(IsNil))
	lowest := pm.nextCol
	c.Assert(lowest.peer, Not(IsNil))
    // THESE THREE TESTS ARE LOGICALLY EQUIVALENT ----------------------
	c.Assert(lowest.peer, Equals, topPeer) // succeeds ...
    c.Assert(xi.SameNodeID(lowest.peer.GetNodeID(), topPeer.GetNodeID()), 
        Equals,true)
    // XXX This fails, but it's a bug in Peer.Equal()
	// c.Assert(topPeer.Equal(lowest.peer), Equals, true) 
    // END LOGICALLY EQUIVALENT -----------------------------------------
	c.Assert(lowest.peer.GetName(), Equals, "top")

	// We expect that bottomPeer will become the lowest with its
	// higher field pointing at topPeer.
	err = pm.AddToPeerMap(bottomPeer)
	c.Assert(err, IsNil)
	lowest = pm.nextCol
	// c.Assert(bottomPeer.Equal(lowest.peer), Equals, true)   // FAILS
	c.Assert(lowest.peer.GetName(), Equals, "bottom") // XXX gets 'top'
}
func (s *XLSuite) TestShallowMap(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_SHALLOW_MAP")
	}
	var pm PeerMap
	c.Assert(pm.nextCol, IsNil)

	peer1 := s.makeAPeer(c, "peer1", 1)
	peer2 := s.makeAPeer(c, "peer2", 2)
	peer3 := s.makeAPeer(c, "peer3", 3)

	err := pm.AddToPeerMap(peer3)
	c.Assert(err, IsNil)
	c.Assert(pm.nextCol, Not(IsNil))
	lowest := pm.nextCol
	c.Assert(lowest.peer, Not(IsNil))
	c.Assert(lowest.peer, Equals, peer3)

	err = pm.AddToPeerMap(peer2)
	c.Assert(err, IsNil)
	c.Assert(pm.nextCol, Not(IsNil))
	lowest = pm.nextCol
	c.Assert(lowest.peer, Not(IsNil))
	c.Assert(lowest.peer, Equals, peer2)

	err = pm.AddToPeerMap(peer1)
	c.Assert(err, IsNil)
	c.Assert(pm.nextCol, Not(IsNil))
	lowest = pm.nextCol
	c.Assert(lowest.peer, Not(IsNil))
	c.Assert(lowest.peer, Equals, peer1)

	rootCell := pm.nextCol
	c.Assert(rootCell.byteVal, Equals, byte(1))
	c.Assert(rootCell.peer.GetName(), Equals, "peer1")
	nextCell := rootCell.thisCol
	c.Assert(nextCell, Not(IsNil)) // FAILS
	c.Assert(nextCell.byteVal, Equals, byte(2))
	nextCell = nextCell.thisCol
	c.Assert(nextCell.byteVal, Equals, byte(3))
}

func (s *XLSuite) TestDeeperMap(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_DEEPER_MAP")
	}
	var pm PeerMap
	c.Assert(pm.nextCol, IsNil)

	peer1   := s.makeAPeer(c, "peer1",  1)
	peer12  := s.makeAPeer(c, "peer12", 1, 2)
	peer123 := s.makeAPeer(c, "peer123",1, 2, 3)

	// add peer123 ================================================
	err := pm.AddToPeerMap(peer123)
	c.Assert(err, IsNil)
	c.Assert(pm.nextCol, Not(IsNil))
	lowest := pm.nextCol
	c.Assert(lowest.peer, Not(IsNil))
	c.Assert(lowest.peer, Equals, peer123)

	// now add peer12 ============================================
	err = pm.AddToPeerMap(peer12)
	c.Assert(err, IsNil)
	c.Assert(pm.nextCol, Not(IsNil))
	col0 := pm.nextCol

    DumpPeerMap(&pm, "after peer123 then peer12 added")

	// column 0 check - expect an empty cell
	c.Assert(col0.thisCol, IsNil)
	c.Assert(col0.peer, IsNil)

	// column 1 check - another empty cell
	col1 := col0.nextCol
	c.Assert(col1, Not(IsNil))
	c.Assert(col1.thisCol, IsNil)
	c.Assert(col1.peer, IsNil)

	// column 2a checks - peer12 with peer123 on the nextCol chain
	col2a := col1.nextCol
	c.Assert(col2a, Not(IsNil))
	c.Assert(col2a.nextCol, IsNil)
	c.Assert(col2a.peer, Not(IsNil))
	c.Assert(col2a.peer.GetName(), Equals, "peer12")

	// column 2b checks
	col2b := col2a.thisCol
	c.Assert(col2b, Not(IsNil))
	c.Assert(col2b.nextCol, IsNil)
	c.Assert(col2b.thisCol, IsNil)
	c.Assert(col2b.peer, Not(IsNil))
	c.Assert(col2b.peer.GetName(), Equals, "peer123")

	// now add peer1 =============================================
	err = pm.AddToPeerMap(peer1)
	c.Assert(err, IsNil)
	c.Assert(pm.nextCol, Not(IsNil))
	col0 = pm.nextCol

    DumpPeerMap(&pm, "after peer123, peer12, then peer1 added")

	// column 0 checks - an empty cell
	c.Assert(col0.peer, IsNil) // FAILS
	c.Assert(col0.thisCol, IsNil)

	// column 1a check -
	col1a := col0.nextCol
	c.Assert(col1a, Not(IsNil))
	c.Assert(col1a.nextCol, IsNil)
	c.Assert(col1a.thisCol, Not(IsNil))
	c.Assert(col1a.peer, Not(IsNil))
	c.Assert(col1a.peer, Equals, peer1)
	c.Assert(col1a.peer.GetName(), Equals, "peer1")

	// column 1b checks - another empty cell
	col1b := col1a.thisCol
	c.Assert(col1b.peer, IsNil)
	c.Assert(col1b.thisCol, IsNil)

	// column 2a checks - peer12 with peer123 on the nextCol chain
	col2a = col1b.nextCol
	c.Assert(col2a, Not(IsNil))
	c.Assert(col2a.nextCol, IsNil)
	c.Assert(col2a.peer, Not(IsNil))
	c.Assert(col2a.peer.GetName(), Equals, "peer12")

	// column 2b checks
	col2b = col2a.thisCol
	c.Assert(col2b, Not(IsNil))
	c.Assert(col2b.nextCol, IsNil)
	c.Assert(col2b.thisCol, IsNil)
	c.Assert(col2b.peer, Not(IsNil))
	c.Assert(col2b.peer.GetName(), Equals, "peer123")

	c.Assert(col0.byteVal, Equals, byte(1))
	c.Assert(col1a.byteVal, Equals, byte(0))
	c.Assert(col1b.byteVal, Equals, byte(2))
	c.Assert(col2a.byteVal, Equals, byte(0))
	c.Assert(col2b.byteVal, Equals, byte(3))

	// add 123, then 1, then 12 ----------------------------------

	// XXX STUB XXX

}

func (s *XLSuite) addAPeer(c *C, pm *PeerMap, peer *xn.Peer) {
	err := pm.AddToPeerMap(peer)
	c.Assert(err, IsNil)
	c.Assert(pm.nextCol, Not(IsNil))
}
func (s *XLSuite) findAPeer(c *C, pm *PeerMap, peer *xn.Peer) {
    nodeID := peer.GetNodeID()
    d := nodeID.Value()
    c.Assert(d, Not(IsNil))
    p := pm.FindPeer(d)
    // DEBUG
    if p == nil {
        fmt.Printf("can't find a match for %d.%d.%d.%d\n", d[0],d[1],d[2],d[3])
    }
    // END
    c.Assert(p, Not(IsNil))
    nodeIDBack := p.GetNodeID()
    c.Assert( xi.SameNodeID(nodeID, nodeIDBack), Equals, true )

}
func (s *XLSuite) TestFindPeer(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("\nTEST_FIND_PEER")
	}
	var pm PeerMap
	c.Assert(pm.nextCol, IsNil)

	peer1   := s.makeAPeer(c, "peer1",  1)
	peer12  := s.makeAPeer(c, "peer12", 1, 2)
	peer123 := s.makeAPeer(c, "peer123",1, 2, 3)
	peer4   := s.makeAPeer(c, "peer4",  4)
	peer42  := s.makeAPeer(c, "peer42", 4, 2)
	peer423 := s.makeAPeer(c, "peer423",4, 2, 3)

    // XXX BUG: if this are added in reverse order (peer123, peer12, peer1)
    // then tests succeed.  If they are added in ascending order (peer1,
    // peer12, peer123) tests fail: specifically, peer12 and its 
    // preceding nil do not appear on dumps.
    //
    // TODO: randomize order in which peers are added
    s.addAPeer(c, &pm, peer123)
    s.addAPeer(c, &pm, peer12)
    s.addAPeer(c, &pm, peer1)

    DumpPeerMap(&pm, "after adding peer1, peer12, peer123, before peer4")

    s.addAPeer(c, &pm, peer423)
    DumpPeerMap(&pm, "after adding peer423")
    s.addAPeer(c, &pm, peer42)
    DumpPeerMap(&pm, "after adding peer42")
    s.addAPeer(c, &pm, peer4)
    DumpPeerMap(&pm, "after adding peer4")

    // TODO: randomize order in which finding peers is tested
    s.findAPeer(c, &pm, peer1)
    s.findAPeer(c, &pm, peer12) 
    s.findAPeer(c, &pm, peer123)
    
    s.findAPeer(c, &pm, peer42)
    // s.findAPeer(c, &pm, peer423) // KNOWN TO FAIL (peer423 has been lost)
    // s.findAPeer(c, &pm, peer4)   // KNOWN TO FAIL, peer4 hanging off 0th cell
}

// XXX THIS DOES NOT BELONG HERE =================================
// XXX Something similar to this should be in nodeID/nodeID.go
func (s *XLSuite) TestSameNodeID(c *C) {
	peer := s.makeAPeer(c, "foo", 1, 2, 3, 4)
	id := peer.GetNodeID()
	c.Assert(xi.SameNodeID(id, id), Equals, true)
	peer2 := s.makeAPeer(c, "foo", 1, 2, 3, 4, 5)
	id2 := peer2.GetNodeID()
	c.Assert(xi.SameNodeID(id, id2), Equals, false)
}
