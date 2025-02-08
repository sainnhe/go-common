package loadconfig_test

import (
	"os"
	"reflect"
	"testing"
	"time"

	loadconfig "github.com/teamsorghum/go-common/pkg/load_config"
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

		got, err := loadconfig.Load[Config](path, "json")
		if err != nil {
			t.Errorf("LoadConfig failed: %v", err)
		}

		want := &Config{
			Name:  "testapp",
			Port:  8080,
			Debug: false,
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Want %+v, got %+v", want, got)
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

		got, err := loadconfig.Load[Config](path, "yaml")
		if err != nil {
			t.Errorf("LoadConfig failed: %v", err)
		}

		want := &Config{
			Name:  "yamlapp",
			Port:  7070,
			Debug: true,
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Want %+v, got %+v", want, got)
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

		got, err := loadconfig.Load[Config](path, "json")
		if err != nil {
			t.Errorf("LoadConfig failed: %v", err)
		}

		want := &Config{
			Name:  "envapp",
			Port:  9090,
			Debug: true,
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Want %+v, got %+v", want, got)
		}
	})

	t.Run("Load with no file path and type", func(t *testing.T) {
		got, err := loadconfig.Load[Config]("", "")
		if err != nil {
			t.Errorf("LoadConfig failed: %v", err)
		}

		want := &Config{}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Want config %+v, got %+v", want, got)
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

		got, err := loadconfig.Load[Config]("", "")
		if err != nil {
			t.Errorf("LoadConfig failed: %v", err)
		}

		want := &Config{
			Name:  "envonly",
			Port:  6060,
			Debug: false,
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Want %+v, got %+v", want, got)
		}
	})

	t.Run("Non-existent file", func(t *testing.T) {
		_, err := loadconfig.Load[Config]("nonexistent.json", "json")
		if err == nil {
			t.Errorf("Want error for non-existent file, got nil")
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

		_, err = loadconfig.Load[Config](path, "json")
		if err == nil {
			t.Errorf("Want error for invalid JSON content, got nil")
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

		_, err = loadconfig.Load[Config](path, "yaml")
		if err == nil {
			t.Errorf("Want error for invalid YAML content, got nil")
		}
	})

	t.Run("Unsupported file type", func(t *testing.T) {
		path, cleanup, err := createTempFile(``, ".txt")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer cleanup()

		_, err = loadconfig.Load[Config](path, "txt")
		if err == nil {
			t.Errorf("Want error for unsupported file type, got nil")
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
		os.Setenv("NESTED_URL", "https://www.example.com")
		defer os.Unsetenv("NESTED_URL")

		got, err := loadconfig.Load[ConfigWithNested](path, "json")
		if err != nil {
			t.Errorf("LoadConfig failed: %v", err)
		}

		want := &ConfigWithNested{
			Name: "testapp",
			Nested: Nested{
				URL: "https://www.example.com",
			},
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Want nested config %+v, got %+v", want, got)
		}
	})

	t.Run("Pointer fields with environment variables", func(t *testing.T) {
		type Nested struct {
			URL string `json:"url" yaml:"url" env:"NESTED_URL"`
		}

		type ConfigWithPointer struct {
			Name   *string `json:"name" yaml:"name" env:"CONFIG_NAME"`
			Port   *int    `json:"port" yaml:"port" env:"CONFIG_PORT"`
			Nested *Nested `json:"nested" yaml:"nested"`
		}

		os.Setenv("CONFIG_NAME", "pointerapp")
		os.Setenv("CONFIG_PORT", "5050")
		os.Setenv("NESTED_URL", "https://www.example.com")
		defer os.Unsetenv("CONFIG_NAME")
		defer os.Unsetenv("CONFIG_PORT")
		defer os.Unsetenv("NESTED_URL")

		got, err := loadconfig.Load[ConfigWithPointer]("", "")
		if err != nil {
			t.Errorf("LoadConfig failed: %v", err)
		}

		wantName := "pointerapp"
		wantPort := 5050
		wantNestedURL := "https://www.example.com"

		if got.Name == nil || *got.Name != wantName {
			t.Errorf("Want Name %v, got %v", wantName, *got.Name)
		}
		if got.Port == nil || *got.Port != wantPort {
			t.Errorf("Want Port %v, got %v", wantPort, *got.Port)
		}
		if got.Nested == nil || got.Nested.URL != wantNestedURL {
			t.Errorf("Want nested url %v, got %v", wantNestedURL, got.Nested)
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

		got, err := loadconfig.Load[ConfigWithSlice](path, "json")
		if err != nil {
			t.Errorf("LoadConfig failed: %v", err)
		}

		want := &ConfigWithSlice{
			Names: []string{"app1", "app2"},
			Ports: []int{8000, 8001},
			Meta:  map[string]int{"version": 1, "build": 100},
			Flags: map[string]bool{"debug": true, "verbose": false},
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Want %+v, got %+v", want, got)
		}
	})

	t.Run("Unsupported field kinds (channels, funcs)", func(t *testing.T) {
		type ConfigUnsupported struct {
			Ch chan int `json:"-"`
			Fn func()   `json:"-"`
		}

		got, err := loadconfig.Load[ConfigUnsupported]("", "")
		if err != nil {
			t.Errorf("LoadConfig failed: %v", err)
		}

		if got.Ch != nil || got.Fn != nil {
			t.Errorf("Want unsupported fields to remain nil, got %+v", got)
		}
	})

	t.Run("Default values with no config and env vars", func(t *testing.T) {
		type ConfigWithDefaults struct {
			Name  string `json:"name" yaml:"name" default:"defaultapp"`
			Port  int    `json:"port" yaml:"port" default:"8080"`
			Debug bool   `json:"debug" yaml:"debug" default:"true"`
		}

		got, err := loadconfig.Load[ConfigWithDefaults]("", "")
		if err != nil {
			t.Errorf("LoadConfig failed: %v", err)
		}

		want := &ConfigWithDefaults{
			Name:  "defaultapp",
			Port:  8080,
			Debug: true,
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Want %+v, got %+v", want, got)
		}
	})

	t.Run("Default values overridden by config file", func(t *testing.T) {
		type ConfigWithDefaults struct {
			Name  string `json:"name" yaml:"name" default:"defaultapp"`
			Port  int    `json:"port" yaml:"port" default:"8080"`
			Debug bool   `json:"debug" yaml:"debug" default:"false"`
		}

		yamlContent := `
name: configapp
port: 9090
debug: true
`

		path, cleanup, err := createTempFile(yamlContent, ".yaml")
		if err != nil {
			t.Fatalf("Failed to create temp YAML file: %v", err)
		}
		defer cleanup()

		got, err := loadconfig.Load[ConfigWithDefaults](path, "yaml")
		if err != nil {
			t.Errorf("LoadConfig failed: %v", err)
		}

		want := &ConfigWithDefaults{
			Name:  "configapp",
			Port:  9090,
			Debug: true,
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Want %+v, got %+v", want, got)
		}
	})

	t.Run("Environment variables override defaults and config file", func(t *testing.T) {
		type ConfigWithDefaults struct {
			Name  string `json:"name" yaml:"name" env:"CONFIG_NAME" default:"defaultapp"`
			Port  int    `json:"port" yaml:"port" env:"CONFIG_PORT" default:"8080"`
			Debug bool   `json:"debug" yaml:"debug" env:"CONFIG_DEBUG" default:"false"`
		}

		yamlContent := `
name: configapp
port: 9090
debug: false
`

		path, cleanup, err := createTempFile(yamlContent, ".yaml")
		if err != nil {
			t.Fatalf("Failed to create temp YAML file: %v", err)
		}
		defer cleanup()

		os.Setenv("CONFIG_NAME", "envapp")
		os.Setenv("CONFIG_PORT", "7070")
		os.Setenv("CONFIG_DEBUG", "true")
		defer os.Unsetenv("CONFIG_NAME")
		defer os.Unsetenv("CONFIG_PORT")
		defer os.Unsetenv("CONFIG_DEBUG")

		got, err := loadconfig.Load[ConfigWithDefaults](path, "yaml")
		if err != nil {
			t.Errorf("LoadConfig failed: %v", err)
		}

		want := &ConfigWithDefaults{
			Name:  "envapp",
			Port:  7070,
			Debug: true,
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Want %+v, got %+v", want, got)
		}
	})

	t.Run("Invalid default values lead to zero value", func(t *testing.T) {
		type ConfigWithInvalidDefaults struct {
			Port  int  `json:"port" yaml:"port" default:"notanumber"`
			Debug bool `json:"debug" yaml:"debug" default:"notabool"`
		}

		got, err := loadconfig.Load[ConfigWithInvalidDefaults]("", "")
		if err != nil {
			t.Errorf("LoadConfig failed: %v", err)
		}

		want := &ConfigWithInvalidDefaults{
			Port:  0,     // zero value
			Debug: false, // zero value
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Want %+v, got %+v", want, got)
		}
	})

	t.Run("Unexported fields remain unset", func(t *testing.T) {
		type ConfigWithUnexported struct {
			Name     string `json:"name" yaml:"name"`
			version  string `json:"version" yaml:"version" env:"CONFIG_VERSION"`    // nolint
			unexport int    `json:"unexport" yaml:"unexport" env:"CONFIG_UNEXPORT"` // nolint
		}

		jsonContent := `
{
    "name": "testapp",
    "version": "1.2.3",
    "unexport": 42
}
`

		path, cleanup, err := createTempFile(jsonContent, ".json")
		if err != nil {
			t.Fatalf("Failed to create temp JSON file: %v", err)
		}
		defer cleanup()

		os.Setenv("CONFIG_VERSION", "2.3.4")
		os.Setenv("CONFIG_UNEXPORT", "99")
		defer os.Unsetenv("CONFIG_VERSION")
		defer os.Unsetenv("CONFIG_UNEXPORT")

		got, err := loadconfig.Load[ConfigWithUnexported](path, "json")
		if err != nil {
			t.Errorf("LoadConfig failed: %v", err)
		}

		want := &ConfigWithUnexported{
			Name: "testapp",
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Want %+v, got %+v", want, got)
		}
	})

	t.Run("Invalid environment variable values are ignored", func(t *testing.T) {
		type Config struct {
			Port  int  `json:"port" yaml:"port" env:"CONFIG_PORT"`
			Debug bool `json:"debug" yaml:"debug" env:"CONFIG_DEBUG"`
		}

		os.Setenv("CONFIG_PORT", "notanumber")
		os.Setenv("CONFIG_DEBUG", "notabool")
		defer os.Unsetenv("CONFIG_PORT")
		defer os.Unsetenv("CONFIG_DEBUG")

		// Default values
		yamlContent := `
port: 8080
debug: false
`

		path, cleanup, err := createTempFile(yamlContent, ".yaml")
		if err != nil {
			t.Fatalf("Failed to create temp YAML file: %v", err)
		}
		defer cleanup()

		got, err := loadconfig.Load[Config](path, "yaml")
		if err != nil {
			t.Errorf("LoadConfig failed: %v", err)
		}

		// Since environment variables are invalid, values from config file should be used
		want := &Config{
			Port:  8080,
			Debug: false,
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Want %+v, got %+v", want, got)
		}
	})

	t.Run("Embedded structs", func(t *testing.T) {
		type Nested struct {
			URL string `json:"url" yaml:"url" env:"NESTED_URL" default:"http://defaulturl.com"`
		}

		type Config struct {
			Name string `json:"name" yaml:"name" env:"CONFIG_NAME"`

			Nested
		}

		yamlContent := `
name: testapp
url: http://configurl.com
`

		path, cleanup, err := createTempFile(yamlContent, ".yaml")
		if err != nil {
			t.Fatalf("Failed to create temp YAML file: %v", err)
		}
		defer cleanup()

		os.Setenv("NESTED_URL", "http://envurl.com")
		defer os.Unsetenv("NESTED_URL")

		got, err := loadconfig.Load[Config](path, "yaml")
		if err != nil {
			t.Errorf("LoadConfig failed: %v", err)
		}

		want := &Config{
			Name: "testapp",
			Nested: Nested{
				URL: "http://envurl.com",
			},
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Want %+v, got %+v", want, got)
		}
	})

	t.Run("Nil pointers are initialized", func(t *testing.T) {
		type Config struct {
			Name   *string `json:"name" yaml:"name"`
			Nested *struct {
				URL string `json:"url" yaml:"url" default:"defaulturl"`
			} `json:"nested" yaml:"nested"`
		}

		got, err := loadconfig.Load[Config]("", "")
		if err != nil {
			t.Errorf("LoadConfig failed: %v", err)
		}

		if got.Name == nil {
			t.Errorf("Expected Name pointer to be initialized")
		}
		if got.Nested == nil || got.Nested.URL != "defaulturl" {
			t.Errorf("Expected Nested.URL to be 'defaulturl', got %+v", got.Nested)
		}
	})

	t.Run("Custom types are not handled", func(t *testing.T) {
		type MyDuration time.Duration

		type Config struct {
			Timeout MyDuration `json:"timeout" yaml:"timeout" env:"CONFIG_TIMEOUT" default:"10s"`
		}

		os.Setenv("CONFIG_TIMEOUT", "20s")
		defer os.Unsetenv("CONFIG_TIMEOUT")

		got, err := loadconfig.Load[Config]("", "")
		if err != nil {
			t.Errorf("LoadConfig failed: %v", err)
		}

		// Since custom types are not handled by 'setVal', the value of Timeout remains zero value
		if got.Timeout != MyDuration(0) {
			t.Errorf("Expected Timeout to be zero value, got %v", got.Timeout)
		}
	})

	t.Run("Empty strings in config and env vars", func(t *testing.T) {
		type Config struct {
			Name string `json:"name" yaml:"name" env:"CONFIG_NAME" default:"defaultname"`
		}

		yamlContent := `
name: ""
`

		path, cleanup, err := createTempFile(yamlContent, ".yaml")
		if err != nil {
			t.Fatalf("Failed to create temp YAML file: %v", err)
		}
		defer cleanup()

		os.Setenv("CONFIG_NAME", "")
		defer os.Unsetenv("CONFIG_NAME")

		got, err := loadconfig.Load[Config](path, "yaml")
		if err != nil {
			t.Fatalf("LoadConfig failed: %v", err)
		}

		want := &Config{
			Name: "",
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Expecting empty Name, got %+v", got.Name)
		}
	})

	t.Run("Unsupported field kinds (complex numbers)", func(t *testing.T) {
		type Config struct {
			Complex complex128 `json:"complex" yaml:"complex" env:"CONFIG_COMPLEX"`
		}

		yamlContent := `
complex: "(1+2i)"
`

		path, cleanup, err := createTempFile(yamlContent, ".yaml")
		if err != nil {
			t.Fatalf("Failed to create temp YAML file: %v", err)
		}
		defer cleanup()

		_, err = loadconfig.Load[Config](path, "yaml")
		if err == nil {
			t.Errorf("Expected error when loading unsupported type, got none")
		}
	})

	t.Run("No config, env vars, or defaults", func(t *testing.T) {
		type Config struct {
			Name  string
			Port  int
			Debug bool
		}

		got, err := loadconfig.Load[Config]("", "")
		if err != nil {
			t.Errorf("LoadConfig failed: %v", err)
		}

		want := &Config{
			Name:  "",
			Port:  0,
			Debug: false,
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Want %+v, got %+v", want, got)
		}
	})
}
