---
layout: article
title: Tips For Creating Reusable Commands
keywords: "cookoo, Go, golang, chain-of-command, tips, tutorial"
description: ""
permalink: tutorial/tips-creating-commands/
---
# Tips For Creating Reusable Commands
Reusable commands sit at the heart of cookoo. In this tutorial we will cover a number of tips for creating commands.

## Getting Parameters Passed Into A Command

For this example we can take a look at a simple set of commands for a route. The route is:

```go
registry.Route("GET /", "The Homepage").
    Does(MyMessage, "msg").
    Does(ActOnMessage, "out").
    Using("type").WithDefault("test").
    Using("content").From("cxt:msg")
```

And the first command, `MyMessage` is:

```go
func MyMessage(c cookoo.Context, p *cookoo.Params) (interface{}, cookoo.Interrupt) {
    msg := "This is a test"

    return msg, nil
}
```

In this case the name `msg` as the second argument on the `Does` function that includes `MyMessage` is the name associated with the returned value `msg` from the command containing the string "This is a test".

In the second command, `ActOnMessage` the `Using()` function sets the name of the param to be passed in. The `From` function that follows tells cookoo where to get the value from. The string "cxt:msg" represents the `msg` option on the context. `WithDefault()` sets a default value to use instead of specifying where a value comes from.

```go
func ActOnMessage(c cookoo.Context, p *cookoo.Params) (interface{}, cookoo.Interrupt) {
    type, ok := p.Has("type")
    content := p.Get("content", "")

    return true, nil
}
```
Here we look at two ways to retrieve a parameter that had been passed in. First, there is the `p.Has()` function. It returns a value, if one exist, and if the value was present. The second function, `p.Get()` returns a value if one exists. If none exists the second argument contains a default value to use.

### Requiring Parameters

There are cases where data is required to be passed in. There are two useful function to use inside a command.

```go
func ActOnMessage(c cookoo.Context, p *cookoo.Params) (interface{}, cookoo.Interrupt) {
    ok, missing := p.Requires("type", "content")
    ...
}
```
`p.Requires()` takes in a comma separated list of names to check. If one is missing `ok` will be false and `missing` will contain a list of the missing parameters.

```go
func ActOnMessage(c cookoo.Context, p *cookoo.Params) (interface{}, cookoo.Interrupt) {
    ok, missing := p.Requires("type", "content")
    ...
}
```
`Requires()` checks if parameter is declared via `Using()`. `RequiresValue()` takes this a step further and makes sure a value exists as well. For example,

```go
func ActOnMessage(c cookoo.Context, p *cookoo.Params) (interface{}, cookoo.Interrupt) {
    ok, missing := p.RequiresValue("type", "content")
    ...
}
```

## Get Values From The Context

Sometimes there is a case to retrieve values from the context that were not passed in as parameters. There are two functions to help that along.

```go
func Foo(c cookoo.Context, p *cookoo.Params) (interface{}, cookoo.Interrupt) {
    writer, ok = cxt.Has("http.ResponseWriter")
    req := cxt.Get("http.Request", nil).(*http.Request)

    return true, nil
}
```

The functions `c.Has()` and `c.Get()` work the same their counterparts on the `Params`. The difference is they have access to anything on the context.

_Note, the preferred method of passing information into commands is `Params` as they allow for the control of the data as it comes and goes from each command._

## Return An Error
Sometimes there is an error and a command needs to let cookoo know about it. The two types of errors a command can return are `cookoo.FatalError` and `cookoo.RecoverableError`. These are returned as the second value, the `cookoo.Interrupt`. For example,

```go
func Foo(c cookoo.Context, p *cookoo.Params) (interface{}, cookoo.Interrupt) {
    return false, &cookoo.RecoverableError{"Not enough arguments."}
}
```

If you need to stop a chain of commands without there being an error there is a `cookoo.Stop` option as well.

## Reroute A Chain of Commands
If you are executing one chain of commands and want to reroute to a different set of commands on a different route, you can use `cookoo.Reroute`. This is a return option for the second return value. For example,

```go
func Foo(c cookoo.Context, p *cookoo.Params) (interface{}, cookoo.Interrupt) {
    return false, &cookoo.Reroute{"GET /hello"}
}
```

The string passed into `Reroute` is the other route to execute.