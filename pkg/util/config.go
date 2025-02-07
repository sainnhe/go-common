package util

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"

	"github.com/teamsorghum/go-common/pkg/log"
	"gopkg.in/yaml.v2"
)

/*
LoadConfig loads config by read config file and environment variables, where environment variables takes precedence over
config file.

Params:
  - path string: The file path. If this argument is empty, only environment variables will be used.
  - typ string: The file type. Possible values are "json" and "yaml". If this argument is empty, only environment
    variables will be used.

Returns:
  - *Config: The config structure.
  - error: The error occurred during the execution. Nil will be returned if no error occurrs.
*/
func LoadConfig[Config any](path string, typ string) (*Config, error) {
	var cfg Config

	if len(path) > 0 && len(typ) > 0 {
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		switch typ {
		case "json":
			err = json.Unmarshal(content, &cfg)
			if err != nil {
				return nil, err
			}
		case "yaml":
			err = yaml.Unmarshal(content, &cfg)
			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("unsupported type: %s", typ)
		}
	}

	// Now override with environment variables.
	overrideWithEnvVars(&cfg)

	// Print log if possible.
	if log.DefaultLogger != nil {
		log.DefaultLogger.Infof("Loaded config: %+v", cfg)
	} else {
		fmt.Printf("Loaded config: %+v\n", cfg)
	}

	return &cfg, nil
}

func overrideWithEnvVars(s interface{}) {
	val := reflect.ValueOf(s)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		// s must be a non-nil pointer to a struct.
		return
	}
	val = val.Elem()
	if val.Kind() != reflect.Struct {
		return
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		structField := typ.Field(i)
		envTag := structField.Tag.Get("env")

		if !field.CanSet() {
			continue
		}

		// Handle environment variable override if envTag is set
		if envTag != "" {
			envVal := os.Getenv(envTag)
			if envVal != "" {
				setField(field, envVal)
				continue
			}
		}

		// Now handle recursive structs or pointers
		switch field.Kind() {
		case reflect.Struct:
			overrideWithEnvVars(field.Addr().Interface())
		case reflect.Ptr:
			if !field.IsNil() {
				elem := field.Elem()
				if elem.Kind() == reflect.Struct {
					overrideWithEnvVars(field.Interface())
				}
			}
		}
	}
}

func setField(field reflect.Value, envVal string) {
	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}
		setFieldValue(field.Elem(), envVal)
	} else {
		setFieldValue(field, envVal)
	}
}

func setFieldValue(field reflect.Value, envVal string) {
	switch field.Kind() {
	case reflect.String:
		field.SetString(envVal)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intVal, err := strconv.ParseInt(envVal, 10, field.Type().Bits())
		if err == nil {
			field.SetInt(intVal)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintVal, err := strconv.ParseUint(envVal, 10, field.Type().Bits())
		if err == nil {
			field.SetUint(uintVal)
		}
	case reflect.Float32, reflect.Float64:
		floatVal, err := strconv.ParseFloat(envVal, field.Type().Bits())
		if err == nil {
			field.SetFloat(floatVal)
		}
	case reflect.Bool:
		boolVal, err := strconv.ParseBool(envVal)
		if err == nil {
			field.SetBool(boolVal)
		}
	}
}
