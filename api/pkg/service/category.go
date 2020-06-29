package hub

import (
	"context"
	"log"

	category "github.com/tektoncd/hub/api/gen/category"
)

// category service example implementation.
// The example methods log the requests and return zero values.
type categorysrvc struct {
	logger *log.Logger
}

// NewCategory returns the category service implementation.
func NewCategory(logger *log.Logger) category.Service {
	return &categorysrvc{logger}
}

// Get all Categories with their tags sorted by name
func (s *categorysrvc) All(ctx context.Context) (res []*category.Category, err error) {
	s.logger.Print("category.All")
	return
}
