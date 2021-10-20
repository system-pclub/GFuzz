package chrec

import oraclert "gfuzz/pkg/oraclert"

func Hello() {
	ch := make(chan int)
	oraclert.StoreChMakeInfo(ch, 0)
	oraclert.StoreOpInfo("Send", 1)

	ch <- 1
	oraclert.StoreOpInfo("Recv", 2)

	<-ch
	oraclert.StoreOpInfo("Close", 3)

	close(ch)
}
