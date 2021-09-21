package output

import "runtime"

func getChans() map[string]ChanRecord {
	var chans map[string]ChanRecord = make(map[string]ChanRecord)
	for _, chr := range runtime.ChRecord {
		if chr == nil {
			continue
		}
		chans[chr.StrCreation] = ChanRecord{
			ID:        chr.StrCreation,
			Closed:    chr.Closed,
			NotClosed: chr.NotClosed,
			CapBuf:    int(chr.CapBuf),
			PeakBuf:   int(chr.PeakBuf),
		}
	}
	return chans
}
