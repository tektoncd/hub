package codegen

import (
	"bytes"
	"fmt"
	"net/http"
	"slices"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/codegen/service"
	"goa.design/goa/v3/expr"
)

// HTTPServices holds the data computed from the design needed to generate the
// transport code of the services.
var HTTPServices = make(ServicesData)

var (
	// pathInitTmpl is the template used to render path constructors code.
	pathInitTmpl = template.Must(
		template.New("path-init").
			Funcs(template.FuncMap{"goify": codegen.Goify}).
			Parse(readTemplate("path_init", "query_slice_conversion")),
	)
	// requestInitTmpl is the template used to render request constructors.
	requestInitTmpl = template.Must(
		template.New("request-init").
			Funcs(template.FuncMap{
				"goTypeRef": func(dt expr.DataType, svc string) string {
					return service.Services.Get(svc).Scope.GoTypeRef(&expr.AttributeExpr{Type: dt})
				},
				"isAliased": func(dt expr.DataType) bool {
					_, ok := dt.(expr.UserType)
					return ok
				},
				"isWebSocketEndpoint": isWebSocketEndpoint,
			}).
			Parse(readTemplate("request_init")),
	)
)

type (
	// ServicesData encapsulates the data computed from the design.
	ServicesData map[string]*ServiceData

	// ServiceData contains the data used to render the code related to a
	// single service.
	ServiceData struct {
		// Service contains the related service data.
		Service *service.Data
		// Endpoints describes the endpoint data for this service.
		Endpoints []*EndpointData
		// FileServers lists the file servers for this service.
		FileServers []*FileServerData
		// ServerStruct is the name of the HTTP server struct.
		ServerStruct string
		// MountPointStruct is the name of the mount point struct.
		MountPointStruct string
		// ServerInit is the name of the constructor of the server
		// struct.
		ServerInit string
		// MountServer is the name of the mount function.
		MountServer string
		// ServerService is the name of service function.
		ServerService string
		// ClientStruct is the name of the HTTP client struct.
		ClientStruct string
		// ServerBodyAttributeTypes is the list of user types used to
		// define the request, response and error response type
		// attributes in the server code.
		ServerBodyAttributeTypes []*TypeData
		// ClientBodyAttributeTypes is the list of user types used to
		// define the request, response and error response type
		// attributes in the client code.
		ClientBodyAttributeTypes []*TypeData
		// ServerTypeNames records the user type names used to define
		// the endpoint request and response bodies for server code.
		// The type name is used as the key and a bool as the value
		// which if true indicates that the type has been generated
		// in the server package.
		ServerTypeNames map[string]bool
		// ClientTypeNames records the user type names used to define
		// the endpoint request and response bodies for client code.
		// The type name is used as the key and a bool as the value
		// which if true indicates that the type has been generated
		// in the client package.
		ClientTypeNames map[string]bool
		// ServerTransformHelpers is the list of transform functions
		// required by the various server side constructors.
		ServerTransformHelpers []*codegen.TransformFunctionData
		// ClientTransformHelpers is the list of transform functions
		// required by the various client side constructors.
		ClientTransformHelpers []*codegen.TransformFunctionData
		// Scope initialized with all the server and client types.
		Scope *codegen.NameScope
	}

	// EndpointData contains the data used to render the code related to a
	// single service HTTP endpoint.
	EndpointData struct {
		// Method contains the related service method data.
		Method *service.MethodData
		// ServiceName is the name of the service exposing the endpoint.
		ServiceName string
		// ServiceVarName is the goified service name (first letter
		// lowercase).
		ServiceVarName string
		// ServicePkgName is the name of the service package.
		ServicePkgName string
		// Payload describes the method HTTP payload.
		Payload *PayloadData
		// Result describes the method HTTP result.
		Result *ResultData
		// Errors describes the method HTTP errors.
		Errors []*ErrorGroupData
		// Routes describes the possible routes for this endpoint.
		Routes []*RouteData
		// BasicScheme is the basic auth security scheme if any.
		BasicScheme *service.SchemeData
		// HeaderSchemes lists all the security requirement schemes that
		// apply to the method and are encoded in the request header.
		HeaderSchemes service.SchemesData
		// BodySchemes lists all the security requirement schemes that
		// apply to the method and are encoded in the request body.
		BodySchemes service.SchemesData
		// QuerySchemes lists all the security requirement schemes that
		// apply to the method and are encoded in the request query
		// string.
		QuerySchemes service.SchemesData
		// Requirements contains the security requirements for the
		// method.
		Requirements service.RequirementsData

		// server

		// MountHandler is the name of the mount handler function.
		MountHandler string
		// HandlerInit is the name of the constructor function for the
		// http handler function.
		HandlerInit string
		// RequestDecoder is the name of the request decoder function.
		RequestDecoder string
		// ResponseEncoder is the name of the response encoder function.
		ResponseEncoder string
		// ErrorEncoder is the name of the error encoder function.
		ErrorEncoder string
		// MultipartRequestDecoder indicates the request decoder for
		// multipart content type.
		MultipartRequestDecoder *MultipartData
		// ServerWebSocket holds the data to render the server struct which
		// implements the server stream interface.
		ServerWebSocket *WebSocketData
		// SSE holds the data to render the server struct which implements the
		// server stream interface for SSE.
		SSE *SSEData
		// Redirect defines a redirect for the endpoint.
		Redirect *RedirectData

		// client

		// ClientStruct is the name of the HTTP client struct.
		ClientStruct string
		// EndpointInit is the name of the constructor function for the
		// client endpoint.
		EndpointInit string
		// RequestInit is the request builder function.
		RequestInit *InitData
		// RequestEncoder is the name of the request encoder function.
		RequestEncoder string
		// ResponseDecoder is the name of the response decoder function.
		ResponseDecoder string
		// MultipartRequestEncoder indicates the request encoder for
		// multipart content type.
		MultipartRequestEncoder *MultipartData
		// ClientWebSocket holds the data to render the client struct which
		// implements the client stream interface.
		ClientWebSocket *WebSocketData
		// BuildStreamPayload is the name of the function used to create the
		// payload for endpoints that use SkipRequestBodyEncodeDecode.
		BuildStreamPayload string
	}

	// FileServerData lists the data needed to generate file servers.
	FileServerData struct {
		// MountHandler is the name of the mount handler function.
		MountHandler string
		// RequestPaths is the set of HTTP paths to the server.
		RequestPaths []string
		// Root is the root server file path.
		FilePath string
		// Dir is true if the file server servers files under a
		// directory, false if it serves a single file.
		IsDir bool
		// PathParam is the name of the parameter used to capture the
		// path for file servers that serve files under a directory.
		PathParam string
		// Redirect defines a redirect for the endpoint.
		Redirect *RedirectData
		// VarName is the name of the variable that holds the file server.
		VarName string
		// ArgName is the name of the argument used to initialize the
		// file server.
		ArgName string
	}

	// RedirectData lists the data needed to generate a redirect.
	RedirectData struct {
		// URL is the URL that is being redirected to.
		URL string
		// StatusCode is the HTTP status code.
		StatusCode string
	}

	// PayloadData contains the payload information required to generate the
	// transport decode (server) and encode (client) code.
	PayloadData struct {
		// Name is the name of the payload type.
		Name string
		// Ref is the fully qualified reference to the payload type.
		Ref string
		// Request contains the data for the corresponding HTTP request.
		Request *RequestData
		// DecoderReturnValue is a reference to the decoder return value
		// if there is no payload constructor (i.e. if Init is nil).
		DecoderReturnValue string
	}

	// ResultData contains the result information required to generate the
	// transport decode (client) and encode (server) code.
	ResultData struct {
		// Name is the name of the result type.
		Name string
		// Ref is the reference to the result type.
		Ref string
		// IsStruct is true if the result type is a user type defining
		// an object.
		IsStruct bool
		// Inits contains the data required to render the result
		// constructors if any.
		Inits []*InitData
		// Responses contains the data for the corresponding HTTP
		// responses.
		Responses []*ResponseData
		// View is the view used to render the result.
		View string
		// MustInit indicates if a variable holding the result type must be
		// initialized. It is used by server response encoder to initialize
		// the result variable only if there are multiple responses, or the
		// response has a body, a header or a cookie.
		MustInit bool
	}

	// ErrorGroupData contains the error information required to generate
	// the transport decode (client) and encode (server) code for all errors
	// with responses using a given HTTP status code.
	ErrorGroupData struct {
		// StatusCode is the response HTTP status code.
		StatusCode string
		// Errors contains the information for each error.
		Errors []*ErrorData
	}

	// ErrorData contains the error information required to generate the
	// transport decode (client) and encode (server) code.
	ErrorData struct {
		// Name is the error name.
		Name string
		// Ref is a reference to the error type.
		Ref string
		// Response is the error response data.
		Response *ResponseData
	}

	// RequestData describes a request.
	RequestData struct {
		// PathParams describes the information about params that are
		// present in the request path.
		PathParams []*ParamData
		// QueryParams describes the information about the params that
		// are present in the request query string.
		QueryParams []*ParamData
		// Headers contains the HTTP request headers used to build the
		// method payload.
		Headers []*HeaderData
		// Cookies contains the HTTP request cookies used to build the
		// method payload.
		Cookies []*CookieData
		// ServerBody describes the request body type used by server
		// code. The type is generated using pointers for all fields so
		// that it can be validated.
		ServerBody *TypeData
		// ClientBody describes the request body type used by client
		// code. The type does NOT use pointers for every fields since
		// no validation is required.
		ClientBody *TypeData
		// PayloadInit contains the data required to render the
		// payload constructor used by server code if any.
		PayloadInit *InitData
		// PayloadType is the type of the payload.
		PayloadType expr.DataType
		// PayloadAttr sets the request body from the specified payload type
		// attribute. This field is set when the design uses Body("name") syntax
		// to set the request body and the payload type is an object.
		PayloadAttr string
		// MustHaveBody is true if the request body cannot be empty.
		MustHaveBody bool
		// MustValidate is true if the request body or at least one
		// parameter or header requires validation.
		MustValidate bool
		// Multipart if true indicates the request is a multipart
		// request.
		Multipart bool
	}

	// ResponseData describes a response.
	ResponseData struct {
		// StatusCode is the return code of the response.
		StatusCode string
		// Description is the response description.
		Description string
		// Headers provides information about the HTTP response headers.
		Headers []*HeaderData
		// Cookies provides information about the HTTP response cookies.
		Cookies []*CookieData
		// ContentType contains the value of the response
		// "Content-Type" header.
		ContentType string
		// ErrorHeader contains the value of the response "goa-error"
		// header if any.
		ErrorHeader string
		// ServerBody is the type of the response body used by server
		// code, nil if body should be empty. The type does NOT use
		// pointers for all fields. If the method result is a result
		// type and the response data describes a success response, then
		// this field contains a type for every view in the result type.
		// The type name is suffixed with the name of the view (except
		// for "default" view where no suffix is added). A constructor
		// is also generated server side for each view to transform the
		// result type to the corresponding response body type. If
		// method result is not a result type or if the response
		// describes an error response, then this field contains at most
		// one item.
		ServerBody []*TypeData
		// ClientBody is the type of the response body used by client
		// code, nil if body should be empty. The type uses pointers for
		// all fields so they can be validated.
		ClientBody *TypeData
		// Init contains the data required to render the result or error
		// constructor if any.
		ResultInit *InitData
		// TagName is the name of the attribute used to test whether the
		// response is the one to use.
		TagName string
		// TagValue is the value the result attribute named by TagName
		// must have for this response to be used.
		TagValue string
		// TagPointer is true if the tag attribute is a pointer.
		TagPointer bool
		// MustValidate is true if at least one header requires validation.
		MustValidate bool
		// ResultAttr sets the response body from the specified result
		// type attribute. This field is set when the design uses
		// Body("name") syntax to set the response body and the result
		// type is an object.
		ResultAttr string
		// ViewedResult indicates whether the response body type is a
		// result type.
		ViewedResult *service.ViewedResultTypeData
	}

	// InitData contains the data required to render a constructor.
	InitData struct {
		// Name is the constructor function name.
		Name string
		// Description is the function description.
		Description string
		// ServerArgs is the list of constructor arguments for server
		// side code.
		ServerArgs []*InitArgData
		// ClientArgs is the list of constructor arguments for client
		// side code.
		ClientArgs []*InitArgData
		// CLIArgs is the list of arguments that should be initialized
		// from CLI flags. This is used for implicit attributes which
		// as the time of writing is only used for the basic auth
		// username and password.
		CLIArgs []*InitArgData
		// ServerCode is the code that builds the payload from the
		// request on the server when it contains user types.
		ServerCode string
		// ClientCode is the code that builds the payload or result type
		// from the request or response state on the client when it
		// contains user types.
		ClientCode string
		// ReturnTypePkg is the package where the return type is present.
		ReturnTypePkg string
		// ReturnTypeName is the qualified (including the package name)
		// name of the payload, result or error type.
		ReturnTypeName string
		// ReturnTypeRef is the qualified (including the package name)
		// reference to the payload, result or error type.
		ReturnTypeRef string
		// ReturnTypeAttribute is the name of the attribute initialized by this
		// constructor when it only initializes one attribute (i.e. body was
		// defined with Body("name") syntax).
		ReturnTypeAttribute string
		// ReturnIsStruct is true if the payload, result or error type is a struct.
		ReturnIsStruct bool
		// ReturnIsPrimitivePointer indicates whether the payload, result or error
		// type is a primitive pointer.
		ReturnIsPrimitivePointer bool
	}

	// AttributeData contains the information needed to generate the code
	// related to a specific payload or result attribute.
	AttributeData struct {
		// Name is the name of the attribute.
		Name string
		// VarName is the name of the variable that holds the attribute value.
		VarName string
		// Pointer is true if the attribute value is a pointer.
		Pointer bool
		// Required is true if the attribute is required in the parent attribute.
		Required bool
		// Type is the attribute type.
		Type expr.DataType
		// TypeName is the generated attribute type name.
		TypeName string
		// TypeRef is the generated attribute type reference.
		TypeRef string
		// Description is the attribute description as defined in the design.
		Description string
		// FieldName is the name of the data structure field that should
		// be initialized with the value if any.
		FieldName string
		// FieldType is the type of the data structure field that should be
		// initialized with the attribute value or read into the attribute value.
		FieldType expr.DataType
		// FieldPointer if true indicates that the data structure field is a
		// pointer.
		FieldPointer bool
		// DefaultValue is the default value of the attribute if any.
		DefaultValue any
		// Validate contains the validation code for the attribute value if any.
		Validate string
		// Example is an example attribute value
		Example any
	}

	// InitArgData represents a single constructor argument.
	InitArgData struct {
		*AttributeData
		// Reference to the argument, e.g. "&body".
		Ref string
	}

	// RouteData describes a route.
	RouteData struct {
		// Verb is the HTTP method.
		Verb string
		// Path is the fullpath including wildcards.
		Path string
		// PathInit contains the information needed to render and call
		// the path constructor for the route.
		PathInit *InitData
	}

	// Element defines the common fields needed to generate HTTP request and
	// response elements including headers, parameters and cookies.
	Element struct {
		*AttributeData
		// HTTPName is the name of the HTTP element (header name, query string name
		// or cookie name)
		HTTPName string
		// AttributeName is the name of the corresponding attribute.
		AttributeName string
		// StringSlice is true if the attribute type is array of strings.
		StringSlice bool
		// Slice is true if the attribute type is an array.
		Slice bool
	}

	// ParamData describes a HTTP request parameter (query string or path
	// parameter).
	ParamData struct {
		*Element
		// MapStringSlice is true if the param type is a map of string
		// slice.
		MapStringSlice bool
		// Map is true if the param type is a map.
		Map bool
		// MapQueryParams indicates that the query params must be mapped
		// to the entire payload (empty string) or a payload attribute
		// (attribute name).
		MapQueryParams *string
	}

	// HeaderData describes a HTTP request or response header.
	HeaderData struct {
		*Element
		// CanonicalName is the canonical header key.
		CanonicalName string
	}

	// CookieData describes a HTTP request or response cookie.
	CookieData struct {
		*Element
		// MaxAge is the cookie "max-age" attribute.
		MaxAge string
		// Path is the cookie "path" attribute.
		Path string
		// Domain is the cookie "domain" attribute.
		Domain string
		// Secure sets the cookie "secure" attribute to "Secure" if true.
		Secure bool
		// HTTPOnly sets the cookie "http-only" attribute to "HttpOnly" if true.
		HTTPOnly bool
		// SameSite sets the cookie "same-site" attribute to the given value.
		SameSite string
	}

	// TypeData contains the data needed to render a type definition.
	TypeData struct {
		// Name is the type name.
		Name string
		// VarName is the Go type name.
		VarName string
		// Description is the type human description.
		Description string
		// Init contains the data needed to render and call the type
		// constructor if any.
		Init *InitData
		// Def is the type definition Go code.
		Def string
		// Ref is the reference to the type.
		Ref string
		// ValidateDef contains the validation code.
		ValidateDef string
		// ValidateRef contains the call to the validation code.
		ValidateRef string
		// Example is an example value for the type.
		Example any
		// View is the view used to render the (result) type if any.
		View string
	}

	// MultipartData contains the data needed to render multipart
	// encoder/decoder.
	MultipartData struct {
		// FuncName is the name used to generate function type.
		FuncName string
		// InitName is the name of the constructor.
		InitName string
		// VarName is the name of the variable referring to the function.
		VarName string
		// ServiceName is the name of the service.
		ServiceName string
		// MethodName is the name of the method.
		MethodName string
		// Payload is the payload data required to generate
		// encoder/decoder.
		Payload *PayloadData
	}
)

// Get retrieves the transport data for the service with the given name
// computing it if needed. It returns nil if there is no service with the given
// name.
func (d ServicesData) Get(name string) *ServiceData {
	if data, ok := d[name]; ok {
		return data
	}
	service := expr.Root.API.HTTP.Service(name)
	if service == nil {
		return nil
	}
	d[name] = d.analyze(service)
	return d[name]
}

// Endpoint returns the service method transport data for the endpoint with the
// given name, nil if there isn't one.
func (svc *ServiceData) Endpoint(name string) *EndpointData {
	for _, e := range svc.Endpoints {
		if e.Method.Name == name {
			return e
		}
	}
	return nil
}

// analyze creates the data necessary to render the code of the given service.
// It records the user types needed by the service definition in userTypes.
func (ServicesData) analyze(httpSvc *expr.HTTPServiceExpr) *ServiceData {
	svc := service.Services.Get(httpSvc.ServiceExpr.Name)
	scope := codegen.NewNameScope()
	scope.Unique("c") // 'c' is reserved as the client's receiver name.
	scope.Unique("v") // 'v' is reserved as the request builder payload argument name.
	sd := &ServiceData{
		Service:          svc,
		ServerStruct:     "Server",
		MountPointStruct: "MountPoint",
		ServerInit:       "New",
		MountServer:      "Mount",
		ServerService:    "Service",
		ClientStruct:     "Client",
		ServerTypeNames:  make(map[string]bool),
		ClientTypeNames:  make(map[string]bool),
		Scope:            scope,
	}

	for _, s := range httpSvc.FileServers {
		paths := make([]string, len(s.RequestPaths))
		for i, p := range s.RequestPaths {
			idx := strings.LastIndex(p, "/{")
			if idx == 0 {
				paths[i] = "/"
			} else if idx > 0 {
				paths[i] = p[:idx]
			} else {
				paths[i] = p
			}
		}
		var pp string
		if s.IsDir() {
			pp = expr.ExtractHTTPWildcards(s.RequestPaths[0])[0]
		}
		var redirect *RedirectData
		if s.Redirect != nil {
			redirect = &RedirectData{
				URL:        s.Redirect.URL,
				StatusCode: statusCodeToHTTPConst(s.Redirect.StatusCode),
			}
		}
		data := &FileServerData{
			MountHandler: scope.Unique(fmt.Sprintf("Mount%s", codegen.Goify(s.FilePath, true))),
			RequestPaths: paths,
			FilePath:     s.FilePath,
			IsDir:        s.IsDir(),
			PathParam:    pp,
			Redirect:     redirect,
			VarName:      scope.Unique(codegen.Goify(s.FilePath, true)),
			ArgName:      scope.Unique(fmt.Sprintf("fileSystem%s", codegen.Goify(s.FilePath, true))),
		}
		sd.FileServers = append(sd.FileServers, data)
	}

	for _, httpEndpoint := range httpSvc.HTTPEndpoints {
		method := svc.Method(httpEndpoint.MethodExpr.Name)

		var routes []*RouteData
		i := 0
		for _, r := range httpEndpoint.Routes {
			for _, rpath := range r.FullPaths() {
				params := expr.ExtractHTTPWildcards(rpath)
				var (
					init *InitData
				)
				{
					initArgs := make([]*InitArgData, len(params))
					pathParamsObj := expr.AsObject(httpEndpoint.PathParams().Type)
					suffix := ""
					if i > 0 {
						suffix = strconv.Itoa(i + 1)
					}
					i++
					name := fmt.Sprintf("%s%sPath%s", method.VarName, svc.StructName, suffix)
					for j, arg := range params {
						patt := pathParamsObj.Attribute(arg)
						att := makeHTTPType(patt)
						pointer := httpEndpoint.Params.IsPrimitivePointer(arg, true)
						if expr.IsObject(httpEndpoint.MethodExpr.Payload.Type) {
							// Path params may override requiredness, need to check payload.
							pointer = httpEndpoint.MethodExpr.Payload.IsPrimitivePointer(arg, true)
						}
						name := sd.Scope.Name(codegen.Goify(arg, false))
						var vcode string
						if att.Validation != nil {
							ctx := httpContext("", sd.Scope, true, false)
							vcode = codegen.AttributeValidationCode(att, nil, ctx, true, expr.IsAlias(att.Type), name, arg)
						}
						initArgs[j] = &InitArgData{
							Ref: name,
							AttributeData: &AttributeData{
								Name:        arg,
								VarName:     name,
								Description: att.Description,
								FieldName:   codegen.Goify(arg, true),
								FieldType:   patt.Type,
								TypeName:    sd.Scope.GoTypeName(att),
								TypeRef:     sd.Scope.GoTypeRef(att),
								Type:        att.Type,
								Pointer:     pointer,
								Required:    true,
								Example:     att.Example(expr.Root.API.ExampleGenerator),
								Validate:    vcode,
							},
						}
					}

					var buffer bytes.Buffer
					pf := expr.HTTPWildcardRegex.ReplaceAllString(rpath, "/%v")
					err := pathInitTmpl.Execute(&buffer, map[string]any{
						"Args":       initArgs,
						"PathParams": pathParamsObj,
						"PathFormat": pf,
					})
					if err != nil {
						panic(err)
					}
					init = &InitData{
						Name:           name,
						Description:    fmt.Sprintf("%s returns the URL path to the %s service %s HTTP endpoint. ", name, svc.Name, method.Name),
						ServerArgs:     initArgs,
						ClientArgs:     initArgs,
						ReturnTypeName: "string",
						ReturnTypeRef:  "string",
						ServerCode:     buffer.String(),
						ClientCode:     buffer.String(),
					}
				}

				routes = append(routes, &RouteData{
					Verb:     strings.ToUpper(r.Method),
					Path:     rpath,
					PathInit: init,
				})
			}
		}

		payload := buildPayloadData(httpEndpoint, sd)

		var (
			reqs  service.RequirementsData
			hsch  service.SchemesData
			bosch service.SchemesData
			qsch  service.SchemesData
			basch *service.SchemeData
		)
		for _, req := range httpEndpoint.Requirements {
			var rs service.SchemesData
			for _, sch := range req.Schemes {
				s := service.BuildSchemeData(sch, httpEndpoint.MethodExpr)
				rs = rs.Append(s)
				switch s.Type {
				case "Basic":
					basch = s
				default:
					switch s.In {
					case "query":
						qsch = qsch.Append(s)
					case "header":
						hsch = hsch.Append(s)
					default:
						bosch = bosch.Append(s)
					}
				}
			}
			reqs = append(reqs, &service.RequirementData{Schemes: rs, Scopes: req.Scopes})
		}

		var requestEncoder string
		if payload.Request.ClientBody != nil || len(payload.Request.Headers) > 0 || len(payload.Request.QueryParams) > 0 || len(payload.Request.Cookies) > 0 || basch != nil {
			requestEncoder = fmt.Sprintf("Encode%sRequest", method.VarName)
		}

		var requestInit *InitData
		var (
			name       string
			args       []*InitArgData
			payloadRef string
			pkg        string
		)
		{
			name = fmt.Sprintf("Build%sRequest", method.VarName)
			s := codegen.NewNameScope()
			s.Unique("c") // 'c' is reserved as the client's receiver name.
			for _, ca := range routes[0].PathInit.ClientArgs {
				if ca.FieldName != "" {
					ca.VarName = s.Unique(ca.VarName)
					ca.Ref = ca.VarName
					args = append(args, ca)
				}
			}
			pkg = pkgWithDefault(method.PayloadLoc, svc.PkgName)
			if len(routes[0].PathInit.ClientArgs) > 0 && httpEndpoint.MethodExpr.Payload.Type != expr.Empty {
				payloadRef = svc.Scope.GoFullTypeRef(httpEndpoint.MethodExpr.Payload, pkg)
			}
		}
		data := map[string]any{
			"PayloadRef":   payloadRef,
			"HasFields":    expr.IsObject(httpEndpoint.MethodExpr.Payload.Type),
			"ServiceName":  svc.Name,
			"EndpointName": method.Name,
			"Args":         args,
			"PathInit":     routes[0].PathInit,
			"Verb":         routes[0].Verb,
			"IsWebSocket":  httpEndpoint.MethodExpr.IsStreaming() && httpEndpoint.SSE == nil,
		}
		if httpEndpoint.SkipRequestBodyEncodeDecode {
			data["RequestStruct"] = pkg + "." + method.RequestStruct
		}
		var buf bytes.Buffer
		if err := requestInitTmpl.Execute(&buf, data); err != nil {
			panic(err) // bug
		}
		clientArgs := []*InitArgData{{Ref: "v", AttributeData: &AttributeData{Name: "payload", VarName: "v", TypeRef: "any"}}}
		requestInit = &InitData{
			Name:        name,
			Description: fmt.Sprintf("%s instantiates a HTTP request object with method and path set to call the %q service %q endpoint", name, svc.Name, method.Name),
			ClientCode:  buf.String(),
			ClientArgs:  clientArgs,
		}

		ed := &EndpointData{
			Method:          method,
			ServiceName:     svc.Name,
			ServiceVarName:  svc.VarName,
			ServicePkgName:  svc.PkgName,
			Payload:         payload,
			Result:          buildResultData(httpEndpoint, sd),
			Errors:          buildErrorsData(httpEndpoint, sd),
			HeaderSchemes:   hsch,
			BodySchemes:     bosch,
			QuerySchemes:    qsch,
			BasicScheme:     basch,
			Routes:          routes,
			MountHandler:    fmt.Sprintf("Mount%sHandler", method.VarName),
			HandlerInit:     fmt.Sprintf("New%sHandler", method.VarName),
			RequestDecoder:  fmt.Sprintf("Decode%sRequest", method.VarName),
			ResponseEncoder: fmt.Sprintf("Encode%sResponse", method.VarName),
			ErrorEncoder:    fmt.Sprintf("Encode%sError", method.VarName),
			ClientStruct:    "Client",
			EndpointInit:    method.VarName,
			RequestInit:     requestInit,
			RequestEncoder:  requestEncoder,
			ResponseDecoder: fmt.Sprintf("Decode%sResponse", method.VarName),
			Requirements:    reqs,
		}
		if httpEndpoint.MethodExpr.IsStreaming() {
			initWebSocketData(ed, httpEndpoint, sd)
			initSSEData(ed, httpEndpoint, sd)
		}

		if httpEndpoint.MultipartRequest {
			ed.MultipartRequestDecoder = &MultipartData{
				FuncName:    fmt.Sprintf("%s%sDecoderFunc", svc.StructName, method.VarName),
				InitName:    fmt.Sprintf("New%s%sDecoder", svc.StructName, method.VarName),
				VarName:     fmt.Sprintf("%s%sDecoderFn", svc.VarName, method.VarName),
				ServiceName: svc.Name,
				MethodName:  method.Name,
				Payload:     ed.Payload,
			}
			ed.MultipartRequestEncoder = &MultipartData{
				FuncName:    fmt.Sprintf("%s%sEncoderFunc", svc.StructName, method.VarName),
				InitName:    fmt.Sprintf("New%s%sEncoder", svc.StructName, method.VarName),
				VarName:     fmt.Sprintf("%s%sEncoderFn", svc.VarName, method.VarName),
				ServiceName: svc.Name,
				MethodName:  method.Name,
				Payload:     ed.Payload,
			}
		}

		if httpEndpoint.SkipRequestBodyEncodeDecode {
			ed.BuildStreamPayload = scope.Unique("Build" + codegen.Goify(method.Name, true) + "StreamPayload")
		}

		if httpEndpoint.Redirect != nil {
			ed.Redirect = &RedirectData{
				URL:        httpEndpoint.Redirect.URL,
				StatusCode: statusCodeToHTTPConst(httpEndpoint.Redirect.StatusCode),
			}
		}

		sd.Endpoints = append(sd.Endpoints, ed)
	}

	for _, a := range httpSvc.HTTPEndpoints {
		collectUserTypes(a.Body.Type, func(ut expr.UserType) {
			if d := attributeTypeData(ut, true, true, true, sd); d != nil {
				sd.ServerBodyAttributeTypes = append(sd.ServerBodyAttributeTypes, d)
			}
			if d := attributeTypeData(ut, true, false, false, sd); d != nil {
				sd.ClientBodyAttributeTypes = append(sd.ClientBodyAttributeTypes, d)
			}
		})

		if a.MethodExpr.StreamingPayload.Type != expr.Empty {
			collectUserTypes(a.StreamingBody.Type, func(ut expr.UserType) {
				if d := attributeTypeData(ut, true, true, true, sd); d != nil {
					sd.ServerBodyAttributeTypes = append(sd.ServerBodyAttributeTypes, d)
				}
				if d := attributeTypeData(ut, true, false, false, sd); d != nil {
					sd.ClientBodyAttributeTypes = append(sd.ClientBodyAttributeTypes, d)
				}
			})
		}

		if res := a.MethodExpr.Result; res != nil {
			for _, v := range a.Responses {
				collectUserTypes(v.Body.Type, func(ut expr.UserType) {
					// NOTE: ServerBodyAttributeTypes for response body types are
					// collected in buildResponseBodyType because we have to generate
					// body types for each view in a result type.
					if d := attributeTypeData(ut, false, true, false, sd); d != nil {
						sd.ClientBodyAttributeTypes = append(sd.ClientBodyAttributeTypes, d)
					}
				})
			}
		}

		for _, v := range a.HTTPErrors {
			collectUserTypes(v.Response.Body.Type, func(ut expr.UserType) {
				// NOTE: ServerBodyAttributeTypes for error response body types are
				// collected in buildResponseBodyType because we have to generate
				// body types for each view in a result type.
				if d := attributeTypeData(ut, false, true, false, sd); d != nil {
					sd.ClientBodyAttributeTypes = append(sd.ClientBodyAttributeTypes, d)
				}
			})
		}
	}

	return sd
}

// makeHTTPType traverses the attribute recursively and performs these actions:
//
// * removes aliased user type by replacing them with the underlying type.
// * changes unions into structs with Type and Value fields.
func makeHTTPType(att *expr.AttributeExpr) *expr.AttributeExpr {
	if att == nil {
		return nil
	}
	att = expr.DupAtt(att)
	return makeHTTPTypeRecursive(att, make(map[string]struct{}))
}

func makeHTTPTypeRecursive(att *expr.AttributeExpr, seen map[string]struct{}) *expr.AttributeExpr {
	switch dt := att.Type.(type) {
	case expr.UserType:
		if _, ok := dt.(*expr.ResultTypeExpr); !ok && !expr.IsObject(dt) {
			// Aliased user type. Use the underlying aliased type instead of
			// generating new types in the client and server packages
			att.Type = dt.Attribute().Type
			if v := dt.Attribute().Validation; v != nil {
				if att.Validation == nil {
					att.Validation = v
				} else {
					att.Validation.Merge(v)
				}
			}
			att.DefaultValue = dt.Attribute().DefaultValue
			att.UserExamples = dt.Attribute().UserExamples
		}
		if _, ok := seen[dt.ID()]; ok {
			return att
		}
		seen[dt.ID()] = struct{}{}
		dt.SetAttribute(makeHTTPTypeRecursive(dt.Attribute(), seen))
	case *expr.Array:
		dt.ElemType = makeHTTPTypeRecursive(dt.ElemType, seen)
	case *expr.Map:
		dt.KeyType = makeHTTPTypeRecursive(dt.KeyType, seen)
		dt.ElemType = makeHTTPTypeRecursive(dt.ElemType, seen)
	case *expr.Object:
		obj := make(expr.Object, len(*dt))
		for i, nat := range *dt {
			obj[i] = &expr.NamedAttributeExpr{Name: nat.Name, Attribute: makeHTTPTypeRecursive(nat.Attribute, seen)}
		}
		att.Type = &obj
	case *expr.Union:
		att = expr.UnionToObject(att)
	}
	return att
}

// buildPayloadData returns the data structure used to describe the endpoint
// payload including the HTTP request details. It also returns the user types
// used by the request body type recursively if any.
func buildPayloadData(e *expr.HTTPEndpointExpr, sd *ServiceData) *PayloadData {
	e.Body = makeHTTPType(e.Body)
	var (
		payload    = e.MethodExpr.Payload
		svc        = sd.Service
		body       = e.Body.Type
		ep         = svc.Method(e.MethodExpr.Name)
		httpsvrctx = httpContext("", sd.Scope, true, true)
		httpclictx = httpContext("", sd.Scope, true, false)
		pkg        = pkgWithDefault(ep.PayloadLoc, svc.PkgName)
		svcctx     = serviceContext(pkg, sd.Service.Scope)

		request       *RequestData
		mapQueryParam *ParamData
	)
	{
		var (
			serverBodyData = buildRequestBodyType(e.Body, payload, e, true, sd)
			clientBodyData = buildRequestBodyType(e.Body, payload, e, false, sd)
			paramsData     = extractPathParams(e.PathParams(), payload, sd.Scope)
			queryData      = extractQueryParams(e.QueryParams(), payload, sd.Scope)
			headersData    = extractHeaders(e.Headers, payload, svcctx, sd.Scope)
			cookiesData    = extractCookies(e.Cookies, payload, svcctx, sd.Scope)
			origin         string

			mustValidate bool
			mustHaveBody = true
		)
		if e.MapQueryParams != nil {
			var (
				fieldName string
				name      = "query"
				required  = true
				pAtt      = payload
			)
			if n := *e.MapQueryParams; n != "" {
				pAtt = expr.AsObject(payload.Type).Attribute(n)
				required = payload.IsRequired(n)
				name = n
				fieldName = codegen.Goify(name, true)
			}
			varn := codegen.Goify(name, false)
			mapQueryParam = &ParamData{
				MapQueryParams: e.MapQueryParams,
				Map:            expr.AsMap(payload.Type) != nil,
				Element: &Element{
					HTTPName: name,
					AttributeData: &AttributeData{
						Name:         name,
						VarName:      varn,
						FieldName:    fieldName,
						FieldType:    pAtt.Type,
						Required:     required,
						Type:         pAtt.Type,
						TypeName:     sd.Scope.GoTypeName(pAtt),
						TypeRef:      sd.Scope.GoTypeRef(pAtt),
						Validate:     codegen.AttributeValidationCode(pAtt, nil, httpsvrctx, required, expr.IsAlias(pAtt.Type), varn, name),
						DefaultValue: pAtt.DefaultValue,
						Example:      pAtt.Example(expr.Root.API.ExampleGenerator),
					},
				},
			}
			queryData = append(queryData, mapQueryParam)
		}
		if serverBodyData != nil {
			sd.ServerTypeNames[serverBodyData.Name] = false
			sd.ClientTypeNames[serverBodyData.Name] = false
		}
		for _, p := range cookiesData {
			if p.Required || p.Validate != "" || needConversion(p.Type) {
				mustValidate = true
				break
			}
		}
		if !mustValidate {
			for _, p := range paramsData {
				if p.Validate != "" || needConversion(p.Type) {
					mustValidate = true
					break
				}
			}
		}
		if !mustValidate {
			for _, q := range queryData {
				if q.Map || q.Validate != "" || q.Required || needConversion(q.Type) {
					mustValidate = true
					break
				}
			}
		}
		if !mustValidate {
			for _, h := range headersData {
				if h.Validate != "" || h.Required || needConversion(h.Type) {
					mustValidate = true
					break
				}
			}
		}
		if e.Body.Type != expr.Empty {
			// If design uses Body("name") syntax we need to use the
			// corresponding attribute in the result type for body
			// transformation.
			if o, ok := e.Body.Meta["origin:attribute"]; ok {
				origin = o[0]
				if !payload.IsRequired(o[0]) {
					mustHaveBody = false
				}
			}
		}
		request = &RequestData{
			PathParams:   paramsData,
			QueryParams:  queryData,
			Headers:      headersData,
			Cookies:      cookiesData,
			ServerBody:   serverBodyData,
			ClientBody:   clientBodyData,
			PayloadAttr:  codegen.Goify(origin, true),
			PayloadType:  e.MethodExpr.Payload.Type,
			MustHaveBody: mustHaveBody,
			MustValidate: mustValidate,
			Multipart:    e.MultipartRequest,
		}
	}

	var init *InitData
	if needInit(payload.Type) {
		// generate constructor function to transform request body,
		// params, headers and cookies into the method payload type
		var (
			name       string
			desc       string
			isObject   bool
			clientArgs []*InitArgData
			serverArgs []*InitArgData
		)
		n := codegen.Goify(ep.Name, true)
		p := codegen.Goify(ep.Payload, true)
		// Raw payload object has type name prefixed with endpoint name. No need to
		// prefix the type name again.
		if strings.HasPrefix(p, n) {
			p = svc.Scope.HashedUnique(payload.Type, p)
			name = fmt.Sprintf("New%s", p)
		} else {
			name = fmt.Sprintf("New%s%s", n, p)
		}
		desc = fmt.Sprintf("%s builds a %s service %s endpoint payload.",
			name, svc.Name, e.Name())
		isObject = expr.IsObject(payload.Type)
		if body != expr.Empty {
			var (
				svcode string
				cvcode string
			)
			if ut, ok := body.(expr.UserType); ok {
				if val := ut.Attribute().Validation; val != nil {
					svcode = codegen.ValidationCode(ut.Attribute(), ut, httpsvrctx, true, expr.IsAlias(ut), false, "body")
					cvcode = codegen.ValidationCode(ut.Attribute(), ut, httpclictx, true, expr.IsAlias(ut), false, "body")
				}
			}
			serverArgs = []*InitArgData{{
				Ref: sd.Scope.GoVar("body", body),
				AttributeData: &AttributeData{
					Name:     "body",
					VarName:  "body",
					TypeName: sd.Scope.GoTypeName(e.Body),
					TypeRef:  sd.Scope.GoTypeRef(e.Body),
					Type:     body,
					Required: true,
					Example:  e.Body.Example(expr.Root.API.ExampleGenerator),
					Validate: svcode,
				},
			}}
			clientArgs = []*InitArgData{{
				Ref: sd.Scope.GoVar("body", body),
				AttributeData: &AttributeData{
					Name:     "body",
					VarName:  "body",
					TypeName: sd.Scope.GoTypeNameWithDefaults(e.Body),
					TypeRef:  sd.Scope.GoTypeRefWithDefaults(e.Body),
					Type:     body,
					Required: true,
					Example:  e.Body.Example(expr.Root.API.ExampleGenerator),
					Validate: cvcode,
				},
			}}
		}
		var args []*InitArgData
		for _, p := range request.PathParams {
			args = append(args, &InitArgData{
				Ref: p.VarName,
				AttributeData: &AttributeData{
					Name:         p.Name,
					VarName:      p.VarName,
					Description:  p.Description,
					FieldName:    p.FieldName,
					FieldPointer: p.FieldPointer,
					FieldType:    p.FieldType,
					TypeName:     p.TypeName,
					TypeRef:      p.TypeRef,
					Type:         p.Type,
					Pointer:      p.Pointer,
					Required:     p.Required,
					Validate:     p.Validate,
					Example:      p.Example,
				},
			})
		}
		for _, p := range request.QueryParams {
			args = append(args, &InitArgData{
				Ref: p.VarName,
				AttributeData: &AttributeData{
					Name:         p.Name,
					VarName:      p.VarName,
					FieldName:    p.FieldName,
					FieldPointer: p.FieldPointer,
					FieldType:    p.FieldType,
					TypeName:     p.TypeName,
					TypeRef:      p.TypeRef,
					Type:         p.Type,
					Pointer:      p.Pointer,
					Required:     p.Required,
					DefaultValue: p.DefaultValue,
					Validate:     p.Validate,
					Example:      p.Example,
				},
			})
		}
		for _, h := range request.Headers {
			args = append(args, &InitArgData{
				Ref: h.VarName,
				AttributeData: &AttributeData{
					Name:         h.Name,
					VarName:      h.VarName,
					FieldName:    h.FieldName,
					FieldPointer: h.FieldPointer,
					FieldType:    h.FieldType,
					TypeName:     h.TypeName,
					TypeRef:      h.TypeRef,
					Type:         h.Type,
					Pointer:      h.Pointer,
					Required:     h.Required,
					DefaultValue: h.DefaultValue,
					Validate:     h.Validate,
					Example:      h.Example,
				},
			})
		}
		for _, c := range request.Cookies {
			args = append(args, &InitArgData{
				Ref: c.VarName,
				AttributeData: &AttributeData{
					Name:         c.Name,
					VarName:      c.VarName,
					FieldName:    c.FieldName,
					FieldPointer: c.FieldPointer,
					FieldType:    c.FieldType,
					TypeName:     c.TypeName,
					TypeRef:      c.TypeRef,
					Type:         c.Type,
					Pointer:      c.Pointer,
					Required:     c.Required,
					DefaultValue: c.DefaultValue,
					Validate:     c.Validate,
					Example:      c.Example,
				},
			})
		}
		serverArgs = append(serverArgs, args...)
		clientArgs = append(clientArgs, args...)

		var (
			cliArgs []*InitArgData
		)
		for _, r := range ep.Requirements {
			done := false
			for _, sc := range r.Schemes {
				if sc.Type == "Basic" {
					uatt := e.MethodExpr.Payload.Find(sc.UsernameAttr)
					uref := svc.Scope.GoTypeRef(uatt)
					if sc.UsernamePointer {
						uref = "*" + uref
					}
					uarg := &InitArgData{
						Ref: sc.UsernameAttr,
						AttributeData: &AttributeData{
							Name:         sc.UsernameAttr,
							VarName:      sc.UsernameAttr,
							FieldName:    sc.UsernameField,
							FieldPointer: sc.UsernamePointer,
							FieldType:    uatt.Type,
							Description:  uatt.Description,
							Required:     sc.UsernameRequired,
							TypeName:     svc.Scope.GoTypeName(uatt),
							TypeRef:      uref,
							Type:         uatt.Type,
							Pointer:      sc.UsernamePointer,
							Validate:     codegen.ValidationCode(uatt, nil, httpsvrctx, sc.UsernameRequired, expr.IsAlias(uatt.Type), false, sc.UsernameAttr),
							Example:      uatt.Example(expr.Root.API.ExampleGenerator),
						},
					}
					patt := e.MethodExpr.Payload.Find(sc.PasswordAttr)
					pref := svc.Scope.GoTypeRef(patt)
					if sc.PasswordPointer {
						pref = "*" + pref
					}
					parg := &InitArgData{
						Ref: sc.PasswordAttr,
						AttributeData: &AttributeData{
							Name:         sc.PasswordAttr,
							VarName:      sc.PasswordAttr,
							FieldName:    sc.PasswordField,
							FieldPointer: sc.PasswordPointer,
							FieldType:    patt.Type,
							Description:  patt.Description,
							Required:     sc.PasswordRequired,
							TypeName:     svc.Scope.GoTypeName(patt),
							TypeRef:      pref,
							Type:         patt.Type,
							Pointer:      sc.PasswordPointer,
							Validate:     codegen.ValidationCode(patt, nil, httpsvrctx, sc.PasswordRequired, expr.IsAlias(patt.Type), false, sc.PasswordAttr),
							Example:      patt.Example(expr.Root.API.ExampleGenerator),
						},
					}
					cliArgs = []*InitArgData{uarg, parg}
					done = true
					break
				}
			}
			if done {
				break
			}
		}

		var (
			serverCode string
			clientCode string
			err        error
			origin     string
			pointer    bool

			pAtt = payload
		)
		if body != expr.Empty {
			// If design uses Body("name") syntax then need to use payload
			// attribute to transform.
			if o, ok := e.Body.Meta["origin:attribute"]; ok {
				origin = o[0]
				pAtt = expr.AsObject(payload.Type).Attribute(origin)
				pointer = !payload.IsRequired(o[0]) && expr.IsPrimitive(pAtt.Type)
			}

			var (
				helpers []*codegen.TransformFunctionData
			)
			serverCode, helpers, err = unmarshal(e.Body, pAtt, "body", "v", httpsvrctx, svcctx)
			if err == nil {
				sd.ServerTransformHelpers = codegen.AppendHelpers(sd.ServerTransformHelpers, helpers)
			}
			// The client code for building the method payload from a request
			// body is used by the CLI tool to build the payload given to the
			// client endpoint. It differs because the body type there does not
			// use pointers for all fields (no need to validate).
			clientCode, helpers, err = marshal(e.Body, pAtt, "body", "v", httpclictx, svcctx)
			if err == nil {
				sd.ClientTransformHelpers = codegen.AppendHelpers(sd.ClientTransformHelpers, helpers)
			}
		} else if expr.IsArray(payload.Type) || expr.IsMap(payload.Type) {
			if params := expr.AsObject(e.Params.Type); len(*params) > 0 {
				var helpers []*codegen.TransformFunctionData
				serverCode, helpers, err = unmarshal((*params)[0].Attribute, payload, codegen.Goify((*params)[0].Name, false), "v", httpsvrctx, svcctx)
				if err == nil {
					sd.ServerTransformHelpers = codegen.AppendHelpers(sd.ServerTransformHelpers, helpers)
				}
				clientCode, helpers, err = marshal((*params)[0].Attribute, payload, codegen.Goify((*params)[0].Name, false), "v", httpclictx, svcctx)
				if err == nil {
					sd.ClientTransformHelpers = codegen.AppendHelpers(sd.ClientTransformHelpers, helpers)
				}
			}
		}
		if err != nil {
			fmt.Println(err.Error()) // TBD validate DSL so errors are not possible
		}
		init = &InitData{
			Name:                     name,
			Description:              desc,
			ServerArgs:               serverArgs,
			ClientArgs:               clientArgs,
			CLIArgs:                  cliArgs,
			ReturnTypeName:           svc.Scope.GoFullTypeName(payload, pkg),
			ReturnTypeRef:            svc.Scope.GoFullTypeRef(payload, pkg),
			ReturnIsStruct:           isObject,
			ReturnTypeAttribute:      codegen.Goify(origin, true),
			ReturnTypePkg:            pkg,
			ServerCode:               serverCode,
			ClientCode:               clientCode,
			ReturnIsPrimitivePointer: pointer,
		}
	}
	request.PayloadInit = init

	var (
		returnValue string
		name        string
		ref         string
	)
	if payload.Type != expr.Empty {
		name = svc.Scope.GoFullTypeName(payload, pkg)
		ref = svc.Scope.GoFullTypeRef(payload, pkg)
	}
	if init == nil {
		if o := expr.AsObject(e.Params.Type); o != nil && len(*o) > 0 {
			returnValue = codegen.Goify((*o)[0].Name, false)
		} else if o := expr.AsObject(e.Headers.Type); o != nil && len(*o) > 0 {
			returnValue = codegen.Goify((*o)[0].Name, false)
		} else if o := expr.AsObject(e.Cookies.Type); o != nil && len(*o) > 0 {
			returnValue = codegen.Goify((*o)[0].Name, false)
		} else if e.MapQueryParams != nil && *e.MapQueryParams == "" {
			returnValue = mapQueryParam.VarName
		}
	}

	return &PayloadData{
		Name:               name,
		Ref:                ref,
		Request:            request,
		DecoderReturnValue: returnValue,
	}
}

// buildResultData builds the result data for the given service endpoint.
func buildResultData(e *expr.HTTPEndpointExpr, sd *ServiceData) *ResultData {
	var (
		svc    = sd.Service
		ep     = svc.Method(e.MethodExpr.Name)
		pkg    = pkgWithDefault(ep.ResultLoc, svc.PkgName)
		result = e.MethodExpr.Result

		name string
		ref  string
		view string
	)
	view = expr.DefaultView
	if v, ok := result.Meta.Last(expr.ViewMetaKey); ok {
		view = v
	}
	if result.Type != expr.Empty {
		name = svc.Scope.GoFullTypeName(result, pkg)
		ref = svc.Scope.GoFullTypeRef(result, pkg)
	}

	var (
		mustInit  bool
		responses []*ResponseData
	)
	{
		viewed := false
		if ep.ViewedResult != nil {
			result = expr.AsObject(ep.ViewedResult.Type).Attribute("projected")
			viewed = true
		}
		responses = buildResponses(e, result, viewed, sd)
		for _, r := range responses {
			// response has a body, headers, cookies or tag
			if len(r.ServerBody) > 0 || len(r.Headers) > 0 || len(r.Cookies) > 0 || r.TagName != "" {
				mustInit = true
			}
		}
	}
	return &ResultData{
		IsStruct:  expr.IsObject(result.Type),
		Name:      name,
		Ref:       ref,
		Responses: responses,
		View:      view,
		MustInit:  mustInit,
	}
}

// buildResponses builds the response data for all the responses in the endpoint
// expression. The response headers, cookies and body for each response are
// inferred from the method's result expression if not specified explicitly.
//
// viewed parameter indicates if the method result uses views.
func buildResponses(e *expr.HTTPEndpointExpr, result *expr.AttributeExpr, viewed bool, sd *ServiceData) []*ResponseData {
	var (
		responses []*ResponseData

		svc        = sd.Service
		md         = svc.Method(e.Name())
		pkg        = pkgWithDefault(md.ResultLoc, svc.PkgName)
		httpclictx = httpContext("", sd.Scope, false, false)
		scope      = svc.Scope
		svcctx     = serviceContext(pkg, sd.Service.Scope)
	)
	{
		if viewed {
			scope = svc.ViewScope
			svcctx = viewContext(sd.Service.ViewsPkg, sd.Service.ViewScope)
		}
		notag := -1
		for i, resp := range e.Responses {
			resp.Body = expr.DupAtt(resp.Body)
			resp.Body = makeHTTPType(resp.Body)
			if resp.Tag[0] == "" {
				if notag > -1 {
					continue // we don't want more than one response with no tag
				}
				notag = i
			}
			var (
				headersData    []*HeaderData
				cookiesData    []*CookieData
				serverBodyData []*TypeData
				clientBodyData *TypeData
				init           *InitData
				origin         string
				mustValidate   bool

				resAttr = result
			)
			{
				headersData = extractHeaders(resp.Headers, result, svcctx, scope)
				cookiesData = extractCookies(resp.Cookies, result, svcctx, scope)
				if resp.Body.Type != expr.Empty {
					// If design uses Body("name") syntax we need to use the
					// corresponding attribute in the result type for body
					// transformation.
					if o, ok := resp.Body.Meta["origin:attribute"]; ok {
						origin = o[0]
						resAttr = expr.AsObject(resAttr.Type).Attribute(origin)
					}
				}
				if viewed {
					vname := ""
					if origin != "" {
						// Response body is explicitly set to an attribute in the method
						// result type. No need to do any view-based projections server side.
						if sbd := buildResponseBodyType(resp.Body, result, md.ResultLoc, e, true, &vname, sd); sbd != nil {
							serverBodyData = append(serverBodyData, sbd)
						}
					} else if v, ok := e.MethodExpr.Result.Meta.Last(expr.ViewMetaKey); ok {
						// Design explicitly sets the view to render the result.
						// We generate only one server body type which will be rendered
						// using the specified view.
						if sbd := buildResponseBodyType(resp.Body, result, md.ResultLoc, e, true, &v, sd); sbd != nil {
							serverBodyData = append(serverBodyData, sbd)
						}
					} else {
						// If a method result uses views (i.e., a result type), we generate
						// one response body type per view defined in the result type. The
						// generated body type names are suffixed with the name of the view
						// (except for the "default" view). Constructors are also generated
						// to create a view-specific body type from the method result.
						// This makes it possible for the server side to return only the
						// attributes defined in the view in the response (NOTE: a required
						// attribute in the result type may not be present in all its views)
						for _, view := range md.ViewedResult.Views {
							if sbd := buildResponseBodyType(resp.Body, result, md.ResultLoc, e, true, &view.Name, sd); sbd != nil {
								serverBodyData = append(serverBodyData, sbd)
							}
						}
					}
					clientBodyData = buildResponseBodyType(resp.Body, result, md.ResultLoc, e, false, &vname, sd)
				} else {
					if sbd := buildResponseBodyType(resp.Body, result, md.ResultLoc, e, true, nil, sd); sbd != nil {
						serverBodyData = append(serverBodyData, sbd)
					}
					clientBodyData = buildResponseBodyType(resp.Body, result, md.ResultLoc, e, false, nil, sd)
				}
				if clientBodyData != nil {
					sd.ClientTypeNames[clientBodyData.Name] = false
				}
				for _, h := range headersData {
					if h.Validate != "" || h.Required || needConversion(h.Type) {
						mustValidate = true
						break
					}
				}
				for _, c := range cookiesData {
					if c.Validate != "" || c.Required || needConversion(c.Type) {
						mustValidate = true
						break
					}
				}
				if needInit(result.Type) {
					// generate constructor function to transform response body,
					// headers and cookies into the method result type
					var (
						name       string
						desc       string
						code       string
						tname      string
						tref       string
						err        error
						pointer    bool
						clientArgs []*InitArgData
						helpers    []*codegen.TransformFunctionData
					)
					{
						tname = svc.Scope.GoFullTypeName(result, pkg)
						tref = svc.Scope.GoFullTypeRef(result, pkg)
						if viewed {
							tname = svc.ViewScope.GoFullTypeName(result, svc.ViewsPkg)
							tref = svc.ViewScope.GoFullTypeRef(result, svc.ViewsPkg)
						}
						status := codegen.Goify(http.StatusText(resp.StatusCode), true)
						n := codegen.Goify(md.Name, true)
						r := codegen.Goify(md.Result, true)
						// Raw result object has type name prefixed with endpoint name. No need to
						// prefix the type name again.
						if strings.HasPrefix(r, n) {
							r = scope.HashedUnique(result.Type, r)
							name = fmt.Sprintf("New%s%s", r, status)
						} else {
							name = fmt.Sprintf("New%s%s%s", n, r, status)
						}
						desc = fmt.Sprintf("%s builds a %q service %q endpoint result from a HTTP %q response.", name, svc.Name, e.Name(), status)
						if resp.Body.Type != expr.Empty {
							if origin != "" {
								pointer = result.IsPrimitivePointer(origin, true)
							}
							ref := "body"
							if expr.IsObject(resp.Body.Type) {
								ref = "&body"
								pointer = false
							}
							var vcode string
							if ut, ok := resp.Body.Type.(expr.UserType); ok {
								if val := ut.Attribute().Validation; val != nil {
									vcode = codegen.ValidationCode(ut.Attribute(), ut, httpclictx, true, expr.IsAlias(ut), false, "body")
								}
							}
							clientArgs = []*InitArgData{{
								Ref: ref,
								AttributeData: &AttributeData{
									Name:     "body",
									VarName:  "body",
									TypeRef:  sd.Scope.GoTypeRef(resp.Body),
									Validate: vcode,
								},
							}}
							// If the method result is a
							// * result type - we unmarshal the client response body to the
							//   corresponding type in the views package so that view-specific
							//   validation logic can be applied.
							// * user type - we unmarshal the client response body to the
							//   corresponding type in the service package after validating the
							//   response body. Here, the transformation code must
							//   rely on the fact that the required attributes are
							//   set in the response body (otherwise validation
							//   would fail).
							code, helpers, err = unmarshal(resp.Body, resAttr, "body", "v", httpclictx, svcctx)
							if err == nil {
								sd.ClientTransformHelpers = codegen.AppendHelpers(sd.ClientTransformHelpers, helpers)
							}
						} else if expr.IsArray(result.Type) || expr.IsMap(result.Type) {
							if params := expr.AsObject(e.QueryParams().Type); len(*params) > 0 {
								code, helpers, err = unmarshal((*params)[0].Attribute, result, codegen.Goify((*params)[0].Name, false), "v", httpclictx, svcctx)
								if err == nil {
									sd.ClientTransformHelpers = codegen.AppendHelpers(sd.ClientTransformHelpers, helpers)
								}
							}
						}
						if err != nil {
							fmt.Println(err.Error()) // TBD validate DSL so errors are not possible
						}
						for _, h := range headersData {
							clientArgs = append(clientArgs, &InitArgData{
								Ref: h.VarName,
								AttributeData: &AttributeData{
									Name:         h.Name,
									VarName:      h.VarName,
									FieldName:    h.FieldName,
									FieldPointer: h.FieldPointer,
									FieldType:    h.FieldType,
									Required:     h.Required,
									Pointer:      h.Pointer,
									TypeRef:      h.TypeRef,
									Type:         h.Type,
									Validate:     h.Validate,
									Example:      h.Example,
								},
							})
						}
						for _, c := range cookiesData {
							clientArgs = append(clientArgs, &InitArgData{
								Ref: c.VarName,
								AttributeData: &AttributeData{
									Name:         c.Name,
									VarName:      c.VarName,
									FieldName:    c.FieldName,
									FieldPointer: c.FieldPointer,
									FieldType:    c.FieldType,
									Required:     c.Required,
									Pointer:      c.Pointer,
									TypeRef:      c.TypeRef,
									Type:         c.Type,
									Validate:     c.Validate,
									Example:      c.Example,
								},
							})
						}
					}
					init = &InitData{
						Name:                     name,
						Description:              desc,
						ClientArgs:               clientArgs,
						ReturnTypeName:           tname,
						ReturnTypeRef:            tref,
						ReturnIsStruct:           expr.IsObject(result.Type),
						ReturnTypeAttribute:      codegen.Goify(origin, true),
						ReturnTypePkg:            pkg,
						ReturnIsPrimitivePointer: pointer,
						ClientCode:               code,
					}
				}

				var (
					tagName string
					tagVal  string
					tagPtr  bool
				)
				if resp.Tag[0] != "" {
					tagName = codegen.Goify(resp.Tag[0], true)
					tagVal = resp.Tag[1]
					tagPtr = viewed || result.IsPrimitivePointer(resp.Tag[0], true)
				}
				responses = append(responses, &ResponseData{
					StatusCode:   statusCodeToHTTPConst(resp.StatusCode),
					Description:  resp.Description,
					Headers:      headersData,
					Cookies:      cookiesData,
					ContentType:  resp.ContentType,
					ServerBody:   serverBodyData,
					ClientBody:   clientBodyData,
					ResultInit:   init,
					TagName:      tagName,
					TagValue:     tagVal,
					TagPointer:   tagPtr,
					MustValidate: mustValidate,
					ResultAttr:   codegen.Goify(origin, true),
					ViewedResult: md.ViewedResult,
				})
			}
		}
		count := len(responses)
		if notag >= 0 && notag < count-1 {
			// Make sure tagless response is last
			responses[notag], responses[count-1] = responses[count-1], responses[notag]
		}
	}
	return responses
}

// buildErrorsData builds the error data for all the error responses in the
// endpoint expression. The response headers, cookies and body for each response
// are inferred from the method's error expression if not specified explicitly.
func buildErrorsData(e *expr.HTTPEndpointExpr, sd *ServiceData) []*ErrorGroupData {
	var (
		svc        = sd.Service
		ep         = svc.Method(e.MethodExpr.Name)
		httpclictx = httpContext("", sd.Scope, false, false)
	)

	data := make(map[string][]*ErrorData)
	for _, v := range e.HTTPErrors {
		v.Response.Body = makeHTTPType(v.Response.Body)
		var (
			init *InitData
			body = v.Response.Body.Type
		)

		pkg := pkgWithDefault(ep.ErrorLocs[v.Name], svc.PkgName)
		errctx := serviceContext(pkg, sd.Service.Scope)

		if needInit(v.Type) {
			var (
				name     string
				desc     string
				isObject bool
				args     []*InitArgData
			)
			name = fmt.Sprintf("New%s%s", codegen.Goify(ep.Name, true), codegen.Goify(v.ErrorExpr.Name, true))
			desc = fmt.Sprintf("%s builds a %s service %s endpoint %s error.",
				name, svc.Name, e.Name(), v.ErrorExpr.Name)
			if body != expr.Empty {
				isObject = expr.IsObject(body)
				ref := "body"
				if isObject {
					ref = "&body"
				}
				args = []*InitArgData{{
					Ref:           ref,
					AttributeData: &AttributeData{Name: "body", VarName: "body", TypeRef: sd.Scope.GoTypeRef(v.Response.Body)},
				}}
			}
			for _, h := range extractHeaders(v.Response.Headers, v.AttributeExpr, errctx, sd.Scope) {
				args = append(args, &InitArgData{
					Ref: h.VarName,
					AttributeData: &AttributeData{
						Name:         h.Name,
						VarName:      h.VarName,
						FieldName:    h.FieldName,
						FieldPointer: false,
						FieldType:    h.FieldType,
						TypeRef:      h.TypeRef,
						Type:         h.Type,
						Validate:     h.Validate,
						Example:      h.Example,
					},
				})
			}
			for _, c := range extractCookies(v.Response.Cookies, v.AttributeExpr, errctx, sd.Scope) {
				args = append(args, &InitArgData{
					Ref: c.VarName,
					AttributeData: &AttributeData{
						Name:         c.Name,
						VarName:      c.VarName,
						FieldName:    c.FieldName,
						FieldPointer: false,
						FieldType:    c.FieldType,
						TypeRef:      c.TypeRef,
						Type:         c.Type,
						Validate:     c.Validate,
						Example:      c.Example,
					},
				})
			}

			var (
				code   string
				origin string
				err    error
			)
			if body != expr.Empty {
				eAtt := v.AttributeExpr
				// If design uses Body("name") syntax then need to use payload
				// attribute to transform.
				if o, ok := v.Response.Body.Meta["origin:attribute"]; ok {
					origin = o[0]
					eAtt = expr.AsObject(v.ErrorExpr.Type).Attribute(origin)
				}

				var helpers []*codegen.TransformFunctionData
				code, helpers, err = unmarshal(v.Response.Body, eAtt, "body", "v", httpclictx, errctx)
				if err == nil {
					sd.ClientTransformHelpers = codegen.AppendHelpers(sd.ClientTransformHelpers, helpers)
				}
			} else if expr.IsArray(v.Type) || expr.IsMap(v.Type) {
				if params := expr.AsObject(e.QueryParams().Type); len(*params) > 0 {
					var helpers []*codegen.TransformFunctionData
					code, helpers, err = unmarshal((*params)[0].Attribute, v.AttributeExpr, codegen.Goify((*params)[0].Name, false), "v", httpclictx, errctx)
					if err == nil {
						sd.ClientTransformHelpers = codegen.AppendHelpers(sd.ClientTransformHelpers, helpers)
					}
				}
			}
			if err != nil {
				fmt.Println(err.Error()) // TBD validate DSL so errors are not possible
			}

			init = &InitData{
				Name:                name,
				Description:         desc,
				ClientArgs:          args,
				ReturnTypeName:      svc.Scope.GoFullTypeName(v.AttributeExpr, pkg),
				ReturnTypeRef:       svc.Scope.GoFullTypeRef(v.AttributeExpr, pkg),
				ReturnIsStruct:      expr.IsObject(v.Type),
				ReturnTypeAttribute: codegen.Goify(origin, true),
				ReturnTypePkg:       pkg,
				ClientCode:          code,
			}
		}

		var (
			responseData *ResponseData
		)
		{
			var (
				serverBodyData []*TypeData
				clientBodyData *TypeData
			)
			{
				errorLoc := ep.ErrorLocs[v.ErrorExpr.Name]
				if sbd := buildResponseBodyType(v.Response.Body, v.AttributeExpr, errorLoc, e, true, nil, sd); sbd != nil {
					serverBodyData = append(serverBodyData, sbd)
				}
				clientBodyData = buildResponseBodyType(v.Response.Body, v.AttributeExpr, errorLoc, e, false, nil, sd)
				if clientBodyData != nil {
					sd.ClientTypeNames[clientBodyData.Name] = false
					clientBodyData.Description = fmt.Sprintf("%s is the type of the %q service %q endpoint HTTP response body for the %q error.",
						clientBodyData.VarName, svc.Name, e.Name(), v.Name)
					serverBodyData[0].Description = fmt.Sprintf("%s is the type of the %q service %q endpoint HTTP response body for the %q error.",
						serverBodyData[0].VarName, svc.Name, e.Name(), v.Name)
				}
			}

			headers := extractHeaders(v.Response.Headers, v.AttributeExpr, errctx, sd.Scope)
			cookies := extractCookies(v.Response.Cookies, v.AttributeExpr, errctx, sd.Scope)
			var mustValidate bool
			for _, h := range headers {
				if h.Validate != "" || h.Required || needConversion(h.Type) {
					mustValidate = true
					break
				}
			}
			for _, c := range cookies {
				if c.Validate != "" || c.Required || needConversion(c.Type) {
					mustValidate = true
					break
				}
			}
			var contentType string
			if v.Response.ContentType != expr.ErrorResultIdentifier {
				contentType = v.Response.ContentType
			}
			responseData = &ResponseData{
				StatusCode:   statusCodeToHTTPConst(v.Response.StatusCode),
				Headers:      headers,
				ContentType:  contentType,
				Cookies:      cookies,
				ErrorHeader:  v.Name,
				ServerBody:   serverBodyData,
				ClientBody:   clientBodyData,
				ResultInit:   init,
				MustValidate: mustValidate,
			}
		}

		ref := svc.Scope.GoFullTypeRef(v.AttributeExpr, pkg)
		data[ref] = append(data[ref], &ErrorData{
			Name:     v.Name,
			Response: responseData,
			Ref:      ref,
		})
	}
	keys := make([]string, len(data))
	i := 0
	for k := range data {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	var vals []*ErrorGroupData
	for _, k := range keys {
		es := data[k]
		for _, e := range es {
			found := false
			for _, eg := range vals {
				if eg.StatusCode == e.Response.StatusCode {
					eg.Errors = append(eg.Errors, e)
					found = true
					break
				}
			}
			if !found {
				vals = append(vals,
					&ErrorGroupData{
						StatusCode: e.Response.StatusCode,
						Errors:     []*ErrorData{e},
					})
			}
		}
	}
	return vals
}

// buildRequestBodyType builds the TypeData for a request body. The data makes
// it possible to generate a function on the client side that creates the body
// from the service method payload.
//
// body is the HTTP request body
//
// att is the payload attribute
//
// e is the HTTP endpoint expression
//
// svr is true if the function is generated for server side code.
//
// sd is the service data
func buildRequestBodyType(body, att *expr.AttributeExpr, e *expr.HTTPEndpointExpr, svr bool, sd *ServiceData) *TypeData {
	if body.Type == expr.Empty {
		return nil
	}
	var (
		name        string
		varname     string
		desc        string
		def         string
		ref         string
		validateDef string
		validateRef string

		svc     = sd.Service
		httpctx = httpContext("", sd.Scope, true, svr)
		ep      = sd.Service.Method(e.Name())
		pkg     = pkgWithDefault(ep.PayloadLoc, sd.Service.PkgName)
		svcctx  = serviceContext(pkg, sd.Service.Scope)
	)
	name = body.Type.Name()
	ref = sd.Scope.GoTypeRef(body)

	AddMarshalTags(body, make(map[string]struct{}))

	if ut, ok := body.Type.(expr.UserType); ok {
		varname = codegen.Goify(ut.Name(), true)
		def = goTypeDef(sd.Scope, ut.Attribute(), svr, !svr)
		desc = fmt.Sprintf("%s is the type of the %q service %q endpoint HTTP request body.",
			varname, svc.Name, e.Name())
		if svr {
			// generate validation code for unmarshaled type (server-side).
			validateDef = codegen.ValidationCode(ut.Attribute(), ut, httpctx, true, expr.IsAlias(ut), false, "body")
			if validateDef != "" {
				validateRef = fmt.Sprintf("err = Validate%s(&body)", varname)
			}
		}
	} else {
		// Generate validation code first because inline struct validation is removed.
		ctx := codegen.NewAttributeContext(!expr.IsPrimitive(body.Type), false, !svr, "", sd.Scope)
		validateRef = codegen.ValidationCode(body, nil, ctx, true, expr.IsAlias(body.Type), false, "body")
		if svr && expr.IsObject(body.Type) {
			// Body is an explicit object described in the design and in
			// this case the GoTypeRef is an inline struct definition. We
			// want to force all attributes to be pointers because we are
			// generating the server body type pre-validation.
			body.Validation = nil
		}
		varname = sd.Scope.GoTypeRef(body)
		desc = body.Description
	}
	var init *InitData
	if !svr && att.Type != expr.Empty && needInit(body.Type) {
		var (
			name    string
			desc    string
			code    string
			origin  string
			err     error
			helpers []*codegen.TransformFunctionData

			sourceVar = "p"
			svc       = sd.Service
		)
		{
			name = fmt.Sprintf("New%s", codegen.Goify(sd.Scope.GoTypeName(body), true))
			desc = fmt.Sprintf("%s builds the HTTP request body from the payload of the %q endpoint of the %q service.",
				name, e.Name(), svc.Name)
			src := sourceVar
			srcAtt := att
			// If design uses Body("name") syntax then need to use payload attribute
			// to transform.
			if o, ok := body.Meta["origin:attribute"]; ok {
				srcObj := expr.AsObject(att.Type)
				origin = o[0]
				srcAtt = srcObj.Attribute(origin)
				src += "." + codegen.Goify(origin, true)
			}
			code, helpers, err = marshal(srcAtt, body, src, "body", svcctx, httpctx)
			if err != nil {
				fmt.Println(err.Error()) // TBD validate DSL so errors are not possible
			}
			sd.ClientTransformHelpers = codegen.AppendHelpers(sd.ClientTransformHelpers, helpers)
		}
		arg := InitArgData{
			Ref: sourceVar,
			AttributeData: &AttributeData{
				Name:     "payload",
				VarName:  sourceVar,
				TypeRef:  svc.Scope.GoFullTypeRef(att, pkg),
				Type:     att.Type,
				Validate: validateDef,
				Example:  att.Example(expr.Root.API.ExampleGenerator),
			},
		}
		init = &InitData{
			Name:                name,
			Description:         desc,
			ReturnTypeRef:       sd.Scope.GoTypeRef(body),
			ReturnTypeAttribute: codegen.Goify(origin, true),
			ClientCode:          code,
			ClientArgs:          []*InitArgData{&arg},
		}
	}
	return &TypeData{
		Name:        name,
		VarName:     varname,
		Description: desc,
		Def:         def,
		Ref:         ref,
		Init:        init,
		ValidateDef: validateDef,
		ValidateRef: validateRef,
		Example:     body.Example(expr.Root.API.ExampleGenerator),
	}
}

// buildResponseBodyType builds the TypeData for a response body. The data
// makes it possible to generate a function that creates the server response
// body from the service method result/projected result or error.
//
// body is the response (success or error) HTTP body.
//
// att is the result/projected attribute.
//
// svr is true if the function is generated for server side code
//
// view is the view name to add as a suffix to the type name.
func buildResponseBodyType(body, att *expr.AttributeExpr, loc *codegen.Location, e *expr.HTTPEndpointExpr, svr bool, view *string, sd *ServiceData) *TypeData {
	if body.Type == expr.Empty {
		return nil
	}
	var (
		name        string
		varname     string
		desc        string
		def         string
		ref         string
		validateDef string
		validateRef string
		viewName    string
		mustInit    bool

		svc     = sd.Service
		httpctx = httpContext("", sd.Scope, false, svr)
		pkg     = pkgWithDefault(loc, sd.Service.PkgName)
		svcctx  = serviceContext(pkg, sd.Service.Scope)
	)
	// For server code, we project the response body type if the type is a result
	// type and generate a type for each view in the result type. This makes it
	// possible to return only the attributes in the view in the server response.
	if svr && view != nil && *view != "" {
		viewName = *view
		body = expr.DupAtt(body)
		if rt, ok := body.Type.(*expr.ResultTypeExpr); ok {
			var err error
			rt, err = expr.Project(rt, *view)
			if err != nil {
				panic(err)
			}
			body.Type = rt
			sd.ServerTypeNames[rt.Name()] = false
		}
	}

	name = body.Type.Name()
	ref = sd.Scope.GoTypeRef(body)
	mustInit = att.Type != expr.Empty && needInit(body.Type)

	AddMarshalTags(body, make(map[string]struct{}))

	if ut, ok := body.Type.(expr.UserType); ok {
		// response body is a user type.
		varname = codegen.Goify(ut.Name(), true)
		def = goTypeDef(sd.Scope, ut.Attribute(), !svr, svr)
		desc = fmt.Sprintf("%s is the type of the %q service %q endpoint HTTP response body.",
			varname, svc.Name, e.Name())
		if !svr && view == nil {
			// generate validation code for unmarshaled type (client-side).
			validateDef = codegen.ValidationCode(body, ut, httpctx, true, expr.IsAlias(body.Type), false, "body")
			if validateDef != "" {
				target := "&body"
				if expr.IsArray(ut) {
					// result type collection
					target = "body"
				}
				validateRef = fmt.Sprintf("err = Validate%s(%s)", varname, target)
			}
		}
	} else if !expr.IsPrimitive(body.Type) && mustInit {
		// response body is an array or map type.
		name = codegen.Goify(e.Name(), true) + "ResponseBody"
		varname = name
		desc = fmt.Sprintf("%s is the type of the %q service %q endpoint HTTP response body.",
			varname, svc.Name, e.Name())
		def = goTypeDef(sd.Scope, body, !svr, svr)
		validateRef = codegen.ValidationCode(body, nil, httpctx, true, expr.IsAlias(body.Type), false, "body")
	} else {
		// response body is a primitive type. They are used as non-pointers when
		// encoding/decoding responses.
		httpctx = httpContext("", sd.Scope, false, true)
		validateRef = codegen.ValidationCode(body, nil, httpctx, true, expr.IsAlias(body.Type), false, "body")
		varname = sd.Scope.GoTypeRef(body)
		desc = body.Description
	}
	if svr {
		sd.ServerTypeNames[name] = false
		// We collect the server body types need to generate a response body type
		// here because the response body type would be different from the actual
		// type in the HTTPResponseExpr since we projected the body type above.
		// For client side, we don't have to generate a separate body type per
		// view. Hence the client types are collected in "analyze" function.
		collectUserTypes(body.Type, func(ut expr.UserType) {
			if d := attributeTypeData(ut, false, false, true, sd); d != nil {
				sd.ServerBodyAttributeTypes = append(sd.ServerBodyAttributeTypes, d)
			}
		})
	}

	var init *InitData
	if svr && mustInit {
		var (
			name    string
			desc    string
			rtref   string
			code    string
			origin  string
			err     error
			helpers []*codegen.TransformFunctionData

			sourceVar = "res"
			svc       = sd.Service
		)
		{
			var rtname string
			if _, ok := body.Type.(expr.UserType); !ok && !expr.IsPrimitive(body.Type) {
				rtname = codegen.Goify(e.Name(), true) + "ResponseBody"
				rtref = rtname
			} else {
				rtname = codegen.Goify(sd.Scope.GoTypeName(body), true)
				rtref = sd.Scope.GoTypeRef(body)
			}
			name = fmt.Sprintf("New%s", rtname)
			desc = fmt.Sprintf("%s builds the HTTP response body from the result of the %q endpoint of the %q service.",
				name, e.Name(), svc.Name)
			if view != nil {
				svcctx = viewContext(sd.Service.ViewsPkg, sd.Service.ViewScope)
			}
			src := sourceVar
			srcAtt := att
			// If design uses Body("name") syntax then need to use result attribute
			// to transform.
			if o, ok := body.Meta["origin:attribute"]; ok {
				srcObj := expr.AsObject(att.Type)
				origin = o[0]
				srcAtt = srcObj.Attribute(origin)
				src += "." + codegen.Goify(origin, true)
			}
			code, helpers, err = marshal(srcAtt, body, src, "body", svcctx, httpctx)
			if err != nil {
				panic(err) // bug
			}
			sd.ServerTransformHelpers = codegen.AppendHelpers(sd.ServerTransformHelpers, helpers)
		}
		ref := sourceVar
		if view != nil {
			ref += ".Projected"
		}
		tref := svc.Scope.GoFullTypeRef(att, pkg)
		if view != nil {
			tref = svc.ViewScope.GoFullTypeRef(att, svc.ViewsPkg)
		}
		arg := InitArgData{
			Ref: ref,
			AttributeData: &AttributeData{
				Name:     "result",
				VarName:  sourceVar,
				TypeRef:  tref,
				Type:     att.Type,
				Validate: validateDef,
				Example:  att.Example(expr.Root.API.ExampleGenerator),
			},
		}
		init = &InitData{
			Name:                name,
			Description:         desc,
			ReturnTypeRef:       rtref,
			ReturnTypeAttribute: codegen.Goify(origin, true),
			ServerCode:          code,
			ServerArgs:          []*InitArgData{&arg},
		}
	}
	return &TypeData{
		Name:        name,
		VarName:     varname,
		Description: desc,
		Def:         def,
		Ref:         ref,
		Init:        init,
		ValidateDef: validateDef,
		ValidateRef: validateRef,
		Example:     body.Example(expr.Root.API.ExampleGenerator),
		View:        viewName,
	}
}

func extractPathParams(a *expr.MappedAttributeExpr, service *expr.AttributeExpr, scope *codegen.NameScope) []*ParamData {
	var params []*ParamData
	codegen.WalkMappedAttr(a, func(name, elem string, _ bool, c *expr.AttributeExpr) error { // nolint: errcheck
		// The StringSlice field of ParamData must be false for aliased primitive types
		var stringSlice bool
		if arr := expr.AsArray(c.Type); arr != nil {
			stringSlice = arr.ElemType.Type.Kind() == expr.StringKind
		}

		c = makeHTTPType(c)
		var (
			varn = scope.Name(codegen.Goify(name, false))
			arr  = expr.AsArray(c.Type)
			ctx  = serviceContext("", scope)
			ft   = service.Type

			fptr bool
		)
		fieldName := codegen.GoifyAtt(c, name, true)
		if !expr.IsObject(service.Type) {
			fieldName = ""
		} else {
			fptr = service.IsPrimitivePointer(name, true)
			ft = service.Find(name).Type
		}
		params = append(params, &ParamData{
			Map:            false,
			MapStringSlice: false,
			Element: &Element{
				HTTPName:      elem,
				AttributeName: name,
				Slice:         arr != nil,
				StringSlice:   stringSlice,
				AttributeData: &AttributeData{
					Name:         name,
					Description:  c.Description,
					FieldName:    fieldName,
					FieldPointer: fptr,
					FieldType:    ft,
					VarName:      varn,
					Required:     true,
					Type:         c.Type,
					TypeName:     scope.GoTypeName(c),
					TypeRef:      scope.GoTypeRef(c),
					Pointer:      false,
					Validate:     codegen.AttributeValidationCode(c, nil, ctx, true, expr.IsAlias(c.Type), varn, name),
					DefaultValue: c.DefaultValue,
					Example:      c.Example(expr.Root.API.ExampleGenerator),
				},
			},
		})
		return nil
	})

	return params
}

func extractQueryParams(a *expr.MappedAttributeExpr, service *expr.AttributeExpr, scope *codegen.NameScope) []*ParamData {
	var params []*ParamData
	codegen.WalkMappedAttr(a, func(name, elem string, required bool, c *expr.AttributeExpr) error { // nolint: errcheck
		// The StringSlice field of ParamData must be false for aliased primitive types
		var stringSlice bool
		if arr := expr.AsArray(c.Type); arr != nil {
			stringSlice = arr.ElemType.Type.Kind() == expr.StringKind
		}

		c = makeHTTPType(c)
		var (
			varn    = scope.Name(codegen.Goify(name, false))
			arr     = expr.AsArray(c.Type)
			mp      = expr.AsMap(c.Type)
			typeRef = scope.GoTypeRef(c)
			ctx     = serviceContext("", scope)
			ft      = service.Type

			pointer bool
			fptr    bool
		)
		pointer = a.IsPrimitivePointer(name, true)
		if pointer {
			typeRef = "*" + typeRef
		}
		fieldName := codegen.GoifyAtt(c, name, true)
		if !expr.IsObject(service.Type) {
			fieldName = ""
		} else {
			fptr = service.IsPrimitivePointer(name, true)
			ft = service.Find(name).Type
		}
		params = append(params, &ParamData{
			Map: mp != nil,
			MapStringSlice: mp != nil &&
				mp.KeyType.Type.Kind() == expr.StringKind &&
				mp.ElemType.Type.Kind() == expr.ArrayKind &&
				expr.AsArray(mp.ElemType.Type).ElemType.Type.Kind() == expr.StringKind,
			Element: &Element{
				Slice:         arr != nil,
				StringSlice:   stringSlice,
				HTTPName:      elem,
				AttributeName: name,
				AttributeData: &AttributeData{
					Name:         name,
					Description:  c.Description,
					FieldName:    fieldName,
					FieldPointer: fptr,
					FieldType:    ft,
					VarName:      varn,
					Required:     required,
					Type:         c.Type,
					TypeName:     scope.GoTypeName(c),
					TypeRef:      typeRef,
					Pointer:      pointer,
					Validate:     codegen.AttributeValidationCode(c, nil, ctx, required, expr.IsAlias(c.Type), varn, name),
					DefaultValue: c.DefaultValue,
					Example:      c.Example(expr.Root.API.ExampleGenerator),
				},
			},
		})
		return nil
	})

	return params
}

func extractHeaders(a *expr.MappedAttributeExpr, svcAtt *expr.AttributeExpr, svcCtx *codegen.AttributeContext, scope *codegen.NameScope) []*HeaderData {
	var headers []*HeaderData
	codegen.WalkMappedAttr(a, func(name, elem string, required bool, _ *expr.AttributeExpr) error { // nolint: errcheck
		var attr *expr.AttributeExpr
		if attr = svcAtt.Find(name); attr == nil {
			attr = svcAtt
		}
		var hattr *expr.AttributeExpr
		var stringSlice bool
		// The StringSlice field of ParamData must be false for aliased primitive types
		if arr := expr.AsArray(attr.Type); arr != nil {
			stringSlice = arr.ElemType.Type.Kind() == expr.StringKind
		}

		hattr = makeHTTPType(attr)
		var (
			varn    = scope.Name(codegen.Goify(name, false))
			arr     = expr.AsArray(hattr.Type)
			typeRef = scope.GoTypeRef(hattr)
			ft      = attr.Type

			fieldName string
			pointer   bool
			fptr      bool
		)
		pointer = a.IsPrimitivePointer(name, true)
		if expr.IsObject(svcAtt.Type) {
			fieldName = codegen.GoifyAtt(attr, name, true)
			fptr = svcCtx.IsPrimitivePointer(name, svcAtt)
		}
		if pointer {
			typeRef = "*" + typeRef
		}
		headers = append(headers, &HeaderData{
			CanonicalName: http.CanonicalHeaderKey(elem),
			Element: &Element{
				HTTPName:      elem,
				Slice:         arr != nil,
				StringSlice:   stringSlice,
				AttributeName: name,
				AttributeData: &AttributeData{
					Name:         name,
					Description:  hattr.Description,
					FieldName:    fieldName,
					FieldPointer: fptr,
					FieldType:    ft,
					VarName:      varn,
					TypeName:     scope.GoTypeName(hattr),
					TypeRef:      typeRef,
					Required:     required,
					Pointer:      pointer,
					Type:         hattr.Type,
					Validate:     codegen.AttributeValidationCode(hattr, nil, svcCtx, required, expr.IsAlias(hattr.Type), varn, name),
					DefaultValue: hattr.DefaultValue,
					Example:      hattr.Example(expr.Root.API.ExampleGenerator),
				},
			},
		})
		return nil
	})
	return headers
}

func extractCookies(a *expr.MappedAttributeExpr, svcAtt *expr.AttributeExpr, svcCtx *codegen.AttributeContext, scope *codegen.NameScope) []*CookieData {
	var cookies []*CookieData
	codegen.WalkMappedAttr(a, func(name, elem string, required bool, _ *expr.AttributeExpr) error { // nolint: errcheck
		var hattr *expr.AttributeExpr
		if hattr = svcAtt.Find(name); hattr == nil {
			hattr = svcAtt
		}
		hattr = makeHTTPType(hattr)
		var (
			varn    = scope.Name(codegen.Goify(name, false))
			typeRef = scope.GoTypeRef(hattr)
			ft      = svcAtt.Type

			fieldName string
			pointer   bool
			fptr      bool
		)
		pointer = a.IsPrimitivePointer(name, true)
		if expr.IsObject(svcAtt.Type) {
			fieldName = codegen.GoifyAtt(hattr, name, true)
			fptr = svcCtx.IsPrimitivePointer(name, svcAtt)
			ft = svcAtt.Find(name).Type
		}
		if pointer {
			typeRef = "*" + typeRef
		}
		c := &CookieData{
			Element: &Element{
				HTTPName:      elem,
				AttributeName: name,
				AttributeData: &AttributeData{
					Name:         name,
					Description:  hattr.Description,
					FieldName:    fieldName,
					FieldPointer: fptr,
					FieldType:    ft,
					VarName:      varn,
					TypeName:     scope.GoTypeName(hattr),
					TypeRef:      typeRef,
					Required:     required,
					Pointer:      pointer,
					Type:         hattr.Type,
					Validate:     codegen.AttributeValidationCode(hattr, nil, svcCtx, required, expr.IsAlias(hattr.Type), varn, name),
					DefaultValue: hattr.DefaultValue,
					Example:      hattr.Example(expr.Root.API.ExampleGenerator),
				},
			},
		}
		for n, v := range a.Meta {
			switch n {
			case "cookie:max-age":
				c.MaxAge = v[0]
			case "cookie:path":
				c.Path = v[0]
			case "cookie:domain":
				c.Domain = v[0]
			case "cookie:secure":
				c.Secure = v[0] == "Secure"
			case "cookie:http-only":
				c.HTTPOnly = v[0] == "HttpOnly"
			case "cookie:same-site":
				switch v[0] {
				case string(expr.CookieSameSiteLax):
					c.SameSite = "http.SameSiteLaxMode"
				case string(expr.CookieSameSiteStrict):
					c.SameSite = "http.SameSiteStrictMode"
				case string(expr.CookieSameSiteNone):
					c.SameSite = "http.SameSiteNoneMode"
				case string(expr.CookieSameSiteDefault):
					c.SameSite = "http.SameSiteDefaultMode"
				}
			}
		}
		cookies = append(cookies, c)
		return nil
	})
	return cookies
}

// collectUserTypes traverses the given data type recursively and calls back the
// given function for each attribute using a user type.
func collectUserTypes(dt expr.DataType, cb func(expr.UserType), seen ...map[string]struct{}) {
	if dt == expr.Empty {
		return
	}
	var s map[string]struct{}
	if len(seen) > 0 {
		s = seen[0]
	} else {
		s = make(map[string]struct{})
	}
	switch actual := dt.(type) {
	case *expr.Object:
		for _, nat := range *actual {
			collectUserTypes(nat.Attribute.Type, cb, seen...)
		}
	case *expr.Array:
		collectUserTypes(actual.ElemType.Type, cb, seen...)
	case *expr.Map:
		collectUserTypes(actual.KeyType.Type, cb, seen...)
		collectUserTypes(actual.ElemType.Type, cb, seen...)
	case expr.UserType:
		if _, ok := s[actual.ID()]; ok {
			return
		}
		s[actual.ID()] = struct{}{}
		cb(actual)
		collectUserTypes(actual.Attribute().Type, cb, s)
	}
}

func attributeTypeData(ut expr.UserType, req, ptr, server bool, rd *ServiceData) *TypeData {
	if ut == expr.Empty {
		return nil
	}
	seen := rd.ServerTypeNames
	if !server {
		seen = rd.ClientTypeNames
	}
	if _, ok := seen[ut.Name()]; ok {
		return nil
	}
	seen[ut.Name()] = false

	var (
		name        string
		desc        string
		validate    string
		validateRef string

		att  = &expr.AttributeExpr{Type: ut}
		hctx = httpContext("", rd.Scope, req, server)
	)
	name = rd.Scope.GoTypeName(att)
	ctx := "request"
	if !req {
		ctx = "response"
	}
	desc = name + " is used to define fields on " + ctx + " body types."
	if (req || !req && !server) && !expr.IsAlias(ut) {
		// Generate validations for responses client-side and for
		// requests server-side and CLI.
		// Alias types are validated inline in the parent type
		validate = codegen.ValidationCode(ut.Attribute(), ut, hctx, true, expr.IsAlias(ut), false, "body")
	}
	if validate != "" {
		validateRef = fmt.Sprintf("err = Validate%s(v)", name)
	}
	return &TypeData{
		Name:        ut.Name(),
		VarName:     name,
		Description: desc,
		Def:         goTypeDef(rd.Scope, ut.Attribute(), ptr, hctx.UseDefault),
		Ref:         rd.Scope.GoTypeRef(att),
		ValidateDef: validate,
		ValidateRef: validateRef,
		Example:     att.Example(expr.Root.API.ExampleGenerator),
	}
}

// httpContext returns a context for attributes of types used to marshal and
// unmarshal HTTP requests and responses.
//
// pkg is the package name where the body type exists
//
// scope is the named scope
//
// request if true indicates that the type is a request type, else response
// type
//
// svr if true indicates that the type is a server type, else client type
func httpContext(pkg string, scope *codegen.NameScope, request, svr bool) *codegen.AttributeContext {
	marshal := !request && svr || request && !svr
	return codegen.NewAttributeContext(!marshal, false, marshal, pkg, scope)
}

// serviceContext returns an attribute context for service types.
func serviceContext(pkg string, scope *codegen.NameScope) *codegen.AttributeContext {
	return codegen.NewAttributeContext(false, false, true, pkg, scope)
}

// viewContext returns an attribute context for projected types.
func viewContext(pkg string, scope *codegen.NameScope) *codegen.AttributeContext {
	return codegen.NewAttributeContext(true, false, true, pkg, scope)
}

// pkgWithDefault returns the package name of the given location if not nil, def otherwise.
func pkgWithDefault(loc *codegen.Location, def string) string {
	if loc == nil {
		return def
	}
	return loc.PackageName()
}

// unmarshal initializes a data structure defined by target type from a data
// structure defined by source type. The attributes in the source data
// structure are pointers and the attributes in the target data structure that
// have default values are non-pointers. Fields in target type are initialized
// with their default values (if any).
//
// source, target are the attributes used in the transformation
//
// sourceVar, targetVar are the variable names for source and target used in
// the transformation code
//
// sourceCtx, targetCtx are the source and target attribute contexts
func unmarshal(source, target *expr.AttributeExpr, sourceVar, targetVar string, sourceCtx, targetCtx *codegen.AttributeContext) (string, []*codegen.TransformFunctionData, error) {
	return codegen.GoTransform(source, target, sourceVar, targetVar, sourceCtx, targetCtx, "unmarshal", true)
}

// marshal initializes a data structure defined by target type from a data
// structure defined by source type. The fields in the source and target
// data structure use non-pointers for attributes with default values.
//
// source, target are the attributes used in the transformation
//
// sourceVar, targetVar are the variable names for source and target used in
// the transformation code
//
// sourceCtx, targetCtx are the source and target attribute contexts
func marshal(source, target *expr.AttributeExpr, sourceVar, targetVar string, sourceCtx, targetCtx *codegen.AttributeContext) (string, []*codegen.TransformFunctionData, error) {
	return codegen.GoTransform(source, target, sourceVar, targetVar, sourceCtx, targetCtx, "marshal", true)
}

// needConversion returns true if the type needs to be converted from a string.
func needConversion(dt expr.DataType) bool {
	if dt == expr.Empty {
		return false
	}
	switch actual := dt.(type) {
	case expr.Primitive:
		if actual.Kind() == expr.StringKind ||
			actual.Kind() == expr.AnyKind ||
			actual.Kind() == expr.BytesKind {
			return false
		}
		return true
	case *expr.Array:
		return needConversion(actual.ElemType.Type)
	case *expr.Map:
		return needConversion(actual.KeyType.Type) ||
			needConversion(actual.ElemType.Type)
	default:
		return true
	}
}

// AddMarshalTags adds JSON, XML and Form tags to all inline object attributes recursively.
func AddMarshalTags(att *expr.AttributeExpr, seen map[string]struct{}) {
	if ut, ok := att.Type.(expr.UserType); ok {
		if _, ok := seen[ut.Hash()]; ok {
			return // avoid infinite recursions
		}
		seen[ut.Hash()] = struct{}{}
		if expr.IsObject(ut.Attribute().Type) {
			for _, att := range *(expr.AsObject(att.Type)) {
				AddMarshalTags(att.Attribute, seen)
			}
		}
		return
	}
	if expr.IsArray(att.Type) {
		AddMarshalTags(expr.AsArray(att.Type).ElemType, seen)
		return
	}
	if expr.IsMap(att.Type) {
		AddMarshalTags(expr.AsMap(att.Type).KeyType, seen)
		AddMarshalTags(expr.AsMap(att.Type).ElemType, seen)
		return
	}
	if !expr.IsObject(att.Type) {
		return
	}
	// inline object
	for _, natt := range *(expr.AsObject(att.Type)) {
		if natt.Attribute.Meta == nil {
			natt.Attribute.Meta = expr.MetaExpr{}
		}
		ns := []string{natt.Name}
		natt.Attribute.Meta["struct:tag:form"] = ns
		natt.Attribute.Meta["struct:tag:json"] = ns
		natt.Attribute.Meta["struct:tag:xml"] = ns
	}
}

// needInit returns true if and only if the given type is or makes use of user
// types.
func needInit(dt expr.DataType) bool {
	if dt == expr.Empty {
		return false
	}
	switch actual := dt.(type) {
	case expr.Primitive:
		return false
	case *expr.Array:
		return needInit(actual.ElemType.Type)
	case *expr.Map:
		return needInit(actual.KeyType.Type) ||
			needInit(actual.ElemType.Type)
	case *expr.Object:
		for _, nat := range *actual {
			if needInit(nat.Attribute.Type) {
				return true
			}
		}
		return false
	case expr.UserType:
		return true
	default:
		panic(fmt.Sprintf("unknown data type %T", actual)) // bug
	}
}

// upgradeParams returns the data required to render the websocket_upgrade
// template.
func upgradeParams(e *EndpointData, fn string) map[string]any {
	return map[string]any{
		"ViewedResult": e.Method.ViewedResult,
		"Function":     fn,
	}
}

// needDialer returns true if at least one method in the defined services
// uses WebSocket for sending payload or result.
func needDialer(data []*ServiceData) bool {
	return slices.ContainsFunc(data, hasWebSocket)
}
