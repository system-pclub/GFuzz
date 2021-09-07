package wgrec

import "sync"

func Hello() {
	wg := sync.WaitGroup{}

	wg.Add(1)
	wg.Wait()
	wg.Done()
}
