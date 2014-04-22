---
layout: article
title: "Cookoo, The Chain-of-Command Controller for Go"
keywords: "Go, Chain-of-Command, front controller, golang"
description: "Learn about using the Chain-of-Command pattern in Go."
---
# Cookoo, The Chain-of-Command Controller for Go

Cookoo is a high performance reusable controller pattern that can be used to build console (CLI) applications, API servers, middle ware, integration tests, bots, and web applications.

Think of Cookoo as a better front-controller:

* Map a route pattern to a chain of commands
* Each time a route is matched, commands are processed in order
* A `Context` travels from command to command, making it possible to
  pass information along.

Since commands are re-usable, you can begin with existing commands, add
your own, and rapidly assemble a collection of components that can
simply and speed up application development.

## Getting Started with Hello World

Here is a dead-simple Cookoo "Hello World" command line app.

    package main

    import (
        "github.com/Masterminds/cookoo"
        "fmt"
    )

    func main() {

        // Build a new Cookoo app.
        registry, router, context := cookoo.Cookoo()

        // For the route TEST execute a single command HelloWorld.
        // The registry stores the mapping of routes to commands.
        registry.Route("TEST", "A test route").Does(HelloWorld, "hi")

        // Execute the route.
        router.HandleRequest("TEST", context, false)
    }

    func HelloWorld(cxt cookoo.Context, params *cookoo.Params) interface{} {
        fmt.Println("Hello World")
        return true
    }

When creating a new cookoo based application there are tree main parts:

1. _Registry_: contains the mapping of commands to callbacks along with how information is passed around the context.
2. _Router_: handles incoming requests and routes them to the correct callback on the registry. Different types of applications will have different routers. Commands coming in from a REST application or a CLI will happen differently. Cookoo includes routers for REST and console applications.
2. _Context_: an execution context passed through the chain of commands as they are executed. It contains information passed around the application along with access to other useful functions such as data sources and logging.

But Cookoo can be used for much more than just command line apps. We use
it for:

* Web apps
* REST servers
* Sophisticated integration tests
* Complex (think Git) CLIs with multiple subcommands

...and that's just what *we've* done with it so far.

To learn more about cookoo, check out the extensive [Go documentation](http://godoc.org/github.com/Masterminds/cookoo). And we're just getting started on the tutorial documentation. Check out the `doc/` directory in our [Git repository](https://github.com/Masterminds/cookoo) to view our work-in-progress.
