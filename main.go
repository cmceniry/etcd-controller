package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/cmceniry/etcd-controller/conductor"
	"github.com/cmceniry/etcd-controller/driver"
	"github.com/cmceniry/etcd-controller/group"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func addGRPCTLSOptions(cafile, certfile, keyfile string) (grpc.ServerOption, error) {
	c, err := tls.LoadX509KeyPair(certfile, keyfile)
	if err != nil {
		return nil, err
	}
	caPool := x509.NewCertPool()
	caData, err := ioutil.ReadFile(cafile)
	if err != nil {
		return nil, err
	}
	if ok := caPool.AppendCertsFromPEM(caData); !ok {
		return nil, fmt.Errorf("unable to add ca pool")
	}
	return grpc.Creds(credentials.NewTLS(&tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{c},
		ClientCAs:    caPool,
	})), nil
}

func main() {
	mainlog := log.WithFields(log.Fields{
		"component": "main",
	})

	v := map[string]string{
		"ETCDCONTROLLER_NODELISTFILENAME": "/config/node-list.yaml",
		"ETCDCONTROLLER_COMMAND_PORT":     "4270",
		"ETCDCONTROLLER_SERF_PORT":        "4271",
		"ETCDCONTROLLER_NAME":             "test001",
		"ETCDCONTROLLER_IP":               "127.0.0.1",
		"ETCDCONTROLLER_BINARY":           "/usr/local/bin/etcd",
		"ETCDCONTROLLER_DATADIR":          "/var/lib/etcd",
		"ETCDCONTROLLER_PEER_PORT":        "2380",
		"ETCDCONTROLLER_PEER_CA":          "",
		"ETCDCONTROLLER_PEER_CERT":        "",
		"ETCDCONTROLLER_PEER_KEY":         "",
		"ETCDCONTROLLER_CLIENT_PORT":      "2379",
		"ETCDCONTROLLER_CLIENT_CA":        "",
		"ETCDCONTROLLER_CLIENT_CERT":      "",
		"ETCDCONTROLLER_CLIENT_KEY":       "",
	}
	for k := range v {
		if os.Getenv(k) != "" {
			v[k] = os.Getenv(k)
		}
	}
	p := map[string]int{
		"ETCDCONTROLLER_COMMAND_PORT": 4270,
		"ETCDCONTROLLER_SERF_PORT":    4271,
		"ETCDCONTROLLER_PEER_PORT":    2380,
		"ETCDCONTROLLER_CLIENT_PORT":  2379,
	}
	for pn := range p {
		vp, ok := v[pn]
		if !ok {
			mainlog.Infof("Unable to find %s; using default %d", pn, p[pn])
			continue
		}
		val, err := strconv.Atoi(vp)
		if err != nil {
			mainlog.Errorf("Unable to parse %s: %s; using default %d", pn, vp, p[pn])
			continue
		}
		p[pn] = val
	}

	peerTLS := false
	for _, f := range []string{"CA", "CERT", "KEY"} {
		if v["ETCDCONTROLLER_PEER_"+f] != "" {
			peerTLS = true
			break
		}
	}
	if peerTLS {
		for _, f := range []string{"CA", "CERT", "KEY"} {
			if v["ETCDCONTROLLER_PEER_"+f] == "" {
				mainlog.Infof("Disabling Peer TLS: missing ETCDCONTROLLER_PEER_%s definition", f)
				peerTLS = false
			}
		}
	}

	clientTLS := false
	for _, f := range []string{"CA", "CERT", "KEY"} {
		if v["ETCDCONTROLLER_CLIENT_"+f] != "" {
			clientTLS = true
			break
		}
	}
	if clientTLS {
		for _, f := range []string{"CA", "CERT", "KEY"} {
			if v["ETCDCONTROLLER_CLIENT_"+f] == "" {
				mainlog.Infof("Disabling Client TLS: missing ETCDCONTROLLER_CLIENT_%s definition", f)
				clientTLS = false
			}
		}
	}

	var opts []grpc.ServerOption
	if peerTLS {
		cred, err := addGRPCTLSOptions(
			v["ETCDCONTROLLER_PEER_CA"],
			v["ETCDCONTROLLER_PEER_CERT"],
			v["ETCDCONTROLLER_PEER_KEY"],
		)
		if err != nil {
			log.Errorf("Disable GRPC TLS - failed to load: %s", err)
		} else {
			opts = append(opts, cred)
		}
	}
	gserver := grpc.NewServer(opts...)

	dc := driver.SimpleDriverConfig{
		Binary:      v["ETCDCONTROLLER_BINARY"],
		Name:        v["ETCDCONTROLLER_NAME"],
		DataDir:     v["ETCDCONTROLLER_DATADIR"],
		IP:          v["ETCDCONTROLLER_IP"],
		ClientPort:  p["ETCDCONTROLLER_CLIENT_PORT"],
		PeerPort:    p["ETCDCONTROLLER_PEER_PORT"],
		CommandPort: p["ETCDCONTROLLER_COMMAND_PORT"],
		Logger:      log.WithFields(log.Fields{"component": "driver"}),
		ECTLLogger:  log.WithFields(log.Fields{"component": "ectl"}),
		ETCDLogger:  log.WithFields(log.Fields{"component": "etcd"}),
	}
	if peerTLS {
		dc.PeerTLSCA = v["ETCDCONTROLLER_PEER_CA"]
		dc.PeerTLSCert = v["ETCDCONTROLLER_PEER_CERT"]
		dc.PeerTLSKey = v["ETCDCONTROLLER_PEER_KEY"]
	}
	if clientTLS {
		dc.ClientTLSCA = v["ETCDCONTROLLER_CLIENT_CA"]
		dc.ClientTLSCert = v["ETCDCONTROLLER_CLIENT_CERT"]
		dc.ClientTLSKey = v["ETCDCONTROLLER_CLIENT_KEY"]
	}
	s, err := driver.NewSimpleDriver(dc)
	if err != nil {
		panic(err)
	}
	s.RegisterWithGRPCServer(gserver)

	nc := conductor.Config{
		NodeListFilename: v["ETCDCONTROLLER_NODELISTFILENAME"],
		CommandPort:      p["ETCDCONTROLLER_COMMAND_PORT"],
		Logger:           log.WithFields(log.Fields{"component": "conductor"}),
	}
	if peerTLS {
		nc.PeerTLSCA = v["ETCDCONTROLLER_PEER_CA"]
		nc.PeerTLSCert = v["ETCDCONTROLLER_PEER_CERT"]
		nc.PeerTLSKey = v["ETCDCONTROLLER_PEER_KEY"]
	}
	if clientTLS {
		nc.ClientTLSCA = v["ETCDCONTROLLER_CLIENT_CA"]
		nc.ClientTLSCert = v["ETCDCONTROLLER_CLIENT_CERT"]
		nc.ClientTLSKey = v["ETCDCONTROLLER_CLIENT_KEY"]
	}
	c := conductor.NewConductor(nc)
	c.RegisterWithGRPCServer(gserver)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", p["ETCDCONTROLLER_COMMAND_PORT"]))
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			err := gserver.Serve(l)
			if err != nil {
				mainlog.Errorf("grpc serve fail (restarting): %s", err)
			}
		}
	}()

	m, err := group.NewManager(
		group.Config{
			Name:             v["ETCDCONTROLLER_NAME"],
			IP:               v["ETCDCONTROLLER_IP"],
			SerfPort:         p["ETCDCONTROLLER_SERF_PORT"],
			NodeListFilename: v["ETCDCONTROLLER_NODELISTFILENAME"],
			Logger:           log.WithField("component", "grouper"),
			SerfLogger:       log.WithField("component", "serf"),
		},
	)
	if err != nil {
		panic(err)
	}
	m.Run()

	t := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-t.C:
			mainlog.Debugf("TICK!")
			if isCon, err := m.IsConductor(); isCon {
				mainlog.Debug("CONDUCTOR: true")
				if !c.IsRunning() {
					go c.Run()
				}
			} else {
				mainlog.Debugf("CONDUCTOR: false")
				if err != nil {
					mainlog.Infof("Conductor evaluation failed: %s", err)
				}
				if c != nil {
					mainlog.Printf("TODO: Should stop the conductor")
				}
			}
		}
	}
}
