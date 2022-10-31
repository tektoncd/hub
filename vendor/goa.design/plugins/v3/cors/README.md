# CORS Plugin

The `cors` plugin is a [Goa](https://github.com/goadesign/goa/tree/v3) plugin
that makes it possible to define Cross-Origin Resource Sharing (CORS) policies for
the server endpoints.

## Enabling the Plugin

To enable the plugin and make use of the CORS DSL simply import both the `cors` and
the `dsl` packages as follows:

```go
import (
  cors "goa.design/plugins/v3/cors/dsl"
  . "goa.design/goa/v3/dsl"
)
```
The `cors` package exports functions that can be used in the design to configure CORS
options, see below.

## Effects on Code Generation

Enabling the plugin changes the behavior of both the `gen` and `example` commands
of the `goa` tool.

The `gen` command output is modified as follows:

1. A new CORS handler is appended to the HTTP server initialization code.
   This handler is configured to handle the preflight (OPTIONS) request from the client
   (browser) for the applicable endpoints. The handler simply returns a 200 OK
   response containing the CORS headers.
2. All HTTP endpoint handlers are modified to add the CORS headers in the response
   based on the CORS policy definition.

The `example` command output is modified as follows:

1. The example server is initialized with the CORS handler to handle the preflight
   requests.

## Design

This plugin adds the following functions to the goa DSL:

* `Origin` is used in `API` or `Service` DSLs to define the CORS policy that apply
  globally to all the endpoints defined in the design (`API`) or to all the endpoints
  in a service (`Service`).
* Origin specific functions such as `Methods`, `Expose`, `Headers`, `MaxAge`, and
  `Credentials` which are only used in the `Origin` DSL to define CORS headers to
  be set in the response.

The usage and effect of the DSL functions are described in the [Godocs](https://godoc.org/goa.design/plugins/cors/dsl)

Here is an example defining a CORS policy at a service level.

```go
var _ = Service("calc", func() {
  // Sets CORS response headers for requests with Origin header matching the string "localhost"
  cors.Origin("localhost")

  // Sets CORS response headers for requests with Origin header matching strings ending with ".domain.com" (e.g. "my.domain.com")
  cors.Origin("*.domain.com", func() {
    Headers("X-Shared-Secret", "X-Api-Version")
    MaxAge(100)
    Credentials()
  })

  // Sets CORS response headers for requests with any Origin header
  cors.Origin("*")

  // Sets CORS response headers for requests with Origin header matching the regular expression ".*domain.*"
  cors.Origin("/.*domain.*/", func() {
    Headers("*")
    Methods("GET", "POST")
    Expose("X-Time")
  })
})
```

Defining a CORS policy at the API-level is similar to the example above.
