package defaultp

import (
	oraclert "gfuzz/pkg/oraclert"
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
	oraclert.CurrentGoAddValue(a, nil, 0)
	b := 2
	println(a, b)
}
func Hello() {
	m := aaa.Mutex{}

	c := aaa.NewCond(&m)
	oraclert.StoreOpInfo("Broadcast",

		//asdfasdf
		1)

	c.Broadcast()
	oraclert.StoreOpInfo("Signal", 2)

	c.Signal()
	oraclert.StoreOpInfo("Wait",

		//asdfadfas
		3)

	c.Wait()

	w := aaa.WaitGroup{}
	oraclert.StoreOpInfo("Wait", 5)

	w.Wait()

	a := aa{}
	oraclert.StoreOpInfo("Lock", 4)

	a.m.Lock()

}
