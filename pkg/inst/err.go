package inst

import (
	"fmt"
)

type NoPassError struct {
	Name string
}

func (e *NoPassError) Error() string {
	return fmt.Sprintf("pass '%s' does not exist", e.Name)
}

type PassExistedError struct {
	Name string
}

func (e *PassExistedError) Error() string {
	return fmt.Sprintf("pass '%s' already existed", e.Name)
}
