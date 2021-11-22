package log

import (
	"io"
	"log"
	"os"
)

var (
	logWriter  io.Writer // could be io.MultiWriter or simply other writer
	fileWriter io.Writer // log file writer
)

const (
	logFlag = log.Ldate | log.Ltime | log.Lmsgprefix
)

func SetupLogger(logFile string, stdout bool) {
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal(err)
	}
	var writers []io.Writer
	fileWriter = file
	writers = append(writers, fileWriter)
	if stdout {
		writers = append(writers, os.Stdout)
	}
	logWriter = io.MultiWriter(writers...)
	log.SetFlags(logFlag)
	log.SetOutput(logWriter)
}

func DisableStdoutLog() {
	log.SetOutput(fileWriter)
}
func NewLogger(prefix string) *log.Logger {
	return log.New(logWriter, prefix, logFlag)
}
