package utils

import (
	"encoding/json"
	"fmt"
	"time"
)
// RFC850
const (
	RFC850  = "Monday, 02-Jan-06 15:04:05 MST"
	ISO8601 = "2006-01-02 15:04:05"
)
// Log is for export.
type Log struct {
	Lvl     string `json:"level"`
	Time    string `json:"time"`
	Message string `json:"msg"`
}
// Logging is for export.
func Logging(lvl, msg string) {
	t := time.Now()
	log := Log{Lvl: lvl, Time: t.Format(ISO8601), Message: msg}
	jsonData, err := json.Marshal(log)
	if err != nil {
		fmt.Println("[ERROR]: Failed to marshal json :(")
	}
	fmt.Println(string(jsonData))
}
