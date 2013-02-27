package main

import (
	"net"
	"time"
	"fmt"
)

type UserInformation struct {
	Nickname	string
	Channel		string
	IP			net.IP
	LastTime	int64
}

type CurrentState struct {
	Nickname	string
	Channel		string
	WhoMap		map[string]*UserInformation	
}

func NewCurrentState() *CurrentState {
	return &CurrentState{"Anonymous","#foyer",make(map[string]*UserInformation)}
}

func (cs *CurrentState) UpdateWhoMap(packet *BCastMessage){
	name := packet.SenderNick
	info,ok := cs.WhoMap[name]
	if !ok {
		info = new(UserInformation)
		info.Nickname = name
		info.Channel = packet.Channel
		info.IP = packet.IP
		info.LastTime = time.Now().Unix()
		cs.WhoMap[name] = info
	}else{
		info.LastTime = time.Now().Unix()
		if info.Channel!=packet.Channel {
			info.Channel = packet.Channel
		}
		if !info.IP.Equal(packet.IP) {
			info.IP = packet.IP
		}
	}
}

func (cs *CurrentState) String() string {
	result := "#######################\n#Bro-Chat Information:\n#######################\n#Your Name: %v\n#Your Chan: %v\n"
	result = fmt.Sprintf(result,cs.Nickname,cs.Channel)
	for _,info := range cs.WhoMap {
		result += fmt.Sprintf("# %v@%v %v (%v)\n",info.Nickname,info.Channel,info.IP,info.LastTime)
	}
	result += "#######################"
	return result
}
