package main

import (
	"crypto/tls"
	"net"
	"log"
)

type TLSAcceptor struct {
	listener	net.Listener
	conns		chan net.Conn
}

func NewTLSAcceptor(portStr string) (*TLSAcceptor,error) {
	accp := new(TLSAcceptor)
	accp.conns = make(chan net.Conn,100)
	certbs,keybs := GenerateTLSKeys()
	cert,err := tls.X509KeyPair(certbs,keybs)
	if err!=nil {
		return nil,err
	}
	config := new(tls.Config)
	config.Certificates = []tls.Certificate{cert}
	accp.listener,err = tls.Listen("tcp",":"+portStr,config)
	if err!=nil {
		return nil,err
	}
	go func(){
		for{
			conn, err := accp.listener.Accept()
			if err != nil {
				log.Print("ERROR::TLSAcceptor:: ",err)
				continue
			}
			accp.conns <- conn
		}
	}()
	return accp,nil
}

func (accp *TLSAcceptor) GetConn() net.Conn {
	return <-accp.conns
}
