package search

import (
	"fmt"
	xi "github.com/jddixon/xlattice_go/nodeID"
	xn "github.com/jddixon/xlattice_go/node"
)

var _ = fmt.Print

type PeerMap struct {
	// lowest *PeerMapCell
	PeerMapCell
}
type PeerMapCell struct {
	byteVal byte
	pred    *PeerMapCell // predecessor
	nextCol *PeerMapCell // points to a cell with same val for this byte
	thisCol *PeerMapCell // points to a cell with higher val for this col
	peer    *xn.Peer
}

// Add a Peer to the map.  This should be idempotent: adding a Peer
// that is already in the map should have no effect at all.  The cell map
// allows us to efficiently return a reference to a Peer, given its nodeID.

func (m *PeerMap) AddToPeerMap(peer *xn.Peer) (err error) {
	id := peer.GetNodeID().Value()
	byte0 := id[0]

	root := m.nextCol
	if root == nil {
		m.nextCol = &PeerMapCell{byteVal: byte0, pred: &m.PeerMapCell, peer: peer}
	} else {
		err = root.AddAtCell(0, peer, id)
	}
	return
}

// depth is that of cell, with the root cell at column 0, and also the
// index into the id slice.

func (p *PeerMapCell) AddAtCell(depth int, peer *xn.Peer, id []byte) (err error) {
	curByte := id[depth]
	if curByte < p.byteVal {
		// DEBUG
		fmt.Printf("lower, adding %s as lowest, curByte is %d\n",
			peer.GetName(), curByte)
		// END
		p.pred.nextCol = &PeerMapCell{
			byteVal: curByte, pred: p.pred, thisCol: p, peer: peer}

	} else if curByte == p.byteVal {
		// DEBUG
		fmt.Printf("%s has match, curByte is %d\n", peer.GetName(), curByte)
		// END
		// XXX If p.nextCol is nil, we know exactly what to do
		if p.nextCol == nil {
			peer2 := p.peer
			var id2 []byte
			if peer2 != nil {
				id2 = peer2.GetNodeID().Value()
			}
			p.peer = nil

			depth++
			nextByte := id[depth]
			nextByte2 := id2[depth]
			curCell := p
			for nextByte == nextByte2 {
				nextCell := &PeerMapCell{byteVal: nextByte, pred: curCell}
				curCell.nextCol = nextCell
				curCell = nextCell
				depth++
				fmt.Printf("depth := %d\n", depth)
				nextByte = id[depth]
				nextByte2 = id2[depth]
			}
			lastCell := &PeerMapCell{byteVal: nextByte, peer: peer}
			lastCell2 := &PeerMapCell{byteVal: nextByte2, peer: peer2}
			if nextByte < nextByte2 {
				curCell.nextCol = lastCell
				lastCell.pred = curCell
				lastCell.thisCol = lastCell2
				lastCell2.pred = lastCell
			} else {
				curCell.nextCol = lastCell2
				lastCell2.pred = curCell
				lastCell2.thisCol = lastCell
				lastCell.pred = lastCell2
			}
			fmt.Printf("END OF CHAINLET, peers %s and %s\n",
				peer.GetName(), peer2.GetName()) // GEEP
		} else {
			// we are guaranteed that this is not nil
			curCell := p.nextCol
			depth++
			var nextB byte
			// skip any cells with matching values
			for nextB = id[depth]; nextB == curCell.byteVal; nextB = id[depth] {
				if curCell.nextCol == nil {
					break
				}
				curCell = curCell.nextCol
				depth++
			}
			newCell := &PeerMapCell{byteVal: nextB, peer: peer}
			if nextB < curCell.byteVal {
				// case C
				pred := curCell.pred
				pred.nextCol = newCell
				newCell.thisCol = curCell
				curCell.pred = newCell
			} else {
				// case D
				curCell.thisCol = newCell
				newCell.pred = curCell
			} // FOO
		}

	} else { // curByte > p.byteVal
		// DEBUG
		fmt.Printf("CALLING AddThisCol: adding %s as higher, curByte is %d\n",
			peer.GetName(), curByte)
		p.AddThisCol(id, depth, peer)      
		// END
	}

	return
}

// The nodeID of the peer being added has the same value for the byte at
// this depth.  id and peer represent the new peer being added, where id
// is the byte slice for its nodeID and peer is a reference to that.
// id2 and peer2 represent any pre-existing value.
func (p *PeerMapCell) AddMatchingToDepth(depth int,
	id, id2 []byte, peer, peer2 *xn.Peer) (err error) {

    // XXX SHOULD NEVER SEE THIS, but do see it
    fmt.Printf("ADD_MATCHING_TO_DEPTH, depth %d, peer %s\n",
        depth, peer.GetName())

	// The byte string id has matched the chain up to this point.
	// We examine the next byte in id and the byte value for the next
	// cell in the chain.
	depth += 1
	nextByte := id[depth]

	if p.nextCol == nil {
		if peer2 == nil {
			p.nextCol = &PeerMapCell{nextByte, p, nil, nil, peer}
		} else {
			nextByte2 := id2[depth]
			if nextByte == nextByte2 {
				fmt.Printf("Case 1b1, %s\n", peer.GetName())
				p.nextCol = &PeerMapCell{nextByte, p, nil, nil, nil}
				p.nextCol.AddMatchingToDepth(depth, id, id2, peer, peer2)
			} else {
				nextCell := &PeerMapCell{nextByte, nil, nil, nil, peer}
				nextCell2 := &PeerMapCell{nextByte2, nil, nil, nil, peer2}
				if nextByte < nextByte2 {
					fmt.Printf("Case 1b2a, %s\n", peer.GetName())
					nextCell.thisCol = nextCell2
					p.nextCol = nextCell
					nextCell.pred = p
					nextCell2.pred = nextCell
				} else {
					fmt.Printf("Case 1b2b, %s\n", peer.GetName())
					nextCell2.thisCol = nextCell
					p.nextCol = nextCell2
					nextCell2.pred = p
					nextCell.pred = nextCell2
				}
			}
		}
	} else {
		// XXX doesn't handle peer2
		curCell := p.nextCol
		if nextByte < curCell.byteVal {
			// DEBUG
			var nextPeerStr string
			if curCell.peer == nil {
				nextPeerStr = "<nil>"
			} else {
				nextPeerStr = curCell.peer.GetName()
			}
			fmt.Printf("CASE 2a: %s => %s\n", peer.GetName(), nextPeerStr)
			// END
			p.peer = peer
			p.nextCol = &PeerMapCell{nextByte, p, nil, curCell, peer2}

		} else if nextByte == curCell.byteVal {
			fmt.Printf("CASE 2b, %s\n", peer.GetName()) // DEBUG
			peer2 := curCell.peer
			var id2 []byte
			if peer2 != nil {
				id2 = peer2.GetNodeID().Value()
			}
			curCell.peer = nil
			curCell.AddMatchingToDepth(depth, id, id2, peer, peer2)

		} else { // nextByte > curCell.byteVal
			fmt.Printf("CASE 2c, %s\n", peer.GetName()) // DEBUG
			curCell.AddThisCol(id, 0, peer)
		}
	}
	return
} // GEEP

func (p *PeerMapCell) AddThisCol(id []byte, depth int, peer *xn.Peer) (
	err error) {

	nextByte := id[depth]
	if p.thisCol == nil {
		p.thisCol = &PeerMapCell{nextByte, p, nil, nil, peer}
	} else {
		curHigher := p.thisCol // current higher value
		if nextByte < curHigher.byteVal {
			curHigher = &PeerMapCell{nextByte, p, nil, curHigher, peer}
		} else if nextByte == curHigher.byteVal {
			// WORKING HERE
			curHigher.AddMatchingToDepth(0, id, nil, peer, nil)
		} else { // nextByte > curHigher.byteVal
			curHigher.AddThisCol(id, 0, peer)
		}
	}
	return
}

// At any particular depth, a match is possible only if (a) peer for the
// cell is not nil and (b) we have a byte-wise match

func (m *PeerMap) FindPeer(id []byte) (peer *xn.Peer) {
	mapCell := m.nextCol
    fmt.Printf("FindPeer for %d.%d.%d.%d\n", id[0], id[1], id[2], id[3])

	for depth := 0; depth < len(id); depth++ {
		myVal := id[depth]
        fmt.Printf("    FindPeer: depth %d, val %d\n", depth, myVal)
        if mapCell == nil {
            fmt.Printf("    Internal error: nil mapCell at depth %d\n", depth)
            return nil
        }
		if myVal > mapCell.byteVal {
            for mapCell.thisCol != nil {
                mapCell = mapCell.thisCol
                if myVal == mapCell.byteVal {
                    goto maybeEqual
                } else if myVal > mapCell.byteVal {
                    return nil
                }
            }
            fmt.Printf("    depth %d, %d < %d returning NIL\n", 
                    depth, myVal, mapCell.byteVal)
		    return nil
		}
        maybeEqual:
        if myVal == mapCell.byteVal {
			if mapCell.nextCol == nil {
                myNodeID, err := xi.NewNodeID(id)
                if err != nil {
                    fmt.Printf("    FindPeer: NewNodeID returns %v", err)
                    return nil
                }
                if mapCell.peer != nil {
                    fmt.Printf("    peer is %s\n", mapCell.peer.GetName())
                    if xi.SameNodeID(myNodeID, mapCell.peer.GetNodeID()) {
                        fmt.Printf("    *MATCH* on %s\n", 
                            mapCell.peer.GetName())
                        return mapCell.peer
                    }
                }
			} else {
                fmt.Printf("    RIGHT, so depth := %d\n", depth + 1)
				mapCell = mapCell.nextCol
				continue
			}

		} else {
			// myVal < mapCell.byteVal
            fmt.Printf("    myval %d > cell's %d\n", myVal, mapCell.byteVal)
            return nil
		}
	}
	return
}
