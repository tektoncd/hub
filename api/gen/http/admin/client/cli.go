// Code generated by goa v3.7.5, DO NOT EDIT.
//
// admin HTTP client CLI support package
//
// Command:
// $ goa gen github.com/tektoncd/hub/api/design

package client

import (
	"encoding/json"
	"fmt"

	admin "github.com/tektoncd/hub/api/gen/admin"
	goa "goa.design/goa/v3/pkg"
)

// BuildUpdateAgentPayload builds the payload for the admin UpdateAgent
// endpoint from CLI flags.
func BuildUpdateAgentPayload(adminUpdateAgentBody string, adminUpdateAgentToken string) (*admin.UpdateAgentPayload, error) {
	var err error
	var body UpdateAgentRequestBody
	{
		err = json.Unmarshal([]byte(adminUpdateAgentBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"name\": \"abc\",\n      \"scopes\": [\n         \"catalog-refresh\",\n         \"agent:create\"\n      ]\n   }'")
		}
		if body.Scopes == nil {
			err = goa.MergeErrors(err, goa.MissingFieldError("scopes", "body"))
		}
		if err != nil {
			return nil, err
		}
	}
	var token string
	{
		token = adminUpdateAgentToken
	}
	v := &admin.UpdateAgentPayload{
		Name: body.Name,
	}
	if body.Scopes != nil {
		v.Scopes = make([]string, len(body.Scopes))
		for i, val := range body.Scopes {
			v.Scopes[i] = val
		}
	}
	v.Token = token

	return v, nil
}

// BuildRefreshConfigPayload builds the payload for the admin RefreshConfig
// endpoint from CLI flags.
func BuildRefreshConfigPayload(adminRefreshConfigToken string) (*admin.RefreshConfigPayload, error) {
	var token string
	{
		token = adminRefreshConfigToken
	}
	v := &admin.RefreshConfigPayload{}
	v.Token = token

	return v, nil
}
