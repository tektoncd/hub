package resource

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tektoncd/hub/api/gen/resource"
	"github.com/tektoncd/hub/api/pkg/testutils"
)

func TestQuery_DefaultLimit(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	resourceSvc := New(tc)
	payload := &resource.QueryPayload{Name: "", Type: "", Limit: 100}
	all, err := resourceSvc.Query(context.Background(), payload)
	assert.NoError(t, err)
	assert.Equal(t, 5, len(all))
}

func TestQuery_ByLimit(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	resourceSvc := New(tc)
	payload := &resource.QueryPayload{Name: "", Type: "", Limit: 2}
	all, err := resourceSvc.Query(context.Background(), payload)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(all))
	assert.Equal(t, "tekton", all[0].Name)
}

func TestQuery_ByName(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	resourceSvc := New(tc)
	payload := &resource.QueryPayload{Name: "tekton", Type: "", Limit: 100}
	all, err := resourceSvc.Query(context.Background(), payload)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(all))
	assert.Equal(t, "0.2", all[0].LatestVersion.Version)
}

func TestQuery_ByPartialName(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	resourceSvc := New(tc)
	payload := &resource.QueryPayload{Name: "build", Type: "", Limit: 100}
	all, err := resourceSvc.Query(context.Background(), payload)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(all))
}

func TestQuery_ByType(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	resourceSvc := New(tc)
	payload := &resource.QueryPayload{Name: "", Type: "pipeline", Limit: 100}
	all, err := resourceSvc.Query(context.Background(), payload)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(all))
}

func TestQuery_ByNameAndType(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	resourceSvc := New(tc)
	payload := &resource.QueryPayload{Name: "build", Type: "pipeline", Limit: 100}
	all, err := resourceSvc.Query(context.Background(), payload)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(all))
	assert.Equal(t, "build-pipeline", all[0].Name)
}

func TestQuery_NotFoundError(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	resourceSvc := New(tc)
	payload := &resource.QueryPayload{Name: "foo", Type: "", Limit: 100}
	_, err := resourceSvc.Query(context.Background(), payload)
	assert.Error(t, err)
	assert.EqualError(t, err, "Resource not found")
}

func TestList_ByLimit(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	resourceSvc := New(tc)
	payload := &resource.ListPayload{Limit: 3}
	all, err := resourceSvc.List(context.Background(), payload)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(all))
	assert.Equal(t, "tekton", all[0].Name)
}

func TestVersionsByID(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	resourceSvc := New(tc)
	payload := &resource.VersionsByIDPayload{ID: 1}
	all, err := resourceSvc.VersionsByID(context.Background(), payload)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(all.Versions))
	assert.Equal(t, "0.2", all.Latest.Version)
}

func TestVersionsByID_NotFoundError(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	resourceSvc := New(tc)
	payload := &resource.VersionsByIDPayload{ID: 11}
	_, err := resourceSvc.VersionsByID(context.Background(), payload)
	assert.Error(t, err)
	assert.EqualError(t, err, "Resource not found")
}

func TestByTypeNameVersion(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	resourceSvc := New(tc)
	payload := &resource.ByTypeNameVersionPayload{Type: "task", Name: "tkn", Version: "0.1"}
	res, err := resourceSvc.ByTypeNameVersion(context.Background(), payload)
	assert.NoError(t, err)
	assert.Equal(t, "0.1", res.Version)
}

func TestByTypeNameVersion_NoResourceWithName(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	resourceSvc := New(tc)
	payload := &resource.ByTypeNameVersionPayload{Type: "task", Name: "foo", Version: "0.1"}
	_, err := resourceSvc.ByTypeNameVersion(context.Background(), payload)
	assert.Error(t, err)
	assert.EqualError(t, err, "Resource not found")
}

func TestByTypeNameVersion_ResourceVersionNotFound(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	resourceSvc := New(tc)
	payload := &resource.ByTypeNameVersionPayload{Type: "task", Name: "tekton", Version: "0.9"}
	_, err := resourceSvc.ByTypeNameVersion(context.Background(), payload)
	assert.Error(t, err)
	assert.EqualError(t, err, "Resource not found")
}
