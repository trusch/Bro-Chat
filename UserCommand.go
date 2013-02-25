package main

import (
	"strings"
	"errors"
)

type UserCommandType int

const (
	UC_SEND_MSG		= iota
	UC_CHANGE_NICK
	UC_CHANGE_CHAN
	UC_LEAVE
)

type UserCommand struct {
	Type	UserCommandType
	Payload	string
}

func NewUserCommand(payload string) (*UserCommand,error) {
	cmd := new(UserCommand)
	var err error
	parts := strings.SplitN(payload," ",2)
	firstWord := parts[0]
	if len(firstWord)>0 && firstWord[0]=='/' {
		switch{
			case strings.HasPrefix(firstWord,"/j"): {
				if len(parts)<2 {
					err = errors.New("malformed command")
					break
				}
				cmd.Type = UC_CHANGE_CHAN
				cmd.Payload = parts[1]
			}
			case strings.HasPrefix(firstWord,"/n"): {
				if len(parts)<2 {
					err = errors.New("malformed command")
					break
				}
				cmd.Type = UC_CHANGE_NICK
				cmd.Payload = parts[1]
			}
			case strings.HasPrefix(firstWord,"/l")||strings.HasPrefix(firstWord,"/q"): {
				cmd.Type = UC_LEAVE
			}
		}
	}
	if err!=nil {
		return nil,err
	}
	if cmd.Type==UC_SEND_MSG {
		cmd.Payload = payload
	}
	return cmd,nil
}

