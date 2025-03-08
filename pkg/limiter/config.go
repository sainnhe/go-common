// -------------------------------------------------------------------------------------------
// Copyright (c) Team Sorghum. All rights reserved.
// Licensed under the GPL v3 License. See LICENSE in the project root for license information.
// -------------------------------------------------------------------------------------------

package limiter

// RateLimitConfig defines the rate limit config.
type RateLimitConfig struct {
	// Enable indicates whether rate limit is enabled.
	Enable bool `json:"enable" yaml:"enable" env:"LIMITER_RATE_LIMIT_ENABLE" default:"true"`

	// Prefix is the prefix for keys, which can be used to describe current business and avoid conflicts.
	Prefix string `json:"prefix" yaml:"prefix" env:"LIMITER_RATE_LIMIT_PREFIX" default:"*"`

	// Limit is the limit of requests in a given time window.
	Limit int `json:"limit" yaml:"limit" env:"LIMITER_RATE_LIMIT_LIMIT" default:"1"`

	// WindowMs is the time window for measurement in milliseconds.
	WindowMs int `json:"window_ms" yaml:"window_ms" env:"LIMITER_RATE_LIMIT_WINDOW_MS" default:"1000"`
}

// PeakShavingConfig defines the peak shaving config.
type PeakShavingConfig struct {
	// Enable indicates whether peak shaving is enabled.
	Enable bool `json:"enable" yaml:"enable" env:"LIMITER_PEAK_SHAVING_ENABLE" default:"true"`

	// Prefix is the prefix for keys, which can be used to describe current business and avoid conflicts.
	Prefix string `json:"prefix" yaml:"prefix" env:"LIMITER_PEAK_SHAVING_PREFIX" default:"*"`

	// Limit is the limit of requests in a given time window.
	Limit int `json:"limit" yaml:"limit" env:"LIMITER_PEAK_SHAVING_LIMIT" default:"1"`

	// WindowMs is the time window for measurement in milliseconds.
	WindowMs int `json:"window_ms" yaml:"window_ms" env:"LIMITER_PEAK_SHAVING_WINDOW_MS" default:"1000"`

	// MaxAttempts is the maximum number of attempts.
	MaxAttempts int `json:"max_attempts" yaml:"max_attempts" env:"LIMITER_PEAK_SHAVING_MAX_ATTEMPTS" default:"10"`

	// AttemptIntervalMs is the interval between each attempt in milliseconds.
	AttemptIntervalMs int `json:"attempt_interval_ms" yaml:"attempt_interval_ms" env:"LIMITER_PEAK_SHAVING_ATTEMPT_INTERVAL_MS" default:"500"` // nolint:lll
}
