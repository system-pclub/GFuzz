package gofmt

import (
	"bytes"
	"io"
	"log"
	"os/exec"
)

func HasSyntaxError(goSrcFile string) bool {
	cmd := exec.Command("go", "fmt", goSrcFile)
	var out bytes.Buffer
	w := io.MultiWriter(&out, log.Writer())
	cmd.Stdout = w
	cmd.Stderr = w

	err := cmd.Run()

	return err != nil
}
