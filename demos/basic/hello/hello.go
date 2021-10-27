package hello

import "time"

func GetMessage() string {
	//location, _ := time.LoadLocation("NZ")
	t := time.Now()
	msg := t.Local().Format("15:04:05 2006/01/02")
	return "The local time is: " + msg
}
