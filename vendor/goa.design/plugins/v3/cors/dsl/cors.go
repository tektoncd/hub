package dsl

import (
	"strings"

	"goa.design/goa/v3/eval"
	goaexpr "goa.design/goa/v3/expr"
	"goa.design/plugins/v3/cors/expr"

	// Register code generators for the CORS plugin
	_ "goa.design/plugins/v3/cors"
)

// Origin defines the CORS policy for a given origin. The origin can use a wildcard prefix
// such as "https://*.mydomain.com". The special value "*" defines the policy for all origins
// (in which case there should be only one Origin DSL in the parent resource).
// The origin can also be a regular expression in which case it must be wrapped with "/".
//
// Origin must appear in API or Service Expression.
//
// Origin accepts an origin string as the first argument and
// an optional DSL function as the second argument.
//
// Optionally, you can specify the name of an environment variable instead, prefixed by a "$".
// The value you store in that environment variable follows the same rules as the above
// and can similarly be a regular expression.
//
// Example:
//
//	import cors "goa.design/plugins/v3/cors"
//
//	var _ = API("calc", func() {
//	    cors.Origin("http://swagger.goa.design", func() { // Define CORS policy, may be prefixed with "*" wildcard
//	        cors.Headers("X-Shared-Secret")  // One or more authorized headers, use "*" to authorize all
//	        cors.Methods("GET", "POST")      // One or more authorized HTTP methods
//	        cors.Expose("X-Time")            // One or more headers exposed to clients
//	        cors.MaxAge(600)                 // How long to cache a preflight request response
//	        cors.Credentials()               // Sets Access-Control-Allow-Credentials header
//	    })
//
//	    cors.Origin("$ORIGIN") // Simple example to demonstrate using an environment variable
//	})
//
//	var _ = Service("calculator", func() {
//	    cors.Origin("/(api|swagger)[.]goa[.]design/") // Define CORS policy with a regular expression
//
//	    Method("add", func() {
//	        Description("Add two operands")
//	        Payload(Operands)
//	        Error(ErrBadRequest, ErrorResult)
//	    })
//	})
func Origin(origin string, args ...interface{}) {
	o := &expr.OriginExpr{Origin: origin}
	if strings.HasPrefix(origin, "/") && strings.HasSuffix(origin, "/") {
		o.Regexp = true
		o.Origin = strings.Trim(origin, "/")
	}
	if strings.HasPrefix(origin, "$") {
		o.EnvVar = true
		o.Origin = strings.Trim(origin, "$")
	}

	var dsl func()
	{
		if len(args) > 0 {
			if d, ok := args[len(args)-1].(func()); ok {
				dsl = d
			}
		}
	}
	if dsl != nil {
		if !eval.Execute(dsl, o) {
			return
		}
	}

	current := eval.Current()
	switch actual := current.(type) {
	case *goaexpr.APIExpr:
		expr.Root.APIOrigins[origin] = o
	case *goaexpr.ServiceExpr:
		{
			s := actual.Name
			if _, ok := expr.Root.ServiceOrigins[s]; !ok {
				expr.Root.ServiceOrigins[s] = make(map[string]*expr.OriginExpr)
			}
			expr.Root.ServiceOrigins[s][origin] = o
		}
	default:
		eval.IncompatibleDSL()
		return
	}
	o.Parent = current
}

// Methods sets the origin allowed methods.
//
// Methods must be used in an Origin expression.
//
// Example:
//
//	Origin("http://swagger.goa.design", func() {
//	    Methods("GET", "POST")           // One or more authorized HTTP methods
//	})
func Methods(vals ...string) {
	switch o := eval.Current().(type) {
	case *expr.OriginExpr:
		o.Methods = append(o.Methods, vals...)
	default:
		eval.IncompatibleDSL()
	}
}

// Expose sets the origin exposed headers.
//
// Expose must appear in an Origin expression.
//
// Example:
//
//	Origin("http://swagger.goa.design", func() {
//	    Expose("X-Time")               // One or more headers exposed to clients
//	})
func Expose(vals ...string) {
	switch o := eval.Current().(type) {
	case *expr.OriginExpr:
		o.Exposed = append(o.Exposed, vals...)
	default:
		eval.IncompatibleDSL()
	}
}

// Headers sets the authorized headers. "*" authorizes all headers.
//
// Headers must be used in an Origin expression.
//
// Example:
//
//	Origin("http://swagger.goa.design", func() {
//	    Headers("X-Shared-Secret")
//	})
//
//	Origin("http://swagger.goa.design", func() {
//	    Headers("*")
//	})
func Headers(vals ...string) {
	switch o := eval.Current().(type) {
	case *expr.OriginExpr:
		o.Headers = append(o.Headers, vals...)
	default:
		eval.IncompatibleDSL()
	}
}

// MaxAge sets the cache expiry for preflight request responses.
//
// MaxAge must be used in an Origin expression.
//
// Example:
//
//	Origin("http://swagger.goa.design", func() {
//	    MaxAge(600)            // How long to cache a preflight request response
//	})
func MaxAge(val uint) {
	switch o := eval.Current().(type) {
	case *expr.OriginExpr:
		o.MaxAge = val
	default:
		eval.IncompatibleDSL()
	}
}

// Credentials sets the allow credentials response header.
//
// Credentials must be used in an Origin expression.
//
// Example:
//
//	Origin("http://swagger.goa.design", func() {
//	    Credentials()            // Sets Access-Control-Allow-Credentials header
//	})
func Credentials() {
	switch o := eval.Current().(type) {
	case *expr.OriginExpr:
		o.Credentials = true
	default:
		eval.IncompatibleDSL()
	}
}
