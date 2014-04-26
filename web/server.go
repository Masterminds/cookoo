package web

import (
	"github.com/Masterminds/cookoo"
	"log"
	"net/http"
	"runtime"
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
//      "github.com/Masterminds/cookoo/src/cookoo"
//      "github.com/Masterminds/cookoo/src/cookoo/web"
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
//    	web.Serve(reg, router, cookoo.SyncContext(cxt))
//    }
//
// Note that we synchronize the context before passing it into Serve(). This
// is optional because each handler gets its own copy of the context already.
// However, if commands pass the context to goroutines, the context ought to be
// synchronized to avoid race conditions.
//
// Note that copies of the context are not synchronized with each other.
// So by declaring the context synchronized here, you
// are not therefore synchronizing across handlers.
func Serve(reg *cookoo.Registry, router *cookoo.Router, cxt cookoo.Context) {
	defer shutdown(router, cxt)

	addr := cxt.Get("server.Address", ":8080").(string)

	handler := NewCookooHandler(reg, router, cxt)
	http.Handle("/", handler)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		cxt.Logf("error", "Caught error while serving: %s", err)
		if router.HasRoute("@crash") {
			router.HandleRequest("@crash", cxt, false)
		}
	}
	// TODO: Need to figure out how to trap signals here, instead of outside.
}
func shutdown(router *cookoo.Router, cxt cookoo.Context) {
	log.Print("Shutdown")
	if router.HasRoute("@shutdown") {
		router.HandleRequest("@shutdown", cxt, false)
	}
}

// The handler for Cookoo.
// You way use this handler in your own web apps, or you can use
// the Serve() function to create and manage a handler for you.
type CookooHandler struct {
	Registry    *cookoo.Registry
	Router      *cookoo.Router
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

// ServeHTTP is the Cookoo request handling function.
//
// This is capable of handling HTTP and HTTPS requests.
func (h *CookooHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	// Trap panics and make them 500 errors:
	defer func() {
		// fmt.Printf("Deferred function executed for path %s\n", req.URL.Path)
		if err := recover(); err != nil {
			//log.Printf("FOUND ERROR: %v", err)
			h.BaseContext.Logf("error", "CookooHandler trapped a panic: %v", err)

			// Buffer for a stack trace.
			stack := make([]byte, 8192)
			size := runtime.Stack(stack, false)
			h.BaseContext.Logf("error", "Stack: %s", stack)

			if size == 8192 {
				h.BaseContext.Logf("error", "<truncated stack trace at 8192 bytes>")
			}

			http.Error(res, "An internal error occurred.", http.StatusInternalServerError)
		}
	}()
	// First we need to clone the context so we have a mutable copy.
	cxt := h.BaseContext.Copy()

	cxt.Put("http.Request", req)
	cxt.Put("http.ResponseWriter", res)

	// Next, we add the datasources for URL and Query params.
	h.addDatasources(cxt, req)

	// Find the route
	path := req.Method + " " + req.URL.Path

	cxt.Logf("info", "Handling request for %s\n", path)

	// If a route matches, run it.
	err := h.Router.HandleRequest(path, cxt, true)
	if err != nil {
		switch err.(type) {

		// For a 404, we bail.
		case *cookoo.RouteError:
			if h.Router.HasRoute("@404") {
				h.Router.HandleRequest("@404", cxt, false)
			} else {
				http.NotFound(res, req)
			}
			return
		// For any other, we go to a 500.
		case *cookoo.FatalError:
			cxt.Logf("error", "Fatal Error on route '%s': %s", path, err)
		default:
			cxt.Logf("error", "Untagged error on route '%s': %v (%T)", path, err, err)
		}

		if h.Router.HasRoute("@500") {
			cxt.Put("error", err)
			h.Router.HandleRequest("@500", cxt, false)
		} else {
			// Passing the error back to the client is a bad default.
			//http.Error(res, err.Error(), http.StatusInternalServerError)
			http.Error(res, "Internal error processing the request.", http.StatusInternalServerError)
		}
	}
}
