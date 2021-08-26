package main

import (
	"log"

	flags "github.com/jessevdk/go-flags"
)

var opts struct {
	Passes []string `long:"passes" description:"A list of passes you want to use in this instrumentation"`
	Dir    string   `long:"dir" description:"instrument all go source code under this directory"`
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
