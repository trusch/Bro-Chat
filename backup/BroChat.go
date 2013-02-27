package main

import (
	"log"
	"os"
	"fmt"
)

type BroChat struct {
	state			*CurrentState

	bCastSender		*BCastSender
	bCastReceiver	*BCastReceiver
	
	tlsAcceptor		*TLSAcceptor
	
	userInterface	UI	
}

func NewBroChat(bCastIP string, bCastPort,tlsPort int) (*BroChat,error) {
	bc := new(BroChat)
	var err error
	bc.bCastSender,err = NewBCastSender(bCastIP,bCastPort)
	if err!=nil {
		return nil,err
	}
	bc.bCastReceiver,err = NewBCastReceiver(bCastPort)
	if err!=nil {
		return nil,err
	}
	bc.state = NewCurrentState()
	
	bc.userInterface = NewGtkUI()
	
	go func(){
		for{
			userInput := bc.userInterface.GetInput()
			bc.ProcessUserInput(userInput)
		}
	}()
	
	go func(){
		for{
			msg := bc.ReadNextMessage()
			bc.state.UpdateWhoMap(msg)
			if bc.state.Channel == msg.Channel {
				bc.userInterface.ProcessOutput(msg.String())
			}
		}
	}()
	return bc,nil
}

func (bc *BroChat) ProcessUserInput(input string){
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
			bc.userInterface.ProcessOutput(bc.state.String()+"\n")
			foundCommand = true
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

