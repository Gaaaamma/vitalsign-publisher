package main

import (
	"flag"
	"time"
	"vitalsign-publisher/common"
	"vitalsign-publisher/server"

	"github.com/fatih/color"
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
		color.Red("server.ServerStart: FAIL - gRPC ServerStart isn't serving")
		return
	}

	for {
		vsp.MuRpn.Lock()
		color.Cyan("%v %v", common.TimeNow(), vsp.RPNs)
		vsp.MuRpn.Unlock()
		time.Sleep((*period) * time.Millisecond)
	}
}
