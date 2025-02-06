package trafficlimit

// Config defines the traffic limit config model.
type Config struct {
	// Prefix is the prefix of the key used in cache, which could be used to describe the current business.
	Prefix string `json:"prefix" yaml:"prefix" env:"TrafficLimitPrefix"`
	// RateLimit is the rate limit config.
	RateLimit *RateLimit `json:"rate_limit" yaml:"rate_limit" env:"TrafficLimitRateLimit"`
	// PeakShaving is the peak shaving config.
	PeakShaving *PeakShaving `json:"peak_shaving" yaml:"peak_shaving" env:"TrafficLimitPeakShaving"`
}

// RateLimit defines the rate limit config.
type RateLimit struct {
	// Enable indicates whether rate limit is enabled.
	Enable bool `json:"enable" yaml:"enable" env:"TrafficLimitRateLimitEnable"`
	// QPS is the QPS of rate limit.
	QPS int64 `json:"qps" yaml:"qps" env:"TrafficLimitRateLimitQPS"`
}

// PeakShaving defines the peak shaving config.
type PeakShaving struct {
	// Enable indicates whether peak shaving is enabled.
	Enable bool `json:"enable" yaml:"enable" env:"TrafficLimitPeakShavingEnable"`
	// QPS is the QPS of peak shaving.
	QPS int64 `json:"qps" yaml:"qps" env:"TrafficLimitPeakShavingQPS"`
	// MaxAttempts is the max number of attempts.
	MaxAttempts int `json:"max_attempts" yaml:"max_attempts" env:"TrafficLimitPeakShavingMaxAttempts"`
	// AttemptIntervalMs is the interval between each attempt.
	AttemptIntervalMs int `json:"attempt_interval_ms" yaml:"attempt_interval_ms" env:"TrafficLimitPeakShavingAttemptIntervalMs"` // nolint:lll
}
