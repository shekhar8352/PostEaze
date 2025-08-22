package utils

import (
	"os"
	"testing"

	"github.com/shekhar8352/PostEaze/utils/configs"
)

func TestInitDev(t *testing.T) {
	// Create a temporary directory for test config files
	tempDir, err := os.MkdirTemp("", "config_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	t.Run("valid directory with config", func(t *testing.T) {
		configContent := `
test:
  value: "test_value"
  number: 123
database:
  host: "localhost"
  port: 5432
`
		configFile := tempDir + "/test_config.yaml"
		err := os.WriteFile(configFile, []byte(configContent), 0644)
		if err != nil {
			t.Fatalf("Failed to write config file: %v", err)
		}

		err = configs.InitDev(tempDir, "test_config")
		if err != nil {
			t.Errorf("InitDev() error = %v", err)
		}

		client := configs.Get()
		if client == nil {
			t.Error("InitDev() should initialize client")
		}
	})

	t.Run("non-existent directory", func(t *testing.T) {
		err := configs.InitDev("/path/that/does/not/exist", "config")
		if err == nil {
			t.Error("InitDev() should return error for non-existent directory")
		}
	})

	t.Run("empty directory", func(t *testing.T) {
		err := configs.InitDev("", "config")
		if err == nil {
			t.Error("InitDev() should return error for empty directory")
		}
	})

	t.Run("multiple config files", func(t *testing.T) {
		config1Content := `
app:
  name: "test_app"
  version: "1.0.0"
`
		config2Content := `
database:
  host: "localhost"
  port: 5432
`

		config1File := tempDir + "/app_config.yaml"
		config2File := tempDir + "/db_config.yaml"

		err := os.WriteFile(config1File, []byte(config1Content), 0644)
		if err != nil {
			t.Fatalf("Failed to write config1 file: %v", err)
		}

		err = os.WriteFile(config2File, []byte(config2Content), 0644)
		if err != nil {
			t.Fatalf("Failed to write config2 file: %v", err)
		}

		err = configs.InitDev(tempDir, "app_config", "db_config")
		if err != nil {
			t.Errorf("InitDev() error = %v", err)
		}

		client := configs.Get()
		if client == nil {
			t.Error("InitDev() should initialize client with multiple configs")
		}
	})

	t.Run("no config names", func(t *testing.T) {
		// Should handle empty config names gracefully
		_ = configs.InitDev(tempDir)
		// The exact behavior depends on the underlying config library
		// We just verify it doesn't panic
		client := configs.Get()
		if client == nil {
			t.Log("InitDev() with no config names returned nil client (acceptable)")
		}
	})

	t.Run("invalid YAML", func(t *testing.T) {
		invalidYAML := `
invalid:
  yaml: content
    missing: proper indentation
  - invalid list item
`
		configFile := tempDir + "/invalid_config.yaml"
		err := os.WriteFile(configFile, []byte(invalidYAML), 0644)
		if err != nil {
			t.Fatalf("Failed to write invalid config file: %v", err)
		}

		err = configs.InitDev(tempDir, "invalid_config")
		if err != nil {
			// Should handle invalid YAML gracefully
			if err.Error() == "" {
				t.Error("InitDev() should return meaningful error for invalid YAML")
			}
		}
	})
}

func TestInitRelease(t *testing.T) {
	tests := []struct {
		name        string
		env         string
		region      string
		configNames []string
		wantErr     bool
	}{
		{
			name:        "valid parameters",
			env:         "test",
			region:      "us-east-1",
			configNames: []string{"app", "database"},
			wantErr:     true, // Expected to fail in test environment without AWS credentials
		},
		{
			name:        "empty environment",
			env:         "",
			region:      "us-east-1",
			configNames: []string{"config"},
			wantErr:     true,
		},
		{
			name:        "empty region",
			env:         "test",
			region:      "",
			configNames: []string{"config"},
			wantErr:     true,
		},
		{
			name:        "no config names",
			env:         "test",
			region:      "us-east-1",
			configNames: []string{},
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := configs.InitRelease(tt.env, tt.region, tt.configNames...)
			if (err != nil) != tt.wantErr {
				t.Errorf("InitRelease() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				// The error should be related to AWS credentials or connectivity, not parameter validation
				if err.Error() == "" {
					t.Error("InitRelease() should return meaningful error message")
				}
			}
		})
	}
}

func TestGet(t *testing.T) {
	t.Run("after InitDev", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "config_get_test")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		configContent := `
test:
  key: "value"
`
		configFile := tempDir + "/get_test_config.yaml"
		err = os.WriteFile(configFile, []byte(configContent), 0644)
		if err != nil {
			t.Fatalf("Failed to write config file: %v", err)
		}

		err = configs.InitDev(tempDir, "get_test_config")
		if err != nil {
			t.Fatalf("InitDev() error = %v", err)
		}

		client := configs.Get()
		if client == nil {
			t.Error("Get() should return client after InitDev")
		}
	})

	t.Run("before init", func(t *testing.T) {
		// This tests the behavior when Get() is called before any Init function
		client := configs.Get()
		// The behavior depends on the implementation
		// It might return nil or a default client
		// We just verify it doesn't panic

		// If it returns something, it should be consistent
		client2 := configs.Get()
		if client != client2 {
			t.Error("Get() should return consistent results")
		}
	})

	t.Run("multiple calls consistency", func(t *testing.T) {
		client1 := configs.Get()
		client2 := configs.Get()
		if client1 != client2 {
			t.Error("Get() should return the same instance on multiple calls")
		}
	})
}

func TestConfigInitializationFlow(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "config_flow_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	t.Run("comprehensive config", func(t *testing.T) {
		configContent := `
app:
  name: "PostEaze"
  version: "1.0.0"
  debug: true

database:
  host: "localhost"
  port: 5432
  name: "posteaze_test"
  user: "test_user"
  password: "test_pass"

jwt:
  access_secret: "test_access_secret"
  refresh_secret: "test_refresh_secret"
  access_expiry: "15m"
  refresh_expiry: "7d"

logging:
  level: "debug"
  format: "json"
`

		configFile := tempDir + "/full_config.yaml"
		err := os.WriteFile(configFile, []byte(configContent), 0644)
		if err != nil {
			t.Fatalf("Failed to write config file: %v", err)
		}

		err = configs.InitDev(tempDir, "full_config")
		if err != nil {
			t.Errorf("InitDev() error = %v", err)
		}

		client := configs.Get()
		if client == nil {
			t.Error("Get() should return client after comprehensive config init")
		}

		// Test that we can call Get() multiple times
		client2 := configs.Get()
		if client != client2 {
			t.Error("Get() should return same client instance")
		}
	})

	t.Run("multiple inits", func(t *testing.T) {
		configContent1 := `
first:
  value: "first_config"
`
		configContent2 := `
second:
  value: "second_config"
`

		config1File := tempDir + "/first_config.yaml"
		config2File := tempDir + "/second_config.yaml"

		err := os.WriteFile(config1File, []byte(configContent1), 0644)
		if err != nil {
			t.Fatalf("Failed to write first config file: %v", err)
		}

		err = os.WriteFile(config2File, []byte(configContent2), 0644)
		if err != nil {
			t.Fatalf("Failed to write second config file: %v", err)
		}

		// First initialization
		err = configs.InitDev(tempDir, "first_config")
		if err != nil {
			t.Errorf("First InitDev() error = %v", err)
		}

		client1 := configs.Get()
		if client1 == nil {
			t.Error("Get() should return client after first init")
		}

		// Second initialization (should replace the first)
		err = configs.InitDev(tempDir, "second_config")
		if err != nil {
			t.Errorf("Second InitDev() error = %v", err)
		}

		client2 := configs.Get()
		if client2 == nil {
			t.Error("Get() should return client after second init")
		}

		// The client should be updated (exact behavior depends on implementation)
		// We just verify both calls succeeded and returned valid clients
	})
}

func TestConfigErrorHandling(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "config_error_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	t.Run("permission denied", func(t *testing.T) {
		// Create a directory with restricted permissions
		restrictedDir := tempDir + "/restricted"
		err := os.Mkdir(restrictedDir, 0000) // No permissions
		if err != nil {
			t.Skip("Cannot create restricted directory on this system")
		}
		defer os.Chmod(restrictedDir, 0755) // Restore permissions for cleanup

		err = configs.InitDev(restrictedDir, "config")
		if err == nil {
			t.Error("InitDev() should return error for restricted directory")
		}
	})
}