package log

import (
	"fmt"
	"time"
)

const dateTime_layout = "2006-01-02T15:04:05-0700"

func Info(format string, a... interface{}) {
	now := time.Now()
	pre := fmt.Sprintf("[INFO] - %s - ", now.Format(dateTime_layout))
	fmt.Println(pre + format, a)
}

func Warn(format string, a... interface{}) {
	now := time.Now()
	pre := fmt.Sprintf("[WARN] - %s ", now.Format(dateTime_layout))
	fmt.Println(pre + format, a)
}

func Error(format string, a... interface{}) {
	now := time.Now()
	pre := fmt.Sprintf("[ERROR] - %s ", now.Format(dateTime_layout))
	fmt.Println(pre + format, a)
}