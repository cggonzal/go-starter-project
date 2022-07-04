package logger

import (
	"log"
	"os"
)

var (
	Logger *log.Logger
)

func InitLogger() *log.Logger {
	logger := log.New(os.Stdout, "logger: ", log.LstdFlags|log.Llongfile)
	return logger
}
