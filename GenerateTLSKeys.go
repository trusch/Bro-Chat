package main

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

// Generate a self-signed X.509 certificate for a TLS server. Outputs to
// 'cert.pem' and 'key.pem' and will overwrite existing files.


import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log"
	"math/big"
	"time"
	"bytes"
)

func GenerateTLSKeys() ([]byte,[]byte){
	hostName := "127.0.0.1"
	priv, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Fatalf("failed to generate private key: %s", err)
		return nil,nil
	}

	now := time.Now()

	template := x509.Certificate{
		SerialNumber: new(big.Int).SetInt64(0),
		Subject: pkix.Name{
			CommonName:   hostName,
			Organization: []string{"Acme Co"},
		},
		NotBefore: now.Add(-5 * time.Minute).UTC(),
		NotAfter:  now.AddDate(1, 0, 0).UTC(), // valid for 1 year.

		SubjectKeyId: []byte{1, 2, 3, 4},
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		log.Fatalf("Failed to create certificate: %s", err)
		return nil,nil
	}
	
	var certOut *bytes.Buffer
	
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	log.Print("written cert.pem\n")

	
	var keyOut *bytes.Buffer

	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
	log.Print("written key.pem\n")

	return certOut.Bytes(),keyOut.Bytes()
}

