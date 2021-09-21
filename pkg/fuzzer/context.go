package fuzzer

//
type fuzzerContext struct {
	execCh chan *ExecTask // task for worker to run
}

func newFuzzerContext() *fuzzerContext {
	return &fuzzerContext{}
}
