package main

import (
	"flag"
	"fmt"
	"slices"
	"time"
	"vitalsign-publisher/common"
	"vitalsign-publisher/config"
	"vitalsign-publisher/mongodb"
	"vitalsign-publisher/mqtt"
	"vitalsign-publisher/request"
	"vitalsign-publisher/server"

	"github.com/fatih/color"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	colUser := client.Database(conf.MongoDB.Database).Collection(conf.MongoDB.User)
	colEcg := client.Database(conf.MongoDB.Database).Collection(conf.MongoDB.Ecg)
	colVital := client.Database(conf.MongoDB.Database).Collection(conf.MongoDB.Vital)
	colBp := client.Database(conf.MongoDB.Database).Collection(conf.MongoDB.BP)
	colHR := client.Database(conf.MongoDB.Database).Collection(conf.MongoDB.HR)
	colVO2 := client.Database(conf.MongoDB.Database).Collection(conf.MongoDB.VO2)
	colCO := client.Database(conf.MongoDB.Database).Collection(conf.MongoDB.CO)

	// Serving Loop
	for {
		// ==================== Serving RPNs ====================
		vsp.MuRpn.Lock()
		rpns := slices.Clone(vsp.RPNs)
		vsp.MuRpn.Unlock()

		color.Cyan("%v RPNs: %v", common.TimeNow(), rpns)
		for _, rpn := range rpns {
			rpnPublish := mqtt.RPNPublish{}

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

			// Step2. Packing data from mongoDB for each patient
			// rpnUploading := mqtt.RPNPublish{}
			for _, p := range data.Patients_list {
				checker := map[string]float64{
					"VitalSign": 0,
					"RT_ECG":    0,
					"BP":        0,
					"HR":        0,
					"VO2":       0,
					"CO":        0,
				}
				record := mqtt.RPNPatientPublish{UserID: p.Patient_CodeID}

				// Step2-1. Get VitalSign data
				filter := bson.M{"Patient_CodeID": p.Patient_CodeID}
				if err = colVital.FindOne(ctx, filter).Decode(&record); err != nil {
					color.Red("Get VitalSign Data %s: %s", p.Patient_CodeID, err)
					continue
				}
				checker["VitalSign"] += 1

				// Step2-2. Get Ecg standard 12 lead data
				// Query ecg.user
				userData := mongodb.User{}
				filter = bson.M{"userId": p.Patient_CodeID}
				if err = colUser.FindOne(ctx, filter).Decode(&userData); err != nil {
					color.Red("Get User Data %s: %s", p.Patient_CodeID, err)
					continue
				}
				record.Lasttime_12lead = userData.LastTime12Lead
				record.Lasttime_bp = userData.LastTimeBP

				// Query ecg.ecgdata12
				filter = bson.M{
					"userId": p.Patient_CodeID,
					"time": bson.M{
						"$gt":  userData.LastTime12Lead - 5,
						"$lte": userData.LastTime12Lead,
					},
				}
				if cursor, err := colEcg.Find(ctx, filter); err != nil {
					color.Red("Get RT_ECG %s: %s", p.Patient_CodeID, err)
				} else {
					for cursor.Next(ctx) {
						ecg := mongodb.RT_ECG{}
						if err = cursor.Decode(&ecg); err != nil {
							color.Red("Cursor decoded ecg %s: %s", p.Patient_CodeID, err)
							continue
						}
						record.Lead2 = append(record.Lead2, mqtt.PublishECG{UserId: ecg.UserID, Time: ecg.Timestamp, II: ecg.II})
						checker["RT_ECG"] += 1
					}
					err = cursor.Close(ctx)
					if err != nil {
						color.Red("Close RT_ECG %s: %s", p.Patient_CodeID, err)
					}
				}

				// Step2-3. Get BP data
				filter = bson.M{
					"userId": p.Patient_CodeID,
					"time": bson.M{
						"$gt":  userData.LastTimeBP - 5,
						"$lte": userData.LastTimeBP,
					},
				}
				if cursor, err := colBp.Find(ctx, filter); err != nil {
					color.Red("Get Bp %s: %s", p.Patient_CodeID, err)
				} else {
					for cursor.Next(ctx) {
						bp := mongodb.Bp{}
						if err = cursor.Decode(&bp); err != nil {
							color.Red("Cursor decoded bp %s: %s", p.Patient_CodeID, err)
							continue
						}
						record.Bp = append(record.Bp, bp)
						checker["BP"] += 1
					}
					err = cursor.Close(ctx)
					if err != nil {
						color.Red("Close BP %s: %s", p.Patient_CodeID, err)
					}
				}

				// Step2-4. Get HR data
				filter = bson.M{"Patient_CodeID": p.Patient_CodeID}
				options := options.FindOne().SetSort(bson.D{{Key: "timestamp", Value: -1}})
				hr := mongodb.Rehb_HR{}
				if err := colHR.FindOne(ctx, filter, options).Decode(&hr); err != nil {
					color.Yellow("Get HR %s: %s", p.Patient_CodeID, err)
					record.HR = -1
				} else {
					record.HR = hr.Value
					checker["HR"] = float64(hr.Value)
				}

				// Step2-5. Get VO2 data
				vo2 := mongodb.Rehb_VO2{}
				if err := colVO2.FindOne(ctx, filter, options).Decode(&vo2); err != nil {
					color.Yellow("Get VO2 %s: %s", p.Patient_CodeID, err)
					record.VO2 = -1.0
				} else {
					record.VO2 = vo2.Value
					checker["VO2"] = vo2.Value
				}

				// Step2-6. Get CO data
				co := mongodb.Rehb_CO{}
				if err := colCO.FindOne(ctx, filter, options).Decode(&co); err != nil {
					color.Yellow("Get CO %s: %s", p.Patient_CodeID, err)
					record.CO = -1.0
				} else {
					record.CO = co.Value
					checker["CO"] = co.Value
				}

				// Step3. Pack record to rpnPublish
				common.DataChecker(p.Patient_CodeID, checker)
				rpnPublish.Patients = append(rpnPublish.Patients, record)
			}
			// Step4. Publish data to MQTT broker

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
