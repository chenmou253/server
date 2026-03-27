package logger

import (
	"log"
	"os"
)

func New() *log.Logger {
	return log.New(os.Stdout, "[go-admin] ", log.LstdFlags|log.Lshortfile)
}
