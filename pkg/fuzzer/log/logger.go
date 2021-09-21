package log

import (
	"io"
	"log"
	"os"
)

var (
	logWriter io.Writer
)

const (
	logFlag = log.Ldate | log.Ltime | log.Lmsgprefix
)

func SetupLogger(logFile string) {
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal(err)
	}
	logWriter = io.MultiWriter(file, os.Stdout)
	log.SetFlags(logFlag)
	log.SetOutput(logWriter)
}

func NewLogger(prefix string) *log.Logger {
	return log.New(logWriter, prefix, logFlag)
}
