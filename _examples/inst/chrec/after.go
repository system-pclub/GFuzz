package chrec

import oraclert "gfuzz/pkg/oraclert"

type A struct {
	aa chan struct{}
}

func Hello() {
	ch := oraclert.StoreChMakeInfo(make(chan int), 1).(chan int)
	a := oraclert.StoreChMakeInfo(make(chan struct{}), 2).(chan struct{})
	b := &A{
		aa: oraclert.StoreChMakeInfo(make(chan struct{}), 3).(chan struct{}),
	}
	oraclert.StoreOpInfo("Send", 4)
	ch <- 1
	oraclert.StoreOpInfo("Recv", 5)

	<-ch
	oraclert.StoreOpInfo("Close", 6)

	close(ch)
}
