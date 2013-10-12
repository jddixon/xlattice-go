package reg

// xlattice_go/reg/admin_client.go

import (
	"crypto/rsa"
	"fmt"
	xi "github.com/jddixon/xlattice_go/nodeID"
	xt "github.com/jddixon/xlattice_go/transport"
	"io"
)

var _ = fmt.Print

// In the current implementation, AdminClient's purpose is to register
// clusters, so it sets up a session with the registry, (calls SessionSetup),
// identifies itself with a Client/ClientOK sequence, uses Create/CreateReply
// to register the cluster, and then ends with Bye/Ack.  It then returns
// the clusterID and size to the caller.

// As implemented so far, this is an ephemeral client, meaning that it
// neither saves nor restores its Node; keys and such are generated for
// each instance.

type AdminClient struct {
	// In this implementation, AdminClient is a one-shot, launched
	// to create a single cluster

	ClientNode
}

func NewAdminClient(
	serverName string, serverID *xi.NodeID, serverEnd xt.EndPointI,
	serverCK, serverSK *rsa.PublicKey,
	clusterName string, size, epCount int) (ac *AdminClient, err error) {

	cn, err := NewClientNode("admin", "", ATTR_ADMIN, //  lfs string, attrs uint64,
		serverName, serverID, serverEnd,
		serverCK, serverSK, //  *rsa.PublicKey,
		clusterName, nil, // nil is clientID
		size,
		epCount,
		nil) // nil = []xt.EndPointI

	if err == nil {
		// Run() fills in clusterID
		ac = &AdminClient{
			ClientNode: *cn,
		}
	}
	return // FOO
}

// Start the client running in separate goroutine, so that this function
// is non-blocking.

func (ac *AdminClient) Run() (err error) {

	cn := &ac.ClientNode

	go func() {
		var (
			version1 uint32
		)
		clientName := cn.name
		cnx, version2, err := cn.SessionSetup(version1)
		_ = version2 // not yet used
		if err == nil {
			err = cn.ClientAndOK()
		}
		if err == nil {
			err = cn.CreateAndReply()
		}
		if err == nil {
			err = cn.ByeAndAck()
		}
		// END OF RUN ===============================================
		if cnx != nil {
			cnx.Close()
		}
		// DEBUG
		fmt.Printf("user client %s run complete ", clientName)
		if err != nil && err != io.EOF {
			fmt.Printf("- ERROR: %v", err)
		}
		fmt.Println("")
		// END

		cn.err = err
		cn.doneCh <- true
	}()
	return
}
