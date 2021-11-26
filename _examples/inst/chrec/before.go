package chrec

type A struct {
	aa chan struct{}
}

func Hello() {
	ch := make(chan int)
	a := make(chan struct{})
	c := make(chan struct{}, 4)
	b := &A{
		aa: make(chan struct{}),
	}
	ch <- 1

	<-ch

	close(ch)
}
