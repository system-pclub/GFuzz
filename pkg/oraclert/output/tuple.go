package output

import "runtime"

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
