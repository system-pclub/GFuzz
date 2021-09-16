package fuzz

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Record struct {
	MapTupleRecord map[int]int
	MapChanRecord  map[string]ChanRecord
}

func (r *Record) GetChanRecords() []ChanRecord {
	values := make([]ChanRecord, 0, len(r.MapChanRecord))

	for _, v := range r.MapChanRecord {
		values = append(values, v)
	}

	return values
}

type ChanRecord struct {
	ChID      string
	Closed    bool
	NotClosed bool
	CapBuf    int
	PeakBuf   int
}

const (
	RecordSplitter = "-----"
)

func EmptyRecord() *Record {
	return &Record{}
}

func ParseRecordFile(content string) (*Record, error) {
	var err error
	text := strings.Split(content, "\n")
	if len(text) == 0 {
		return nil, errors.New("Record is empty")
	}

	newRecord := &Record{
		MapTupleRecord: make(map[int]int),
		MapChanRecord:  make(map[string]ChanRecord),
	}

	indexSplitter := -1
	for i, eachline := range text {
		if eachline == RecordSplitter {
			indexSplitter = i
			break
		}
		if eachline == "" {
			continue
		}
		vecStr := strings.Split(eachline, ":")
		if len(vecStr) != 2 {
			return nil, fmt.Errorf("tuple in record has incorrect format: %s, at line %d", eachline, i)
		}
		var tuple, count int
		if tuple, err = strconv.Atoi(vecStr[0]); err != nil {
			return nil, fmt.Errorf("tuple in record has incorrect format: %s, at line %d", vecStr, i)
		}
		if count, err = strconv.Atoi(vecStr[1]); err != nil {
			return nil, fmt.Errorf("tuple in record has incorrect format: %s, at line %d", vecStr, i)
		}
		newRecord.MapTupleRecord[tuple] = count
	}

	if indexSplitter == -1 {
		return nil, fmt.Errorf("doesn't find RecordSplitter in record. Full text: %s", text)
	}

	for i := indexSplitter + 1; i < len(text); i++ {
		eachline := text[i]
		if eachline == "" {
			continue
		}
		//filename:linenum:closedBit:notClosedBit:capBuf:peakBuf
		vecStr := strings.Split(eachline, ":")
		if len(vecStr) != 6 {
			return nil, fmt.Errorf("channel in record has incorrect format: %s, at line %d", eachline, i)
		}
		chRecord := ChanRecord{}
		chRecord.ChID = vecStr[0] + ":" + vecStr[1]
		if vecStr[2] == "0" {
			chRecord.Closed = false
		} else {
			chRecord.Closed = true
		}
		if vecStr[3] == "0" {
			chRecord.NotClosed = false
		} else {
			chRecord.NotClosed = true
		}
		if chRecord.CapBuf, err = strconv.Atoi(vecStr[4]); err != nil {
			return nil, fmt.Errorf("line of channel in record has incorrect format: %s, at line %d", eachline, i)
		}
		if chRecord.PeakBuf, err = strconv.Atoi(vecStr[5]); err != nil {
			return nil, fmt.Errorf("line of channel in record has incorrect format: %s, at line %d", eachline, i)
		}
		newRecord.MapChanRecord[chRecord.ChID] = chRecord
	}
	return newRecord, nil
}

func UpdateMainRecord(mainRecord, curRecord Record) Record {

	// Update a tuple if it doesn't exist or it exists but a better count is observed
	for tuple, count := range curRecord.MapTupleRecord {
		mainCount, exist := mainRecord.MapTupleRecord[tuple]
		if exist {
			if mainCount < count {
				mainRecord.MapTupleRecord[tuple] = count
			}
		} else {
			mainRecord.MapTupleRecord[tuple] = count
		}
	}

	for chID, chRecord := range curRecord.MapChanRecord {
		mainChRecord, exist := mainRecord.MapChanRecord[chID]
		if exist {
			// Update an existing channel's status
			if mainChRecord.Closed == false && chRecord.Closed == true {
				mainChRecord.Closed = true
			}
			if mainChRecord.NotClosed == false && chRecord.NotClosed == true {
				mainChRecord.NotClosed = true
			}
			if mainChRecord.PeakBuf < chRecord.PeakBuf {
				mainChRecord.PeakBuf = chRecord.PeakBuf
			}
			mainChRecord.CapBuf = chRecord.CapBuf

			mainRecord.MapChanRecord[chID] = mainChRecord
		} else {
			// Update a new chan if it doesn't exist
			mainRecord.MapChanRecord[chID] = chRecord
		}
	}

	return mainRecord
}
