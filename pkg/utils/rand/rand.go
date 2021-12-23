package rand

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func GetRandomWithMax(max int) int {
	mutateMethod, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		fmt.Println("Crypto/rand returned non-nil errors: ", err)
	}
	return int(mutateMethod.Int64())
}