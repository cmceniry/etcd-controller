package ectl

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
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
	Logger                   *log.Entry
	ETCDLogger               *log.Entry
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
	Config     ETCDConfig
	command    *exec.Cmd
	mux        sync.Mutex
	logger     *log.Entry
	etcdLogger *log.Entry
}

func New(c ETCDConfig) (*ETCDProcess, error) {
	e := ETCDProcess{
		Config:     c,
		logger:     c.Logger,
		etcdLogger: c.ETCDLogger,
	}
	if e.logger == nil {
		l := log.New()
		l.SetOutput(ioutil.Discard)
		e.logger = l.WithFields(log.Fields{"component": "none"})
	}
	if e.etcdLogger == nil {
		l := log.New()
		l.SetOutput(ioutil.Discard)
		e.etcdLogger = l.WithFields(log.Fields{"component": "none"})
	}
	return &e, nil
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
	cmd.Stdout = e.etcdLogger.Writer()
	cmd.Stderr = e.etcdLogger.Writer()
	err := cmd.Start()
	if err != nil {
		return err
	}
	e.command = cmd
	go e.wait()
	return nil
}

// StartServer starts the ETCD server normally (i.e. without any join/
// init/etc options)
func (e *ETCDProcess) StartServer() error {
	return e.start(e.Config.buildEnvironment())
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
		peers = append(peers, p+"="+u)
	}
	startOpts := append(
		e.Config.buildEnvironment(),
		"ETCD_INITIAL_CLUSTER_STATE=existing",
		"ETCD_INITIAL_CLUSTER="+strings.Join(peers, ","),
	)
	err := e.start(startOpts)
	if err != nil {
		return false, err
	}
	return true, nil
}

// IsInitialized returns if it is ready to run
func (e *ETCDProcess) IsInitialized() bool {
	s, err := os.Lstat(filepath.Join(e.Config.DataDir, "member"))
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
	e.logger.Debugf("e.command.ProcessState: %#v", e.command.ProcessState)
	if e.command.ProcessState == nil {
		return true
	}
	e.logger.Debugf("e.command.ProcessState.Exited(): %#v", e.command.ProcessState.Exited())
	return !e.command.ProcessState.Exited()
}

func (e *ETCDProcess) getLocalURL() string {
	// TODO: Probably put some intelligence here to convert to 127.0.0.1?
	us := strings.Split(e.Config.AdvertiseClientURLs, ",")
	return us[0]
}

// GetHealth shows the status of this node
func (e *ETCDProcess) GetHealth() (bool, error) {
	if e.command == nil {
		return false, fmt.Errorf("no etcd server running")
	}
	var tlsConfig *tls.Config
	if e.Config.ClientCertAuth {
		cert, err := tls.LoadX509KeyPair(e.Config.CertFile, e.Config.KeyFile)
		if err != nil {
			return false, err
		}
		caPool := x509.NewCertPool()
		caData, err := ioutil.ReadFile(e.Config.CAFile)
		if err != nil {
			return false, err
		}
		if ok := caPool.AppendCertsFromPEM(caData); !ok {
			return false, fmt.Errorf("unable to load ca certs")
		}
		tlsConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      caPool,
			ServerName:   "127.0.0.1",
		}
	}
	url := e.getLocalURL()
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{url},
		DialTimeout: 5 * time.Second,
		TLS:         tlsConfig,
	})
	if err != nil {
		return false, err
	}
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	resp, err := client.Status(ctx, url)
	cancel()
	e.logger.Debugf("StatusResponse: %#v", resp)
	e.logger.Debugf("StatusError: %s", err)
	switch err {
	case nil, rpctypes.ErrPermissionDenied:
		return true, nil
	case rpctypes.ErrTimeout:
		return false, nil
	default:
		return false, err
	}
}
