---
layout: article
title: "Cookoo, The Chain-of-Command Controller for Go"
keywords: "Go, Chain-of-Command, front controller, golang"
description: "Learn about using the Chain-of-Command pattern in Go."
---
# Cookoo, The Chain-of-Command Controller for Go

Cookoo is a high performance reusable controller pattern that can be used to build console (CLI) applications, API servers, middle ware, bots, and web applications.

In its simplest form, you have a set of re-usable commands. For a callback a series of commands are executed. Each command receives information from a context and can put information back into the context.

## Getting Started with Hello World

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

1. Registry: contains the mapping of commands to callbacks along with how information is passed around the context.
2. Router: handles incoming requests and routes them to the correct callback on the registry. Different types of applications will have different routers. Commands coming in from a REST application or a CLI will happen differently. Cookoo includes routers for REST and console applications.
2. Context: an execution context passed through the chain of commands as they are executed. It contains information passed around the application along with access to other useful functions such as data sources and logging.

To learn more about cookoo please see the [documentation](http://godoc.org/github.com/Masterminds/cookoo).