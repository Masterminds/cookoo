package web

import (
	"io"
	"github.com/masterminds/cookoo"
	"net/http"
	"fmt"
)

// Common web-oriented commands

func Flush (cxt cookoo.Context, params *cookoo.Params) (interface{}, cookoo.Interrupt) {
	ok, _ := params.Requires("writer")

	if (!ok) {
		return false, nil
	}

	out := params.Get("writer", nil).(http.ResponseWriter)
	content := params.Get("content", "").(string)
	contentType := params.Get("contentType", "plain/text; charset=utf-8").(string)

	//out.Header(

	fmt.Fprintf(out, "%s: %s\n", http.CanonicalHeaderKey("content-type"), contentType)
	io.WriteString(out, content)

	return true, nil
}

func ServerInfo(cxt cookoo.Context, params *cookoo.Params) (interface{}, cookoo.Interrupt) {
	req := cxt.Get("http.Request").(http.Request)
	out := cxt.Get("http.ResponseWriter").(http.ResponseWriter)

	out.Header().Add("X-Foo", "Bar")
	out.Header().Add("Content-type", "text/markdown; charset=utf-8")

	fmt.Fprintf(out, "Request: %+v\n", req)
	fmt.Fprintf(out, "\n\n\nResponse: %+v\n", out)
	return true, nil
}
