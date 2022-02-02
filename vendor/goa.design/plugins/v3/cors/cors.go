package cors

import (
	"regexp"
	"strings"
)

// MatchOrigin returns true if the given Origin header value matches the
// origin specification.
// Spec can be one of:
// - a plain string identifying an origin. eg http://swagger.goa.design
// - a plain string containing a wildcard. eg *.goa.design
// - the special string * that matches every host
func MatchOrigin(origin, spec string) bool {
	if spec == "*" {
		return true
	}

	// Check regular expression
	if strings.HasPrefix(spec, "/") && strings.HasSuffix(spec, "/") {
		stripped := strings.Trim(spec, "/")
		r := regexp.MustCompile(stripped)
		return MatchOriginRegexp(origin, r)
	}

	if !strings.Contains(spec, "*") {
		return origin == spec
	}
	parts := strings.SplitN(spec, "*", 2)
	if !strings.HasPrefix(origin, parts[0]) {
		return false
	}
	if !strings.HasSuffix(origin, parts[1]) {
		return false
	}
	return true
}

// MatchOriginRegexp returns true if the given Origin header value matches the
// origin specification.
// Spec must be a valid regex
func MatchOriginRegexp(origin string, spec *regexp.Regexp) bool {
	return spec.Match([]byte(origin))
}
