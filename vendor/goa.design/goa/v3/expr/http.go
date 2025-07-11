package expr

import (
	"regexp"
)

type (
	// HTTPExpr contains the API level HTTP specific expressions.
	HTTPExpr struct {
		// Path is the common request path prefix to all the service
		// HTTP endpoints.
		Path string
		// Params defines the HTTP request path and query parameters
		// common to all the API endpoints.
		Params *MappedAttributeExpr
		// Headers defines the HTTP request headers common to to all
		// the API endpoints.
		Headers *MappedAttributeExpr
		// Cookies defines the HTTP request cookies common to to all
		// the API endpoints.
		Cookies *MappedAttributeExpr
		// Consumes lists the mime types supported by the API
		// controllers.
		Consumes []string
		// Produces lists the mime types generated by the API
		// controllers.
		Produces []string
		// Services contains the services created by the DSL.
		Services []*HTTPServiceExpr
		// Errors lists the error HTTP responses.
		Errors []*HTTPErrorExpr
		// SSE contains the Server-Sent Events configuration for all
		// streaming endpoints in the API.
		SSE *HTTPSSEExpr
	}
)

// HTTPWildcardRegex is the regular expression used to capture path
// parameters.
var HTTPWildcardRegex = regexp.MustCompile(`/{\*?([a-zA-Z0-9_]+)}`)

// ExtractHTTPWildcards returns the names of the wildcards that appear in
// a HTTP path.
func ExtractHTTPWildcards(path string) []string {
	matches := HTTPWildcardRegex.FindAllStringSubmatch(path, -1)
	wcs := make([]string, len(matches))
	for i, m := range matches {
		wcs[i] = m[1]
	}
	return wcs
}

// Service returns the service with the given name if any.
func (h *HTTPExpr) Service(name string) *HTTPServiceExpr {
	for _, res := range h.Services {
		if res.Name() == name {
			return res
		}
	}
	return nil
}

// ServiceFor creates a new or returns the existing service definition for the
// given service.
func (h *HTTPExpr) ServiceFor(s *ServiceExpr) *HTTPServiceExpr {
	if res := h.Service(s.Name); res != nil {
		return res
	}
	res := &HTTPServiceExpr{
		ServiceExpr: s,
	}
	h.Services = append(h.Services, res)
	return res
}

// EvalName returns the name printed in case of evaluation error.
func (*HTTPExpr) EvalName() string {
	return "API HTTP"
}

// Finalize initializes Consumes and Produces with defaults if not set.
func (h *HTTPExpr) Finalize() {
	if len(h.Consumes) == 0 {
		h.Consumes = []string{"application/json", "application/xml", "application/gob"}
	}
	if len(h.Produces) == 0 {
		h.Produces = []string{"application/json", "application/xml", "application/gob"}
	}
}
