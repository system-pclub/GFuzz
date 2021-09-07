package inst

type GoOpType string

const (
	ChSend  = "chsend"
	ChMake  = "chmake"
	ChRecv  = "chrecv"
	ChClose = "chclose"
)

type GoOp struct {
	Type GoOpType
	ID   uint16
}

type GoSrcStats struct {
	Ops []GoOp
}
