Cookoo
======

A chain-of-command framework written in Go

## Usage

- Go get Cookoo.

Use it as follows:

~~~go
package main

import (
	"cookoo"
	"fmt"
)

func main() {

	// Build a new Cookoo app.
	registry, router, context := cookoo.Cookoo()

	// Fill the registry.
  registry.Route("TEST", "A test route").Does(SomeCommandFunc, "a") //...

	// Execute the route.
	router.HandleRequest("TEST", context, false)
}

func HelloWorld(cxt cookoo.Context, params cookoo.Params) {
	fmt.Println("Hello World")
	return
}

~~~
