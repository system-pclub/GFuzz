package hash

import (
	"crypto/sha256"
	"fmt"
)

type Hashable interface {
	Hash() string
}

func AsSha256(o interface{}) string {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%v", o)))

	return fmt.Sprintf("%x", h.Sum(nil))
}
