package chrec

import oraclert "gfuzz/pkg/oraclert"

func Hello() {
	ch := make(chan int)
	oraclert.StoreChMakeInfo(ch, 1)
	oraclert.StoreOpInfo("Send", 2)

	ch <- 1
	oraclert.StoreOpInfo("Recv", 3)

	<-ch
	oraclert.StoreOpInfo("Close", 4)

	close(ch)
}
