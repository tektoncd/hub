package expr

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"goa.design/goa/v3/eval"
	"goa.design/goa/v3/expr"
)

type (
	// OriginExpr describes a CORS policy.
	OriginExpr struct {
		// Origin is the origin string.
		Origin string
		// Methods is the list of authorized HTTP methods.
		Methods []string
		// Exposed is the list of headers exposed to clients.
		Exposed []string
		// Headers is the list of authorized headers, "*" authorizes all.
		Headers []string
		// MaxAge is the duration to cache a preflight request response.
		MaxAge uint
		// Credentials sets Access-Control-Allow-Credentials header in the
		// response.
		Credentials bool
		// Regexp tells whether the Origin string is a regular expression.
		Regexp bool
		// Parent expression, ServiceExpr or APIExpr.
		Parent eval.Expression
	}
)

// Origins returns the origin expressions (sorted alphabetically
// by origin string) for the given service.
func Origins(svc string) []*OriginExpr {
	origins := make(map[string]*OriginExpr)
	for s, no := range Root.ServiceOrigins {
		if s == svc {
			for n, o := range no {
				origins[n] = o
			}
		}
	}
	for n, o := range Root.APIOrigins {
		if _, ok := origins[n]; !ok {
			// Include API level origins not found in service level
			origins[n] = o
		}
	}
	names := make([]string, 0, len(origins))
	for n := range origins {
		names = append(names, n)
	}
	sort.Strings(names)
	oexps := make([]*OriginExpr, 0, len(names))
	for _, n := range names {
		oexps = append(oexps, origins[n])
	}
	return oexps
}

// PreflightPaths returns the paths that should handle OPTIONS requests
// for the given service.
func PreflightPaths(svc string) []string {
	var paths []string
	s := expr.Root.API.HTTP.Service(svc)
	if s == nil {
		return paths
	}
	for _, e := range s.HTTPEndpoints {
		for _, r := range e.Routes {
			if r.Method == "OPTIONS" {
				continue
			}
			fps := r.FullPaths()
			for _, fp := range fps {
				found := false
				for _, p := range paths {
					if fp == p {
						found = true
						break
					}
				}
				if !found {
					paths = append(paths, fp)
				}
			}
		}
	}
	for _, fs := range s.FileServers {
		fps := fs.RequestPaths
		for _, fp := range fps {
			found := false
			for _, p := range paths {
				if fp == p {
					found = true
					break
				}
			}
			if !found {
				paths = append(paths, fp)
			}
		}
	}
	return paths
}

// EvalName returns the generic expression name used in error messages.
func (o *OriginExpr) EvalName() string {
	var suffix string
	if o.Parent != nil {
		suffix = fmt.Sprintf(" of %s", o.Parent.EvalName())
	}
	return "CORS" + suffix
}

// Validate ensures the origin expression is valid.
func (o *OriginExpr) Validate() *eval.ValidationErrors {
	verr := new(eval.ValidationErrors)
	if !o.Regexp && strings.Count(o.Origin, "*") > 1 {
		verr.Add(o, "invalid origin, can only contain one wildcard character")
	}
	if o.Regexp {
		_, err := regexp.Compile(o.Origin)
		if err != nil {
			verr.Add(o, "invalid origin, should be a valid regular expression")
		}
	}
	return verr
}
