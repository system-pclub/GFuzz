package interest

import (
	"gfuzz/pkg/oraclert/output"
	"gfuzz/pkg/selefcm"
)

func IsEfcmCovered(efcms []selefcm.SelEfcm, records []output.SelectRecord) bool {
	records_cache := make(map[string]map[int]struct{})
	for _, r := range records {
		if _, exist := records_cache[r.ID]; !exist {
			records_cache[r.ID] = make(map[int]struct{})
		}

		records_cache[r.ID][int(r.Chosen)] = struct{}{}
	}
	for _, e := range efcms {
		cases, exist := records_cache[e.ID]
		// if a enforcement, not reflect on the outputs, return false
		if !exist {
			return false
		}
		if _, exist := cases[int(e.Case)]; !exist {

			return false
		}

	}
	return true
}
