package msg

import (
	"bytes"
	//"code.google.com/p/go.crypto/sha3"
    "golang.org/x/crypto/sha3"

	"code.google.com/p/goprotobuf/proto"
	"encoding/binary"
	"fmt"
	xr "github.com/jddixon/rnglib_go"
	. "gopkg.in/check.v1"
)

var _ = proto.Marshal
var _ = fmt.Println

func (d *XLSuite) TestXLatticePkt(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_XLATTICE_PKT")
	}

	rng := xr.MakeSimpleRNG()

	myMsgN := uint64(rng.Int63())
	for myMsgN == 0 { // must not be zero
		myMsgN = uint64(rng.Int63())
	}

	id := make([]byte, 32) // sha3 length
	rng.NextBytes(id)      // random bytes

	seqBuf := new(bytes.Buffer)
	binary.Write(seqBuf, binary.LittleEndian, myMsgN)

	msgLen := 64 + rng.Intn(64)
	msg := make([]byte, msgLen)
	rng.NextBytes(msg) // fill with rubbish

	salt := make([]byte, 8)
	rng.NextBytes(salt) // still more rubbish

	digest := sha3.NewKeccak256()
	digest.Write(id)
	digest.Write(seqBuf.Bytes())
	digest.Write(msg)
	digest.Write([]byte(salt))

	hash := digest.Sum(nil)

	// XXX This does not adhere to the rules: it has no Cmd field;
	// since it has a payload it must be a Put, and so the id is
	// also required and the Hash field should be a Sig instead, right?
	var pkt = XLatticeMsg{
		MsgN:    &myMsgN,
		Payload: msg,
		Salt:    salt,
		Hash:    hash,
	}

	// In each of these cases, the test proves that the field
	// was present; otherwise the 'empty' value (zero, nil, etc)
	// would have been returned.
	msgNOut := pkt.GetMsgN()
	c.Assert(msgNOut, Equals, myMsgN)

	msgOut := pkt.GetPayload()
	d.compareByteSlices(c, msgOut, msg)

	saltOut := pkt.GetSalt()
	d.compareByteSlices(c, saltOut, salt)

	hashOut := pkt.GetHash()
	d.compareByteSlices(c, hashOut, hash)
}

func (d *XLSuite) compareByteSlices(c *C, a []byte, b []byte) {
	c.Assert(len(a), Equals, len(b))
	for i := 0; i < len(b); i++ {
		c.Assert(a[i], Equals, b[i])
	}
}
