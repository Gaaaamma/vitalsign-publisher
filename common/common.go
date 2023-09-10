package common

import (
	"fmt"
	"time"
)

func TimeNow() string {
	time := time.Now()
	return fmt.Sprintf("%04d/%02d/%02d %02d:%02d:%02d", time.Year(), time.Month(), time.Day(),
		time.Hour(), time.Minute(), time.Second())
}
