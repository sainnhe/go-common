package encoding

import (
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"strconv"

	"gopkg.in/yaml.v2"
)

var (
	// ErrLoadConfigNotStruct indicates an error that the Config is not a struct.
	ErrLoadConfigNotStruct = errors.New("config must be a struct")

	// ErrLoadConfigUnsupportedType indicates an error that the type is unsupported.
	ErrLoadConfigUnsupportedType = errors.New("unsupported type")
)

/*
LoadConfig loads config by reading the config content and environment variables.

The Config generic should be a struct and supports 4 struct tags:

 1. "json": Used to mark JSON fields.
 2. "yaml": Used to mark YAML fields.
 3. "env": Used to mark environment variable fields.
 4. "default": Used to mark the default value of a field.

The "env" and "default" tag is parsed using [strconv] for basic data types, and [json.Unmarshal] for arrays, slices,
maps and structs.

The parsing process is as follows:

 1. Initialize a Config struct literal and assign default values to corresponding fields.
 2. Unmarshal the config content to this struct literal. This will override the default values assigned in the previous
    step if such field exists in the config content.
 3. Read the environment variables and assign it to corresponding fields. This will override the values assigned in the
    previous step if such environment variable exists.

Params:
  - content []byte: The config content. For example, you can use [os.ReadFile] to read the content from a local file.
    If this argument is nil or an empty slice, only default values and environment variables will be used.
  - typ [Type]: The config type. If [TypeNil] is passed, only default values and environment variables will be used.

Returns:
  - *Config: The config struct.
  - error: The error occurred during the execution, which may be [ErrLoadConfigNotStruct],
    [ErrLoadConfigUnsupportedType] or other runtime errors.
*/
func LoadConfig[Config any](content []byte, typ Type) (*Config, error) {
	var cfg Config

	// Config must be a struct
	if reflect.ValueOf(cfg).Kind() != reflect.Struct {
		return nil, ErrLoadConfigNotStruct
	}

	// Initialize nil pointers
	initNilPointers(reflect.ValueOf(&cfg).Elem())

	// Override default values
	overrideWithDefaultValues(&cfg)

	// Load config content
	if len(content) > 0 {
		switch typ {
		case TypeNil:
			break
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
			return nil, ErrLoadConfigUnsupportedType
		}
	}

	// Override with environment variables.
	overrideWithEnvVars(&cfg)

	return &cfg, nil
}

// initNilPointers initializes nil pointers recursively.
func initNilPointers(val reflect.Value) {

	switch val.Kind() {
	case reflect.Pointer:
		if val.IsNil() {
			if !val.CanSet() {
				return
			}
			elemType := val.Type().Elem()
			newVal := reflect.New(elemType)
			val.Set(newVal)
			initNilPointers(newVal.Elem())
		} else {
			initNilPointers(val.Elem())
		}
	case reflect.Struct:
		for i := range val.NumField() {
			initNilPointers(val.Field(i))
		}
	case reflect.Array, reflect.Slice:
		for i := range val.Len() {
			initNilPointers(val.Index(i))
		}
	case reflect.Map:
		for _, key := range val.MapKeys() {
			initNilPointers(val.MapIndex(key))
		}
	}
}

func overrideWithDefaultValues(cfg any) {
	cfgVal := reflect.ValueOf(cfg).Elem()

	// Iterate over each field
	for i := range cfgVal.NumField() {
		val := cfgVal.Field(i)
		defaultTag := cfgVal.Type().Field(i).Tag.Get("default")

		// Handle default value override if defaultVal is set
		if len(defaultTag) > 0 {
			if val.Kind() == reflect.Pointer {
				setVal(val.Elem(), defaultTag)
			} else {
				setVal(val, defaultTag)
			}
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
	for i := range cfgVal.NumField() {
		val := cfgVal.Field(i)
		envTag := cfgVal.Type().Field(i).Tag.Get("env")

		// Handle environment variable override if envTag is set
		if len(envTag) != 0 {
			envVal := os.Getenv(envTag)
			if len(envVal) != 0 {
				if val.Kind() == reflect.Pointer {
					setVal(val.Elem(), envVal)
				} else {
					setVal(val, envVal)
				}
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
	if !field.CanSet() {
		return
	}
	switch field.Kind() {
	case reflect.String:
		field.SetString(val)
	case reflect.Bool:
		boolVal, err := strconv.ParseBool(val)
		if err == nil {
			field.SetBool(boolVal)
		}
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
	case reflect.Complex64, reflect.Complex128:
		complexVal, err := strconv.ParseComplex(val, field.Type().Bits())
		if err == nil {
			field.SetComplex(complexVal)
		}
	case reflect.Slice, reflect.Array, reflect.Map, reflect.Struct:
		target := reflect.New(field.Type()).Interface()
		if json.Unmarshal([]byte(val), target) == nil {
			field.Set(reflect.ValueOf(target).Elem())
		}
	}
}
