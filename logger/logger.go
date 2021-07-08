package logger

import (
	"log"
	"os"
	"x-net.at/idp/helpers"
)

var (
	Warning   *log.Logger
	Info      *log.Logger
	Error     *log.Logger
	debugging *log.Logger
)

func Init() {
	Info = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	Warning = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	debugging = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime)
}

func Debug(message string) {
	if helpers.Config.Debug {
		debugging.Println(message)
	}
}
