package main

import (
	"flag"
	"log"
	"time"
	"vitalsign-publisher/server"
)

var (
	port   = flag.Int("port", 50051, "The server port")
	period = flag.Duration("period", 5000, "vitalsign-publisher working period(second)")
)

func main() {
	flag.Parse()

	vsp := &server.VSP{}
	serving := make(chan bool)

	go server.ServerStart(vsp, *port, serving)
	if !<-serving {
		log.Fatal("gRPC ServerStart FAIL")
		return
	}

	for {
		vsp.Mutex.Lock()
		log.Println(vsp.RPNs)
		vsp.Mutex.Unlock()
		time.Sleep((*period) * time.Millisecond)
	}
}
