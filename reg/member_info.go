package reg

// xlattice_go/reg/member_info.go

// This file contains functions and structures used to describe
// and manage the cluster data managed by the registry.

import (
	"crypto/rsa"
	"encoding/hex"
	"fmt"
	xc "github.com/jddixon/xlattice_go/crypto"
	xn "github.com/jddixon/xlattice_go/node"
	xi "github.com/jddixon/xlattice_go/nodeID"
	"strings"
)

var _ = fmt.Print

type MemberInfo struct {
	Attrs       uint64   //  bit flags are defined in const.go
	MyEnds      []string // serialized EndPointI
	xn.BaseNode          // name and ID must be unique
}

func NewMemberInfo(name string, id *xi.NodeID,
	commsPubKey, sigPubKey *rsa.PublicKey, attrs uint64, myEnds []string) (
	member *MemberInfo, err error) {

	// all attrs bits are zero by default

	base, err := xn.NewBaseNode(name, id, commsPubKey, sigPubKey, nil)
	if err == nil {
		member = &MemberInfo{
			Attrs:    attrs,
			MyEnds:   myEnds,
			BaseNode: *base,
		}
	}
	return
}

// Create the MemberInfo corresponding to the token passed.

func NewMemberInfoFromToken(token *XLRegMsg_Token) (
	m *MemberInfo, err error) {

	var nodeID *xi.NodeID
	if token == nil {
		err = NilToken
	} else {
		nodeID, err = xi.New(token.GetID())
		if err == nil {
			ck, err := xc.RSAPubKeyFromWire(token.GetCommsKey())
			if err == nil {
				sk, err := xc.RSAPubKeyFromWire(token.GetSigKey())
				if err == nil {
					m, err = NewMemberInfo(token.GetName(), nodeID,
						ck, sk, token.GetAttrs(), token.GetMyEnds())
				}
			}
		}
	}
	return
}

// Return the XLRegMsg_Token corresponding to this cluster member.
func (cm *MemberInfo) Token() (token *XLRegMsg_Token, err error) {

	var ckBytes, skBytes []byte

	ck := cm.GetCommsPublicKey()
	// DEBUG
	if ck == nil {
		fmt.Printf("MemberInfo.Token: %s commsPubKey is nil\n", cm.GetName())
	}
	// END
	ckBytes, err = xc.RSAPubKeyToWire(ck)
	if err == nil {
		skBytes, err = xc.RSAPubKeyToWire(cm.GetSigPublicKey())
		if err == nil {
			name := cm.GetName()
			token = &XLRegMsg_Token{
				Name:     &name,
				Attrs:    &cm.Attrs,
				ID:       cm.GetNodeID().Value(),
				CommsKey: ckBytes,
				SigKey:   skBytes,
				MyEnds:   cm.MyEnds,
			}
		}
	}
	return
}

// EQUAL ////////////////////////////////////////////////////////////

func (cm *MemberInfo) Equal(any interface{}) bool {

	if any == cm {
		return true
	}
	if any == nil {
		return false
	}
	switch v := any.(type) {
	case *MemberInfo:
		_ = v
	default:
		return false
	}
	other := any.(*MemberInfo) // type assertion
	if cm.Attrs != other.Attrs {
		return false
	}
	if cm.MyEnds == nil {
		if other.MyEnds != nil {
			return false
		}
	} else {
		if other.MyEnds == nil {
			return false
		}
		if len(cm.MyEnds) != len(other.MyEnds) {
			return false
		}
		for i := 0; i < len(cm.MyEnds); i++ {
			if cm.MyEnds[i] != other.MyEnds[i] {
				return false
			}
		}
	}
	// WARNING: panics without the ampersand !
	return cm.BaseNode.Equal(&other.BaseNode)
}

// SERIALIZATION ////////////////////////////////////////////////////

func (cm *MemberInfo) Strings() (ss []string) {
	ss = []string{"memberInfo {"}
	bns := cm.BaseNode.Strings()
	for i := 0; i < len(bns); i++ {
		ss = append(ss, "    "+bns[i])
	}
	ss = append(ss, fmt.Sprintf("    attrs: 0x%016x", cm.Attrs))
	ss = append(ss, "    endPoints {")
	for i := 0; i < len(cm.MyEnds); i++ {
		ss = append(ss, "        "+cm.MyEnds[i])
	}
	ss = append(ss, "    }")
	ss = append(ss, "}")
	return
}

func (cm *MemberInfo) String() string {
	return strings.Join(cm.Strings(), "\n")
}
func collectAttrs(cm *MemberInfo, ss []string) (rest []string, err error) {
	rest = ss
	line := xn.NextNBLine(&rest) // trims
	// attrs line looks like "attrs: 0xHHHH..." where H is a hex digit
	if strings.HasPrefix(line, "attrs: 0x") {
		var val []byte
		var attrs uint64
		line := line[9:]
		val, err = hex.DecodeString(line)
		if err == nil {
			if len(val) != 8 {
				err = WrongNumberOfBytesInAttrs
			} else {
				for i := 0; i < 8; i++ {
					// assume little-endian ; but printf has put
					// high order bytes first - ie, it's big-endian
					attrs |= uint64(val[i]) << uint(8*(7-i))
				}
				cm.Attrs = attrs
			}
		}
	} else {
		err = BadAttrsLine
	}
	return
}
func collectMyEnds(cm *MemberInfo, ss []string) (rest []string, err error) {
	rest = ss
	line := xn.NextNBLine(&rest)
	if line == "endPoints {" {
		for {
			line = strings.TrimSpace(rest[0]) // peek
			if line == "}" {
				break
			}
			line = xn.NextNBLine(&rest)
			// XXX NO CHECK THAT THIS IS A VALID ENDPOINT
			cm.MyEnds = append(cm.MyEnds, line)
		}
		line = xn.NextNBLine(&rest)
		if line != "}" {
			err = MissingClosingBrace
		}
	} else {
		err = MissingEndPointsSection
	}
	return
}
func ParseMemberInfo(s string) (
	cm *MemberInfo, rest []string, err error) {

	ss := strings.Split(s, "\n")
	return ParseMemberInfoFromStrings(ss)
}

func ParseMemberInfoFromStrings(ss []string) (
	cm *MemberInfo, rest []string, err error) {

	bn, rest, err := xn.ParseBNFromStrings(ss, "memberInfo")
	if err == nil {
		cm = &MemberInfo{BaseNode: *bn}
		rest, err = collectAttrs(cm, rest)
		if err == nil {
			rest, err = collectMyEnds(cm, rest)
		}
		if err == nil {
			// expect and consume a closing brace
			line := xn.NextNBLine(&rest)
			if line != "}" {
				err = MissingClosingBrace
			}
		}
	}
	return
}
