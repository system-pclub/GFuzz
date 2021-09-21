package output

import (
	"gfuzz/pkg/oraclert/config"
	"gfuzz/pkg/selefcm"
)

type Output struct {
	Tuples         map[uint32]uint32     `json:"tuple"`
	Channels       map[string]ChanRecord `json:"channels"`
	Ops            []string              `json:"ops"`
	SelEfcmHistory []selefcm.SelEfcm     `json:"sel_efcm_history"`
}

type ChanRecord struct {
	ID        string `json:"id"`
	Closed    bool   `json:"closed"`
	NotClosed bool   `json:"not_closed"`
	CapBuf    int    `json:"cap_buf"`
	PeakBuf   int    `json:"peak_buf"`
}

func GenerateOracleRtOutput(config *config.Config) *Output {
	output := &Output{
		Tuples: getTuples(),
	}

	if config.DumpSelEfcmHistory {

	}

	return output
}

func Serialize() {

}

func Deserialize() {

}
