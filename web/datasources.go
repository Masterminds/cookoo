// Extra datasources for Web servers.
package web

import (
	"net/url"
	"net/http"
)

// Get the query parameters by name.
type QueryParameterDatasource struct {
	Parameters url.Values
}
// The datasource for URLs.
// This datasource knows the following items:
// - url: the URL struct
// - scheme: The scheme of the URL as a string
// - opaque: The opaque identifier
// - user: A *Userinfo
// - host: The string hostname
// - path: The entire path
// - rawquery: The query string, not decoded.
// - fragment: The fragment string.
// - query: The array of Query parameters. Usually it is better to use the
//   'query:foo' syntax.
type URLDatasource struct {
	URL *url.URL
}

// Access to name/value pairs in POST/PUT form data from the body.
// This will attempt to access form data supplied in the HTTP request's body.
// If the MIME type is not correct or if there is no POST data, no data will
// be made available.
//
// Parsing is lazy: No form data is parsed until it is requested.
type FormValuesDatasource struct {
	req *http.Request
}

func (f *FormValuesDatasource) Init(req *http.Request) *FormValuesDatasource {
	f.req = req
	return f
}

// The return value will always be a string or nil.
// To match the interface, we use interface{}.
func (f *FormValuesDatasource) Value(name string) interface{} {
	return f.req.PostFormValue(name)
}

func (d *QueryParameterDatasource) Init(vals url.Values) *QueryParameterDatasource {
	d.Parameters = vals
	return d
}

func (d *QueryParameterDatasource) Value(name string) interface{} {
	return d.Parameters.Get(name)
}


func (d *URLDatasource) Init(parsedUrl *url.URL) *URLDatasource {
	d.URL = parsedUrl
	return d
}

func (d *URLDatasource) Value(name string) interface{} {
	switch name {
	case "host", "Host":
		return d.URL.Host
	case "path", "Path":
		return d.URL.Path
	case "url", "URL", "Url":
		return d.URL
	case "user", "User":
		return d.URL.Path
	case "scheme", "Scheme":
		return d.URL.Scheme
	case "rawquery", "RawQuery":
		return d.URL.RawQuery
	case "query", "Query":
		return d.URL.Query()
	case "fragment", "Fragment":
		return d.URL.Fragment
	case "opaque", "Opaque":
		return d.URL.Opaque
	}
	return nil
}

