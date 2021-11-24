package wgrec

import "sync"

type RandomStruct struct {
	wg sync.WaitGroup
}

func Hello() {
	wg := sync.WaitGroup{}
	r := RandomStruct{}
	wg.Add(1)
	wg.Wait()
	wg.Done()

	r.wg.Wait()
}
