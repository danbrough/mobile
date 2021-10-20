package hello


import "time"

func GetMessage() string {
	//location, _ := time.LoadLocation("NZ")
	t := time.Now()
	msg := t.Local().Format("2006/01/02 15:04:05")
	return "Hello World at " + msg
}
