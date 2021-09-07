package oracle

import gooracle "gooracle"

func TestHello() {
	oracleEntry := gooracle.BeforeRun()
	defer gooracle.AfterRun(oracleEntry)
	println("hello")
}
