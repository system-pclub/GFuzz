package fuzz

import "testing"

func TestParseInputFileHappy(t *testing.T) {
	input, err := DeserializeOracleRtInput("PrintInput\n0\nabc.go:23:4:1")
	if err != nil {
		t.Fail()
	}
	if input.Note != "PrintInput" {
		t.Fail()
	}
	if len(input.VecSelect) != 1 {
		t.Fail()
	}

}

func TestParseInputFileShouldFail(t *testing.T) {
	_, err := DeserializeOracleRtInput("PrintInput\nabc.go:23:4:1")
	if err == nil {
		t.Fail()
	}
}

func TestSelectInputHappy(t *testing.T) {
	input, err := ParseSelectInput("abc.go:23:4:1")
	if err != nil {
		t.Fail()
	}
	if input.StrFileName != "abc.go" {
		t.Fail()
	}
	if input.IntLineNum != 23 {
		t.Fail()
	}

	if input.IntNumCase != 4 {
		t.Fail()
	}

	if input.IntPrioCase != 1 {
		t.Fail()
	}
}

func TestSelectInputShouldFail(t *testing.T) {
	_, err := ParseSelectInput("abc.go:23:4")
	if err == nil {
		t.Fail()
	}

	_, err = ParseSelectInput("abc.go")
	if err == nil {
		t.Fail()
	}
}
