// Code generated by goa v3.7.5, DO NOT EDIT.
//
// resource HTTP server encoders and decoders
//
// Command:
// $ goa gen github.com/tektoncd/hub/api/design

package server

import (
	"context"
	"net/http"
	"strconv"

	resource "github.com/tektoncd/hub/api/gen/resource"
	resourceviews "github.com/tektoncd/hub/api/gen/resource/views"
	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
)

// EncodeQueryResponse returns an encoder for responses returned by the
// resource Query endpoint.
func EncodeQueryResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res, _ := v.(*resource.QueryResult)
		w.Header().Set("Location", res.Location)
		w.WriteHeader(http.StatusFound)
		return nil
	}
}

// DecodeQueryRequest returns a decoder for requests sent to the resource Query
// endpoint.
func DecodeQueryRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			name       string
			catalogs   []string
			categories []string
			kinds      []string
			tags       []string
			platforms  []string
			limit      uint
			match      string
			err        error
		)
		nameRaw := r.URL.Query().Get("name")
		if nameRaw != "" {
			name = nameRaw
		}
		catalogs = r.URL.Query()["catalogs"]
		categories = r.URL.Query()["categories"]
		kinds = r.URL.Query()["kinds"]
		tags = r.URL.Query()["tags"]
		platforms = r.URL.Query()["platforms"]
		{
			limitRaw := r.URL.Query().Get("limit")
			if limitRaw == "" {
				limit = 1000
			} else {
				v, err2 := strconv.ParseUint(limitRaw, 10, strconv.IntSize)
				if err2 != nil {
					err = goa.MergeErrors(err, goa.InvalidFieldTypeError("limit", limitRaw, "unsigned integer"))
				}
				limit = uint(v)
			}
		}
		matchRaw := r.URL.Query().Get("match")
		if matchRaw != "" {
			match = matchRaw
		} else {
			match = "contains"
		}
		if !(match == "exact" || match == "contains") {
			err = goa.MergeErrors(err, goa.InvalidEnumValueError("match", match, []interface{}{"exact", "contains"}))
		}
		if err != nil {
			return nil, err
		}
		payload := NewQueryPayload(name, catalogs, categories, kinds, tags, platforms, limit, match)

		return payload, nil
	}
}

// EncodeVersionsByIDResponse returns an encoder for responses returned by the
// resource VersionsByID endpoint.
func EncodeVersionsByIDResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res, _ := v.(*resource.VersionsByIDResult)
		w.Header().Set("Location", res.Location)
		w.WriteHeader(http.StatusFound)
		return nil
	}
}

// DecodeVersionsByIDRequest returns a decoder for requests sent to the
// resource VersionsByID endpoint.
func DecodeVersionsByIDRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			id  uint
			err error

			params = mux.Vars(r)
		)
		{
			idRaw := params["id"]
			v, err2 := strconv.ParseUint(idRaw, 10, strconv.IntSize)
			if err2 != nil {
				err = goa.MergeErrors(err, goa.InvalidFieldTypeError("id", idRaw, "unsigned integer"))
			}
			id = uint(v)
		}
		if err != nil {
			return nil, err
		}
		payload := NewVersionsByIDPayload(id)

		return payload, nil
	}
}

// EncodeByCatalogKindNameVersionResponse returns an encoder for responses
// returned by the resource ByCatalogKindNameVersion endpoint.
func EncodeByCatalogKindNameVersionResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res, _ := v.(*resource.ByCatalogKindNameVersionResult)
		w.Header().Set("Location", res.Location)
		w.WriteHeader(http.StatusFound)
		return nil
	}
}

// DecodeByCatalogKindNameVersionRequest returns a decoder for requests sent to
// the resource ByCatalogKindNameVersion endpoint.
func DecodeByCatalogKindNameVersionRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			catalog string
			kind    string
			name    string
			version string
			err     error

			params = mux.Vars(r)
		)
		catalog = params["catalog"]
		kind = params["kind"]
		if !(kind == "task" || kind == "pipeline") {
			err = goa.MergeErrors(err, goa.InvalidEnumValueError("kind", kind, []interface{}{"task", "pipeline"}))
		}
		name = params["name"]
		version = params["version"]
		if err != nil {
			return nil, err
		}
		payload := NewByCatalogKindNameVersionPayload(catalog, kind, name, version)

		return payload, nil
	}
}

// EncodeByVersionIDResponse returns an encoder for responses returned by the
// resource ByVersionId endpoint.
func EncodeByVersionIDResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res, _ := v.(*resource.ByVersionIDResult)
		w.Header().Set("Location", res.Location)
		w.WriteHeader(http.StatusFound)
		return nil
	}
}

// DecodeByVersionIDRequest returns a decoder for requests sent to the resource
// ByVersionId endpoint.
func DecodeByVersionIDRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			versionID uint
			err       error

			params = mux.Vars(r)
		)
		{
			versionIDRaw := params["versionID"]
			v, err2 := strconv.ParseUint(versionIDRaw, 10, strconv.IntSize)
			if err2 != nil {
				err = goa.MergeErrors(err, goa.InvalidFieldTypeError("versionID", versionIDRaw, "unsigned integer"))
			}
			versionID = uint(v)
		}
		if err != nil {
			return nil, err
		}
		payload := NewByVersionIDPayload(versionID)

		return payload, nil
	}
}

// EncodeByCatalogKindNameResponse returns an encoder for responses returned by
// the resource ByCatalogKindName endpoint.
func EncodeByCatalogKindNameResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res, _ := v.(*resource.ByCatalogKindNameResult)
		w.Header().Set("Location", res.Location)
		w.WriteHeader(http.StatusFound)
		return nil
	}
}

// DecodeByCatalogKindNameRequest returns a decoder for requests sent to the
// resource ByCatalogKindName endpoint.
func DecodeByCatalogKindNameRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			catalog          string
			kind             string
			name             string
			pipelinesversion *string
			err              error

			params = mux.Vars(r)
		)
		catalog = params["catalog"]
		kind = params["kind"]
		if !(kind == "task" || kind == "pipeline") {
			err = goa.MergeErrors(err, goa.InvalidEnumValueError("kind", kind, []interface{}{"task", "pipeline"}))
		}
		name = params["name"]
		pipelinesversionRaw := r.URL.Query().Get("pipelinesversion")
		if pipelinesversionRaw != "" {
			pipelinesversion = &pipelinesversionRaw
		}
		if pipelinesversion != nil {
			err = goa.MergeErrors(err, goa.ValidatePattern("pipelinesversion", *pipelinesversion, "^\\d+(?:\\.\\d+){0,2}$"))
		}
		if err != nil {
			return nil, err
		}
		payload := NewByCatalogKindNamePayload(catalog, kind, name, pipelinesversion)

		return payload, nil
	}
}

// EncodeByIDResponse returns an encoder for responses returned by the resource
// ById endpoint.
func EncodeByIDResponse(encoder func(context.Context, http.ResponseWriter) goahttp.Encoder) func(context.Context, http.ResponseWriter, interface{}) error {
	return func(ctx context.Context, w http.ResponseWriter, v interface{}) error {
		res, _ := v.(*resource.ByIDResult)
		w.Header().Set("Location", res.Location)
		w.WriteHeader(http.StatusFound)
		return nil
	}
}

// DecodeByIDRequest returns a decoder for requests sent to the resource ById
// endpoint.
func DecodeByIDRequest(mux goahttp.Muxer, decoder func(*http.Request) goahttp.Decoder) func(*http.Request) (interface{}, error) {
	return func(r *http.Request) (interface{}, error) {
		var (
			id  uint
			err error

			params = mux.Vars(r)
		)
		{
			idRaw := params["id"]
			v, err2 := strconv.ParseUint(idRaw, 10, strconv.IntSize)
			if err2 != nil {
				err = goa.MergeErrors(err, goa.InvalidFieldTypeError("id", idRaw, "unsigned integer"))
			}
			id = uint(v)
		}
		if err != nil {
			return nil, err
		}
		payload := NewByIDPayload(id)

		return payload, nil
	}
}

// marshalResourceviewsResourceDataViewToResourceDataResponseBodyWithoutVersion
// builds a value of type *ResourceDataResponseBodyWithoutVersion from a value
// of type *resourceviews.ResourceDataView.
func marshalResourceviewsResourceDataViewToResourceDataResponseBodyWithoutVersion(v *resourceviews.ResourceDataView) *ResourceDataResponseBodyWithoutVersion {
	res := &ResourceDataResponseBodyWithoutVersion{
		ID:         *v.ID,
		Name:       *v.Name,
		Kind:       *v.Kind,
		HubURLPath: *v.HubURLPath,
		Rating:     *v.Rating,
	}
	if v.Catalog != nil {
		res.Catalog = marshalResourceviewsCatalogViewToCatalogResponseBodyMin(v.Catalog)
	}
	if v.Categories != nil {
		res.Categories = make([]*CategoryResponseBody, len(v.Categories))
		for i, val := range v.Categories {
			res.Categories[i] = marshalResourceviewsCategoryViewToCategoryResponseBody(val)
		}
	}
	if v.LatestVersion != nil {
		res.LatestVersion = marshalResourceviewsResourceVersionDataViewToResourceVersionDataResponseBodyWithoutResource(v.LatestVersion)
	}
	if v.Tags != nil {
		res.Tags = make([]*TagResponseBody, len(v.Tags))
		for i, val := range v.Tags {
			res.Tags[i] = marshalResourceviewsTagViewToTagResponseBody(val)
		}
	}
	if v.Platforms != nil {
		res.Platforms = make([]*PlatformResponseBody, len(v.Platforms))
		for i, val := range v.Platforms {
			res.Platforms[i] = marshalResourceviewsPlatformViewToPlatformResponseBody(val)
		}
	}

	return res
}

// marshalResourceviewsCatalogViewToCatalogResponseBodyMin builds a value of
// type *CatalogResponseBodyMin from a value of type *resourceviews.CatalogView.
func marshalResourceviewsCatalogViewToCatalogResponseBodyMin(v *resourceviews.CatalogView) *CatalogResponseBodyMin {
	res := &CatalogResponseBodyMin{
		ID:   *v.ID,
		Name: *v.Name,
		Type: *v.Type,
	}

	return res
}

// marshalResourceviewsCategoryViewToCategoryResponseBody builds a value of
// type *CategoryResponseBody from a value of type *resourceviews.CategoryView.
func marshalResourceviewsCategoryViewToCategoryResponseBody(v *resourceviews.CategoryView) *CategoryResponseBody {
	res := &CategoryResponseBody{
		ID:   *v.ID,
		Name: *v.Name,
	}

	return res
}

// marshalResourceviewsResourceVersionDataViewToResourceVersionDataResponseBodyWithoutResource
// builds a value of type *ResourceVersionDataResponseBodyWithoutResource from
// a value of type *resourceviews.ResourceVersionDataView.
func marshalResourceviewsResourceVersionDataViewToResourceVersionDataResponseBodyWithoutResource(v *resourceviews.ResourceVersionDataView) *ResourceVersionDataResponseBodyWithoutResource {
	res := &ResourceVersionDataResponseBodyWithoutResource{
		ID:                  *v.ID,
		Version:             *v.Version,
		DisplayName:         *v.DisplayName,
		Deprecated:          v.Deprecated,
		Description:         *v.Description,
		MinPipelinesVersion: *v.MinPipelinesVersion,
		RawURL:              *v.RawURL,
		WebURL:              *v.WebURL,
		UpdatedAt:           *v.UpdatedAt,
		HubURLPath:          *v.HubURLPath,
	}
	if v.Platforms != nil {
		res.Platforms = make([]*PlatformResponseBody, len(v.Platforms))
		for i, val := range v.Platforms {
			res.Platforms[i] = marshalResourceviewsPlatformViewToPlatformResponseBody(val)
		}
	}

	return res
}

// marshalResourceviewsPlatformViewToPlatformResponseBody builds a value of
// type *PlatformResponseBody from a value of type *resourceviews.PlatformView.
func marshalResourceviewsPlatformViewToPlatformResponseBody(v *resourceviews.PlatformView) *PlatformResponseBody {
	res := &PlatformResponseBody{
		ID:   *v.ID,
		Name: *v.Name,
	}

	return res
}

// marshalResourceviewsTagViewToTagResponseBody builds a value of type
// *TagResponseBody from a value of type *resourceviews.TagView.
func marshalResourceviewsTagViewToTagResponseBody(v *resourceviews.TagView) *TagResponseBody {
	res := &TagResponseBody{
		ID:   *v.ID,
		Name: *v.Name,
	}

	return res
}
