package status

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/ikawaha/goahttpcheck"
	"github.com/stretchr/testify/assert"
	"github.com/tektoncd/hub/api/gen/http/status/server"
	"github.com/tektoncd/hub/api/gen/status"
)

func TestOk_http(t *testing.T) {

	checker := goahttpcheck.New()
	checker.Mount(
		server.NewStatusHandler,
		server.MountStatusHandler,
		status.NewStatusEndpoint(New()),
	)

	checker.Test(t, http.MethodGet, "/").Check().
		HasStatus(http.StatusOK).Cb(func(r *http.Response) {
		b, readErr := ioutil.ReadAll(r.Body)
		assert.NoError(t, readErr)
		defer r.Body.Close()

		var jsonMap map[string]interface{}
		marshallErr := json.Unmarshal([]byte(b), &jsonMap)
		assert.NoError(t, marshallErr)

		assert.Equal(t, "ok", jsonMap["status"])
	})
}
