package fuzzer

import (
	"gfuzz/pkg/gexec"
	"math/rand"
	"time"
)

func Shuffle(vals []gexec.Executable) {
	r := rand.New(rand.NewSource(time.Now().Unix()))

	for n := len(vals); n > 0; n-- {
		randIndex := r.Intn(n)
		vals[n-1], vals[randIndex] = vals[randIndex], vals[n-1]
	}
}
