package overlay

// xlattice_go/overlay/addr_range.go

import (
	"encoding/hex"
	"errors"
	"fmt"
)

// An address range as the term is used in XLattice and on the global
// internet.  The range is defined by a prefix, the number of significant
// bits in that prefix, and the number of bits in a valid address.  So
// an ipv4 10/8 address range would be represented as [10], 8, 32.
type AddrRange struct {
	prefix    []byte // all addresses in the range match this prefix
	prefixLen uint   // number of leading bits, the 8 in 10/8
	addrLen   uint   // in bits: 32 for ipv4, 64 for ipv6
}

// XXX UNSAFE: should copy
func (r *AddrRange) Prefix()	[]byte	{	return r.prefix		}
func (r *AddrRange) PrefixLen() uint	{	return r.prefixLen	}
func (r *AddrRange) AddrLen()	uint	{	return r.addrLen	}

func NewAddrRange(prefix []byte, prefixLen uint, addrLen uint) (*AddrRange, error) {
	if prefix == nil {
		return nil, errors.New("IllegalArgument: nil prefix")
	}
	byteLen := uint(len(prefix))
	if prefixLen > 8*byteLen {
		return nil, errors.New("IllegalArgument: too few bits in prefix")
	}
	if prefixLen > addrLen {
		return nil, errors.New("IllegalArgument: prefix too long for addr len")
	}

	return &AddrRange{prefix, prefixLen, addrLen}, nil
}

func NewV4AddrRange(prefix []byte, prefixLen uint) (*AddrRange, error) {
	return NewAddrRange(prefix, prefixLen, uint(32))
}

func NewV6AddrRange(prefix []byte, prefixLen uint) (*AddrRange, error) {
	return NewAddrRange(prefix, prefixLen, uint(64))
}

func (r *AddrRange) Equal (any interface{}) bool {
	if any == r		{ return true  }
	if any == nil	{ return false }
	switch v := any.(type) {
	case *AddrRange:
		_ = v
	default:
		return false
	}
	other := any.(*AddrRange)
	if len(r.prefix) != len(other.prefix)		{ return false }
	for i := 0; i < len(r.prefix); i++ {
		if r.prefix[i] != other.prefix[i]		{ return false }
	}
	if r.prefixLen != other.prefixLen			{ return false }
	if r.addrLen   != other.addrLen				{ return false }
	return true
}
func (r *AddrRange) String() string {
	// XXX SOMETHING OF A HACK, because the AddrLen is absent
	return fmt.Sprintf("%s/%u", hex.EncodeToString(r.prefix), r.prefixLen)
}
