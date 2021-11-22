package output

import (
	"encoding/json"
)

type Output struct {
	Tuples   map[uint32]uint32     `json:"tuple"`
	Channels map[string]ChanRecord `json:"channels"`
	Ops      []uint16              `json:"ops"`
	Selects  []SelectRecord        `json:"selects"`
}

type ChanRecord struct {
	ID        string `json:"id"`
	Closed    bool   `json:"closed"`
	NotClosed bool   `json:"not_closed"`
	CapBuf    int    `json:"cap_buf"`
	PeakBuf   int    `json:"peak_buf"`
}

type SelectRecord struct {
	ID     string `json:"id"`
	Cases  uint   `json:"cases"`
	Chosen uint   `json:"chosen"`
}

func Serialize(o *Output) ([]byte, error) {
	if o == nil {
		return []byte{}, nil
	}
	return json.Marshal(o)
}

func Deserialize(data []byte) (*Output, error) {
	o := Output{}
	err := json.Unmarshal(data, &o)
	if err != nil {
		return nil, err
	}
	return &o, nil

}

type BySelectID []SelectRecord

func (a BySelectID) Len() int           { return len(a) }
func (a BySelectID) Less(i, j int) bool { return a[i].ID < a[j].ID }
func (a BySelectID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
