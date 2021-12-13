package defaultp

import (
	aaa "sync"

	_ "github.com/go-kit/log"
)

type aa struct {
	m aaa.Mutex
}

func (_ *aa) abcde() {
	println(3)
}

func (a *aa) abcd() {
	b := 2
	println(a, b)
}
func Hello() {
	m := aaa.Mutex{}

	c := aaa.NewCond(&m)

	c.Broadcast()

	c.Signal()

	//asdfasdf
	c.Wait()
	//asdfadfas

	w := aaa.WaitGroup{}

	w.Wait()

	a := aa{}

	a.m.Lock()

}
