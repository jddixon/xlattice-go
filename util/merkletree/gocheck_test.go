package merkletree

import (
	. "gopkg.in/check.v1"
	"testing"
)

// IF USING gocheck, need a file like this in each package=directory.

func Test(t *testing.T) { TestingT(t) }

type XLSuite struct{}

var _ = Suite(&XLSuite{})

// TEST CONSTANTS BELOW THIS LINE
const (
	VERBOSITY = 1
)
