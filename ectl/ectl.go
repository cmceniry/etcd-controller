package ectl

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
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
	mux     sync.Mutex
}

func (e *ETCDProcess) wait() {
	e.command.Wait()
	e.mux.Lock()
	defer e.mux.Unlock()
	e.command = nil
}

func (e *ETCDProcess) start(env []string) error {
	e.mux.Lock()
	defer e.mux.Unlock()
	if e.command != nil {
		return fmt.Errorf("etcd already running")
	}
	cmd := exec.Command(e.Config.Binary)
	cmd.Env = env
	cmd.Stdin = nil
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		return err
	}
	e.command = cmd
	go e.wait()
	return nil
}

// StopServer stops the ETCD server
func (e *ETCDProcess) StopServer() error {
	e.mux.Lock()
	defer e.mux.Unlock()
	if e.command == nil {
		return fmt.Errorf("etcd not running")
	}
	err := e.command.Process.Signal(syscall.SIGTERM)
	if err != nil {
		return fmt.Errorf("terminate failed: %s", err)
	}
	return nil
}

// StartInitial starts an empty ETCD server
func (e *ETCDProcess) StartInitial() error {
	startOpts := append(
		e.Config.buildEnvironment(),
		"ETCD_INITIAL_CLUSTER_STATE=new",
	)
	err := e.start(startOpts)
	if err != nil {
		return err
	}
	return nil
}

// JoinCluster starts ETCD server with the parameters to join an existing
// cluster.
func (e *ETCDProcess) JoinCluster(peerURLs map[string]string) (bool, error) {
	peers := []string{}
	for p, u := range peerURLs {
		peers = append(peers, p + "=" + u)
	}
	startOpts := append(
		e.Config.buildEnvironment(),
		"ETCD_INITIAL_CLUSTER_STATE=existing",
		"ETCD_INITIAL_CLUSTER=" + strings.Join(peers, ","),
	)
	err := e.start(startOpts)
	if err != nil {
		return false, err
	}
	return true, nil
}

// IsInitialized returns if it is ready to run
func (e *ETCDProcess) IsInitialized() bool {
	s, err := os.Lstat(e.Config.DataDir)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// IsRunning returns if there is a process running
func (e *ETCDProcess) IsRunning() bool {
	e.mux.Lock()
	defer e.mux.Unlock()
	if e.command == nil {
		return false
	}
	fmt.Printf("e.command.ProcessState: %#v\n", e.command.ProcessState)
	if e.command.ProcessState == nil {
		return true
	}
	fmt.Printf("e.command.ProcessState.Exited(): %#v\n", e.command.ProcessState.Exited())
	return !e.command.ProcessState.Exited()
}

// GetHealth shows the status of this node
func (e *ETCDProcess) GetHealth() (bool, error) {
	if e.command == nil {
		return false, fmt.Errorf("no etcd server running")
	}
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
