# ZapLogger Plugin

The `ZapLogger` plugin is a [Goa v3](https://github.com/goadesign/goa/tree/v3) plugin
that adapt the basic logger to use the [zap](https://github.com/uber-go/zap) library.

## Enabling the Plugin

To enable the plugin import it in your design.go file using the blank identifier `_` as follows:

```go

package design

import . "goa.design/goa/v3/http/design"
import . "goa.design/goa/v3/http/dsl"
import _ "goa.design/plugins/v3/zaplogger" // Enables the plugin

var _ = API("...

```

and generate as usual:

```bash
goa gen PACKAGE
goa example PACKAGE
```

where `PACKAGE` is the Go import path of the design package.
