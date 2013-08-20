package web

import(
	"github.com/masterminds/cookoo"
	"net/http"
	"fmt"
)

func Serve (reg *cookoo.Registry, router *cookoo.Router, cxt cookoo.Context) {

	handler := NewCookooHandler(reg, router, cxt)
	http.Handle("/", handler)
	http.ListenAndServe(":8080", nil)
}

func dummy(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, "Hi Annabelle")
}

// The handler for Cookoo.
// You way use this handler in your own web apps, or you can use
// the Serve() function to create and manage a handler for you.
type CookooHandler struct {
	Registry *cookoo.Registry
	Router *cookoo.Router
	BaseContext cookoo.Context
}

func NewCookooHandler(reg *cookoo.Registry, router *cookoo.Router, cxt cookoo.Context) *CookooHandler {
	handler := new(CookooHandler)
	handler.Registry = reg
	handler.Router = router
	handler.BaseContext = cxt

	return handler
}

func (h *CookooHandler) addDatasources(cxt cookoo.Context, req *http.Request) {
	parsedURL := req.URL
	urlDS := new(URLDatasource).Init(parsedURL)
	queryDS := new(QueryParameterDatasource).Init(parsedURL.Query())
	formDS := new(FormValuesDatasource).Init(req)

	cxt.AddDatasource("url", urlDS)
	cxt.AddDatasource("query", queryDS)
	cxt.AddDatasource("q", queryDS)
	cxt.AddDatasource("post", formDS)
}

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
	path := req.URL.Path;

	fmt.Printf("Handling request for %s\n", path)

	// If a route matches, run it.
	if h.Router.HasRoute(path) {
		h.Router.HandleRequest(path, cxt, true)


	// Else if there is a custom 404 handler, run it.
	} else if h.Router.HasRoute("@404") {
		// Taint mode is false for error paths.
		h.Router.HandleRequest("@404", cxt, false)

	// Else run the default 404 handler.
	} else {
		http.NotFound(res, req)
	}
}

