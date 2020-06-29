package category

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/ikawaha/goahttpcheck"
	"github.com/stretchr/testify/assert"
	category "github.com/tektoncd/hub/api/gen/category"
	server "github.com/tektoncd/hub/api/gen/http/category/server"
	"github.com/tektoncd/hub/api/pkg/testutils"
)

func TestCategories_List_Http(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	checker := goahttpcheck.New()
	checker.Mount(
		server.NewListHandler,
		server.MountListHandler,
		category.NewListEndpoint(New(tc)))

	checker.Test(t, http.MethodGet, "/categories").Check().
		HasStatus(http.StatusOK).Cb(func(r *http.Response) {
		b, readErr := ioutil.ReadAll(r.Body)
		assert.NoError(t, readErr)
		defer r.Body.Close()

		var jsonMap []map[string]interface{}
		marshallErr := json.Unmarshal([]byte(b), &jsonMap)
		assert.NoError(t, marshallErr)

		assert.Equal(t, 3, len(jsonMap))
		assert.Equal(t, "abc", jsonMap[0]["name"])
	})
}
