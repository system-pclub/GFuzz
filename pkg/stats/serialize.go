package stats

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func Serialize(s *stats) ([]byte, error) {
	j, err := json.Marshal(struct {
		Stats map[string]uint64
	}{
		Stats: s.stats,
	})
	if err != nil {
		return nil, err
	}
	return j, nil
}

func Deserialize(data []byte) (*stats, error) {
	b := struct {
		Stats map[string]uint64
	}{}
	err := json.Unmarshal(data, &b)
	if err != nil {
		return nil, err
	}
	return &stats{
		stats: b.Stats,
	}, nil
}

func ToFile(s *stats, dstFile string) error {
	file, err := os.OpenFile(dstFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	defer file.Close()
	if err != nil {
		return err
	}
	data, err := Serialize(s)
	if err != nil {
		return err
	}
	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func FromFile(file string) (*stats, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return Deserialize(data)
}
