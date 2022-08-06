package logger

import (
	"log"
	"os"
)

// Logger should be used for all logging in the project
var (
	Logger *log.Logger
)

func InitLogger() {
	Logger = log.New(os.Stdout, "logger: ", log.LstdFlags|log.Llongfile)
}
