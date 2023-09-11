package main

import (
	"flag"
	"fmt"
	"slices"
	"time"
	"vitalsign-publisher/common"
	"vitalsign-publisher/config"
	"vitalsign-publisher/request"
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

	// Serving Loop
	for {
		// ==================== Serving RPNs ====================
		vsp.MuRpn.Lock()
		rpns := slices.Clone(vsp.RPNs)
		vsp.MuRpn.Unlock()

		color.Cyan("%v RPNs: %v", common.TimeNow(), rpns)
		for _, rpn := range rpns {
			// Step1. Get patient list from API server
			data, err := request.GetPatientsByRPN(rpn)
			if err != nil {
				color.Red("%s", err)
				continue
			}

			// Step1-1 Filter rpn without patients
			if len(data.Patients_list) == 0 {
				color.Yellow("RPN %s: no any binding patient", rpn.Id)
				continue
			}

			// Step1-2 Filter patient that status == 4 (Uploading differential leads)
			color.HiGreen("%+v", data.Patients_list)

			// Step2. Packing data from mongoDB for each patient

			// Step2-1. Get VitalSign data (Default value is set if data didn't exist)

			// Step2-2. Get Ecg standard 12 lead data

			// Step2-3. Get BP data

			// Step3. Push data to MQTT broker

		}

		// ==================== Serving Patients ====================
		vsp.MuPatient.Lock()
		patients := slices.Clone(vsp.Patients)
		vsp.MuPatient.Unlock()

		color.Cyan("%v Patients: %v", common.TimeNow(), patients)

		fmt.Println()
		time.Sleep((*period) * time.Millisecond)
	}
}
