package main

import (
	"net"
	"log"
)

type BCastReceiver struct {
	conn 	*net.UDPConn
	msgs	chan *BCastMessage	
}

func (r *BCastReceiver) GetMessage() *BCastMessage {
	return <-r.msgs
}

func NewBCastReceiver(port int) (*BCastReceiver,error){
	result := new(BCastReceiver)
	result.msgs = make(chan *BCastMessage,64)
	conn, err := net.ListenUDP("udp4", &net.UDPAddr{
	        IP:   net.IPv4(0, 0, 0, 0),
	        Port: port,
		})
	if err!=nil {
		return nil,err
	}
	result.conn = conn 
	go func(){
	    data := make([]byte, (1<<16)-1)  //maximal udp packet size
		for {
		    read, addr , err := result.conn.ReadFromUDP(data[0:])
		    if err!=nil {
		    	log.Print("BCastReceiver:: ",err,addr)
		    	continue
		    }
		    msg := BCastMessageFromBytes(data[:read])
		    msg.IP = addr.IP
		    result.msgs <- msg
		}
	}()
	return result,nil
}
