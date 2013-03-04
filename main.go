package main

import (
	"flag"
	"log"
)

var Nickname *string = flag.String("nick","Anonymous","Your Nickname")
var BCastIP	 *string = flag.String("ip","","The broadcast IP to use")
var Port	 *int	 = flag.Int("port",1234,"The port used to broadcast")


func main(){
	flag.Parse()
	_,err := NewBroChat(*BCastIP,*Port)
	if err!=nil {
		log.Print(err)
		return
	}
	log.Print("All started!")
	select{}
}
