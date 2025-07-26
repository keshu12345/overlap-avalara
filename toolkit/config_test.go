package toolkit

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test configuration struct
type TestConfig struct {
	Server struct {
		Port int    `mapstructure:"port"`
		Host string `mapstructure:"host"`
	} `mapstructure:"server"`
	Database struct {
		URL      string `mapstructure:"url"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
	} `mapstructure:"database"`
	Debug bool `mapstructure:"debug"`
}

func TestNewConfig_BasicYAMLConfig(t *testing.T) {
	configContent := `
server:
  port: 8080
  host: localhost
database:
  url: postgresql://localhost:5432/testdb
  username: testuser
  password: testpass
debug: true
`

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	var config TestConfig
	err = NewConfig(&config, configPath, "")

	assert.NoError(t, err)
	assert.Equal(t, 8080, config.Server.Port)
	assert.Equal(t, "localhost", config.Server.Host)
	assert.Equal(t, "postgresql://localhost:5432/testdb", config.Database.URL)
	assert.Equal(t, "testuser", config.Database.Username)
	assert.Equal(t, "testpass", config.Database.Password)
	assert.True(t, config.Debug)
}

func TestNewConfig_BasicJSONConfig(t *testing.T) {
	configContent := `{
  "server": {
    "port": 9090,
    "host": "0.0.0.0"
  },
  "database": {
    "url": "mysql://localhost:3306/testdb",
    "username": "jsonuser",
    "password": "jsonpass"
  },
  "debug": false
}`

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	var config TestConfig
	err = NewConfig(&config, configPath, "")

	assert.NoError(t, err)
	assert.Equal(t, 9090, config.Server.Port)
	assert.Equal(t, "0.0.0.0", config.Server.Host)
	assert.Equal(t, "mysql://localhost:3306/testdb", config.Database.URL)
	assert.Equal(t, "jsonuser", config.Database.Username)
	assert.Equal(t, "jsonpass", config.Database.Password)
	assert.False(t, config.Debug)
}

func TestNewConfig_WithEnvironmentVariables(t *testing.T) {
	configContent := `
server:
  port: 8080
  host: localhost
database:
  url: postgresql://localhost:5432/testdb
  username: testuser
  password: testpass
debug: true
`

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	os.Setenv("TEST_PORT", "9000")
	os.Setenv("TEST_HOST", "production.com")
	os.Setenv("TEST_DEBUG", "false")
	defer func() {
		os.Unsetenv("TEST_PORT")
		os.Unsetenv("TEST_HOST")
		os.Unsetenv("TEST_DEBUG")
	}()

	envVars := map[string]string{
		"server.port": "TEST_PORT",
		"server.host": "TEST_HOST",
		"debug":       "TEST_DEBUG",
	}

	var config TestConfig
	err = NewConfig(&config, configPath, "", envVars)

	assert.NoError(t, err)
	assert.Equal(t, 9000, config.Server.Port)             // Overridden by env var
	assert.Equal(t, "production.com", config.Server.Host) // Overridden by env var
	assert.False(t, config.Debug)                         // Overridden by env var
	assert.Equal(t, "postgresql://localhost:5432/testdb", config.Database.URL)
}

func TestNewConfig_WithOverrideConfig(t *testing.T) {
	// Base config
	baseConfigContent := `
server:
  port: 8080
  host: localhost
database:
  url: postgresql://localhost:5432/testdb
  username: testuser
  password: testpass
debug: true
`

	// Override config
	overrideConfigContent := `
server:
  port: 9090
database:
  password: overridepass
debug: false
`

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	overridePath := filepath.Join(tmpDir, "override.yaml")

	err := os.WriteFile(configPath, []byte(baseConfigContent), 0644)
	require.NoError(t, err)

	err = os.WriteFile(overridePath, []byte(overrideConfigContent), 0644)
	require.NoError(t, err)

	var config TestConfig
	err = NewConfig(&config, configPath, overridePath)

	assert.NoError(t, err)
	assert.Equal(t, 9090, config.Server.Port)                                  // Overridden
	assert.Equal(t, "localhost", config.Server.Host)                           // From base config
	assert.Equal(t, "postgresql://localhost:5432/testdb", config.Database.URL) // From base config
	assert.Equal(t, "testuser", config.Database.Username)                      // From base config
	assert.Equal(t, "overridepass", config.Database.Password)                  // Overridden
	assert.False(t, config.Debug)                                              // Overridden
}

func TestNewConfig_MultipleEnvVarMaps(t *testing.T) {
	configContent := `
server:
  port: 8080
  host: localhost
debug: true
`

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	// Set environment variables
	os.Setenv("SERVER_PORT", "3000")
	os.Setenv("SERVER_HOST", "example.com")
	os.Setenv("APP_DEBUG", "false")
	defer func() {
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("SERVER_HOST")
		os.Unsetenv("APP_DEBUG")
	}()

	envVars1 := map[string]string{
		"server.port": "SERVER_PORT",
		"server.host": "SERVER_HOST",
	}

	envVars2 := map[string]string{
		"debug": "APP_DEBUG",
	}

	var config TestConfig
	err = NewConfig(&config, configPath, "", envVars1, envVars2)

	assert.NoError(t, err)
	assert.Equal(t, 3000, config.Server.Port)
	assert.Equal(t, "example.com", config.Server.Host)
	assert.False(t, config.Debug)
}

func TestNewConfig_FileNotFound(t *testing.T) {
	var config TestConfig
	err := NewConfig(&config, "/nonexistent/config.yaml", "")

	assert.Error(t, err)

	assert.Contains(t, err.Error(), "Config File")
	assert.Contains(t, err.Error(), "Not Found")
}

func TestNewConfig_InvalidYAMLFormat(t *testing.T) {
	invalidContent := `
server:
  port: 8080
  host: localhost
database:
  url: postgresql://localhost:5432/testdb
  username: testuser
    password: testpass  # Invalid indentation
debug: true
`

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "invalid.yaml")

	err := os.WriteFile(configPath, []byte(invalidContent), 0644)
	require.NoError(t, err)

	var config TestConfig
	err = NewConfig(&config, configPath, "")

	assert.Error(t, err)
}

func TestNewConfig_InvalidJSONFormat(t *testing.T) {
	invalidContent := `{
  "server": {
    "port": 9090,
    "host": "0.0.0.0"
  },
  "database": {
    "url": "mysql://localhost:3306/testdb",
    "username": "jsonuser",
    "password": "jsonpass"
  },
  "debug": false,  // Invalid trailing comma
}`

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "invalid.json")

	err := os.WriteFile(configPath, []byte(invalidContent), 0644)
	require.NoError(t, err)

	var config TestConfig
	err = NewConfig(&config, configPath, "")

	assert.Error(t, err)
}

func TestNewConfig_OverrideFileNotFound(t *testing.T) {
	configContent := `
server:
  port: 8080
  host: localhost
debug: true
`

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	var config TestConfig
	err = NewConfig(&config, configPath, "/nonexistent/override.yaml")

	assert.Error(t, err)
}

func TestNewConfig_EmptyOverridePath(t *testing.T) {
	configContent := `
server:
  port: 8080
  host: localhost
debug: true
`

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	var config TestConfig
	err = NewConfig(&config, configPath, "") // Empty override path

	assert.NoError(t, err)
	assert.Equal(t, 8080, config.Server.Port)
	assert.Equal(t, "localhost", config.Server.Host)
	assert.True(t, config.Debug)
}

func TestNewConfig_UnmarshalError(t *testing.T) {
	// Config with incompatible types
	configContent := `
server:
  port: "not_a_number"  # This should be an integer
  host: localhost
debug: true
`

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	var config TestConfig
	err = NewConfig(&config, configPath, "")

	assert.Error(t, err)
}

func TestNewConfig_TOMLConfig(t *testing.T) {
	configContent := `
debug = true

[server]
port = 7070
host = "toml.host"

[database]
url = "postgresql://toml:5432/db"
username = "tomluser"
password = "tomlpass"
`

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.toml")

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	var config TestConfig
	err = NewConfig(&config, configPath, "")

	assert.NoError(t, err)
	assert.Equal(t, 7070, config.Server.Port)
	assert.Equal(t, "toml.host", config.Server.Host)
	assert.Equal(t, "postgresql://toml:5432/db", config.Database.URL)
	assert.Equal(t, "tomluser", config.Database.Username)
	assert.Equal(t, "tomlpass", config.Database.Password)
	assert.True(t, config.Debug)
}

func TestNewConfig_EnvironmentVariableBindingError(t *testing.T) {
	configContent := `
server:
  port: 8080
debug: true
`

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	envVars := map[string]string{
		"server.port": "",
	}

	var config TestConfig
	err = NewConfig(&config, configPath, "", envVars)

	assert.NoError(t, err)
}

func TestNewConfig_RelativeConfigPath(t *testing.T) {
	configContent := `
server:
  port: 5555
  host: relative.host
debug: false
`

	err := os.WriteFile("test_config.yaml", []byte(configContent), 0644)
	require.NoError(t, err)
	defer os.Remove("test_config.yaml")

	var config TestConfig
	err = NewConfig(&config, "test_config.yaml", "")

	assert.NoError(t, err)
	assert.Equal(t, 5555, config.Server.Port)
	assert.Equal(t, "relative.host", config.Server.Host)
	assert.False(t, config.Debug)
}
