// Package loadconfig loads config from a file and environment variables.
package loadconfig

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"

	"github.com/teamsorghum/go-common/pkg/log"
	"gopkg.in/yaml.v2"
)

// Type is the type of the config.
type Type int

const (
	// TypeNil means no type, which results in only environment variables being used to load the config.
	TypeNil Type = 0
	// TypeJSON is the JSON type.
	TypeJSON Type = 1
	// TypeYAML is the YAML type.
	TypeYAML Type = 2
)

/*
Load loads config by read the config content and environment variables.

The Config generic supports 4 struct tags:

 1. "json": Used to mark JSON fields.
 2. "yaml": Used to mark YAML fields.
 3. "env": Used to mark environment variable fields. The value is parsed using strconv.
 4. "default": Used to mark the default value of a field. The value is parsed using strconv.

The parsing process is as follows:

 1. Initialize a Config struct literal and assign default values to corresponding fields.
 2. Unmarshal the config content to this struct literal. This will override the default values assigned in the previous
    step.
 3. Read the environment variables and assign it to corresponding fields. This will override the values assigned in the
    previous step.

Params:
  - content []byte: The config content. For example, you can use os.ReadFile() to read the content from a local file.
    If this argument is nil or an empty slice, only environment variables will be used.
  - typ Type: The config type. If TypeNil is passed, only environment variables will be used.

Returns:
  - *Config: The config struct.
  - error: The error occurred during the execution. Nil will be returned if no error occurrs.
*/
func Load[Config any](content []byte, typ Type) (*Config, error) {
	var cfg Config

	// Config must be a Struct
	if reflect.ValueOf(cfg).Kind() != reflect.Struct {
		return nil, errors.New("config must be a struct")
	}

	// Initialize nil pointers
	initNilPointers(reflect.ValueOf(&cfg).Elem())

	// Override default values
	overrideWithDefaultValues(&cfg)

	// Load config content
	if len(content) > 0 && typ != TypeNil {
		switch typ {
		case TypeJSON:
			err := json.Unmarshal(content, &cfg)
			if err != nil {
				return nil, err
			}
		case TypeYAML:
			err := yaml.Unmarshal(content, &cfg)
			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("unsupported type: %+v", typ)
		}
	}

	// Override with environment variables.
	overrideWithEnvVars(&cfg)

	log.GetDefault().Infof("Loaded config: %+v", cfg)

	return &cfg, nil
}

func initNilPointers(v reflect.Value) {
	switch v.Kind() {
	case reflect.Pointer:
		if v.IsNil() {
			elemType := v.Type().Elem()
			newVal := reflect.New(elemType)
			v.Set(newVal)
			initNilPointers(newVal.Elem())
		} else {
			initNilPointers(v.Elem())
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			initNilPointers(field)
		}
	}
}

func overrideWithDefaultValues(cfg any) {
	cfgVal := reflect.ValueOf(cfg).Elem()

	// Iterate over each field
	for i := 0; i < cfgVal.NumField(); i++ {
		val := cfgVal.Field(i)
		defaultTag := cfgVal.Type().Field(i).Tag.Get("default")

		// Skip if can't set
		if !val.CanSet() {
			continue
		}

		// Handle default value override if defaultTag is set
		if defaultTag != "" {
			if val.Kind() == reflect.Pointer {
				setVal(val.Elem(), defaultTag)
			} else {
				setVal(val, defaultTag)
			}
			continue
		}

		// Now handle recursive structs or pointers
		switch val.Kind() {
		case reflect.Struct:
			overrideWithDefaultValues(val.Addr().Interface())
		case reflect.Pointer:
			if val.Elem().Kind() == reflect.Struct {
				overrideWithDefaultValues(val.Interface())
			}
		}
	}
}

func overrideWithEnvVars(cfg any) {
	cfgVal := reflect.ValueOf(cfg).Elem()

	// Iterate over each field
	for i := 0; i < cfgVal.NumField(); i++ {
		val := cfgVal.Field(i)
		envTag := cfgVal.Type().Field(i).Tag.Get("env")

		// Skip if can't set
		if !val.CanSet() {
			continue
		}

		// Handle environment variable override if envTag is set
		if envTag != "" {
			envVal := os.Getenv(envTag)
			if envVal != "" {
				if val.Kind() == reflect.Pointer {
					setVal(val.Elem(), envVal)
				} else {
					setVal(val, envVal)
				}
				continue
			}
		}

		// Now handle recursive structs or pointers
		switch val.Kind() {
		case reflect.Struct:
			overrideWithEnvVars(val.Addr().Interface())
		case reflect.Pointer:
			if val.Elem().Kind() == reflect.Struct {
				overrideWithEnvVars(val.Interface())
			}
		}
	}
}

func setVal(field reflect.Value, val string) {
	switch field.Kind() {
	case reflect.String:
		field.SetString(val)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intVal, err := strconv.ParseInt(val, 10, field.Type().Bits())
		if err == nil {
			field.SetInt(intVal)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintVal, err := strconv.ParseUint(val, 10, field.Type().Bits())
		if err == nil {
			field.SetUint(uintVal)
		}
	case reflect.Float32, reflect.Float64:
		floatVal, err := strconv.ParseFloat(val, field.Type().Bits())
		if err == nil {
			field.SetFloat(floatVal)
		}
	case reflect.Bool:
		boolVal, err := strconv.ParseBool(val)
		if err == nil {
			field.SetBool(boolVal)
		}
	}
}
