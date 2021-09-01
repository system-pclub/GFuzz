package chrec

import gooracle "gooracle"

func Hello() {
	ch := make(chan int)
	gooracle.StoreChMakeInfo(ch, 0)
	gooracle.StoreOpInfo("Send", 1)

	ch <- 1
	gooracle.StoreOpInfo("Recv", 2)

	<-ch
	gooracle.StoreOpInfo("Close", 3)

	close(ch)
}
