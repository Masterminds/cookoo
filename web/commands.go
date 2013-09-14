package web

import (
	"github.com/masterminds/cookoo"
	"html/template"
	"net/http"
	"strings"
	"bytes"
	"fmt"
	"log"
	"io"
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
//
// Returns
// - boolean true
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

// Render an HTML template.
//
// This uses the `html/template` system built into Go to render data into a writer.
//
// Params:
// - template (required): A string with the full path to the template.
// - values: An interface{} with the values to be passed to the template. If 
//   this is not specified, the contents of the Context are passed as a map[string]interface{}.
//   Note that datasources, in this model, are not accessible to the template.
// - writer: The writer that data should be sent to. By default, this will create a new
//   Buffer and put it into the context. (If no Writer was passed in, the returned writer
//   is actually a bytes.Buffer.) To flush the contents directly to the client, you can
//   use `.Using('writer').From('http.ResponseWriter')`.
//
// Returns
// - An io.Writer. The template's contents have already been written into the writer.
//
// Example:
//
//	reg.Route("GET /html", "Test HTML").
//		Does(cookoo.AddToContext, "_").
//			Using("Title").WithDefault("Hello World").
//			Using("Body").WithDefault("This is the body.").
//		Does(web.RenderHTML, "render").
//			Using("template").WithDefault("index.html").
//		Does(web.Flush, "_").
//			Using("contentType").WithDefault("text/html").
//			Using("content").From("cxt:render")
//
// In the example above, we do three things:
// - Add Title and Body to the context. For the template rendered, it will see these as
//   {{.Title}} and {{.Body}}.
// - Render the template located in a local file called "index.html"
// - Flush the result out to the client. This gives you a chance to add any additional headers.
func RenderHTML(cxt cookoo.Context, params *cookoo.Params) (interface{}, cookoo.Interrupt) {
	ok, missing := params.Requires("template")
	if !ok {
		return nil, &cookoo.FatalError{"Missing params: " + strings.Join(missing, ", ")}
	}
	
	var buf bytes.Buffer
	out := params.Get("writer", &buf).(io.Writer)
	tplName := params.Get("template", nil).(string)
	vals := params.Get("values", cxt.AsMap())

	//out := cxt.Get("http.ResponseWriter", nil).(http.ResponseWriter)

	tpl, err := template.ParseFiles(tplName)
	if err != nil {
		return nil, &cookoo.FatalError{"Failed to parse template(s): " + tplName}
	}
	err = tpl.Execute(out, vals)
	if err != nil {
		log.Printf("Recoverable error parsing template: %s", err)
		// XXX: This outputs partially completed templates. Is this what we want?
		io.WriteString(out, "Template error. The error has been logged.")
		return out, &cookoo.RecoverableError{"Template failed to completely render."}
	}
	return out, nil
}

// Get the server info for this request.
//
// This assumes that `http.Request` and `http.ResponseWriter` are in the context, which
// they are by default.
//
// Returns:
// - boolean true
func ServerInfo(cxt cookoo.Context, params *cookoo.Params) (interface{}, cookoo.Interrupt) {
	req := cxt.Get("http.Request", nil).(*http.Request)
	out := cxt.Get("http.ResponseWriter", nil).(http.ResponseWriter)

	out.Header().Add("X-Foo", "Bar")
	out.Header().Add("Content-type", "text/plain; charset=utf-8")

	fmt.Fprintf(out, "Request:\n %+v\n", req)
	fmt.Fprintf(out, "\n\n\nResponse:\n%+v\n", out)
	return true, nil
}
