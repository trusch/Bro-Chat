package main

import(
	"net"
	"crypto/tls"
	"log"
)

type TLSConnector struct {
	port	int
	config	*tls.Config
	conns	map[string]*BroConn
}

func (c *TLSConnector)Connect(nick string,ip net.IP) *BroConn{
	conn := c.conns[nick]
	if conn==nil {
		log.Print("Must create new BroConn")
		conn = NewBroConn(ip,c.port,c.config)
		c.conns[nick] = conn
	}
	return conn
}

func (c *TLSConnector)SetConn(nick string,conn *BroConn){
	c.conns[nick] = conn
}

func NewTLSConnector(port int,cfg *tls.Config) *TLSConnector {
	result := new(TLSConnector)
	result.port = port
	result.config = cfg
	result.conns = make(map[string]*BroConn)
	return result
}
