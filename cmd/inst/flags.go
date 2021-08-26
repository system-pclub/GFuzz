package main

import (
	"log"

	flags "github.com/jessevdk/go-flags"
)

var opts struct {
	Passes []string `long:"passes" description:"A list of passes you want to use in this instrumentation"`
	File   string   `long:"file" required:"true"`
}

func parseFlags() {
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Panic(err)
	}
}
