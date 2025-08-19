package config_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/juninhoitabh/clob-go/internal/infra/config"
)

func TestGetEnv(t *testing.T) {
	os.Unsetenv("API_HOST")

	cfg := config.LoadConfig()

	assert.Equal(t, "localhost", cfg.ApiHost, "Deve usar o valor padrão quando a variável de ambiente não está definida")

	t.Setenv("API_HOST", "custom-host")

	cfg = config.LoadConfig()

	assert.Equal(t, "custom-host", cfg.ApiHost, "Deve usar o valor da variável de ambiente quando definida")

	os.Unsetenv("API_HOST")
}

func TestLoadConfig(t *testing.T) {
	t.Setenv("API_HOST", "test-host")
	t.Setenv("API_PORT", "8080")
	t.Setenv("ENVIRONMENT", "staging")

	cfg := config.LoadConfig()

	assert.Equal(t, "test-host", cfg.ApiHost)
	assert.Equal(t, "8080", cfg.ApiPort)
	assert.Equal(t, "staging", cfg.Environment)

	os.Unsetenv("API_HOST")
	os.Unsetenv("API_PORT")
	os.Unsetenv("ENVIRONMENT")
}

func TestInit(t *testing.T) {
	t.Setenv("API_HOST", "init-test-host")
	t.Setenv("API_PORT", "9090")
	t.Setenv("ENVIRONMENT", "production")

	config.Init()

	assert.NotNil(t, config.EnvConfigInstance)
	assert.Equal(t, "init-test-host", config.EnvConfigInstance.ApiHost)
	assert.Equal(t, "9090", config.EnvConfigInstance.ApiPort)
	assert.Equal(t, "production", config.EnvConfigInstance.Environment)

	os.Unsetenv("API_HOST")
	os.Unsetenv("API_PORT")
	os.Unsetenv("ENVIRONMENT")
}

func TestInit_UnknownEnvironment(t *testing.T) {
	t.Setenv("ENVIRONMENT", "unknown")

	config.Init()

	assert.Equal(t, "unknown", config.EnvConfigInstance.Environment)

	os.Unsetenv("ENVIRONMENT")
}

func TestConfig_DefaultValues(t *testing.T) {
	os.Unsetenv("API_HOST")
	os.Unsetenv("API_PORT")
	os.Unsetenv("ENVIRONMENT")

	cfg := config.LoadConfig()

	assert.Equal(t, "localhost", cfg.ApiHost)
	assert.Equal(t, "3000", cfg.ApiPort)
	assert.Equal(t, "development", cfg.Environment)
}
