package main

type UI interface {
	Init(sender *BCastSender,receiver *BCastReceiver) error
	Run() error
}



