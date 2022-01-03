package hello

import (
	"time"
)

func Test(data []int8) {
	for _, b := range data {
		println("B is", b)
	}
}

func Test2(data []byte) {
	var data2 []int8
	data2 = make([]int8, len(data))
	for i, b := range data {
		data2[i] = int8(b)
	}
	Test(data2)
}

func GetMessage() string {
	//location, _ := time.LoadLocation("NZ")

	t := time.Now()
	msg := t.Local().Format("15:04:05 2006/01/02")
	return "The local time is: " + msg
}
