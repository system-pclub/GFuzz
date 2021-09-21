package config

import (
	"encoding/json"
	"gfuzz/pkg/selefcm"
)

// Config contains all information need to be passed to fuzz target (consumed by oracle runtime)
type Config struct {
	// SelEfcm, select enforcement
	SelEfcm            selefcm.SelEfcmConfig `json:"selefcm"`
	DumpSelEfcmHistory bool                  `json:"dump_sel_efcm_history"`
}

func Serialize(l *Config) ([]byte, error) {
	if l == nil {
		return []byte{}, nil
	}

	return json.Marshal(l)
}

func Deserilize(data []byte) (*Config, error) {
	l := Config{}
	err := json.Unmarshal(data, &l)
	if err != nil {
		return nil, err
	}
	return &l, nil
}
