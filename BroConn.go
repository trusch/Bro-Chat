package main

import (
	"crypto/tls"
	"net"
	"fmt"
	"errors"
	"crypto/sha1"
	"encoding/base64"
)

type BroConn struct {
	*tls.Conn
	ip		net.IP
	port	int
	config	*tls.Config
	keyhash string
}

func(bc *BroConn) Fingerprint() string {
	return bc.keyhash
}

func(bc *BroConn) connect() error{
	conn,err := tls.Dial("tcp",fmt.Sprintf("%v:%v",bc.ip,bc.port),bc.config)
	if err!=nil {
		return err
	}
	bc.Conn=conn
	return bc.handshake()
}

func(bc *BroConn) handshake() error {
	err := bc.Handshake()
	if err!=nil {
		return err
	}
	certs := bc.ConnectionState().PeerCertificates
	if certs==nil || len(certs)==0 {
		return errors.New("no certificate submitted")
	}
	raw := certs[0].Raw
	hash := sha1.New().Sum(raw)
	b64 := base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")
	bc.keyhash = b64.EncodeToString(hash)
	return nil
}

func (bc *BroConn) Read(b []byte) (n int, err error) {
	if bc.Conn==nil {
		err = bc.connect()
		if err!=nil {
			return 0,err
		}
	}
	n,err = bc.Conn.Read(b)
	if err!=nil {
		err = bc.connect()
		if err!=nil {
			return 0,err
		}
		n,err = bc.Conn.Read(b)
		if err!=nil {
			return 0,err
		}
	}
	return n,nil
}

func (bc *BroConn) Write(b []byte) (n int, err error) {
	if bc.Conn==nil {
		err = bc.connect()
		if err!=nil {
			return 0,err
		}
	}
	n,err = bc.Conn.Write(b)
	if err!=nil {
		err = bc.connect()
		if err!=nil {
			return 0,err
		}
		n,err = bc.Conn.Write(b)
		if err!=nil {
			return 0,err
		}
	}
	return n,nil
}

func (bc *BroConn) SetConn(conn *tls.Conn) error {
	bc.Conn = conn
	bc.ip = net.ParseIP(conn.RemoteAddr().String())
	return bc.handshake()
}

func NewBroConn(ip net.IP,port int,config *tls.Config) *BroConn {
	bc := new(BroConn)
	bc.ip = ip
	bc.port = port
	bc.config = config
	return bc
}
