package log

// Config defines the log config model.
type Config struct {
	// Type is the type of logger. Currently support slog.
	Type string `json:"type" yaml:"type"`
	// Level is the log level. Possible values are: debug, info, warn, error
	Level string `json:"level" yaml:"level"`
	// FilePath is the file path to store logs.
	FilePath string `json:"file_path" yaml:"file_path"`
}
