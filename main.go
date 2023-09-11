package main

import (
	"flag"
	"fmt"
	"slices"
	"time"
	"vitalsign-publisher/common"
	"vitalsign-publisher/config"
	"vitalsign-publisher/mongodb"
	"vitalsign-publisher/request"
	"vitalsign-publisher/server"

	"github.com/fatih/color"
	"go.mongodb.org/mongo-driver/bson"
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

	// Mongo database setting
	ctx, _, client, err := mongodb.GetMongoClient()
	if err != nil {
		panic(err)
	}

	// colUser := client.Database(conf.MongoDB.Database).Collection(conf.MongoDB.User)
	// colRaw := client.Database(conf.MongoDB.Database).Collection(conf.MongoDB.Raw)
	// colEcg := client.Database(conf.MongoDB.Database).Collection(conf.MongoDB.Ecg)
	colVital := client.Database(conf.MongoDB.Database).Collection(conf.MongoDB.Vital)
	// colBp := client.Database(conf.MongoDB.Database).Collection(conf.MongoDB.BP)
	// colHR := client.Database(conf.MongoDB.Database).Collection(conf.MongoDB.HR)
	// colVO2 := client.Database(conf.MongoDB.Database).Collection(conf.MongoDB.VO2)
	// colCO := client.Database(conf.MongoDB.Database).Collection(conf.MongoDB.CO)

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

			// Step1-1 Filter patient that status == 4 (Uploading differential leads)
			data.Patients_list = slices.DeleteFunc(data.Patients_list, func(p request.Patient) bool {
				return p.Patient_Status != 4
			})

			// Step1-2 Filter rpn without patients
			if len(data.Patients_list) == 0 {
				color.Yellow("RPN %s: no any binding patient or uploading data (patient.status != 4)", rpn.Id)
				continue
			}

			color.HiGreen("%+v", data.Patients_list)

			// Step2. Packing data from mongoDB for each patient
			vital := mongodb.VitalSign{}
			filter := bson.M{"Patient_CodeID": "NCTU0000"}
			err = colVital.FindOne(ctx, filter).Decode(&vital)
			if err != nil {
				color.Red("%s", err)
			}
			color.HiGreen("%+v", vital)

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
