package main

import (
	"flag"
	"log"
	"os"
)

var Nickname *string = flag.String("nick","Anonymous","Your Nickname")
var BCastIP	 *string = flag.String("ip","","The broadcast IP to use")
var Port	 *int	 = flag.Int("port",1234,"The port used to broadcast")


func main(){
	flag.Parse()
	
	sender,err := NewBCastSender(*BCastIP,*Port)
	if err!=nil {
		log.Print("can't setup Broadcast Sender")
		log.Print(err)
		os.Exit(1)
	}
	receiver,err := NewBCastReceiver(*Port)
	if err!=nil {
		log.Print("can't setup Broadcast Receiver")
		log.Print(err)
		os.Exit(1)
	}
	
	ui := NewGtkUI(*Nickname)
	ui.Init(sender,receiver)
	ui.Run()
	select{}
}
