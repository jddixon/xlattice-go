package node

// xlattice_go/node/bn_map_test.go

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	xi "github.com/jddixon/xlattice_go/nodeID"
	xr "github.com/jddixon/xlattice_go/rnglib"
	xt "github.com/jddixon/xlattice_go/transport"
	. "launchpad.net/gocheck"
)

var _ = fmt.Print

func (s *XLSuite) makePubKey(c *C) *rsa.PublicKey {
	key, err := rsa.GenerateKey(rand.Reader, 512) // 512 because cheaper
	c.Assert(err, IsNil)
	return &key.PublicKey
}

func (s *XLSuite) makePeerGivenID(c *C, rng *xr.PRNG, name string,
	id []byte) (member *Peer) {

	nodeID, err := xi.New(id)
	c.Assert(err, IsNil)

	ck := s.makePubKey(c)
	sk := s.makePubKey(c)
	port := 1024 + rng.Intn(256*256-1024)
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	myEnd, err := xt.NewTcpEndPoint(addr)
	c.Assert(err, IsNil)
	ctor, err := xt.NewTcpConnector(myEnd)
	c.Assert(err, IsNil)
	ctors := []xt.ConnectorI{ctor}

	member, err = NewPeer(name, nodeID, ck, sk, nil, ctors)
	c.Assert(err, IsNil)
	return
}
func (s *XLSuite) makeTopAndBottomBNI(c *C, rng *xr.PRNG) (
	topBNI, bottomBNI BaseNodeI) {

	// top contains  a slice of SHA1_LEN 0xFF as NodeID
	t := make([]byte, SHA1_LEN)
	for i := 0; i < SHA1_LEN; i++ {
		t[i] = byte(0xff)
	}
	topBNI = s.makePeerGivenID(c, rng, "top", t)

	// bottom contains  a slice of SHA1_LEN zeroes as NodeID
	b := make([]byte, SHA1_LEN)
	bottomBNI = s.makePeerGivenID(c, rng, "bottom", b)

	return topBNI, bottomBNI
}

// Create a randomish Peer for use as a BNI, assigning it a
// nodeID based upon the variable-length list of ints passed

func (s *XLSuite) makeABNI(c *C, rng *xr.PRNG, name string, id ...int) (
	bni BaseNodeI) {

	t := make([]byte, SHA1_LEN)
	for i := 0; i < len(id); i++ {
		t[i] = byte(id[i])
	}
	bni = s.makePeerGivenID(c, rng, name, t)
	c.Assert(bni, Not(IsNil))
	return
}
func (s *XLSuite) TestBNIMapTools(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_BNI_MAP_TOOLS")
	}
	rng := xr.MakeSimpleRNG()
	threeBaseNode := s.makeABNI(c, rng, "threeBaseNode", 1, 2, 3)
	nodeID := threeBaseNode.GetNodeID()
	value := nodeID.Value()
	c.Assert(threeBaseNode.GetName(), Equals, "threeBaseNode")
	c.Assert(value[0], Equals, byte(1))
	c.Assert(value[1], Equals, byte(2))
	c.Assert(value[2], Equals, byte(3))
	for i := 3; i < SHA1_LEN; i++ {
		c.Assert(value[i], Equals, byte(0))
	}

}
func (s *XLSuite) TestTopBottomBNMap(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_TOP_BOTTOM_MAP")
	}

	var pm BNIMap
	c.Assert(pm.NextCol, IsNil)

	rng := xr.MakeSimpleRNG()
	topBNI, bottomBNI := s.makeTopAndBottomBNI(c, rng)
	err := pm.AddToBNIMap(topBNI)
	c.Assert(err, IsNil)
	c.Assert(pm.NextCol, Not(IsNil))
	lowest := pm.NextCol
	c.Assert(lowest.CellNode, Not(IsNil))

	// THESE THREE TESTS ARE LOGICALLY EQUIVALENT ----------------------
	c.Assert(lowest.CellNode, Equals, topBNI)
	c.Assert(xi.SameNodeID(lowest.CellNode.GetNodeID(), topBNI.GetNodeID()),
		Equals, true)
	c.Assert(topBNI.Equal(lowest.CellNode), Equals, true)
	// END LOGICALLY EQUIVALENT -----------------------------------------

	c.Assert(lowest.CellNode.GetName(), Equals, "top")

	// We expect that bottomBNI will become the lowest with its
	// higher field pointing at topBNI.
	err = pm.AddToBNIMap(bottomBNI)
	c.Assert(err, IsNil)
	lowest = pm.NextCol

	c.Assert(bottomBNI.Equal(lowest.CellNode), Equals, true)
	c.Assert(lowest.CellNode.GetName(), Equals, "bottom")
}
func (s *XLSuite) TestShallowBNMap(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_SHALLOW_MAP")
	}
	var pm BNIMap
	c.Assert(pm.NextCol, IsNil)

	rng := xr.MakeSimpleRNG()
	baseNode1 := s.makeABNI(c, rng, "baseNode1", 1)
	baseNode2 := s.makeABNI(c, rng, "baseNode2", 2)
	baseNode3 := s.makeABNI(c, rng, "baseNode3", 3)

	// ADD BNI 3 ---------------------------------------------------
	err := pm.AddToBNIMap(baseNode3)
	c.Assert(err, IsNil)
	c.Assert(pm.NextCol, Not(IsNil))
	cell3 := pm.NextCol
	c.Assert(cell3.ByteVal, Equals, byte(3))
	c.Assert(cell3.CellNode, Not(IsNil))
	c.Assert(cell3.CellNode.GetName(), Equals, baseNode3.GetName())

	// INSERT BNI 2 ------------------------------------------------
	err = pm.AddToBNIMap(baseNode2)
	c.Assert(err, IsNil)
	c.Assert(pm.NextCol, Not(IsNil))
	cell2 := pm.NextCol
	c.Assert(cell2.ByteVal, Equals, byte(2)) // FAILS, is 3
	c.Assert(cell2.ThisCol.ByteVal, Equals, byte(3))
	c.Assert(cell2.CellNode, Not(IsNil))
	c.Assert(cell2.CellNode.GetName(), Equals, baseNode2.GetName()) // FAILS

	// DumpBNIMap(&pm, "dump of shallow map, baseNodes 3 and 2")

	// INSERT BNI 1 ------------------------------------------------
	err = pm.AddToBNIMap(baseNode1)
	c.Assert(err, IsNil)
	c.Assert(pm.NextCol, Not(IsNil))
	cell1 := pm.NextCol
	c.Assert(cell1.ByteVal, Equals, byte(1))
	c.Assert(cell1.CellNode, Not(IsNil))
	c.Assert(cell1.CellNode.GetName(), Equals, baseNode1.GetName())

	// DumpBNIMap(&pm, "dump of shallow map, baseNodes 3,2,1")

	rootCell := pm.NextCol
	c.Assert(rootCell.ByteVal, Equals, byte(1))
	c.Assert(rootCell.CellNode.GetName(), Equals, "baseNode1")
	nextCell := rootCell.ThisCol
	c.Assert(nextCell, Not(IsNil))
	c.Assert(nextCell.ByteVal, Equals, byte(2))
	nextCell = nextCell.ThisCol
	c.Assert(nextCell.ByteVal, Equals, byte(3))
}

func (s *XLSuite) TestDeeperBNMap(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_DEEPER_MAP")
	}
	var pm BNIMap
	c.Assert(pm.NextCol, IsNil)

	rng := xr.MakeSimpleRNG()
	baseNode1 := s.makeABNI(c, rng, "baseNode1", 1)
	baseNode12 := s.makeABNI(c, rng, "baseNode12", 1, 2)
	baseNode123 := s.makeABNI(c, rng, "baseNode123", 1, 2, 3)

	// add baseNode123 ================================================
	err := pm.AddToBNIMap(baseNode123)
	c.Assert(err, IsNil)
	c.Assert(pm.NextCol, Not(IsNil))
	lowest := pm.NextCol
	c.Assert(lowest.CellNode, Not(IsNil))
	c.Assert(lowest.CellNode, Equals, baseNode123)

	// now add baseNode12 ============================================
	err = pm.AddToBNIMap(baseNode12)
	c.Assert(err, IsNil)
	c.Assert(pm.NextCol, Not(IsNil))
	col0 := pm.NextCol

	// DumpBNIMap(&pm, "after baseNode123 then baseNode12 added")

	// column 0 check - expect an empty cell
	c.Assert(col0.ThisCol, IsNil)
	c.Assert(col0.CellNode, IsNil)

	// column 1 check - another empty cell
	col1 := col0.NextCol
	c.Assert(col1, Not(IsNil))
	c.Assert(col1.ThisCol, IsNil)
	c.Assert(col1.CellNode, IsNil)

	// column 2a checks - baseNode12 with baseNode123 on the NextCol chain
	col2a := col1.NextCol
	c.Assert(col2a, Not(IsNil))
	c.Assert(col2a.NextCol, IsNil)
	c.Assert(col2a.CellNode, Not(IsNil))
	c.Assert(col2a.CellNode.GetName(), Equals, "baseNode12")

	// column 2b checks
	col2b := col2a.ThisCol
	c.Assert(col2b, Not(IsNil))
	c.Assert(col2b.NextCol, IsNil)
	c.Assert(col2b.ThisCol, IsNil)
	c.Assert(col2b.CellNode, Not(IsNil))
	c.Assert(col2b.CellNode.GetName(), Equals, "baseNode123")

	// now add baseNode1 =============================================
	err = pm.AddToBNIMap(baseNode1)
	c.Assert(err, IsNil)
	c.Assert(pm.NextCol, Not(IsNil))
	col0 = pm.NextCol

	// DumpBNIMap(&pm, "after baseNode123, baseNode12, then baseNode1 added")

	// column 0 checks - an empty cell
	c.Assert(col0.CellNode, IsNil)
	c.Assert(col0.ThisCol, IsNil)

	// column 1a check -
	col1a := col0.NextCol
	c.Assert(col1a, Not(IsNil))
	c.Assert(col1a.NextCol, IsNil)
	c.Assert(col1a.ThisCol, Not(IsNil))
	c.Assert(col1a.CellNode, Not(IsNil))
	c.Assert(col1a.CellNode, Equals, baseNode1)
	c.Assert(col1a.CellNode.GetName(), Equals, "baseNode1")

	// column 1b checks - another empty cell
	col1b := col1a.ThisCol
	c.Assert(col1b.CellNode, IsNil)
	c.Assert(col1b.ThisCol, IsNil)

	// column 2a checks - baseNode12 with baseNode123 on the NextCol chain
	col2a = col1b.NextCol
	c.Assert(col2a, Not(IsNil))
	c.Assert(col2a.NextCol, IsNil)
	c.Assert(col2a.CellNode, Not(IsNil))
	c.Assert(col2a.CellNode.GetName(), Equals, "baseNode12")

	// column 2b checks
	col2b = col2a.ThisCol
	c.Assert(col2b, Not(IsNil))
	c.Assert(col2b.NextCol, IsNil)
	c.Assert(col2b.ThisCol, IsNil)
	c.Assert(col2b.CellNode, Not(IsNil))
	c.Assert(col2b.CellNode.GetName(), Equals, "baseNode123")

	c.Assert(col0.ByteVal, Equals, byte(1))
	c.Assert(col1a.ByteVal, Equals, byte(0))
	c.Assert(col1b.ByteVal, Equals, byte(2))
	c.Assert(col2a.ByteVal, Equals, byte(0))
	c.Assert(col2b.ByteVal, Equals, byte(3))

	// add 123, then 1, then 12 ----------------------------------

	// XXX STUB XXX

}

func (s *XLSuite) addABNI(c *C, pm *BNIMap, baseNode BaseNodeI) {
	err := pm.AddToBNIMap(baseNode)
	c.Assert(err, IsNil)
	c.Assert(pm.NextCol, Not(IsNil))
}
func (s *XLSuite) findABaseNode(c *C, pm *BNIMap, baseNode BaseNodeI) {
	nodeID := baseNode.GetNodeID()
	d := nodeID.Value()
	c.Assert(d, Not(IsNil))
	p := pm.FindBNI(d)
	// DEBUG
	if p == nil {
		fmt.Printf("can't find a match for %d.%d.%d.%d\n", d[0], d[1], d[2], d[3])
	}
	// END
	c.Assert(p, Not(IsNil))
	nodeIDBack := p.GetNodeID()
	c.Assert(xi.SameNodeID(nodeID, nodeIDBack), Equals, true)

}
func (s *XLSuite) TestFindFlatBaseNodes(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_FIND_FLAT_BNIS")
	}
	var pm BNIMap
	c.Assert(pm.NextCol, IsNil)

	rng := xr.MakeSimpleRNG()
	baseNode1 := s.makeABNI(c, rng, "baseNode1", 1)
	baseNode2 := s.makeABNI(c, rng, "baseNode2", 2)
	baseNode4 := s.makeABNI(c, rng, "baseNode4", 4)
	baseNode5 := s.makeABNI(c, rng, "baseNode5", 5)
	baseNode6 := s.makeABNI(c, rng, "baseNode6", 6)

	// TODO: randomize order in which baseNodes are added

	// ADD 1 AND THEN 5 ---------------------------------------------
	s.addABNI(c, &pm, baseNode1)
	s.addABNI(c, &pm, baseNode5)

	cell1 := pm.NextCol
	c.Assert(cell1.Pred, Equals, &pm.BNIMapCell)
	c.Assert(cell1.NextCol, IsNil)

	cell5 := cell1.ThisCol
	c.Assert(cell5, Not(IsNil)) // FAILS
	c.Assert(cell5.ByteVal, Equals, byte(5))
	c.Assert(cell5.Pred, Equals, cell1)
	c.Assert(cell5.NextCol, IsNil)
	c.Assert(cell5.ThisCol, IsNil)

	// INSERT 4 -----------------------------------------------------
	s.addABNI(c, &pm, baseNode4)

	cell4 := cell1.ThisCol
	c.Assert(cell4.ByteVal, Equals, byte(4))
	c.Assert(cell4.Pred, Equals, cell1)
	c.Assert(cell4.NextCol, IsNil)
	c.Assert(cell4.ThisCol, Equals, cell5)
	c.Assert(cell5.Pred, Equals, cell4)

	// ADD 6 --------------------------------------------------------
	s.addABNI(c, &pm, baseNode6)

	cell6 := cell5.ThisCol
	c.Assert(cell6.ByteVal, Equals, byte(6))
	c.Assert(cell6.Pred, Equals, cell5)
	c.Assert(cell6.NextCol, IsNil)
	c.Assert(cell6.ThisCol, IsNil)

	// INSERT 2 -----------------------------------------------------
	s.addABNI(c, &pm, baseNode2)

	cell2 := cell1.ThisCol
	c.Assert(cell2.ByteVal, Equals, byte(2))
	c.Assert(cell2.Pred, Equals, cell1)
	c.Assert(cell2.NextCol, IsNil)
	c.Assert(cell2.ThisCol, Equals, cell4)
	c.Assert(cell4.Pred, Equals, cell2)

	// DumpBNIMap(&pm, "after adding baseNode2")

	// TODO: randomize order in which finding baseNodes is tested
	s.findABaseNode(c, &pm, baseNode1)
	s.findABaseNode(c, &pm, baseNode2)
	s.findABaseNode(c, &pm, baseNode4)
	s.findABaseNode(c, &pm, baseNode5)
	s.findABaseNode(c, &pm, baseNode6)
}
func (s *XLSuite) TestFindBNI(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_FIND_BNI")
	}
	var pm BNIMap
	c.Assert(pm.NextCol, IsNil)

	rng := xr.MakeSimpleRNG()
	baseNode0123 := s.makeABNI(c, rng, "baseNode0123", 0, 1, 2, 3)
	baseNode1 := s.makeABNI(c, rng, "baseNode1", 1)
	baseNode12 := s.makeABNI(c, rng, "baseNode12", 1, 2)
	baseNode123 := s.makeABNI(c, rng, "baseNode123", 1, 2, 3)
	baseNode4 := s.makeABNI(c, rng, "baseNode4", 4)
	baseNode42 := s.makeABNI(c, rng, "baseNode42", 4, 2)
	baseNode423 := s.makeABNI(c, rng, "baseNode423", 4, 2, 3)
	// baseNode5 := s.makeABNI(c, rng, "baseNode5", 5)
	baseNode6 := s.makeABNI(c, rng, "baseNode6", 6)
	baseNode62 := s.makeABNI(c, rng, "baseNode62", 6, 2)
	baseNode623 := s.makeABNI(c, rng, "baseNode623", 6, 2, 3)

	// TODO: randomize order in which baseNodes are added
	s.addABNI(c, &pm, baseNode123)
	s.addABNI(c, &pm, baseNode12)
	s.addABNI(c, &pm, baseNode1)
	//DumpBNIMap(&pm, "after adding baseNode1, baseNode12, baseNode123, before baseNode4")

	// s.addABNI(c, &pm, baseNode5)
	// DumpBNIMap(&pm, "after adding baseNode5")

	s.addABNI(c, &pm, baseNode4)
	s.addABNI(c, &pm, baseNode42)
	s.addABNI(c, &pm, baseNode423)
	// DumpBNIMap(&pm, "after adding baseNode4, baseNode42, baseNode423")

	s.addABNI(c, &pm, baseNode6)
	// DumpBNIMap(&pm, "after adding baseNode6")
	s.addABNI(c, &pm, baseNode623)
	//DumpBNIMap(&pm, "after adding baseNode623")
	s.addABNI(c, &pm, baseNode62)
	//DumpBNIMap(&pm, "after adding baseNode62")

	s.addABNI(c, &pm, baseNode0123)
	//DumpBNIMap(&pm, "after adding baseNode0123")

	// adding duplicates should have no effect
	s.addABNI(c, &pm, baseNode4)
	s.addABNI(c, &pm, baseNode42)
	s.addABNI(c, &pm, baseNode423)

	// TODO: randomize order in which finding baseNodes is tested
	s.findABaseNode(c, &pm, baseNode0123) // XXX

	s.findABaseNode(c, &pm, baseNode1)
	s.findABaseNode(c, &pm, baseNode12)
	s.findABaseNode(c, &pm, baseNode123)

	s.findABaseNode(c, &pm, baseNode4)
	s.findABaseNode(c, &pm, baseNode42)
	s.findABaseNode(c, &pm, baseNode423)

	s.findABaseNode(c, &pm, baseNode6)
	s.findABaseNode(c, &pm, baseNode62)
	s.findABaseNode(c, &pm, baseNode623)
}