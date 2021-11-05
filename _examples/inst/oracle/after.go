package oracle

func TestHello() {
	oracleEntry := oraclert.BeforeRun()
	defer oraclert.AfterRun(oracleEntry)
	println("hello")
}
