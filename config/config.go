package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Define a struct that matches the structure of your YAML configuration.
type Config struct {
	Setting struct {
		Host           string `yaml:"host"`
		Port           int    `yaml:"port"`
		RpnBase        string `yaml:"rpn_base"`
		PatientBase    string `yaml:"patient_base"`
		HrExpireSecs   int    `yaml:"hrExpireSecs"`
		FloatPrecision int    `yaml:"floatPrecision"`
		LogMode        int    `yaml:"logMode"`
		SleepTime      int    `yaml:"sleepTime"`
	} `yaml:"setting"`

	Api struct {
		Host           string `yaml:"host"`
		Port           int    `yaml:"port"`
		RpnPatientList string `yaml:"rpn_patient_list"`
	} `yaml:"api"`

	Wisepaas struct {
		Host         string `yaml:"host"`
		Port         int    `yaml:"port"`
		Websocket    int    `yaml:"websocket"`
		MqttUser     string `yaml:"mqtt_user"`
		MqttPassword string `yaml:"mqtt_password"`
	} `yaml:"wisepaas"`

	MongoDB struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Database string `yaml:"database"`
		User     string `yaml:"colUser"`
		Raw      string `yaml:"colRaw"`
		Ecg      string `yaml:"colEcg"`
		Vital    string `yaml:"colVital"`
		BP       string `yaml:"colBp"`
		HR       string `yaml:"rehab_hr"`
		VO2      string `yaml:"rehab_VO2"`
		CO       string `yaml:"rehab_CO"`
	} `yaml:"mongoDB"`
}

func GetConfig() Config {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	path := fmt.Sprintf("%s/%s", cwd, "config/config.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	conf := Config{}
	err = yaml.Unmarshal(data, &conf)

	if err != nil {
		panic(err)
	}
	return conf
}
