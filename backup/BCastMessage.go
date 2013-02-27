package main

import (
	"strings"
	"net"
)

type BCastMessage struct {
	SenderNick	string
	Channel		string
	Payload		string
	
	IP			net.IP
}

func NewBCastMessage(senderNick,channel,payload string) *BCastMessage {
	return &BCastMessage{senderNick,channel,payload,nil}
}

func (m *BCastMessage) ToBytes() []byte {
	return []byte(m.SenderNick+"\n"+m.Channel+"\n"+m.Payload)
}

func BCastMessageFromBytes(data []byte) *BCastMessage {
	str := string(data)
	parts := strings.SplitN(str,"\n",3)
	if len(parts)<3 {
		return nil
	}
	return NewBCastMessage(parts[0],parts[1],parts[2])
}

func (m *BCastMessage) String() string {
	return m.SenderNick+"@"+m.Channel+": "+m.Payload
}
