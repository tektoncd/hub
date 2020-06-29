package testutils

import (
	"sync"
	"testing"
)

var setup sync.Once

func Setup(t *testing.T) *TestConfig {
	tc := Config()
	setup.Do(func() {
		if !tc.IsValid() {
			t.Fatalf("Failed to create test configuration object: %s", tc.Error())
		}

		if err := applyMigration(); err != nil {
			t.Fatalf("Failed to apply migration: %s", tc.Error())
		}
	})

	// NOTE: not calling tc.Cleanup as the connection
	//       need to be resued by other tests
	// t.Cleanup(func() {
	// 	 tc.Cleanup()
	// })

	return tc
}
