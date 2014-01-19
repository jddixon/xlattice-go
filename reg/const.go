package reg

// xlattice_go/reg/const.go

const (
	SHA1_LEN = 20 // length in bytes of binary hash
	SHA3_LEN = 32

	BLOCK_SIZE = 4096

	DEFAULT_M = uint(20)
	DEFAULT_K = uint(8)
)

// client attrs bits, also used for member attrs
const (
	ATTR_EPHEMERAL = 1 << iota
	ATTR_ADMIN
	ATTR_SOLO // no related cluster, persists config to LFS
)
