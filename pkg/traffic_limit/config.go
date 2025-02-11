package trafficlimit

// RateLimitConfig defines the rate limit config.
type RateLimitConfig struct {
	// Enable indicates whether rate limit is enabled.
	Enable bool `json:"enable" yaml:"enable" env:"TrafficLimitRateLimitEnable" default:"true"`
	// Prefix is the prefix for keys, which can be used to describe current business and avoid conflicts.
	Prefix string `json:"prefix" yaml:"prefix" env:"TrafficLimitRateLimitPrefix" default:"*"`
	// Limit is the limit of requests in a given time window.
	Limit int `json:"limit" yaml:"limit" env:"TrafficLimitRateLimitLimit" default:"1"`
	// WindowMs is the time window for measurement in milliseconds.
	WindowMs int `json:"window_ms" yaml:"window_ms" env:"TrafficLimitRateLimitWindowMs" default:"1000"`
}

// PeakShavingConfig defines the peak shaving config.
type PeakShavingConfig struct {
	// Enable indicates whether peak shaving is enabled.
	Enable bool `json:"enable" yaml:"enable" env:"TrafficLimitPeakShavingEnable" default:"true"`
	// Prefix is the prefix for keys, which can be used to describe current business and avoid conflicts.
	Prefix string `json:"prefix" yaml:"prefix" env:"TrafficLimitPeakShavingPrefix" default:"*"`
	// Limit is the limit of requests in a given time window.
	Limit int `json:"limit" yaml:"limit" env:"TrafficLimitPeakShavingLimit" default:"1"`
	// WindowMs is the time window for measurement in milliseconds.
	WindowMs int `json:"window_ms" yaml:"window_ms" env:"TrafficLimitPeakShavingWindowMs" default:"1000"`
	// MaxAttempts is the max number of attempts.
	MaxAttempts int `json:"max_attempts" yaml:"max_attempts" env:"TrafficLimitPeakShavingMaxAttempts" default:"10"`
	// AttemptIntervalMs is the interval between each attempt.
	AttemptIntervalMs int `json:"attempt_interval_ms" yaml:"attempt_interval_ms" env:"TrafficLimitPeakShavingAttemptIntervalMs" default:"500"` // nolint:lll
}
