package testutils

import (
	"errors"
	"path"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/tektoncd/hub/api/pkg/app"
)

type TestConfig struct {
	*app.APIConfig
	fixturePath string
	configPath  string
	err         error
}

var _ app.Config = (*TestConfig)(nil)

var once sync.Once
var tc *TestConfig

// config creates the test configuration once and returns the
// same every time the function is called
func Config() *TestConfig {
	once.Do(func() {
		tc = initializeConfig()
	})
	return tc
}

func (tc *TestConfig) Path() string {
	return tc.configPath
}

func (tc *TestConfig) FixturePath() string {
	return tc.fixturePath
}

func (tc *TestConfig) IsValid() bool {
	return tc.APIConfig != nil && tc.err == nil
}

func (tc *TestConfig) Error() error {
	return tc.err
}

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
