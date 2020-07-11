package category

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/ikawaha/goahttpcheck"
	"github.com/stretchr/testify/assert"
	"gotest.tools/v3/golden"

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

		res, err := testutils.FormatJSON(b)
		assert.NoError(t, err)

		golden.Assert(t, res, fmt.Sprintf("%s.golden", t.Name()))
	})
}
