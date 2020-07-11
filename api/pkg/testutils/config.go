package testutils

import (
	"errors"
	"path"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/tektoncd/hub/api/pkg/app"
)

// TestConfig defines configurations required for running tests
// APIConfig contains the db object and logger
// fixturePath is the path to fixture directory which contains test data
// configPath is the path to test config file
// err will have error if occured during initialising the test db connection
type TestConfig struct {
	*app.APIConfig
	fixturePath string
	configPath  string
	err         error
}

var _ app.Config = (*TestConfig)(nil)

var once sync.Once
var tc *TestConfig

// Config creates the test configuration once and returns the
// same every time the function is called
func Config() *TestConfig {
	once.Do(func() {
		tc = initializeConfig()
	})
	return tc
}

// Path return the file path to test config
func (tc *TestConfig) Path() string {
	return tc.configPath
}

// FixturePath return the directory path to fixtures
func (tc *TestConfig) FixturePath() string {
	return tc.fixturePath
}

// IsValid checks if connection to test db is valid or not
func (tc *TestConfig) IsValid() bool {
	return tc.APIConfig != nil && tc.err == nil
}

// Error returns error if occured during initialising connection to test db
func (tc *TestConfig) Error() error {
	return tc.err
}

// initializeConfig compute the path to test config and fixture directory
// then initiate the test db connection and returns a TestConfig Object
func initializeConfig() *TestConfig {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return &TestConfig{err: errors.New("failed to find filename")}
	}

	testDir := filepath.Join(path.Dir(filename), "..", "..", "test")
	configPath := filepath.Join(testDir, "config", "env.test")
	fixturePath := filepath.Join(testDir, "fixtures")
	api, err := app.FromEnvFile(configPath)

	return &TestConfig{
		APIConfig:   api,
		fixturePath: fixturePath,
		configPath:  configPath,
		err:         err,
	}
}
