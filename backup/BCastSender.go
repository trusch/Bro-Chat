package main

import (
	"fmt"
	"net"
	"strings"
	"log"
	"errors"
)


type BCastSender struct {
	conn	net.Conn
	input	chan *BCastMessage
}

func (s *BCastSender) SetupUDPConn(bCastIpAddr string,bCastPort int) error {
	if bCastIpAddr=="" {
		addr,err := GetBCastAddr()
		if err!=nil {
			return err
		}
		bCastIpAddr = addr
	}
	ip := net.ParseIP(bCastIpAddr)
	if ip==nil {
		return errors.New("BCastSender can't parse ip")
	}
	conn,err := net.DialUDP("udp4", nil, &net.UDPAddr{
        	IP:   ip,
        	Port: bCastPort,
		})
	if err!=nil {
		return err
	}
	s.conn = conn
	return nil
}

func NewBCastSender(bCastIpAddr string,port int) (*BCastSender,error) {
	sender := new(BCastSender)
	sender.input = make(chan *BCastMessage,64)
	err := sender.SetupUDPConn(bCastIpAddr,port)
	if err!=nil {
		return nil,err
	}
	go func(){
		for msg := range sender.input {
			payload := msg.ToBytes()
			bs,err := sender.conn.Write(payload)
			if f := bs!=len(payload); f || err!=nil {
				log.Print("ERROR::BCastSender::SendMessage:: writtenBytes!=messageLenght:",f,err)
			}
		}
	}()
	return sender,err
}

func (s *BCastSender) SendMessage(msg *BCastMessage){
	s.input <- msg
}


func GetBCastAddr() (string,error) {
	addrs,err := net.InterfaceAddrs()
	if err!=nil {
		return "",err
	}
	result := ""
	for _,addr := range addrs {
		str := addr.String()
		parts1 := strings.Split(str,"/")
		parts2 := strings.Split(parts1[0],".")
		if parts2[0]!="127" {
			switch(parts1[1]){
				case "24": {
					result = fmt.Sprintf("%v.%v.%v.255",parts2[0],parts2[1],parts2[2])
				}
				case "16": {
					result = fmt.Sprintf("%v.%v.255.255",parts2[0],parts2[1])
				}
				case "8": {
					result = fmt.Sprintf("%v.255.255.255",parts2[0])
				}
			}
		}
		if result!=""{
			break
		}
	}
	return result,nil
}

