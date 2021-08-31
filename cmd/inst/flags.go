package main

import (
	"log"

	flags "github.com/jessevdk/go-flags"
)

var opts struct {
	Passes []string `long:"passes" description:"A list of passes you want to use in this instrumentation"`
	Dir    string   `long:"dir" description:"Instrument all go source files under this directory"`
	File   string   `long:"file" description:"Instrument single go source file"`
	Out    string   `long:"out" description:"Output instrumented golang source file to the given file. Only allow when instrumenting single golang source file"`
	Args   struct {
		Globs []string
	} `positional-arg-name:"globs" positional-args:"yes"`
}

func parseFlags() {
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Panic(err)
	}
}
