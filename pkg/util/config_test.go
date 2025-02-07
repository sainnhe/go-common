package util_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/teamsorghum/go-common/pkg/util"
)

// nolint:gocyclo,paralleltest
func TestLoadConfig(t *testing.T) {
	t.Parallel()
	// Do not use t.Parallel() in t.Run() to avoid environment variable conflicts.

	// Define a sample Config struct with struct tags for JSON, YAML, and env variables.
	type Config struct {
		Name  string `json:"name" yaml:"name" env:"CONFIG_NAME"`
		Port  int    `json:"port" yaml:"port" env:"CONFIG_PORT"`
		Debug bool   `json:"debug" yaml:"debug" env:"CONFIG_DEBUG"`
	}

	// Helper function to create a temporary file with given content.
	createTempFile := func(content string, suffix string) (string, func(), error) {
		file, err := os.CreateTemp("", "config_*"+suffix)
		if err != nil {
			return "", nil, err
		}
		defer file.Close()

		if _, err := file.WriteString(content); err != nil {
			return "", nil, err
		}

		cleanup := func() { os.Remove(file.Name()) }
		return file.Name(), cleanup, nil
	}

	t.Run("Load from JSON file", func(t *testing.T) {
		// Create a temporary JSON config file.
		jsonContent := `
{
	"name": "testapp",
	"port": 8080,
	"debug": false
}
`

		path, cleanup, err := createTempFile(jsonContent, ".json")
		if err != nil {
			t.Fatalf("Failed to create temp JSON file: %v", err)
		}
		defer cleanup()

		cfg, err := util.LoadConfig[Config](path, "json")
		if err != nil {
			t.Errorf("LoadConfig failed: %v", err)
		}

		expected := &Config{
			Name:  "testapp",
			Port:  8080,
			Debug: false,
		}

		if !reflect.DeepEqual(cfg, expected) {
			t.Errorf("Expected config %+v, got %+v", expected, cfg)
		}
	})

	t.Run("Load from YAML file", func(t *testing.T) {
		// Create a temporary YAML config file.
		yamlContent := `
name: yamlapp
port: 7070
debug: true
`

		path, cleanup, err := createTempFile(yamlContent, ".yaml")
		if err != nil {
			t.Fatalf("Failed to create temp YAML file: %v", err)
		}
		defer cleanup()

		cfg, err := util.LoadConfig[Config](path, "yaml")
		if err != nil {
			t.Errorf("LoadConfig failed: %v", err)
		}

		expected := &Config{
			Name:  "yamlapp",
			Port:  7070,
			Debug: true,
		}

		if !reflect.DeepEqual(cfg, expected) {
			t.Errorf("Expected config %+v, got %+v", expected, cfg)
		}
	})

	t.Run("Environment variable overrides", func(t *testing.T) {
		// Set environment variables to override config values.
		os.Setenv("CONFIG_NAME", "envapp")
		os.Setenv("CONFIG_PORT", "9090")
		os.Setenv("CONFIG_DEBUG", "true")

		defer os.Unsetenv("CONFIG_NAME")
		defer os.Unsetenv("CONFIG_PORT")
		defer os.Unsetenv("CONFIG_DEBUG")

		// Create a temporary JSON config file.
		jsonContent := `
{
	"name": "testapp",
	"port": 8080,
	"debug": false
}
`

		path, cleanup, err := createTempFile(jsonContent, ".json")
		if err != nil {
			t.Fatalf("Failed to create temp JSON file: %v", err)
		}
		defer cleanup()

		cfg, err := util.LoadConfig[Config](path, "json")
		if err != nil {
			t.Errorf("LoadConfig failed: %v", err)
		}

		expected := &Config{
			Name:  "envapp",
			Port:  9090,
			Debug: true,
		}

		if !reflect.DeepEqual(cfg, expected) {
			t.Errorf("Expected config %+v, got %+v", expected, cfg)
		}
	})

	t.Run("Load with no file path and type", func(t *testing.T) {
		cfg, err := util.LoadConfig[Config]("", "")
		if err != nil {
			t.Errorf("LoadConfig failed: %v", err)
		}

		expected := &Config{}

		if !reflect.DeepEqual(cfg, expected) {
			t.Errorf("Expected empty config %+v, got %+v", expected, cfg)
		}
	})

	t.Run("Environment variables without config file", func(t *testing.T) {
		// Set environment variables.
		os.Setenv("CONFIG_NAME", "envonly")
		os.Setenv("CONFIG_PORT", "6060")
		os.Setenv("CONFIG_DEBUG", "false")

		defer os.Unsetenv("CONFIG_NAME")
		defer os.Unsetenv("CONFIG_PORT")
		defer os.Unsetenv("CONFIG_DEBUG")

		cfg, err := util.LoadConfig[Config]("", "")
		if err != nil {
			t.Errorf("LoadConfig failed: %v", err)
		}

		expected := &Config{
			Name:  "envonly",
			Port:  6060,
			Debug: false,
		}

		if !reflect.DeepEqual(cfg, expected) {
			t.Errorf("Expected config %+v, got %+v", expected, cfg)
		}
	})

	t.Run("Non-existent file", func(t *testing.T) {
		_, err := util.LoadConfig[Config]("nonexistent.json", "json")
		if err == nil {
			t.Errorf("Expected error for non-existent file, got nil")
		}
	})

	t.Run("Invalid JSON content", func(t *testing.T) {
		invalidJSON := `
{
	"name": "test",
	"port": "invalid"
}
`

		path, cleanup, err := createTempFile(invalidJSON, ".json")
		if err != nil {
			t.Fatalf("Failed to create temp invalid JSON file: %v", err)
		}
		defer cleanup()

		_, err = util.LoadConfig[Config](path, "json")
		if err == nil {
			t.Errorf("Expected error for invalid JSON content, got nil")
		}
	})

	t.Run("Invalid YAML content", func(t *testing.T) {
		invalidYAML := `
name: test
port: invalidPort
debug: yes
`
		path, cleanup, err := createTempFile(invalidYAML, ".yaml")
		if err != nil {
			t.Fatalf("Failed to create temp invalid YAML file: %v", err)
		}
		defer cleanup()

		_, err = util.LoadConfig[Config](path, "yaml")
		if err == nil {
			t.Errorf("Expected error for invalid YAML content, got nil")
		}
	})

	t.Run("Unsupported file type", func(t *testing.T) {
		path, cleanup, err := createTempFile(``, ".txt")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer cleanup()

		_, err = util.LoadConfig[Config](path, "txt")
		if err == nil {
			t.Errorf("Expected error for unsupported file type, got nil")
		}
	})

	t.Run("Nested struct with environment variable override", func(t *testing.T) {
		type Nested struct {
			URL string `json:"url" yaml:"url" env:"NESTED_URL"`
		}

		type ConfigWithNested struct {
			Name   string `json:"name" yaml:"name" env:"CONFIG_NAME"`
			Nested Nested `json:"nested" yaml:"nested"`
		}

		jsonContent := `{
            "name": "testapp",
            "nested": {
                "url": "http://localhost"
            }
        }`

		path, cleanup, err := createTempFile(jsonContent, ".json")
		if err != nil {
			t.Fatalf("Failed to create temp JSON file: %v", err)
		}
		defer cleanup()

		// Set environment variable to override nested struct field.
		os.Setenv("NESTED_URL", "http://env.com")
		defer os.Unsetenv("NESTED_URL")

		cfg, err := util.LoadConfig[ConfigWithNested](path, "json")
		if err != nil {
			t.Errorf("LoadConfig failed: %v", err)
		}

		expected := &ConfigWithNested{
			Name: "testapp",
			Nested: Nested{
				URL: "http://env.com",
			},
		}

		if !reflect.DeepEqual(cfg, expected) {
			t.Errorf("Expected nested config %+v, got %+v", expected, cfg)
		}
	})

	t.Run("Pointer fields with environment variables", func(t *testing.T) {
		type ConfigWithPointer struct {
			Name *string `json:"name" yaml:"name" env:"CONFIG_NAME"`
			Port *int    `json:"port" yaml:"port" env:"CONFIG_PORT"`
		}

		os.Setenv("CONFIG_NAME", "pointerapp")
		os.Setenv("CONFIG_PORT", "5050")
		defer os.Unsetenv("CONFIG_NAME")
		defer os.Unsetenv("CONFIG_PORT")

		cfg, err := util.LoadConfig[ConfigWithPointer]("", "")
		if err != nil {
			t.Errorf("LoadConfig failed: %v", err)
		}

		expectedName := "pointerapp"
		expectedPort := 5050

		if cfg.Name == nil || *cfg.Name != expectedName {
			t.Errorf("Expected Name %v, got %v", expectedName, *cfg.Name)
		}
		if cfg.Port == nil || *cfg.Port != expectedPort {
			t.Errorf("Expected Port %v, got %v", expectedPort, *cfg.Port)
		}
	})

	t.Run("Slice and Map fields", func(t *testing.T) {
		type ConfigWithSlice struct {
			Names []string        `json:"names" yaml:"names"`
			Ports []int           `json:"ports" yaml:"ports"`
			Meta  map[string]int  `json:"meta" yaml:"meta"`
			Flags map[string]bool `json:"flags" yaml:"flags"`
		}

		jsonContent := `{
            "names": ["app1", "app2"],
            "ports": [8000, 8001],
            "meta": {"version": 1, "build": 100},
            "flags": {"debug": true, "verbose": false}
        }`

		path, cleanup, err := createTempFile(jsonContent, ".json")
		if err != nil {
			t.Fatalf("Failed to create temp JSON file: %v", err)
		}
		defer cleanup()

		cfg, err := util.LoadConfig[ConfigWithSlice](path, "json")
		if err != nil {
			t.Errorf("LoadConfig failed: %v", err)
		}

		expected := &ConfigWithSlice{
			Names: []string{"app1", "app2"},
			Ports: []int{8000, 8001},
			Meta:  map[string]int{"version": 1, "build": 100},
			Flags: map[string]bool{"debug": true, "verbose": false},
		}

		if !reflect.DeepEqual(cfg, expected) {
			t.Errorf("Expected config %+v, got %+v", expected, cfg)
		}
	})

	t.Run("Unsupported field kinds (channels, funcs)", func(t *testing.T) {
		type ConfigUnsupported struct {
			Ch chan int `json:"-"`
			Fn func()   `json:"-"`
		}

		cfg, err := util.LoadConfig[ConfigUnsupported]("", "")
		if err != nil {
			t.Errorf("LoadConfig failed: %v", err)
		}

		if cfg.Ch != nil || cfg.Fn != nil {
			t.Errorf("Expected unsupported fields to remain nil, got %+v", cfg)
		}
	})
}
