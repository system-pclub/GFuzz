package oraclert

import (
	"fmt"
	"gfuzz/pkg/oraclert/config"
	"gfuzz/pkg/oraclert/output"
	"os"
	"runtime"
	"sort"
	"strings"
)

// getChans returns all the channels touched during
// runtime and their status/properties
func getChans() map[string]output.ChanRecord {
	var chans map[string]output.ChanRecord = make(map[string]output.ChanRecord)
	for _, chr := range runtime.ChRecord {
		if chr == nil {
			continue
		}
		chans[chr.StrCreation] = output.ChanRecord{
			ID:        chr.StrCreation,
			Closed:    chr.Closed,
			NotClosed: chr.NotClosed,
			CapBuf:    int(chr.CapBuf),
			PeakBuf:   int(chr.PeakBuf),
		}
	}
	return chans
}

func GenerateOracleRtOutput(config *config.Config) *output.Output {
	output := &output.Output{
		Channels: getChans(),
		Tuples:   getTuples(),
		Ops:      getOps(),
	}
	if config.DumpSelects {
		output.Selects = getSelects()
	}

	return output
}

func DumpOracleRtOutput(config *config.Config, outputFile string) error {
	if outputFile == "" {
		return nil
	}
	o := GenerateOracleRtOutput(config)
	bytes, err := output.Serialize(o)
	if err != nil {
		return err
	}
	err = os.WriteFile(outputFile, bytes, 0666)
	if err != nil {
		return err
	}
	return nil
}

// getOps returns all operation ID touched by the program in runtime
func getOps() []uint16 {
	return opRecords
}

// getTuples returns all location touched by the program in runtime
func getTuples() map[uint32]uint32 {
	var tuples map[uint32]uint32 = make(map[uint32]uint32)
	for xorLoc, count := range runtime.TupleRecord {
		if count == 0 {
			continue // no need to record tuple that doesn't show up at all
		}
		tuples[uint32(xorLoc)] = count
	}
	return tuples
}

// getSelects returns all selected cases during the runtime
func getSelects() []output.SelectRecord {
	var selects []output.SelectRecord

	runtime.ProcessSelectInfo(func(selectInfos map[string]runtime.SelectInfo) {
		fmt.Printf("[oraclert]: %d selects\n", len(selectInfos))
		for _, selectInput := range selectInfos {
			// filename:linenum:totalCaseNum:chooseCaseNum
			strFileName := selectInput.StrFileName
			if indexEnter := strings.Index(strFileName, "\n"); indexEnter > -1 {
				strFileName = strFileName[:indexEnter]
			}
			id := fmt.Sprintf("%s:%s", strFileName, selectInput.StrLineNum)
			selects = append(selects, output.SelectRecord{
				ID:     id,
				Cases:  uint(selectInput.IntNumCase),
				Chosen: uint(selectInput.IntPrioCase),
			})

		}
	})

	sort.Sort(output.BySelectID(selects))

	return selects
}
