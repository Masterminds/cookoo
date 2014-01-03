package io

import (
	cio "io"
	"strings"
)

// A Colorizing logger middleman.
//
// This can be used to colorize logs as they pass through to another writer. The
// colorization uses the UNIX-style shell color coding.
//
// TODO: This is a very basic implementation, and could use some TLC.
//
// Example Usage:
//
//  import (
// 		"github.com/Masterminds/cookoo"
// 		"github.com/Masterminds/cookoo/io"
// 		// And other stuff
// 		cio "io"
// 	)
// 	func main() {
// 		reg, router, cxt := cookoo.Cookoo()
// 		clogger := io.NewColorizer(cio.Stdout)
// 		cxt.AddLogger("stdout", clogger)
// 		// etc.
// 	}
//
// Given the above, log messages will be colorized before written
// to `io.Stdout`.
type Colorizer struct {
	writer cio.Writer
}

// Creates a new colorizer that wraps a given io.Writer.
func NewColorizer(writer cio.Writer) *Colorizer {
	c := new(Colorizer)
	c.writer = writer

	return c
}

// Colorizes a message and then passes it to the underlying writer.
func (r *Colorizer) Write(data []byte) (int, error) {

	str := string(data)
	if strings.HasPrefix(str, "error") {
		str = "\033[0;31m" + str + "\033[m"
	} else if strings.HasPrefix(str, "warning") {
		str = "\033[0;33m" + str + "\033[m"
	} else if strings.HasPrefix(str, "info") {
		str = "\033[0;36m" + str + "\033[m"
	}

	return r.writer.Write([]byte(str))
}
