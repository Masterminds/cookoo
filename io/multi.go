package io

import (
	"io"
)

// MultiWriter enables you to have a writer that passes on the writing to one
// of more Writers where the write is duplicated to each Writer. MultiWriter
// is similar to the multiWriter that is part of Go. The difference is
// this MultiWriter allows you to manager the Writers attached to it via CRUD
// operations. To do this you will need to mock the type. For example,
// mw := NewMultiWriter()
// mw.(*MultiWriter).AddWriter("foo", foo)
type MultiWriter struct {
	writers map[string]io.Writer
}

func (t *MultiWriter) Write(p []byte) (n int, err error) {
	for _, w := range t.writers {
		n, err = w.Write(p)
		if err != nil {
			return
		}
		if n != len(p) {
			err = io.ErrShortWrite
			return
		}
	}
	return len(p), nil
}

func (t *MultiWriter) Init() *MultiWriter {
	t.writers = make(map[string]io.Writer)
	return t
}

func (t *MultiWriter) Writer(name string) (io.Writer, bool) {
	value, found := t.writers[name]
	return value, found
}

func (t *MultiWriter) Writers() map[string]io.Writer {
	return t.writers
}

func (t *MultiWriter) AddWriter(name string, writer io.Writer) {
	t.writers[name] = writer
}

func (t *MultiWriter) RemoveWriter(name string) {
	delete(t.writers, name)
}

func NewMultiWriter() io.Writer {
	w := new(MultiWriter).Init()
	return w
}
