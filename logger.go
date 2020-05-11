package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"time"
)

var (
	Logger *log.Logger
)

// Initializing the logger and customizing prefix
func initLogger() {
	f, err := os.OpenFile("run.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	Logger = log.New(f, "", 0)
	Logger.SetPrefix("[" + time.Now().Format("02-Jan-2006 15:04:05 MST") + "]: ")
}

// Logs to log file
// takes generic object and then based on the type of object,
// logs is in appropriate style.
func logIt(val interface{}, console ...bool) {
	// if console is passed, print to stdout as well
	if len(console) != 0 {
		fmt.Println(val)
	}

	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.Map:
		str, _ := json.MarshalIndent(val, "", "  ")
		Logger.Println(string(str))
	default:
		Logger.Println(val)
	}
}
