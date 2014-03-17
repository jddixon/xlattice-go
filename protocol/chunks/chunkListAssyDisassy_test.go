package chunks

// xlattice_go/protocol/chunks/chunkListAssyDisassy_test.go

import (
	"bytes"
	"code.google.com/p/go.crypto/sha3"
	"crypto/rand"
	"crypto/rsa"
	// "encoding/hex"
	"fmt"
	xi "github.com/jddixon/xlattice_go/nodeID"
	xr "github.com/jddixon/xlattice_go/rnglib"
	"github.com/jddixon/xlattice_go/u"
	// xu "github.com/jddixon/xlattice_go/util"
	xf "github.com/jddixon/xlattice_go/util/lfs"
	. "launchpad.net/gocheck"
	"os"
	"path"
)

var _ = fmt.Print

func (s *XLSuite) TestChunkListAssyDisassy(c *C) {
	rng := xr.MakeSimpleRNG()

	// make a slice 3 to 7 chunks long, fill with random-ish data ---
	chunkCount := 3 + rng.Intn(5) // so 3 to 7, inclusive
	lastChunkLen := 1 + rng.Intn(MAX_DATA_BYTES-1)
	dataLen := (chunkCount-1)*MAX_DATA_BYTES + lastChunkLen
	data := make([]byte, dataLen)
	rng.NextBytes(data)

	// calculate datum, the SHA3 hash of the data -------------------
	d := sha3.NewKeccak256()
	d.Write(data)
	hash := d.Sum(nil)
	datum, err := xi.NewNodeID(hash)
	c.Assert(err, IsNil)

	// create tmp if it doesn't exist -------------------------------
	found, err := xf.PathExists("tmp")
	c.Assert(err, IsNil)
	if !found {
		err = os.MkdirAll("tmp", 0755)
		c.Assert(err, IsNil)
	}

	// create scratch subdir with unique name -----------------------
	var pathToU string
	for {
		dirName := rng.NextFileName(8)
		pathToU = path.Join("tmp", dirName)
		found, err = xf.PathExists(pathToU)
		c.Assert(err, IsNil)
		if !found {
			break
		}
	}

	// create a FLAT uDir at that point -----------------------------
	myU, err := u.New(pathToU, u.DIR_FLAT, 0) // 0 means default perm
	c.Assert(err, IsNil)

	// write the test data into uDir --------------------------------
	bytesWritten, key, err := myU.PutData(data, datum.Value())
	c.Assert(err, IsNil)
	c.Assert(bytes.Equal(datum.Value(), key), Equals, true)
	c.Assert(bytesWritten, Equals, int64(dataLen))

	skPriv, err := rsa.GenerateKey(rand.Reader, 1024) // cheap key
	c.Assert(err, IsNil)
	c.Assert(skPriv, NotNil)

	// XXX WORKING HERE
}