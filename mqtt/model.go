package mqtt

import "vitalsign-publisher/mongodb"

type RPNPatientPublish struct {
	UserID          string       `json:"Patient_CodeID"`
	HR              int          `json:"HR"`
	SBP             int          `json:"SBP"`
	DBP             int          `json:"DBP"`
	MAP             int          `json:"MAP"`
	RESPIRATORY     int          `json:"RESPIRATORY"`
	SPO2            int          `json:"SPO2"`
	VO2MAX          int          `json:"VO2MAX"`
	MET             float64      `json:"MET"`
	Lasttime_12lead int          `json:"lasttime_12lead"`
	Lead2           []PublishECG `json:"lead2"`
	Lasttime_bp     int          `json:"lasttime_bp"`
	Bp              []mongodb.Bp `json:"bp"`
	VO2             float64      `json:"vo2"`
	CO              float64      `json:"co"`
}

type PublishECG struct {
	UserId string    `json:"userId"`
	Time   int       `json:"time"`
	II     []float64 `json:"II"`
}
