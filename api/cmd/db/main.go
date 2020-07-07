package main

import (
	"fmt"
	"os"

	"github.com/tektoncd/hub/api/pkg/app"
)

func main() {
	api, err := app.FromEnv()
	if err != nil {
		fmt.Fprintf(os.Stderr, "FATAL: failed to initialise: %s", err)
		os.Exit(1)
	}
	defer api.Cleanup()

	logger := api.Logger()
	if err = Migrate(api); err != nil {
		logger.Errorf("DB initialisation failed !!")
		return
	}
	logger.Info("DB initialisation successful !!")
}
