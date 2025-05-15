package limiter

// Config defines the config model for limiter.
type Config struct {
	// Enable indicates whether to enable limiter.
	Enable bool `json:"enable" yaml:"enable" toml:"enable" xml:"enable" env:"LIMITER_ENABLE" default:"true"`

	// Prefix is the prefix for redis keys. Use different keys in different scenarios to avoid conflicts.
	Prefix string `json:"prefix" yaml:"prefix" toml:"prefix" xml:"prefix" env:"LIMITER_PREFIX" default:"*"`

	// Limit is the limit of request volume within the specified time window.
	Limit int `json:"limit" yaml:"limit" toml:"limit" xml:"limit" env:"LIMITER_LIMIT" default:"1"`

	// WindowMs is the time window for measurement in milliseconds.
	WindowMs int `json:"window_ms" yaml:"window_ms" toml:"window_ms" xml:"window_ms" env:"LIMITER_WINDOW_MS" default:"1000"` // nolint:lll

	// MaxAttempts is the maximum number of attempts.
	// Setting it to 0 will disable peak shaving and degrade the limiter to traditional rate limiter.
	MaxAttempts int `json:"max_attempts" yaml:"max_attempts" toml:"max_attempts" xml:"max_attempts" env:"LIMITER_MAX_ATTEMPTS" default:"0"` // nolint:lll

	// AttemptIntervalMs is the interval between each attempt in milliseconds.
	AttemptIntervalMs int `json:"attempt_interval_ms" yaml:"attempt_interval_ms" toml:"attempt_interval_ms" xml:"attempt_interval_ms" env:"LIMITER_ATTEMPT_INTERVAL_MS" default:"500"` // nolint:lll

	// EnableLog indicates whether to output logs when sleeping and retrying.
	EnableLog bool `json:"enable_log" yaml:"enable_log" toml:"enable_log" xml:"enable_log" env:"LIMITER_ENABLE_LOG" default:"true"` // nolint:lll
}
