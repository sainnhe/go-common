package log

// Config defines the log config model.
type Config struct {
	// Type is the type of logger. Currently support "slog".
	Type string `json:"type" yaml:"type" env:"LogType" default:"slog"`
	// Level is the log level. Possible values are: debug, info, warn, error
	Level string `json:"level" yaml:"level" env:"LogLevel" default:"debug"`
	// File is the log file config.
	File *File `json:"file" yaml:"file"`
}

// File defines the log file config.
type File struct {
	// Path is the file to write logs to. Backup log files will be retained in the same directory.
	Path string `json:"path" yaml:"path" env:"LogFilePath" default:"/tmp/test/log"`
	// MaxSizeMB is the maximum size in megabytes of the log file before it gets rotated.
	MaxSizeMB int `json:"max_size_mb" yaml:"max_size_mb" env:"LogFileMaxSizeMB" default:"100"`
	// MaxBackups is the maximum number of old log files to retain.
	MaxBackups int `json:"max_backups" yaml:"max_backups" env:"LogFileMaxBackups" default:"3"`
}
