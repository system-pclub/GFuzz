package main

import (
	"io"
	"log"
	"os"
)

const (
	GFUZZ_LOG_FILE = "fuzzer.log"
)

func setupLogger(logFile string) {

	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal(err)
	}
	w := io.MultiWriter(file, os.Stdout)
	log.SetOutput(w)
}
