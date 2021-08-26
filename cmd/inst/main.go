package main

import (
	"gfuzz/pkg/inst"
	"gfuzz/pkg/inst/pass"
)

func main() {
	parseFlags()

	reg := inst.NewPassRegistry()

	// register passes
	reg.AddPass(&pass.SelEfcmPass{})
}
