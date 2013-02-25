package main

import (
	"log"
	"fmt"
	"os"
)

type BaseUI struct {
	Nickname	string
	Channel		string
	
	Sender		*BCastSender
	Receiver	*BCastReceiver
}

func (ui *BaseUI) Init(sender *BCastSender,receiver *BCastReceiver) error {
	ui.Sender = sender
	ui.Receiver = receiver
	ui.Channel = "#foyer"
	if ui.Nickname=="" {
		ui.Nickname = "Anonymous"
	}
	ui.SendMessage("joined the channel")
	return nil
}

func (ui *BaseUI) SendMessage(payload string) error {
	msg := NewBCastMessage(ui.Nickname,ui.Channel,payload)
	ui.Sender.SendMessage(msg)
	return nil
}

func (ui *BaseUI) ReadNextMessage() *BCastMessage {
	return ui.Receiver.GetMessage()
}

func (ui *BaseUI) ProcessUserInput(input string){
	foundCommand,err := ui.ProcessPossibleCommand(input)
	if foundCommand && err!=nil {
		log.Print("malformed command. grml...")
		return
	}
	if foundCommand || len(input)<1 {
		return
	}
	ui.SendMessage(input)
}

func (ui *BaseUI) ProcessPossibleCommand(payload string) (foundCommand bool,err error){
	cmd,err := NewUserCommand(payload)
	if err!=nil {
		return
	}
	switch cmd.Type {
		case UC_CHANGE_NICK: {
			ui.UpdateNickname(cmd.Payload)
			foundCommand = true
		}
		case UC_CHANGE_CHAN: {
			ui.JoinChannel(cmd.Payload)
			foundCommand = true
		}
		case UC_LEAVE: {
			ui.SendMessage("leaves the channel")
			os.Exit(0)
		}
	}
	return
}

func (ui *BaseUI) UpdateNickname(newnick string){
	/*
	TODO: Add RegExp Checker here
	*/
	if len(newnick)>0 {
		text := fmt.Sprintf("changed nick to %v",newnick)
		ui.SendMessage(text)
		ui.Nickname = newnick
	}
}

func (ui *BaseUI) JoinChannel(newchan string){
	text1 := "leaves the channel"
	text2 := "joins the channel"
	if ui.Channel!=""{
		ui.SendMessage(text1)
	}
	ui.Channel = newchan
	ui.SendMessage(text2)
}
