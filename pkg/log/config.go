// -------------------------------------------------------------------------------------------
// Copyright (c) Team Sorghum. All rights reserved.
// Licensed under the GPL v3 License. See LICENSE in the project root for license information.
// -------------------------------------------------------------------------------------------

package log

// Config defines the log config model.
type Config struct {
	// Type is the type of logger. Currently support "light", "local" and "otel".
	// The "light" logger outputs logs to stderr, the "local" logger outputs logs to stderr and a local file, and the
	// "otel" logger outputs logs to the global open telemetry logger provider.
	Type string `json:"type" yaml:"type" toml:"type" xml:"type" env:"LOG_TYPE" default:"light"`

	// Level is the log level. Possible values are "debug", "info", "warn" and "error".
	// Note that this config option doesn't effect "otel" logger.
	Level string `json:"level" yaml:"level" toml:"level" xml:"level" env:"LOG_LEVEL" default:"debug"`

	// Local is the local log config.
	Local LocalConfig `json:"local" yaml:"local" toml:"local" xml:"local"`

	// OTel is the otel log config.
	OTel OTelConfig `json:"otel" yaml:"otel" toml:"otel" xml:"otel"`
}

// LocalConfig defines the local log config.
type LocalConfig struct {
	// Path is the file to write logs to. Backup log files will be retained in the same directory.
	Path string `json:"path" yaml:"path" toml:"path" xml:"path" env:"LOG_LOCAL_PATH" default:"/tmp/test/log"`

	// MaxSizeMB is the maximum size in megabytes of the log file before it gets rotated.
	MaxSizeMB int `json:"max_size_mb" yaml:"max_size_mb" toml:"max_size_mb" xml:"max_size_mb" env:"LOG_LOCAL_MAX_SIZE_MB" default:"100"` // nolint:lll

	// MaxBackups is the maximum number of old log files to retain.
	MaxBackups int `json:"max_backups" yaml:"max_backups" toml:"max_backups" xml:"max_backups" env:"LOG_LOCAL_MAX_BACKUPS" default:"3"` // nolint:lll
}

// OTelConfig defines the otel log config.
type OTelConfig struct {
	// Name is the logger name, which is most commonly the package name of the code.
	Name string `json:"name" yaml:"name" toml:"name" xml:"name" env:"LOG_OTEL_NAME" default:""`
}
