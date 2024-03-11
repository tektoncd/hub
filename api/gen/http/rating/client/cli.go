// Code generated by goa v3.15.1, DO NOT EDIT.
//
// rating HTTP client CLI support package
//
// Command:
// $ goa gen github.com/tektoncd/hub/api/design

package client

import (
	"encoding/json"
	"fmt"
	"strconv"

	rating "github.com/tektoncd/hub/api/gen/rating"
	goa "goa.design/goa/v3/pkg"
)

// BuildGetPayload builds the payload for the rating Get endpoint from CLI
// flags.
func BuildGetPayload(ratingGetID string, ratingGetToken string) (*rating.GetPayload, error) {
	var err error
	var id uint
	{
		var v uint64
		v, err = strconv.ParseUint(ratingGetID, 10, strconv.IntSize)
		id = uint(v)
		if err != nil {
			return nil, fmt.Errorf("invalid value for id, must be UINT")
		}
	}
	var token string
	{
		token = ratingGetToken
	}
	v := &rating.GetPayload{}
	v.ID = id
	v.Token = token

	return v, nil
}

// BuildUpdatePayload builds the payload for the rating Update endpoint from
// CLI flags.
func BuildUpdatePayload(ratingUpdateBody string, ratingUpdateID string, ratingUpdateToken string) (*rating.UpdatePayload, error) {
	var err error
	var body UpdateRequestBody
	{
		err = json.Unmarshal([]byte(ratingUpdateBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"rating\": 0\n   }'")
		}
		if body.Rating < 0 {
			err = goa.MergeErrors(err, goa.InvalidRangeError("body.rating", body.Rating, 0, true))
		}
		if body.Rating > 5 {
			err = goa.MergeErrors(err, goa.InvalidRangeError("body.rating", body.Rating, 5, false))
		}
		if err != nil {
			return nil, err
		}
	}
	var id uint
	{
		var v uint64
		v, err = strconv.ParseUint(ratingUpdateID, 10, strconv.IntSize)
		id = uint(v)
		if err != nil {
			return nil, fmt.Errorf("invalid value for id, must be UINT")
		}
	}
	var token string
	{
		token = ratingUpdateToken
	}
	v := &rating.UpdatePayload{
		Rating: body.Rating,
	}
	v.ID = id
	v.Token = token

	return v, nil
}
