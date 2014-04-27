package node

// xlattice_go/node/peer_test.go

import (
	"encoding/hex"
	"fmt"
	xc "github.com/jddixon/xlattice_go/crypto"
	xo "github.com/jddixon/xlattice_go/overlay"
	"github.com/jddixon/xlattice_go/rnglib"
	xt "github.com/jddixon/xlattice_go/transport"
	. "gopkg.in/check.v1"
	"strings"
)

// available:
//		func makeNodeID(rng *rnglib.PRNG) *NodeID

func (s *XLSuite) addAString(slice *[]string, str string) *[]string {
	*slice = append(*slice, str)
	return slice
}
func (s *XLSuite) TestPeerSerialization(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_PEER_SERIALIZATION")
	}
	rng := rnglib.MakeSimpleRNG()

	// this is just a lazy way of building a peer
	name := rng.NextFileName(4)
	nid, err := makeNodeID(rng)
	c.Assert(err, Equals, nil)

	lfs := "tmp/" + hex.EncodeToString(nid.Value())
	node, err := NewNew(name, nid, lfs)
	c.Assert(err, Equals, nil)

	// harvest its keys
	ck := &node.commsKey.PublicKey
	ckSSH, err := xc.RSAPubKeyToDisk(ck)
	c.Assert(err, Equals, nil)
	sk := &node.sigKey.PublicKey
	skSSH, err := xc.RSAPubKeyToDisk(sk)
	c.Assert(err, Equals, nil)

	// the other bits necessary
	port := 1024 + rng.Intn(1024)
	addr := fmt.Sprintf("1.2.3.4:%d", port)
	ep, err := xt.NewTcpEndPoint(addr)
	c.Assert(err, Equals, nil)
	ctor, err := xt.NewTcpConnector(ep)
	c.Assert(err, Equals, nil)
	overlay, err := xo.DefaultOverlay(ep)
	c.Assert(err, Equals, nil)
	oSlice := []xo.OverlayI{overlay}
	ctorSlice := []xt.ConnectorI{ctor}
	peer, err := NewPeer(name, nid, ck, sk, oSlice, ctorSlice)
	c.Assert(err, Equals, nil)
	c.Assert(peer, Not(Equals), nil)

	// build the expected serialization

	// BaseNode
	var bns []string
	s.addAString(&bns, "peer {")
	s.addAString(&bns, fmt.Sprintf("    name: %s", name))
	s.addAString(&bns, fmt.Sprintf("    nodeID: %s", nid.String()))
	s.addAString(&bns, fmt.Sprintf("    commsPubKey: %s", ckSSH))
	s.addAString(&bns, fmt.Sprintf("    sigPubKey: %s", skSSH))
	s.addAString(&bns, fmt.Sprintf("    overlays {"))
	for i := 0; i < len(oSlice); i++ {
		s.addAString(&bns, fmt.Sprintf("        %s", oSlice[i].String()))
	}
	s.addAString(&bns, fmt.Sprintf("    }"))

	// Specific to Peer
	s.addAString(&bns, fmt.Sprintf("    connectors {"))
	for i := 0; i < len(ctorSlice); i++ {
		s.addAString(&bns, fmt.Sprintf("        %s", ctorSlice[i].String()))
	}
	s.addAString(&bns, fmt.Sprintf("    }")) // closes connectors
	s.addAString(&bns, fmt.Sprintf("}"))     // closes peer
	myVersion := strings.Join(bns, "\n")

	serialized := peer.String()
	c.Assert(serialized, Equals, myVersion)

	backAgain, rest, err := ParsePeer(serialized)
	c.Assert(err, Equals, nil)
	c.Assert(len(rest), Equals, 0)
	reserialized := backAgain.String()
	c.Assert(reserialized, Equals, serialized)

}
