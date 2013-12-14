package chunks

import (
	e "errors"
)

var (
	BadDatum                    = e.New("bad datum - doesn't match content hash")
	BadDatumLength              = e.New("bad datum - length must be 32")
	EmptyTitle                  = e.New("empty title parameter")
	MismatchedPublicPrivateKeys = e.New("public and private keys don't match")
	NilData                     = e.New("nil data parameter")
	NilDatum                    = e.New("nil datum (content hash) parameter")
	NilReader                   = e.New("nil io.Reader")
	NilRSAPrivKey               = e.New("nil RSA private key parameter")
	NilSubClass                 = e.New("nil subClass parameter")
	NoDigSig                    = e.New("no digital signature: the list is not signed")
	NoNthItem                   = e.New("no Nth item")
	ZeroLengthInput             = e.New("zero length input")
)