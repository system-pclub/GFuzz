package chrec

func Hello() {
	ch := make(chan int)

	ch <- 1

	<-ch

	close(ch)
}
