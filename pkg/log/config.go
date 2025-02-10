package log

// Config defines the log config model.
type Config struct {
	// Type is the type of logger. Currently support "light", "slog" and "loki".
	Type string `json:"type" yaml:"type" env:"LogType" default:"light"`
	// Level is the log level. Possible values are "debug", "info", "warn", "error" and "fatal".
	Level string `json:"level" yaml:"level" env:"LogLevel" default:"debug"`
	// Slog is the slog config.
	Slog *Slog `json:"slog" yaml:"slog"`
	// Loki is the loki config.
	Loki *Loki `json:"loki" yaml:"loki"`
}

// Slog defines the slog config.
type Slog struct {
	// Path is the file to write logs to. Backup log files will be retained in the same directory.
	Path string `json:"path" yaml:"path" env:"LogSlogPath" default:"/tmp/test/log"`
	// MaxSizeMB is the maximum size in megabytes of the log file before it gets rotated.
	MaxSizeMB int `json:"max_size_mb" yaml:"max_size_mb" env:"LogSlogMaxSizeMB" default:"100"`
	// MaxBackups is the maximum number of old log files to retain.
	MaxBackups int `json:"max_backups" yaml:"max_backups" env:"LogSlogMaxBackups" default:"3"`
}

// Loki defines the loki config.
type Loki struct {
	// URL is the URL of Loki server endpoint.
	URL string `json:"url" yaml:"url" env:"LogLokiURL" default:"http://localhost:3100/loki/api/v1/push"`
	// TenantID is the tenant id.
	TenantID string `json:"tenant_id" yaml:"tenant_id" env:"LogLokiTenantID" default:"fake"`
	// TimeoutSec is the timeout duration in seconds.
	TimeoutSec int `json:"timeout_sec" yaml:"timeout_sec" env:"LogLokiTimeoutSec" default:"10"`
	// ExternalLabels is the external labels in JSON map format.
	ExternalLabels string `json:"external_labels" yaml:"external_labels" env:"LogLokiExternalLabels" default:"{\"env\": \"dev\"}"` // nolint:lll
}
