package gofmt

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
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

func GoModDownload(dir string) error {
	if _, err := os.Stat(path.Join(dir, "go.mod")); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("no go.mod presented in '%s'", dir)
	}
	cmd := exec.Command("go", "mod", "download")
	if dir != "" {
		cmd.Dir = dir
	}
	return cmd.Run()
}
