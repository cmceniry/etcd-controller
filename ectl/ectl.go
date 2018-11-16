package ectl

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/etcdserver/api/v3rpc/rpctypes"
)

// This is the main package for running the etcd process

// ETCDConfig is the config items for etcd
type ETCDConfig struct {
	Binary                   string
	Name                     string
	DataDir                  string
	AdvertiseClientURLs      string
	ListenClientURLs         string
	ClientCertAuth           bool
	CAFile                   string
	CertFile                 string
	KeyFile                  string
	InitialAdvertisePeerURLs string
	ListenPeerURLs           string
	PeerClientCertAuth       bool
	PeerCAFile               string
	PeerCertFile             string
	PeerKeyFile              string
}

func (c ETCDConfig) buildEnvironment() []string {
	return []string{
		"ETCD_NAME=" + c.Name,
		"ETCD_DATA_DIR=" + c.DataDir,

		"ETCD_ADVERTISE_CLIENT_URLS=" + c.AdvertiseClientURLs,
		"ETCD_LISTEN_CLIENT_URLS=" + c.ListenClientURLs,
		"ETCD_CLIENT_CERT_AUTH=" + fmt.Sprintf("%t", c.ClientCertAuth),
		"ETCD_CA_FILE=" + c.CAFile,
		"ETCD_CERT_FILE=" + c.CertFile,
		"ETCD_KEY_FILE=" + c.KeyFile,

		"ETCD_INITIAL_ADVERTISE_PEER_URLS=" + c.InitialAdvertisePeerURLs,
		"ETCD_LISTEN_PEER_URLS=" + c.ListenPeerURLs,
		"ETCD_PEER_CLIENT_CERT_AUTH=" + fmt.Sprintf("%t", c.PeerClientCertAuth),
		"ETCD_PEER_CA_FILE=" + c.PeerCAFile,
		"ETCD_PEER_CERT_FILE=" + c.PeerCertFile,
		"ETCD_PEER_KEY_FILE=" + c.PeerKeyFile,
	}
}

// ETCDProcess is the main struct that is used to pass etcd commands around
//
// ETCD_INITIAL_CLUSTER and ETCD_INITIAL_CLUSTER_STATE are handled internal
// to this struct
//
type ETCDProcess struct {
	Config  ETCDConfig
	command *exec.Cmd
}

// Stop stops the ETCD server
func (e *ETCDProcess) Stop() error {
	if e.command == nil {
		return nil
	}
	// TODO: timeout
	return e.command.Process.Kill()
}

// StartInitial starts an empty ETCD server
func (e *ETCDProcess) StartInitial() error {
	cmd := exec.Command(e.Config.Binary)
	cmd.Env = append(
		e.Config.buildEnvironment(),
		// "ETCD_INITIAL_CLUSTER=none",
		"ETCD_INITIAL_CLUSTER_STATE=new",
	)
	cmd.Stdin = nil
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Start()
	if err != nil {
		return err
	}
	e.command = cmd
	return nil
}

// JoinCluster starts ETCD server with the parameters to join an existing
// cluster.
func (e *ETCDProcess) JoinCluster(peerURLs map[string]string) (bool, error) {
	peers := []string{}
	for p, u := range peerURLs {
		peers = append(peers, p + "=" + u)
	}
	cmd := exec.Command(e.Config.Binary)
	cmd.Env = append(
		e.Config.buildEnvironment(),
		"ETCD_INITIAL_CLUSTER_STATE=existing",
		"ETCD_INITIAL_CLUSTER=" + strings.Join(peers, ","),
	)
	cmd.Stdin = nil
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Start()
	if err != nil {
		return false, err
	}
	e.command = cmd
	return true, nil
}

// GetHealth shows the status of this node
func (e *ETCDProcess) GetHealth() (bool, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return false, err
	}
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	resp, err := client.Status(ctx, "http://localhost:2379")
	cancel()
	if err == nil {
		fmt.Printf("%#v\n", resp)
	}
	switch err {
	case nil, rpctypes.ErrPermissionDenied:
		return true, nil
	case rpctypes.ErrTimeout:
		return false, nil
	default:
		return false, err
	}
}
