// customLogger should be used for all logging in the project
package customLogger

import (
	"log"
	"os"
)

var (
	logger *log.Logger
)

// initialize logger. Needs to be called once at the beginning of the program.
func InitLogger() {
	logger = log.New(os.Stdout, "logger: ", log.LstdFlags|log.Llongfile)
}

// returns logger that should be used for all logging in the project
func GetLogger() *log.Logger {
	return logger
}
