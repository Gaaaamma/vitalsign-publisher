package common

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

func TimeNow() string {
	time := time.Now()
	return fmt.Sprintf("%04d/%02d/%02d %02d:%02d:%02d", time.Year(), time.Month(), time.Day(),
		time.Hour(), time.Minute(), time.Second())
}

func DataChecker(id string, m map[string]float64) {
	success := true
	for k, v := range m {
		if (k == "RT_ECG" || k == "BP") && v < 5 {
			success = false
			break
		} else if (k == "HR" || k == "VO2" || k == "CO") && v < 0 {
			success = false
			break
		}
	}
	if success {
		color.Green("%s: (%s=%.0f) (%s=%.0fs) (%s=%.0fs) (%s=%.0f) (%s=%.1f) (%s=%.1f)", id, "VitalSign", m["VitalSign"], "RT_ECG", m["RT_ECG"], "BP", m["BP"], "HR", m["HR"], "VO2", m["VO2"], "CO", m["CO"])
	} else {
		color.Yellow("%s: (%s=%.0f) (%s=%.0fs) (%s=%.0fs) (%s=%.0f) (%s=%.1f) (%s=%.1f)", id, "VitalSign", m["VitalSign"], "RT_ECG", m["RT_ECG"], "BP", m["BP"], "HR", m["HR"], "VO2", m["VO2"], "CO", m["CO"])
	}
}
