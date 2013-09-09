package main

import (
	"flag"
	"fmt"
	"github.com/jddixon/xlattice_go/reg"
	xt "github.com/jddixon/xlattice_go/transport"
	"os"
	"path"
)

func Usage() {
	fmt.Printf("Usage: %s [OPTIONS]\n", os.Args[0])
	fmt.Printf("where the options are:\n")
	flag.PrintDefaults()
}

const (
	DEFAULT_ADDR = "127.0.0.1"
	DEFAULT_NAME = "xlReg"
	DEFAULT_LFS  = "/var/XLReg"
	DEFAULT_PORT = 55555
)

var (
	// these need to be referenced as pointers
	address  = flag.String("p", DEFAULT_ADDR, "registry IP address")
	justShow = flag.Bool("j", false, "display option settings and exit")
	lfs      = flag.String("lfs", DEFAULT_LFS, "path to work directory")
	name     = flag.String("p", DEFAULT_NAME, "registry name")
	port     = flag.Int("p", DEFAULT_PORT, "listening port")
	testing  = flag.Bool("T", false, "test run")
	verbose  = flag.Bool("v", false, "be talkative")
)

func init() {
	fmt.Println("init() invocation") // DEBUG
}

func main() {
	flag.Usage = Usage
	flag.Parse()

	// FIXUPS ///////////////////////////////////////////////////////
	if *testing {
		if *lfs == DEFAULT_LFS || *lfs == "" {
			*lfs = "./tmp"
		} else {
			*lfs = path.Join("tmp", *lfs)
		}
	}
	addrAndPort := fmt.Sprintf("%s:%d", *address, port)
	endPoint, err := xt.NewTcpEndPoint( addrAndPort )
	if err != nil { 
		fmt.Printf("not a valid endPoint: %s\n", addrAndPort)
		// XXX STUB XXX 
	}

	// SANITY CHECKS ////////////////////////////////////////////////

	// DISPLAY FLAGS ////////////////////////////////////////////////
	if *verbose || *justShow {
		fmt.Printf("address                = %v\n", *address)
		fmt.Printf("endPoint               = %v\n", *endPoint)
		fmt.Printf("justShow               = %v\n", *justShow)
		fmt.Printf("lfs                    = %s\n", *lfs)
		fmt.Printf("name                   = %s\n", *name)
		fmt.Printf("port                   = %d\n", *port)
		fmt.Printf("testing                = %v\n", *testing)
		fmt.Printf("verbose                = %v\n", *verbose)
	}
	if *justShow {
		return
	}
	// SET UP OPTIONS ///////////////////////////////////////////////
	var opt reg.RegOptions
	opt.Lfs = *lfs
	opt.Port = *port
	opt.Testing = *testing
	opt.Verbose = *verbose

	r, err := setup(&opt)
	if err == nil {
		err = serve(r)
	}
	_ = err
}
func setup(opt *reg.RegOptions) (r *reg.RegNode, err error) {
	// If LFS/.xlattice/config exists, we load that.  Otherwise we
	// create a node.  In either case we force the node to listen on
	// the designated port

	// XXX STUB XXX

	r, err = reg.New(opt.Name, opt.Lfs, 
		nil, nil, nil, // opt.Id, opt.CKey, opt.SKey, 
		nil, 
		opt.EndPoint)

	return r, err
}
func serve(r *reg.RegNode) (err error) {

	// XXX STUB XXX

	// ORDERLY SHUTDOWN /////////////////////////////////////////////

	return
}
