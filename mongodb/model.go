package mongodb

type User struct {
	UserID            string `bson:"userId"`
	Status            int    `bson:"Status"`
	PublishStatus     int    `bson:"publishStatus"`
	Transfer3To12     int    `bson:"transfer3To12"`
	SubStartTime      []int  `bson:"subStartTime"`
	LastTime3Lead     int    `bson:"lasttime_3lead"`
	LastTime12Lead    int    `bson:"lasttime_12lead"`
	LastTimeSCG       int    `bson:"lasttime_scg"`
	LastTimeGCG       int    `bson:"lasttime_gcg"`
	LastTimeMCG       int    `bson:"lasttime_mcg"`
	LastTime3Denoised int    `bson:"lasttime_3denoised"`
	Visible           int    `bson:"visible"`
	DenoiseON         int    `bson:"Denoise_ON"`
	AFON              int    `bson:"AF_ON"`
	MION              int    `bson:"MI_ON"`
	HFON              int    `bson:"HF_ON"`
	VFON              int    `bson:"VF_ON"`
	AFDetect          int    `bson:"AF_detect"`
	MIDetect          int    `bson:"MI_detect"`
	VFDetect          int    `bson:"VF_detect"`
	HFDetect          int    `bson:"HF_detect"`
	HFONRef           int    `bson:"HF_ON_ref"`
	AFONRef           int    `bson:"AF_ON_ref"`
	MIONRef           int    `bson:"MI_ON_ref"`
	VFONRef           int    `bson:"VF_ON_ref"`
	LastTimeBP        int    `bson:"lasttime_bp"`
}

type RT_ECG struct {
	UserID    string    `bson:"userId"`
	Timestamp int       `bson:"time"`
	I         []float64 `bson:"I"`
	II        []float64 `bson:"II"`
	III       []float64 `bson:"III"`
	AVR       []float64 `bson:"aVR"`
	AVL       []float64 `bson:"aVL"`
	AVF       []float64 `bson:"aVF"`
	V1        []float64 `bson:"V1"`
	V2        []float64 `bson:"V2"`
	V3        []float64 `bson:"V3"`
	V4        []float64 `bson:"V4"`
	V5        []float64 `bson:"V5"`
	V6        []float64 `bson:"V6"`
}

type RT_VitalSign struct {
	UserID      string  `bson:"Patient_CodeID"`
	HR          int     `bson:"HR"`
	SBP         int     `bson:"SBP"`
	DBP         int     `bson:"DBP"`
	RESPIRATORY int     `bson:"RESPIRATORY"`
	SPO2        int     `bson:"SPO2"`
	VO2MAX      int     `bson:"VO2MAX"`
	MET         float64 `bson:"MET"`
}

type Bp struct {
	UserID    string    `bson:"userId"`
	Timestamp int       `bson:"time"`
	Value     []float64 `bson:"bp"`
}

type Rehb_HR struct {
	UserID    string `bson:"Patient_CodeID"`
	Timestamp int    `bson:"timestamp"`
	Value     int    `bson:"heart_rate"`
}

type Rehb_CO struct {
	UserID    string  `bson:"Patient_CodeID"`
	Timestamp int     `bson:"timestamp"`
	Value     float64 `bson:"CO"`
}

type Rehb_VO2 struct {
	UserID    string `bson:"Patient_CodeID"`
	Timestamp int    `bson:"timestamp"`
	Value     int    `bson:"VO2"`
}
