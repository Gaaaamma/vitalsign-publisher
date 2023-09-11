package mqtt

import "vitalsign-publisher/mongodb"

type RPNPublish struct {
	Patients []RPNPatientPublish
}
type RPNPatientPublish struct {
	UserID          string
	HR              int
	SBP             int
	DBP             int
	MAP             int
	RESPIRATORY     int
	SPO2            int
	VO2MAX          int
	MET             float64
	Lasttime_12lead int
	Lead2           []PublishECG
	Lasttime_bp     int
	Bp              []mongodb.Bp

	Rehab_hr  []mongodb.Rehb_HR
	Rehab_VO2 []mongodb.Rehb_VO2
	Rehab_CO  []mongodb.Rehb_CO
}

type PublishECG struct {
	UserId string
	Time   int
	II     []float64
}
