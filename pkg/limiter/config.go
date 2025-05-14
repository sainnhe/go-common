package limiter

// RateLimitConfig defines the rate limit config.
type RateLimitConfig struct {
	// Enable indicates whether rate limit is enabled.
	Enable bool `json:"enable" yaml:"enable" toml:"enable" xml:"enable" env:"LIMITER_RATE_LIMIT_ENABLE" default:"true"`

	// Prefix is the prefix for keys, which can be used to describe current business and avoid conflicts.
	Prefix string `json:"prefix" yaml:"prefix" toml:"prefix" xml:"prefix" env:"LIMITER_RATE_LIMIT_PREFIX" default:"*"`

	// Limit is the limit of requests in a given time window.
	Limit int `json:"limit" yaml:"limit" toml:"limit" xml:"limit" env:"LIMITER_RATE_LIMIT_LIMIT" default:"1"`

	// WindowMs is the time window for measurement in milliseconds.
	WindowMs int `json:"window_ms" yaml:"window_ms" toml:"window_ms" xml:"window_ms" env:"LIMITER_RATE_LIMIT_WINDOW_MS" default:"1000"` // nolint:lll
}

// PeakShavingConfig defines the peak shaving config.
type PeakShavingConfig struct {
	// Enable indicates whether peak shaving is enabled.
	Enable bool `json:"enable" yaml:"enable" toml:"enable" xml:"enable" env:"LIMITER_PEAK_SHAVING_ENABLE" default:"true"`

	// Prefix is the prefix for keys, which can be used to describe current business and avoid conflicts.
	Prefix string `json:"prefix" yaml:"prefix" toml:"prefix" xml:"prefix" env:"LIMITER_PEAK_SHAVING_PREFIX" default:"*"`

	// Limit is the limit of requests in a given time window.
	Limit int `json:"limit" yaml:"limit" toml:"limit" xml:"limit" env:"LIMITER_PEAK_SHAVING_LIMIT" default:"1"`

	// WindowMs is the time window for measurement in milliseconds.
	WindowMs int `json:"window_ms" yaml:"window_ms" toml:"window_ms" xml:"window_ms" env:"LIMITER_PEAK_SHAVING_WINDOW_MS" default:"1000"` // nolint:lll

	// MaxAttempts is the maximum number of attempts.
	MaxAttempts int `json:"max_attempts" yaml:"max_attempts" toml:"max_attempts" xml:"max_attempts" env:"LIMITER_PEAK_SHAVING_MAX_ATTEMPTS" default:"10"` // nolint:lll

	// AttemptIntervalMs is the interval between each attempt in milliseconds.
	AttemptIntervalMs int `json:"attempt_interval_ms" yaml:"attempt_interval_ms" toml:"attempt_interval_ms" xml:"attempt_interval_ms" env:"LIMITER_PEAK_SHAVING_ATTEMPT_INTERVAL_MS" default:"500"` // nolint:lll
}
