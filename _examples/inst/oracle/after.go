package oracle

import oraclert "gfuzz/pkg/oraclert"

func TestHello() {
	oracleEntry := oraclert.BeforeRun()
	defer oraclert.AfterRun(oracleEntry)
	println("hello")
}
