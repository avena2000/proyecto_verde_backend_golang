package database

import (
	"log"
	"time"
)

type DBLogger struct{}

func (l DBLogger) Print(v ...interface{}) {
	if v[0] == "sql" {
		log.Printf("SQL [%s] %s %v", v[2], v[3], time.Duration(v[1].(float64)))
	}
}