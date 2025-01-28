package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	tests := map[string]struct {
		configYAML  string
		assertError assert.ErrorAssertionFunc
		expected    *Config
	}{
		"Valid config": {
			configYAML: `
port: "8080"
mongodb:
  uri: "mongodb://localhost:27017"
  database: "testdb"`,
			assertError: assert.NoError,
			expected: &Config{
				Port: "8080",
				MonogoDB: MongoDB{
					URI:      "mongodb://localhost:27017",
					Database: "testdb",
				},
			},
		},
		"Missing required field": {
			configYAML: `
mongodb:
  uri: "mongodb://localhost:27017"
  database: "testdb"`,
			assertError: assert.Error,
		},
		"Invalid yaml": {
			configYAML: `
port: 8080
mongodb: {
`,
			assertError: assert.Error,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Create temporary config file
			tmpDir := t.TempDir()
			configPath := filepath.Join(tmpDir, "config.yaml")

			err := os.WriteFile(configPath, []byte(tt.configYAML), 0644)
			require.NoError(t, err)

			got, err := LoadConfig(configPath)
			tt.assertError(t, err)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	_, err := LoadConfig("nonexistent.yaml")
	assert.Error(t, err)
}
