package web

import (
	"io"
	"github.com/masterminds/cookoo"
	"net/http"
	"fmt"
)

// Common web-oriented commands

// Send content to output.
//
// If no writer is specified, this will attempt to write to whatever is in the
// Context with the key "http.ResponseWriter". If no suitable writer is found, it will
// not write to anything at all.
//
// Params:
// - writer: A Writer of some sort. This will try to write to the HTTP response if no writer
//   is specified.
// - content: The content to write as a body.
// - contentType: The content type header (e.g. text/html). Default is text/plain
func Flush (cxt cookoo.Context, params *cookoo.Params) (interface{}, cookoo.Interrupt) {
	// ok, _ := params.Requires("writer")

	//if (!ok) {
	//	return false, nil
	//}

	//out := params.Get("writer", nil).(http.ResponseWriter)
	writer, ok := params.Has("writer")
	if writer == nil {
		writer, ok = cxt.Has("http.ResponseWriter")
		if !ok {
			return false, nil
		}
	}
	out := writer.(io.Writer)

	content := params.Get("content", "").(string)
	contentType := params.Get("contentType", "text/plain; charset=utf-8").(string)

	fmt.Fprintf(out, "%s: %s\n", http.CanonicalHeaderKey("content-type"), contentType)
	io.WriteString(out, content)

	return true, nil
}

func ServerInfo(cxt cookoo.Context, params *cookoo.Params) (interface{}, cookoo.Interrupt) {
	req := cxt.Get("http.Request").(*http.Request)
	out := cxt.Get("http.ResponseWriter").(http.ResponseWriter)

	out.Header().Add("X-Foo", "Bar")
	out.Header().Add("Content-type", "text/plain; charset=utf-8")

	fmt.Fprintf(out, "Request:\n %+v\n", req)
	fmt.Fprintf(out, "\n\n\nResponse:\n%+v\n", out)
	return true, nil
}
