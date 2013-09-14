package web

import (
	"github.com/masterminds/cookoo"
	//cookoo "../"
	"net/http"
	"fmt"
	"log"
)

// Create a new Cookoo web server.
//
// Important details:
//
// - A URIPathResolver is used for resolving request names.
// - The following datasources are added to the Context:
//   * url: A URLDatasource (Provides access to parts of the URL)
//   * path: A PathDatasource (Provides access to parts of a path. E.g. "/foo/bar")
//   * query: A QueryParameterDatasource (Provides access to URL query parameters.)
//   * post: A FormValuesDatasource (Provides access to form data or the body of a request.)
// - The following context variables are set:
//   * http.Request: A pointer to the http.Request object
//   * http.ResponseWriter: The response writer.
//   * server.Address: The server's address and port (NOT ALWAYS PRESENT)
// - The handler includes logic to redirect "not found" errors to a path named "@404" if present.
//
// Context Params:
//
// - server.Address: If this key exists in the context, it will be used to determine the host/port the
//   server runes on. EXPERIMENTAL. Default is ":8080".
//
// Example:
//
//    package main
//
//    import (
//      //This is the path to Cookoo
//      "github.com/masterminds/cookoo/src/cookoo"
//      "github.com/masterminds/cookoo/src/cookoo/web"
//      "fmt"
//    )
//
//    func main() {
//      // Build a new Cookoo app.
//      handler, registry, router, context := cookoo.Cookoo()
//
//      // Fill the registry.
//      registry.Route("GET /", "The index").Does(web.Flush, "example").
//      	Using("content").WithDefault("Hello World")
//      	Using("writer").From(
//
//    	// Create a server
//    	web.Serve(reg, router, cxt)
//    }
//
func Serve(reg *cookoo.Registry, router *cookoo.Router, cxt cookoo.Context) {

	addr := cxt.Get("server.Address", ":8080").(string)

	handler := NewCookooHandler(reg, router, cxt)
	http.Handle("/", handler)
	http.ListenAndServe(addr, nil)
}

// The handler for Cookoo.
// You way use this handler in your own web apps, or you can use
// the Serve() function to create and manage a handler for you.
type CookooHandler struct {
	Registry *cookoo.Registry
	Router *cookoo.Router
	BaseContext cookoo.Context
}

// Create a new Cookoo HTTP handler.
//
// This will create an HTTP hanlder, but will not automatically attach it to a server. Implementors
// can take the handler and attach it to an existing HTTP server wiht http.HandleFunc() or
// http.ListenAndServe().
//
// For simple web servers, using this package's Serve() function may be the easier route.
//
// Important details:
//
// - A URIPathResolver is used for resolving request names.
// - The following datasources are added to the Context:
//   * url: A URLDatasource (Provides access to parts of the URL)
//   * path: A PathDatasource (Provides access to parts of a path. E.g. "/foo/bar")
//   * query: A QueryParameterDatasource (Provides access to URL query parameters.)
//   * post: A FormValuesDatasource (Provides access to form data or the body of a request.)
// - The following context variables are set:
//   * http.Request: A pointer to the http.Request object
//   * http.ResponseWriter: The response writer.
//   * server.Address: The server's address and port (NOT ALWAYS PRESENT)
func NewCookooHandler(reg *cookoo.Registry, router *cookoo.Router, cxt cookoo.Context) *CookooHandler {
	handler := new(CookooHandler)
	handler.Registry = reg
	handler.Router = router
	handler.BaseContext = cxt

	// Use the URI oriented request resolver in this package.
	resolver := new(URIPathResolver)
	resolver.Init(reg)
	router.SetRequestResolver(resolver)

	return handler
}

// Adds the built-in HTTP-specific datasources.
func (h *CookooHandler) addDatasources(cxt cookoo.Context, req *http.Request) {
	parsedURL := req.URL
	urlDS := new(URLDatasource).Init(parsedURL)
	queryDS := new(QueryParameterDatasource).Init(parsedURL.Query())
	formDS := new(FormValuesDatasource).Init(req)
	pathDS := new(PathDatasource).Init(parsedURL.Path)

	cxt.AddDatasource("url", urlDS)
	cxt.AddDatasource("query", queryDS)
	// cxt.AddDatasource("q", queryDS)
	cxt.AddDatasource("post", formDS)
	cxt.AddDatasource("path", pathDS)
}

// The Cookoo request handling function.
func (h *CookooHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	// Trap panics and make them 500 errors:
	defer func() {
		// fmt.Printf("Deferred function executed for path %s\n", req.URL.Path)
		if err := recover(); err != nil {
			fmt.Printf("FOUND ERROR: %v", err)
			http.Error(res, "An internal error occurred.", http.StatusInternalServerError)
		}
	}()
	// First we need to clone the context so we have a mutable copy.
	cxt := h.BaseContext.Copy()

	cxt.Add("http.Request", req)
	cxt.Add("http.ResponseWriter", res)

	// Next, we add the datasources for URL and Query params.
	h.addDatasources(cxt, req)

	// Find the route
	path := req.Method + " " + req.URL.Path;

	fmt.Printf("Handling request for %s\n", path)

	// If a route matches, run it.
	if h.Router.HasRoute(path) {
		err := h.Router.HandleRequest(path, cxt, true)

		if err != nil {
			fatal, ok := err.(*cookoo.FatalError)
			if !ok {
				log.Printf("Unknown error: %v %T", err, err)
			} else {
				log.Printf("Fatal Error: %s", fatal)
			}
		}

	// Else if there is a custom 404 handler, run it.
	} else if h.Router.HasRoute("@404") {
		// Taint mode is false for error paths.
		h.Router.HandleRequest("@404", cxt, false)

	// Else run the default 404 handler.
	} else {
		http.NotFound(res, req)
	}
}

