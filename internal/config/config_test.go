package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	// Set required environment variables
	os.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/testdb")
	defer os.Unsetenv("DATABASE_URL")

	cfg, err := Load()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// Test database config
	assert.Equal(t, "postgres://user:pass@localhost:5432/testdb", cfg.Database.URL)
	assert.Equal(t, 100, cfg.Database.MaxConnections)
	assert.Equal(t, 10, cfg.Database.MaxIdle)

	// Test storage config
	assert.Equal(t, "local", cfg.Storage.Type)
	assert.Equal(t, "./uploads", cfg.Storage.Path)

	// Test server config
	assert.Equal(t, "8080", cfg.Server.Port)

	// Test worker config
	assert.Equal(t, 4, cfg.Worker.Concurrency)
	assert.Equal(t, 5*time.Second, cfg.Worker.PollInterval)

	// Test logging config
	assert.Equal(t, "info", cfg.Logging.Level)
	assert.Equal(t, "json", cfg.Logging.Format)
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config",
			cfg: &Config{
				Database: DatabaseConfig{URL: "postgres://localhost/db"},
				Storage:  StorageConfig{Type: "local", Path: "/tmp"},
				Logging:  LoggingConfig{Level: "info", Format: "json"},
			},
			wantErr: false,
		},
		{
			name: "missing database URL",
			cfg: &Config{
				Database: DatabaseConfig{URL: ""},
				Storage:  StorageConfig{Type: "local", Path: "/tmp"},
				Logging:  LoggingConfig{Level: "info", Format: "json"},
			},
			wantErr: true,
			errMsg:  "DATABASE_URL is required",
		},
		{
			name: "invalid storage type",
			cfg: &Config{
				Database: DatabaseConfig{URL: "postgres://localhost/db"},
				Storage:  StorageConfig{Type: "invalid", Path: "/tmp"},
				Logging:  LoggingConfig{Level: "info", Format: "json"},
			},
			wantErr: true,
			errMsg:  "STORAGE_TYPE must be 'local' or 's3'",
		},
		{
			name: "missing storage path for local",
			cfg: &Config{
				Database: DatabaseConfig{URL: "postgres://localhost/db"},
				Storage:  StorageConfig{Type: "local", Path: ""},
				Logging:  LoggingConfig{Level: "info", Format: "json"},
			},
			wantErr: true,
			errMsg:  "STORAGE_PATH is required when STORAGE_TYPE is 'local'",
		},
		{
			name: "invalid log level",
			cfg: &Config{
				Database: DatabaseConfig{URL: "postgres://localhost/db"},
				Storage:  StorageConfig{Type: "local", Path: "/tmp"},
				Logging:  LoggingConfig{Level: "invalid", Format: "json"},
			},
			wantErr: true,
			errMsg:  "LOG_LEVEL must be one of: debug, info, warn, error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetEnv(t *testing.T) {
	os.Setenv("TEST_VAR", "test_value")
	defer os.Unsetenv("TEST_VAR")

	assert.Equal(t, "test_value", getEnv("TEST_VAR", "default"))
	assert.Equal(t, "default", getEnv("NONEXISTENT_VAR", "default"))
}

func TestGetEnvInt(t *testing.T) {
	os.Setenv("TEST_INT", "42")
	defer os.Unsetenv("TEST_INT")

	assert.Equal(t, 42, getEnvInt("TEST_INT", 10))
	assert.Equal(t, 10, getEnvInt("NONEXISTENT_INT", 10))

	os.Setenv("INVALID_INT", "not_a_number")
	defer os.Unsetenv("INVALID_INT")
	assert.Equal(t, 10, getEnvInt("INVALID_INT", 10))
}

func TestGetEnvDuration(t *testing.T) {
	os.Setenv("TEST_DURATION", "5s")
	defer os.Unsetenv("TEST_DURATION")

	assert.Equal(t, 5*time.Second, getEnvDuration("TEST_DURATION", 10*time.Second))
	assert.Equal(t, 10*time.Second, getEnvDuration("NONEXISTENT_DURATION", 10*time.Second))

	os.Setenv("INVALID_DURATION", "not_a_duration")
	defer os.Unsetenv("INVALID_DURATION")
	assert.Equal(t, 10*time.Second, getEnvDuration("INVALID_DURATION", 10*time.Second))
}
