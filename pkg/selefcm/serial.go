package selefcm

import "encoding/json"

func Serialize(l *SelEfcmInput) ([]byte, error) {
	if l == nil {
		return []byte{}, nil
	}

	return json.Marshal(l)
}

func Deserilize(data []byte) (*SelEfcmInput, error) {
	l := SelEfcmInput{}
	err := json.Unmarshal(data, &l)
	if err != nil {
		return nil, err
	}
	return &l, nil
}
