package main

import (
	"flag"
	"time"
	"vitalsign-publisher/common"
	"vitalsign-publisher/config"
	"vitalsign-publisher/server"

	"github.com/fatih/color"
)

var (
	conf   = config.GetConfig()
	port   = flag.Int("port", conf.Setting.Port, "The server port")
	period = flag.Duration("period", time.Duration(conf.Setting.SleepTime), "vitalsign-publisher working period(second)")
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
		color.Cyan("%v RPNs: %v", common.TimeNow(), vsp.RPNs)
		vsp.MuRpn.Unlock()

		vsp.MuPatient.Lock()
		color.Cyan("%v Patients: %v", common.TimeNow(), vsp.Patients)
		vsp.MuPatient.Unlock()
		time.Sleep((*period) * time.Millisecond)
	}
}
