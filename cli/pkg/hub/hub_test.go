package hub

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetURL(t *testing.T) {

	hub := &client{}
	err := hub.SetURL("http://localhost:80000")
	assert.NoError(t, err)

	err = hub.SetURL("localhost:8000")
	assert.NoError(t, err)

	err = hub.SetURL("http://80.80.79.9:80")
	assert.NoError(t, err)
}

func TestSetURL_InvalidCase(t *testing.T) {

	hub := &client{}
	err := hub.SetURL("abc")
	assert.Error(t, err)
	assert.EqualError(t, err, "parse \"abc\": invalid URI for request")
}
