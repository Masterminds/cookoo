package web

import (
	//"io"
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
// - content: The content to write as a body. If this is a byte[], it is sent unchanged. Otherwise.
//   we first try to convert to a string, then pass it into a writer.
// - contentType: The content type header (e.g. text/html). Default is text/plain
// - responseCode: Integer HTTP Response Code: Default is `http.StatusOK`.
// - headers: a map[string]string of HTTP headers. The keys will be run through 
//   http.CannonicalHeaderKey()
//
// Note that this is optimized for writing from strings or arrays, not Readers. For larger
// objects, you may find it more efficient to use a different command.
func Flush (cxt cookoo.Context, params *cookoo.Params) (interface{}, cookoo.Interrupt) {

	// Make sure we have a place to write this stuff.
	writer, ok := params.Has("writer")
	if writer == nil {
		writer, ok = cxt.Has("http.ResponseWriter")
		if !ok {
			return false, nil
		}
	}
	out := writer.(http.ResponseWriter)

	// Get the rest of the info.
	code := params.Get("responseCode", http.StatusOK).(int)
	header := out.Header()
	contentType := params.Get("contentType", "text/plain; charset=utf-8").(string)

	// Prepare the content.
	var content []byte
	rawContent, ok := params.Has("content")
	if !ok {
		// No content. Send nothing in the body.
		content = []byte("")
	} else if byteContent, ok := rawContent.([]byte); ok {
		// Got a byte[]; add it as is.
		content = byteContent
	} else {
		// Use the formatter to convert to a string, and then
		// cast it to bytes.
		content = []byte(fmt.Sprintf("%v", rawContent))
	}

	// Add headers:
	header.Set(http.CanonicalHeaderKey("content-type"), contentType)
	headerO, ok := params.Has("headers")
	if ok {
		headers := headerO.(map[string]string)
		for k, v := range headers {
			header.Add(http.CanonicalHeaderKey(k), v)
		}
	}

	// Send the headers.
	out.WriteHeader(code)

	//io.WriteString(out, content)
	out.Write(content)

	return true, nil
}

func ServerInfo(cxt cookoo.Context, params *cookoo.Params) (interface{}, cookoo.Interrupt) {
	req := cxt.Get("http.Request", nil).(*http.Request)
	out := cxt.Get("http.ResponseWriter", nil).(http.ResponseWriter)

	out.Header().Add("X-Foo", "Bar")
	out.Header().Add("Content-type", "text/plain; charset=utf-8")

	fmt.Fprintf(out, "Request:\n %+v\n", req)
	fmt.Fprintf(out, "\n\n\nResponse:\n%+v\n", out)
	return true, nil
}
