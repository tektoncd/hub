package httpcheck

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/itchyny/gojq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// body //////////////////////////////////////////////////////////

// WithBody adds the []byte data to the body.
func (tt *Tester) WithBody(body []byte) *Tester {
	tt.request.Body = io.NopCloser(bytes.NewReader(body))
	tt.request.ContentLength = int64(len(body))
	return tt
}

// HasBody checks if the body is equal to provided []byte data.
func (tt *Tester) HasBody(expected []byte) *Tester {
	body, err := io.ReadAll(tt.response.Body)
	require.NoError(tt.t, err)
	tt.response.Body.Close()
	defer func(body []byte) {
		tt.response.Body = io.NopCloser(bytes.NewReader(body))
	}(body)

	assert.Equal(tt.t, expected, body)
	return tt
}

// MustHasBody checks if the body is equal to provided []byte data.
func (tt *Tester) MustHasBody(expected []byte) *Tester {
	body, err := io.ReadAll(tt.response.Body)
	require.NoError(tt.t, err)
	tt.response.Body.Close()
	defer func(body []byte) {
		tt.response.Body = io.NopCloser(bytes.NewReader(body))
	}(body)

	require.Equal(tt.t, expected, body)
	return tt
}

// ContainsBody checks if the body contains provided [] byte data.
func (tt *Tester) ContainsBody(segment []byte) *Tester {
	body, err := io.ReadAll(tt.response.Body)
	require.NoError(tt.t, err)
	tt.response.Body.Close()
	defer func(body []byte) {
		tt.response.Body = io.NopCloser(bytes.NewReader(body))
	}(body)

	if !bytes.Contains(body, segment) {
		assert.Fail(tt.t, fmt.Sprintf("%#v does not contain %#v", body, segment))
	}
	return tt
}

// MustContainsBody checks if the body contains provided [] byte data.
func (tt *Tester) MustContainsBody(segment []byte) *Tester {
	body, err := io.ReadAll(tt.response.Body)
	require.NoError(tt.t, err)
	tt.response.Body.Close()
	defer func(body []byte) {
		tt.response.Body = io.NopCloser(bytes.NewReader(body))
	}(body)

	if !bytes.Contains(body, segment) {
		require.Fail(tt.t, fmt.Sprintf("%#v does not contain %#v", body, segment))
	}
	return tt
}

// NotContainsBody checks if the body does not contain provided [] byte data.
func (tt *Tester) NotContainsBody(segment []byte) *Tester {
	body, err := io.ReadAll(tt.response.Body)
	require.NoError(tt.t, err)
	tt.response.Body.Close()
	defer func(body []byte) {
		tt.response.Body = io.NopCloser(bytes.NewReader(body))
	}(body)

	if bytes.Contains(body, segment) {
		assert.Fail(tt.t, fmt.Sprintf("%#v contains %#v", body, segment))
	}
	return tt
}

// MustNotContainsBody checks if the body does not contain provided [] byte data.
func (tt *Tester) MustNotContainsBody(segment []byte) *Tester {
	body, err := io.ReadAll(tt.response.Body)
	require.NoError(tt.t, err)
	tt.response.Body.Close()
	defer func(body []byte) {
		tt.response.Body = io.NopCloser(bytes.NewReader(body))
	}(body)

	if bytes.Contains(body, segment) {
		require.Fail(tt.t, fmt.Sprintf("%#v contains %#v", body, segment))
	}
	return tt
}

// WithString adds the string to the body.
func (tt *Tester) WithString(body string) *Tester {
	tt.request.Body = io.NopCloser(strings.NewReader(body))
	tt.request.ContentLength = int64(len(body))
	return tt
}

// HasString converts the response to a string type and then compares it with the given string.
func (tt *Tester) HasString(expected string) *Tester {
	body, err := io.ReadAll(tt.response.Body)
	require.NoError(tt.t, err)
	tt.response.Body.Close()
	defer func(body []byte) {
		tt.response.Body = io.NopCloser(bytes.NewReader(body))
	}(body)

	assert.Equal(tt.t, expected, string(body))
	return tt
}

// MustHasString converts the response to a string type and then compares it with the given string.
func (tt *Tester) MustHasString(expected string) *Tester {
	body, err := io.ReadAll(tt.response.Body)
	require.NoError(tt.t, err)
	tt.response.Body.Close()
	defer func(body []byte) {
		tt.response.Body = io.NopCloser(bytes.NewReader(body))
	}(body)

	require.Equal(tt.t, expected, string(body))
	return tt
}

// ContainsString converts the response to a string type and then checks it containing the given string.
func (tt *Tester) ContainsString(substr string) *Tester {
	body, err := io.ReadAll(tt.response.Body)
	require.NoError(tt.t, err)
	tt.response.Body.Close()
	defer func(body []byte) {
		tt.response.Body = io.NopCloser(bytes.NewReader(body))
	}(body)

	assert.Contains(tt.t, string(body), substr)
	return tt
}

// MustContainsString converts the response to a string type and then checks it containing the given string.
func (tt *Tester) MustContainsString(substr string) *Tester {
	body, err := io.ReadAll(tt.response.Body)
	require.NoError(tt.t, err)
	tt.response.Body.Close()
	defer func(body []byte) {
		tt.response.Body = io.NopCloser(bytes.NewReader(body))
	}(body)

	require.Contains(tt.t, string(body), substr)
	return tt
}

// NotContainsString converts the response to a string type and then checks if it does not
// contain the given string.
func (tt *Tester) NotContainsString(substr string) *Tester {
	body, err := io.ReadAll(tt.response.Body)
	require.NoError(tt.t, err)
	tt.response.Body.Close()
	defer func(body []byte) {
		tt.response.Body = io.NopCloser(bytes.NewReader(body))
	}(body)

	assert.NotContains(tt.t, string(body), substr)
	return tt
}

// MustNotContainsString converts the response to a string type and then checks if it does not
// contain the given string.
func (tt *Tester) MustNotContainsString(substr string) *Tester {
	body, err := io.ReadAll(tt.response.Body)
	require.NoError(tt.t, err)
	tt.response.Body.Close()
	defer func(body []byte) {
		tt.response.Body = io.NopCloser(bytes.NewReader(body))
	}(body)

	require.NotContains(tt.t, string(body), substr)
	return tt
}

func (tt *Tester) MatchesJSONQuery(q string) *Tester {
	body, err := io.ReadAll(tt.response.Body)
	require.NoError(tt.t, err)
	tt.response.Body.Close()
	defer func(body []byte) {
		tt.response.Body = io.NopCloser(bytes.NewReader(body))
	}(body)
	var in any
	require.NoError(tt.t, json.Unmarshal(body, &in), "failed to unmarshal json: %s", string(body))
	jq, err := gojq.Parse(q)
	require.NoError(tt.t, err, "failed to parse query %q: %s", q)
	it := jq.Run(in)
	var detect bool
	for {
		v, ok := it.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			require.NoError(tt.t, err, "query %q does not match: %s", q, string(body))
		}
		if v != nil {
			detect = true
		}
	}
	assert.True(tt.t, detect, "query %q does not match: %s", q, string(body))
	return tt
}

func (tt *Tester) NotMatchesJSONQuery(q string) *Tester {
	body, err := io.ReadAll(tt.response.Body)
	require.NoError(tt.t, err)
	tt.response.Body.Close()
	defer func(body []byte) {
		tt.response.Body = io.NopCloser(bytes.NewReader(body))
	}(body)
	var in any
	require.NoError(tt.t, json.Unmarshal(body, &in))
	jq, err := gojq.Parse(q)
	require.NoError(tt.t, err, "failed to parse query %q: %s", q)
	it := jq.Run(in)
	var detect bool
	for {
		v, ok := it.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			require.NoError(tt.t, err, "query %q does not match: %s", q, string(body))
		}
		if v != nil {
			detect = true
		}
	}
	assert.False(tt.t, detect, "query %q does not match: %s", q, string(body))
	return tt
}
