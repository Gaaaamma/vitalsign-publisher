package request

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"vitalsign-publisher/config"
	"vitalsign-publisher/protos"
)

var conf = config.GetConfig()

type PatientList struct {
	Date          string    `json:"Date"`
	Patients_list []Patient `json:"Patients_list"`
}

type Patient struct {
	Alarm_Status     int    `json:"Alarm_Status"`
	Battery          int    `json:"Battery"`
	HR               int    `json:"HR"`
	HR_Timestamp     string `json:"HR_Timestamp"`
	Bed_Id           string `json:"Bed_Id"`
	Device_Id        int    `json:"Device_Id"`
	MRN              string `json:"MRN"`
	Patient_CodeID   string `json:"Patient_CodeID"`
	Patient_EnName   string `json:"Patient_EnName"`
	Patient_Name     string `json:"Patient_Name"`
	Patient_Status   int    `json:"Patient_Status"`
	RPN_Id           string `json:"RPN_Id"`
	Room_Id          int    `json:"Room_Id"`
	Attending_Doctor string `json:"Attending_Doctor"`
	AI_Model         struct {
		MI  int `json:"MI"`
		VF  int `json:"VF"`
		AF  int `json:"AF"`
		HF  int `json:"HF"`
		VT  int `json:"VT"`
		PVC int `json:"PVC"`
	} `json:"AI_Model"`
}

func GetPatientsByRPN(rpn *protos.RPN) (PatientList, error) {
	// Get patient list from API server
	url := fmt.Sprintf("http://%s:%d%s%s", conf.Api.Host, conf.Api.Port, conf.Api.RpnPatientList, rpn.Id)
	resp, err := http.Get(url)
	if err != nil {
		return PatientList{}, fmt.Errorf("RPN %s: %s", rpn.Id, err)
	}
	defer resp.Body.Close()

	read, err := io.ReadAll(resp.Body)
	if err != nil {
		return PatientList{}, fmt.Errorf("RPN %s - ReadAll: %s", rpn.Id, err)
	}

	// Parsing
	data := PatientList{}
	err = json.Unmarshal(read, &data)
	if err != nil {
		return PatientList{}, fmt.Errorf("RPN %s - Unmarshal: %s", rpn.Id, err)
	}
	return data, nil
}
