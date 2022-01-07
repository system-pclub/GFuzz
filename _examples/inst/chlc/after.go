package defaultp

import (
	"fmt"
	oraclert "gfuzz/pkg/oraclert"
)

func getChannel() chan string {
	return make(chan string)
}

func useOfCh(c interface{}) {

}
func Hello() {
	ch1 := make(chan int)
	ch2 := make(chan struct{})
	ch3 := getChannel()
	ch4 := ch1

	for {
		ch5 := make(chan int)
		go func() {
			oraclert.CurrentGoAddCh(ch4)
			oraclert.CurrentGoAddCh(ch2)
			oraclert.CurrentGoAddCh(ch3)
			oraclert.CurrentGoAddCh(ch5)
			fmt.Printf("many many code here")
			fmt.Printf("many many code here")
			fmt.Printf("many many code here")
			fmt.Printf("many many code here")
			fmt.Printf("many many code here")

			ch5 <- 1
			useOfCh(ch3)
			select {
			case <-ch2:
			case ch4 <- 3:
			}
		}()
	}

}
