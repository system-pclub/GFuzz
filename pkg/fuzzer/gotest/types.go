package gotest

type GoTest struct {
	// If test should be triggered from compiled binary file
	Bin string

	// Test function name
	Func string

	// Which package the test function located (if bin has set, package is bin file name since test bin is compiled by package level)
	Package string
}
