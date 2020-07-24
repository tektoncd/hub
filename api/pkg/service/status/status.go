package status

import (
	"context"

	"github.com/tektoncd/hub/api/gen/status"
)

// status service implementation.
type service struct{}

// New returns the status service implementation.
func New() status.Service {
	return &service{}
}

// Return status 'ok' when the server has started successfully
func (s *service) Status(ctx context.Context) (res *status.StatusResult, err error) {

	res = &status.StatusResult{
		Status: "ok",
	}
	return res, nil
}
