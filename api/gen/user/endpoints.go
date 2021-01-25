// Code generated by goa v3.2.2, DO NOT EDIT.
//
// user endpoints
//
// Command:
// $ goa gen github.com/tektoncd/hub/api/design

package user

import (
	"context"

	goa "goa.design/goa/v3/pkg"
	"goa.design/goa/v3/security"
)

// Endpoints wraps the "user" service endpoints.
type Endpoints struct {
	RefreshAccessToken goa.Endpoint
}

// NewEndpoints wraps the methods of the "user" service with endpoints.
func NewEndpoints(s Service) *Endpoints {
	// Casting service to Auther interface
	a := s.(Auther)
	return &Endpoints{
		RefreshAccessToken: NewRefreshAccessTokenEndpoint(s, a.JWTAuth),
	}
}

// Use applies the given middleware to all the "user" service endpoints.
func (e *Endpoints) Use(m func(goa.Endpoint) goa.Endpoint) {
	e.RefreshAccessToken = m(e.RefreshAccessToken)
}

// NewRefreshAccessTokenEndpoint returns an endpoint function that calls the
// method "RefreshAccessToken" of service "user".
func NewRefreshAccessTokenEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*RefreshAccessTokenPayload)
		var err error
		sc := security.JWTScheme{
			Name:           "jwt",
			Scopes:         []string{"rating:read", "rating:write", "agent:create", "catalog:refresh", "config:refresh", "refresh:token"},
			RequiredScopes: []string{"refresh:token"},
		}
		ctx, err = authJWTFn(ctx, p.RefreshToken, &sc)
		if err != nil {
			return nil, err
		}
		return s.RefreshAccessToken(ctx, p)
	}
}
