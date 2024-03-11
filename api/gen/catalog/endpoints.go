// Code generated by goa v3.15.1, DO NOT EDIT.
//
// catalog endpoints
//
// Command:
// $ goa gen github.com/tektoncd/hub/api/design

package catalog

import (
	"context"

	goa "goa.design/goa/v3/pkg"
	"goa.design/goa/v3/security"
)

// Endpoints wraps the "catalog" service endpoints.
type Endpoints struct {
	Refresh      goa.Endpoint
	RefreshAll   goa.Endpoint
	CatalogError goa.Endpoint
}

// NewEndpoints wraps the methods of the "catalog" service with endpoints.
func NewEndpoints(s Service) *Endpoints {
	// Casting service to Auther interface
	a := s.(Auther)
	return &Endpoints{
		Refresh:      NewRefreshEndpoint(s, a.JWTAuth),
		RefreshAll:   NewRefreshAllEndpoint(s, a.JWTAuth),
		CatalogError: NewCatalogErrorEndpoint(s, a.JWTAuth),
	}
}

// Use applies the given middleware to all the "catalog" service endpoints.
func (e *Endpoints) Use(m func(goa.Endpoint) goa.Endpoint) {
	e.Refresh = m(e.Refresh)
	e.RefreshAll = m(e.RefreshAll)
	e.CatalogError = m(e.CatalogError)
}

// NewRefreshEndpoint returns an endpoint function that calls the method
// "Refresh" of service "catalog".
func NewRefreshEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req any) (any, error) {
		p := req.(*RefreshPayload)
		var err error
		sc := security.JWTScheme{
			Name:           "jwt",
			Scopes:         []string{"rating:read", "rating:write", "agent:create", "catalog:refresh", "config:refresh", "refresh:token"},
			RequiredScopes: []string{"catalog:refresh"},
		}
		ctx, err = authJWTFn(ctx, p.Token, &sc)
		if err != nil {
			return nil, err
		}
		res, err := s.Refresh(ctx, p)
		if err != nil {
			return nil, err
		}
		vres := NewViewedJob(res, "default")
		return vres, nil
	}
}

// NewRefreshAllEndpoint returns an endpoint function that calls the method
// "RefreshAll" of service "catalog".
func NewRefreshAllEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req any) (any, error) {
		p := req.(*RefreshAllPayload)
		var err error
		sc := security.JWTScheme{
			Name:           "jwt",
			Scopes:         []string{"rating:read", "rating:write", "agent:create", "catalog:refresh", "config:refresh", "refresh:token"},
			RequiredScopes: []string{"catalog:refresh"},
		}
		ctx, err = authJWTFn(ctx, p.Token, &sc)
		if err != nil {
			return nil, err
		}
		return s.RefreshAll(ctx, p)
	}
}

// NewCatalogErrorEndpoint returns an endpoint function that calls the method
// "CatalogError" of service "catalog".
func NewCatalogErrorEndpoint(s Service, authJWTFn security.AuthJWTFunc) goa.Endpoint {
	return func(ctx context.Context, req any) (any, error) {
		p := req.(*CatalogErrorPayload)
		var err error
		sc := security.JWTScheme{
			Name:           "jwt",
			Scopes:         []string{"rating:read", "rating:write", "agent:create", "catalog:refresh", "config:refresh", "refresh:token"},
			RequiredScopes: []string{"catalog:refresh"},
		}
		ctx, err = authJWTFn(ctx, p.Token, &sc)
		if err != nil {
			return nil, err
		}
		return s.CatalogError(ctx, p)
	}
}
