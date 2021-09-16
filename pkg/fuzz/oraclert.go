package fuzz

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// OracleRtInput contains all information need to be passed to fuzz target (consumed by oracle runtime)
type OracleRtInput struct {
	// First line of input file, string `PrintInput` for recording input (used by gooracle),
	// otherwise it is just a placeholder for the first line of file to make sure file has correct format.
	Note string

	// How many milliseconds a select will wait for the prioritized case
	SelectDelayMS int

	// Select choice need to be forced during runtime
	VecSelect []SelectInput
}

type SelectInput struct {
	StrFileName string
	IntLineNum  int
	IntNumCase  int
	IntPrioCase int
}

func ShouldSkipInput(fuzzCtx *FuzzContext, exec string, i *OracleRtInput) bool {
	if i == nil {
		return true
	}

	// drop it if list of selects has been run more than 3 times
	if GetCountsOfSelects(exec, i.VecSelect) > 5 {
		log.Printf("drop %s since the select inputs run more than 3 times", i)
		return false
	}

	// drop it if it's source already timeout more than 3 times
	fuzzCtx.timeoutTargetsLock.RLock()
	timeoutCnt, exist := fuzzCtx.timeoutTargets[exec]
	fuzzCtx.timeoutTargetsLock.RUnlock()
	if exist {
		if timeoutCnt > 3 {
			log.Printf("drop %s since it has timeout more than 3 times", i)
			return true
		}
	}

	return false
}

func GetHashOfSelects(selects []SelectInput) string {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%v", selects)))
	return fmt.Sprintf("%x", h.Sum(nil))
}

const (
	NotePrintInput string = "PrintInput"
)

func DumpToFile(i *OracleRtInput, dstFile string) error {
	out, err := os.Create(dstFile)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = out.WriteString(i.String())
	if err != nil {
		return err
	}
	return nil
}

func (i *OracleRtInput) String() string {
	retStr := ""

	// The first line is a note. Could be empty.
	if i.Note == NotePrintInput {
		retStr += NotePrintInput + "\n"
	} else {
		retStr += "\n"
	}

	// The second line is how many milliseconds to wait
	retStr += strconv.Itoa(i.SelectDelayMS) + "\n"

	// From the third line, each line corresponds to a select
	for _, selectInput := range i.VecSelect {
		// filename:linenum:totalCaseNum:chooseCaseNum
		str := selectInput.StrFileName + ":" + strconv.Itoa(selectInput.IntLineNum)
		str += ":" + strconv.Itoa(selectInput.IntNumCase)
		str += ":" + strconv.Itoa(selectInput.IntPrioCase)
		str += "\n"
		retStr += str
	}
	return retStr
}

func (i *OracleRtInput) DeepCopy() *OracleRtInput {
	newInput := &OracleRtInput{
		Note:          i.Note,
		SelectDelayMS: i.SelectDelayMS,
		VecSelect:     []SelectInput{},
	}
	for _, selectInput := range i.VecSelect {
		newInput.VecSelect = append(newInput.VecSelect, copySelectInput(selectInput))
	}
	return newInput
}

func HashOfRecord(record *Record) string {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%v", record))) // TODO: we may need to replace `record` with `StrOfRecord(record)`
	return fmt.Sprintf("%x", h.Sum(nil))
}

func FindRecordHashInSlice(recordHash string, recordHashSlice []string) bool {
	for _, searchRecordHash := range recordHashSlice {
		if recordHash == searchRecordHash {
			return true
		}
	}
	return false
}

func DeserializeOracleRtInput(content string) (*OracleRtInput, error) {
	var err error

	lines := strings.Split(content, "\n")
	if len(lines) < 2 {
		return nil, errors.New("Input has less than 2 lines")
	}

	newInput := &OracleRtInput{
		Note:          lines[0],
		SelectDelayMS: 0,
		VecSelect:     []SelectInput{},
	}

	strDelayMS := lines[1]
	newInput.SelectDelayMS, err = strconv.Atoi(strDelayMS)
	if err != nil {
		return nil, err
	}

	// Skip line 1 (PrintInput) and line 2 (select time out)
	selectInfos := lines[2:]
	for _, eachLine := range selectInfos {
		if eachLine == "" {
			continue
		}
		selectInput, err := ParseSelectInput(eachLine)
		if err != nil {
			return nil, err
		}

		newInput.VecSelect = append(newInput.VecSelect, *selectInput)
	}

	return newInput, nil
}

// ParseSelectInput parses the each select in input file
// which has format filename:linenum:totalCaseNum:chooseCaseNum
func ParseSelectInput(line string) (*SelectInput, error) {
	var err error
	selectInput := SelectInput{}
	vecStr := strings.Split(line, ":")
	if len(vecStr) != 4 {
		return nil, fmt.Errorf("expect number of components: 4, actual: %d", len(vecStr))
	}
	selectInput.StrFileName = vecStr[0]
	if selectInput.IntLineNum, err = strconv.Atoi(vecStr[1]); err != nil {
		return nil, fmt.Errorf("incorrect format at line number")
	}
	if selectInput.IntNumCase, err = strconv.Atoi(vecStr[2]); err != nil {
		return nil, fmt.Errorf("incorrect format at number of cases")
	}
	if selectInput.IntPrioCase, err = strconv.Atoi(vecStr[3]); err != nil {
		return nil, fmt.Errorf("incorrect format at priority case")
	}
	return &selectInput, nil
}

func copySelectInput(sI SelectInput) SelectInput {
	return SelectInput{
		StrFileName: sI.StrFileName,
		IntLineNum:  sI.IntLineNum,
		IntNumCase:  sI.IntNumCase,
		IntPrioCase: sI.IntPrioCase,
	}
}
