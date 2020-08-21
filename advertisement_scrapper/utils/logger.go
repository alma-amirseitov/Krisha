package utils

import (
	"encoding/json"
	"fmt"
	"time"
)
const (
	ISO8601 = "2006-01-02 15:04:05"
)

type Log struct {
	Lvl     string `json:"level"`
	Time    string `json:"time"`
	Message string `json:"msg"`
}

func Logging(lvl, msg string) {
	t := time.Now()
	log := Log{Lvl: lvl, Time: t.Format(ISO8601), Message: msg}
	jsonData, err := json.Marshal(log)
	if err != nil {
		fmt.Println("[ERROR]: Failed to marshal json")
	}
	fmt.Println(string(jsonData))
}
