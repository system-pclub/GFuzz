package fuzz

import (
	"gfuzz/pkg/fuzz/exec"
)

type Handler interface {
	Handle(fc *Context, i *exec.Input, o *exec.Output) ([]*exec.Input, error)
}

type ScoreStrategy interface {
	Score(fc *Context, i *exec.Input, o *exec.Output) (int, error)
}
