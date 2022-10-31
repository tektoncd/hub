package cors

import (
	"path/filepath"
	"regexp"
	"strings"

	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/codegen/service"
	"goa.design/goa/v3/eval"
	httpcodegen "goa.design/goa/v3/http/codegen"

	"goa.design/plugins/v3/cors/expr"
)

// ServicesData holds the all the ServiceData indexed by service name.
var ServicesData = make(map[string]*ServiceData)

type (
	// ServiceData contains the data necessary to generate origin handlers
	ServiceData struct {
		// Name is the name of the service.
		Name string
		// Origins is a list of origin expressions defined in API and service levels.
		Origins []*expr.OriginExpr
		// OriginHandler is the name of the handler function that sets CORS headers.
		OriginHandler string
		// PreflightPaths is the list of paths that should handle OPTIONS requests.
		PreflightPaths []string
		// Endpoint is the CORS endpoint data.
		Endpoint *httpcodegen.EndpointData
	}
)

// Register the plugin Generator functions.
func init() {
	codegen.RegisterPlugin("cors", "gen", nil, Generate)
	codegen.RegisterPlugin("cors-example", "example", nil, TweakExample)
}

// Generate produces server code that handle preflight requests and updates
// the HTTP responses with the appropriate CORS headers.
func Generate(genpkg string, roots []eval.Root, files []*codegen.File) ([]*codegen.File, error) {
	for _, f := range files {
		serverCORS(f)
	}
	return files, nil
}

// TweakExample handles the special case where a service only has file servers
// and no method in which case the Goa generator generate a Mount method that
// does not take a second argument but this plugin generates one that does. The
// second argument is the actual HTTP server which is needed so it can be
// configured with the CORS endpoint. So this method simply removes the special
// case from the Goa template generating the example.
func TweakExample(genpkg string, roots []eval.Root, files []*codegen.File) ([]*codegen.File, error) {
	re := regexp.MustCompile("{{ if .Endpoints }}(.+){{ end }}")
	for _, f := range files {
		for _, t := range f.SectionTemplates {
			if t.Name == "server-http-init" {
				t.Source = re.ReplaceAllString(t.Source, "$1")
			}
		}
	}
	return files, nil
}

// buildServiceData builds the data needed to render the CORS handlers.
func buildServiceData(svc string) *ServiceData {
	preflights := expr.PreflightPaths(svc)
	routes := make([]*httpcodegen.RouteData, len(preflights))
	for i, p := range preflights {
		routes[i] = &httpcodegen.RouteData{Verb: "OPTIONS", Path: p}
	}

	return &ServiceData{
		Name:           svc,
		Origins:        expr.Origins(svc),
		PreflightPaths: preflights,
		OriginHandler:  "Handle" + codegen.Goify(svc, true) + "Origin",
		Endpoint: &httpcodegen.EndpointData{
			Method: &service.MethodData{
				VarName: "CORS",
			},
			MountHandler: "MountCORSHandler",
			HandlerInit:  "NewCORSHandler",
			Routes:       routes,
		},
	}
}

// serverCORS updates the HTTP server file to handle preflight paths and
// adds the required CORS headers to the response.
func serverCORS(f *codegen.File) {
	if filepath.Base(f.Path) != "server.go" {
		return
	}

	var svcData *ServiceData
	for _, s := range f.Section("server-struct") {

		data, ok := s.Data.(*httpcodegen.ServiceData)
		if !ok { // other transport, e.g. gRPC
			continue
		}

		codegen.AddImport(f.SectionTemplates[0],
			&codegen.ImportSpec{Path: "goa.design/plugins/v3/cors"})

		if d, ok := ServicesData[data.Service.Name]; !ok {
			svcData = buildServiceData(data.Service.Name)
			ServicesData[data.Service.Name] = svcData
		} else {
			svcData = d
		}
		for _, o := range svcData.Origins {
			if o.Regexp {
				codegen.AddImport(f.SectionTemplates[0],
					&codegen.ImportSpec{Path: "regexp"})
				break
			}
		}
		for _, o := range svcData.Origins {
			if o.EnvVar {
				codegen.AddImport(f.SectionTemplates[0],
					&codegen.ImportSpec{Path: "os"})
				codegen.AddImport(f.SectionTemplates[0],
					&codegen.ImportSpec{Path: "strings"})
				break
			}
		}
		data.Endpoints = append(data.Endpoints, svcData.Endpoint)
		fm := codegen.TemplateFuncs()
		f.SectionTemplates = append(f.SectionTemplates, &codegen.SectionTemplate{
			Name:    "mount-cors",
			Source:  mountCORST,
			Data:    svcData,
			FuncMap: fm,
		})
		f.SectionTemplates = append(f.SectionTemplates, &codegen.SectionTemplate{
			Name:    "cors-handler-init",
			Source:  corsHandlerInitT,
			Data:    svcData,
			FuncMap: fm,
		})
		fm["join"] = strings.Join
		f.SectionTemplates = append(f.SectionTemplates, &codegen.SectionTemplate{
			Name:    "handle-cors",
			Source:  handleCORST,
			Data:    svcData,
			FuncMap: fm,
		})
	}
	for _, s := range f.Section("server-init") {
		s.Source = strings.Replace(s.Source,
			`e.{{ .Method.VarName }}, mux, {{ if .MultipartRequestDecoder }}{{ .MultipartRequestDecoder.InitName }}(mux, {{ .MultipartRequestDecoder.VarName }}){{ else }}decoder{{ end }}, encoder, errhandler, formatter{{ if isWebSocketEndpoint . }}, upgrader, configurer.{{ .Method.VarName }}Fn{{ end }})`,
			`{{ if ne .Method.VarName "CORS" }}e.{{ .Method.VarName }}, mux, {{ if .MultipartRequestDecoder }}{{ .MultipartRequestDecoder.InitName }}(mux, {{ .MultipartRequestDecoder.VarName }}){{ else }}decoder{{ end }}, encoder, errhandler, formatter{{ if isWebSocketEndpoint . }}, upgrader, configurer.{{ .Method.VarName }}Fn{{ end }}{{ end }})`,
			-1)
	}
	for _, s := range f.Section("server-handler") {
		s.Source = strings.Replace(s.Source, "h.(http.HandlerFunc)", svcData.OriginHandler+"(h).(http.HandlerFunc)", -1)
	}
	for _, s := range f.Section("server-files") {
		s.Source = strings.Replace(s.Source, "h.ServeHTTP", svcData.OriginHandler+"(h).ServeHTTP", -1)
	}
}

// Data: ServiceData
var corsHandlerInitT = `{{ printf "%s creates a HTTP handler which returns a simple 200 response." .Endpoint.HandlerInit | comment }}
func {{ .Endpoint.HandlerInit }}() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})
}
`

// Data: ServiceData
var mountCORST = `{{ printf "%s configures the mux to serve the CORS endpoints for the service %s." .Endpoint.MountHandler .Name | comment }}
func {{ .Endpoint.MountHandler }}(mux goahttp.Muxer, h http.Handler) {
	h = {{ .OriginHandler }}(h)
	{{- range $p := .PreflightPaths }}
	mux.Handle("OPTIONS", "{{ $p }}", h.ServeHTTP)
	{{- end }}
}
`

// Data: ServiceData
var handleCORST = `{{ printf "%s applies the CORS response headers corresponding to the origin for the service %s." .OriginHandler .Name | comment }}
func {{ .OriginHandler }}(h http.Handler) http.Handler {
{{- range $i, $policy := .Origins }}
	{{- if $policy.EnvVar }}
	originStr{{$i}}, present := os.LookupEnv({{ printf "%q" $policy.Origin }})
	if !present {
		panic("CORS origin environment variable \"{{ $policy.Origin }}\" not set!")
	}
	{{- end }}
	{{- if $policy.Regexp }}
	spec{{$i}} := regexp.MustCompile({{ printf "%q" $policy.Origin }})
	{{- end }}
{{- end }}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	if origin == "" {
		// Not a CORS request
		h.ServeHTTP(w, r)
		return
	}
	{{- range $i, $policy := .Origins }}
		{{- if $policy.Regexp }}
	if cors.MatchOriginRegexp(origin, spec{{$i}}) {
		{{- else }}
	{{- if $policy.EnvVar }}
	if cors.MatchOrigin(origin, originStr{{$i}}) {
	{{- else }}
	if cors.MatchOrigin(origin, {{ printf "%q" $policy.Origin }}) {
	{{- end }}
		{{- end }}
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Vary", "Origin")
			{{- if $policy.Exposed }}
		w.Header().Set("Access-Control-Expose-Headers", "{{ join $policy.Exposed ", " }}")
			{{- end }}
			{{- if gt $policy.MaxAge 0 }}
		w.Header().Set("Access-Control-Max-Age", "{{ $policy.MaxAge }}")
			{{- end }}
			{{- if eq $policy.Credentials true }}
		w.Header().Set("Access-Control-Allow-Credentials", "{{ $policy.Credentials }}")
			{{- end }}
		if acrm := r.Header.Get("Access-Control-Request-Method"); acrm != "" {
			// We are handling a preflight request
				{{- if $policy.Methods }}
			w.Header().Set("Access-Control-Allow-Methods", "{{ join $policy.Methods ", " }}")
				{{- end }}
				{{- if $policy.Headers }}
			w.Header().Set("Access-Control-Allow-Headers", "{{ join $policy.Headers ", " }}")
				{{- end }}
		}
		h.ServeHTTP(w, r)
		return
	}
	{{- end }}
	h.ServeHTTP(w, r)
	return
  })
}
`
