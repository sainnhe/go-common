package limiter

// RateLimitConfig defines the rate limit config.
type RateLimitConfig struct {
	// Enable indicates whether rate limit is enabled.
	Enable bool `json:"enable" yaml:"enable" env:"LimiterRateLimitEnable" default:"true"`

	// Prefix is the prefix for keys, which can be used to describe current business and avoid conflicts.
	Prefix string `json:"prefix" yaml:"prefix" env:"LimiterRateLimitPrefix" default:"*"`

	// Limit is the limit of requests in a given time window.
	Limit int `json:"limit" yaml:"limit" env:"LimiterRateLimitLimit" default:"1"`

	// WindowMs is the time window for measurement in milliseconds.
	WindowMs int `json:"window_ms" yaml:"window_ms" env:"LimiterRateLimitWindowMs" default:"1000"`
}

// PeakShavingConfig defines the peak shaving config.
type PeakShavingConfig struct {
	// Enable indicates whether peak shaving is enabled.
	Enable bool `json:"enable" yaml:"enable" env:"LimiterPeakShavingEnable" default:"true"`

	// Prefix is the prefix for keys, which can be used to describe current business and avoid conflicts.
	Prefix string `json:"prefix" yaml:"prefix" env:"LimiterPeakShavingPrefix" default:"*"`

	// Limit is the limit of requests in a given time window.
	Limit int `json:"limit" yaml:"limit" env:"LimiterPeakShavingLimit" default:"1"`

	// WindowMs is the time window for measurement in milliseconds.
	WindowMs int `json:"window_ms" yaml:"window_ms" env:"LimiterPeakShavingWindowMs" default:"1000"`

	// MaxAttempts is the maximum number of attempts.
	MaxAttempts int `json:"max_attempts" yaml:"max_attempts" env:"LimiterPeakShavingMaxAttempts" default:"10"`

	// AttemptIntervalMs is the interval between each attempt in milliseconds.
	AttemptIntervalMs int `json:"attempt_interval_ms" yaml:"attempt_interval_ms" env:"LimiterPeakShavingAttemptIntervalMs" default:"500"` // nolint:lll
}
