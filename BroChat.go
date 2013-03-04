package main

import (
	"log"
	"os"
	"fmt"
	"net"
	"strings"
)

type BroChat struct {
	state			*CurrentState

	bCastSender		*BCastSender
	bCastReceiver	*BCastReceiver
	
	tlsAcceptor		*TLSAcceptor
	tlsConnector	*TLSConnector
	
	userInterface	UI	
}

func NewBroChat(bCastIP string, port int) (*BroChat,error) {
	bc := new(BroChat)
	var err error
	
	bc.bCastSender,err = NewBCastSender(bCastIP,port)
	if err!=nil {
		return nil,err
	}
	
	bc.bCastReceiver,err = NewBCastReceiver(port)
	if err!=nil {
		return nil,err
	}
	
	bc.tlsAcceptor,err = NewTLSAcceptor(port)
	if err!=nil {
		return nil,err
	}
	
	bc.tlsConnector = NewTLSConnector(port,bc.tlsAcceptor.GetConfig())
	
	bc.state = NewCurrentState()
	bc.userInterface = NewGtkUI()
	
	/*
	Get Input from UI and process it
	*/
	go func(){
		for{
			userInput := bc.userInterface.GetInput()
			bc.ProcessUserInput(userInput)
		}
	}()
	
	/*
	Get BroadcastMessages and process them
	(recognize and display to ui)
	*/
	go func(){
		for{
			msg := bc.ReadNextMessage()
			bc.state.UpdateWhoMap(msg)
			if bc.state.Channel == msg.Channel {
				bc.userInterface.ProcessOutput(msg.String())
			}
		}
	}()
	
	/*
	Accept Whisper-Connections and run handler
	*/
	go func(){
		for{
			conn := bc.tlsAcceptor.GetConn()
			go func(){
				defer conn.Close()
				ip := conn.RemoteAddr().(*net.TCPAddr).IP
				nick := bc.state.GetNickToIP(ip)
				log.Print("Accepted BroConn from ",ip," nick: ",nick," fingerprint: ",conn.Fingerprint())
				bc.tlsConnector.SetConn(nick,conn)
				buff := make([]byte,1024)
				for{
					bs,err := conn.Read(buff[0:])
					if err!=nil {
						log.Print("finished connection: ",err)
						break
					}
					msg := string(buff[:bs])
					log.Print("received (tls) ",msg)
					bc.userInterface.ProcessOutput(nick+" <(private)>: "+msg)					
				}
			}()
		}
	}()
	return bc,nil
}

func (bc *BroChat) ProcessUserInput(input string){
	//check for command. if no command: broadcast.
	foundCommand,err := bc.ProcessPossibleCommand(input)
	if foundCommand && err!=nil {
		log.Print("malformed command. grml...")
		return
	}
	if foundCommand || len(input)<1 {
		return
	}
	bc.SendMessage(input)
}

func (bc *BroChat) ProcessPossibleCommand(payload string) (foundCommand bool,err error){
	cmd,err := NewUserCommand(payload)
	if err!=nil {
		return
	}
	switch cmd.Type {
		case UC_CHANGE_NICK: {
			bc.UpdateNickname(cmd.Payload)
			foundCommand = true
		}
		case UC_CHANGE_CHAN: {
			bc.JoinChannel(cmd.Payload)
			foundCommand = true
		}
		case UC_LEAVE: {
			bc.SendMessage("leaves the channel")
			os.Exit(0)
		}
		case UC_INFO: {
			bc.userInterface.ProcessOutput(bc.state.String())
			foundCommand = true
		}
		case UC_WHISPER: {
			foundCommand = true
			words := strings.SplitN(cmd.Payload," ",2)
			if len(words)<2 {
				break
			}
			target := words[0]
			info,ok := bc.state.WhoMap[target]
			if !ok {
				bc.userInterface.ProcessOutput("Don't know who you mean.\n")
				break
			}
			conn := bc.tlsConnector.Connect(target,info.IP)
			bs,err := conn.Write([]byte(words[1]))
			if err!=nil {
				log.Print("fail...",err,bs)
			}
			log.Print("finished whispering")
		}
		case UC_HELP: {
			foundCommand = true
			text := []string{
				"# Bro-Chat Help:",
				"# Commands:",
				"# /nick <new nickname>",
				"# /join <new channel>",
				"# /info",
				"# /whisper <nickname> <message>",
				"# /quit",
			}
			for _,line := range text {
				bc.userInterface.ProcessOutput(line)
			}
		}
	}
	return
}

func (bc *BroChat) UpdateNickname(newnick string){
	/*
	TODO: Add RegExp Checker here
	*/
	if len(newnick)>0 {
		text := fmt.Sprintf("changed nick to %v",newnick)
		bc.SendMessage(text)
		bc.state.Nickname = newnick
	}
}

func (bc *BroChat) JoinChannel(newchan string){
	text1 := "leaves the channel"
	text2 := "joins the channel"
	if bc.state.Channel!=""{
		bc.SendMessage(text1)
	}
	bc.state.Channel = newchan
	bc.SendMessage(text2)
}

func (bc *BroChat) SendMessage(payload string) error {
	msg := NewBCastMessage(bc.state.Nickname,bc.state.Channel,payload+"\n")
	bc.bCastSender.SendMessage(msg)
	return nil
}

func (bc *BroChat) ReadNextMessage() *BCastMessage {
	return bc.bCastReceiver.GetMessage()
}

