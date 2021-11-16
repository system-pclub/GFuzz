package main

import (
	"os"

	flags "github.com/jessevdk/go-flags"
)

var opts struct {
	Passes   []string `long:"pass" description:"A list of passes you want to use in this instrumentation"`
	Dir      string   `long:"dir" description:"Instrument all go source files under this directory"`
	File     string   `long:"file" description:"Instrument single go source file"`
	Out      string   `long:"out" description:"Output instrumented golang source file to the given file. Only allow when instrumenting single golang source file"`
	StatsOut string   `long:"statsOut" description:"Output statistics"`
	Args     struct {
		Globs []string
	} `positional-arg-name:"globs" positional-args:"yes"`
	Version         bool   `long:"version" description:"Print version and exit"`
	IgnoreSyntaxErr bool   `long:"ignore-syntax-err"`
	Parallel        int    `long:"parallel" default:"5"`
	CPUProfile      string `long:"cpuprofile"`
}

func parseFlags() {
	if _, err := flags.Parse(&opts); err != nil {
		switch flagsErr := err.(type) {
		case flags.ErrorType:
			if flagsErr == flags.ErrHelp {
				os.Exit(0)
			}
			os.Exit(1)
		default:
			os.Exit(1)
		}
	}
}
