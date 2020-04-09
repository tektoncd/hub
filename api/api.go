package hub

import (
	"context"
	"log"

	api "github.com/tektoncd/hub/api/gen/api"
)

// api service example implementation.
// The example methods log the requests and return zero values.
type apisrvc struct {
	logger *log.Logger
}

// NewAPI returns the api service implementation.
func NewAPI(logger *log.Logger) api.Service {
	return &apisrvc{logger}
}

// Get all tasks and pipelines.
func (s *apisrvc) List(ctx context.Context) (err error) {
	s.logger.Print("api.list")
	return
}
