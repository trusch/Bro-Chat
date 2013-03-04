package main

import (
	"crypto/tls"
	"net"
	"log"
	"fmt"
)

type TLSAcceptor struct {
	listener	net.Listener
	conns		chan *tls.Conn
	config		*tls.Config
	port		int
}

func NewTLSAcceptor(port int) (*TLSAcceptor,error) {
	accp := new(TLSAcceptor)
	accp.port = port
	portStr := fmt.Sprintf("%v",port)
	accp.conns = make(chan *tls.Conn,100)
	certbs,keybs := GenerateTLSKeys()
	cert,err := tls.X509KeyPair(certbs,keybs)
	if err!=nil {
		return nil,err
	}
	config := new(tls.Config)
	config.Certificates = []tls.Certificate{cert}
	config.ClientAuth = tls.RequestClientCert
	config.InsecureSkipVerify = true
	accp.config = config
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
			accp.conns <- conn.(*tls.Conn)
		}
	}()
	return accp,nil
}

func (accp *TLSAcceptor) GetConn() *BroConn {
	conn :=  <-accp.conns
	log.Print("accepted new conn. Build BroConn...")
	broConn := NewBroConn(nil,accp.port,accp.config)
	err := broConn.SetConn(conn)
	if err!=nil {
		log.Print(err)
		return nil
	}
	return broConn
}

func (accp *TLSAcceptor) GetConfig() *tls.Config {
	return accp.config
}
