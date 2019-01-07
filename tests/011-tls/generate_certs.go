package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"
	"net"
	"os"
	"time"
)

func generateKeyAndCert(name, ipstr string, signer *x509.Certificate, signerkey *rsa.PrivateKey) (*x509.Certificate, *rsa.PrivateKey) {
	var ips []net.IP
	if ipstr != "" {
		ips = []net.IP{net.ParseIP(ipstr), net.IPv4(127, 0, 0, 1)}
	}
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	template := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: name},
		NotBefore:             time.Now().Truncate(24 * time.Hour),
		NotAfter:              time.Now().Truncate(24 * time.Hour).Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
		IPAddresses:           ips,
	}
	if signer == nil || signerkey == nil {
		template.IsCA = true
		template.KeyUsage |= x509.KeyUsageCertSign
		signer = template
		signerkey = key
	}
	der, _ := x509.CreateCertificate(
		rand.Reader,
		template,
		signer,
		&key.PublicKey,
		signerkey,
	)
	cert, err := x509.ParseCertificate(der)
	if err != nil {
		panic(err)
	}
	return cert, key
}

func saveKeyAndCert(prefix string, cert *x509.Certificate, key *rsa.PrivateKey) {
	certPem := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw})
	if err := ioutil.WriteFile(prefix+".crt", certPem, 0600); err != nil {
		panic(err)
	}
	keyBytes := x509.MarshalPKCS1PrivateKey(key)
	keyPem := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: keyBytes})
	if err := ioutil.WriteFile(prefix+".key", keyPem, 0600); err != nil {
		panic(err)
	}
}

func genAndSave(name, ipstr, prefix string, scert *x509.Certificate, skey *rsa.PrivateKey) (*x509.Certificate, *rsa.PrivateKey) {
	c, k := generateKeyAndCert(name, ipstr, scert, skey)
	saveKeyAndCert(prefix, c, k)
	return c, k
}

func main() {
	ipbase := os.Args[1]
	peerCACert, peerCAKey := genAndSave("Peer CA", "", "./config/peer-ca", nil, nil)
	genAndSave("node1 peer", ipbase+".101", "./config/node1-peer", peerCACert, peerCAKey)
	genAndSave("node2 peer", ipbase+".102", "./config/node2-peer", peerCACert, peerCAKey)
	genAndSave("node3 peer", ipbase+".103", "./config/node3-peer", peerCACert, peerCAKey)
	genAndSave("controller peer", "", "./config/controller-peer", peerCACert, peerCAKey)

	clientCACert, clientCAKey := genAndSave("Client CA", "", "./config/client-ca", nil, nil)
	genAndSave("node1 client", ipbase+".101", "./config/node1-client", clientCACert, clientCAKey)
	genAndSave("node2 client", ipbase+".102", "./config/node2-client", clientCACert, clientCAKey)
	genAndSave("node3 client", ipbase+".103", "./config/node3-client", clientCACert, clientCAKey)
	genAndSave("controller client", "", "./config/controller-client", clientCACert, clientCAKey)

	fmt.Printf("Certificates generated\n")
}
