package category

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tektoncd/hub/api/pkg/testutils"
)

func TestCategory_List(t *testing.T) {
	tc := testutils.Setup(t)
	testutils.LoadFixtures(t, tc.FixturePath())

	category := New(tc)
	all, err := category.List(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, 3, len(all))
	assert.Equal(t, 2, len(all[0].Tags))
	assert.Equal(t, "abc", all[0].Name) // categories are sorted by name
}
