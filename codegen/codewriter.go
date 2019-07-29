package codegen

import (
	"io"
)

type CodeWriter struct {
	io.Writer
	IndentationPrefix string
	indentation int
}

func (w *CodeWriter) WriteLine(line string) (n int, err error) {
	prefix := ""
	for i :=0; i < w.indentation; i++ {
		prefix += w.IndentationPrefix
	}
	return w.Write([]byte(line))
}