package defaultp

func getChannel() chan string {
	return make(chan string)
}
func Hello() {
	ch1 := make(chan int)
	ch2 := make(chan struct{})
	ch3 := getChannel()
	ch4 := ch1

	for {
		ch5 := make(chan int)
		go func() {
			ch5 <- 1
		}()
	}

}
