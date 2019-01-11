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
			fmt.Printf("Unable to find %s; using default %d\n", pn, p[pn])
			continue
		}
		val, err := strconv.Atoi(vp)
		if err != nil {
			fmt.Printf("Unable to parse %s: %s; using default %d\n", pn, vp, p[pn])
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
				fmt.Printf("Disabling Peer TLS: missing ETCDCONTROLLER_PEER_%s definition\n", f)
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
				fmt.Printf("Disabling Client TLS: missing ETCDCONTROLLER_CLIENT_%s definition\n", f)
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
			fmt.Printf("Disable GRPC TLS - failed to load: %s", err)
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
				fmt.Printf("grpc serve fail: %s\n", err)
			}
		}
	}()

	m, err := group.NewManager(
		group.Config{
			Name:             v["ETCDCONTROLLER_NAME"],
			IP:               v["ETCDCONTROLLER_IP"],
			SerfPort:         p["ETCDCONTROLLER_SERF_PORT"],
			NodeListFilename: v["ETCDCONTROLLER_NODELISTFILENAME"],
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
			fmt.Printf("main TICK!\n")
			if isCon, notCon := m.IsConductor(); isCon {
				fmt.Printf("IS CONDUCTOR\n")
				if !c.IsRunning() {
					go c.Run()
				}
			} else {
				fmt.Printf("NOT CONDUCTOR: %s\n", notCon)
				if c != nil {
					fmt.Printf("TODO: Should stop conductor\n")
				}
			}
		}
	}
}
