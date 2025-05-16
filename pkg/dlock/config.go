package dlock

// Config defines the config model for dlock.
type Config struct {
	// Prefix is the prefix for redis keys. Use different keys in different scenarios to avoid conflicts.
	Prefix string `json:"prefix" yaml:"prefix" toml:"prefix" xml:"prefix" env:"DLOCK_PREFIX" default:"dlock"`

	// ExpireMs x indicates how long before a key expires.
	ExpireMs int64 `json:"expire_ms" yaml:"expire_ms" toml:"expire_ms" xml:"expire_ms" env:"DLOCK_EXPIRE_MS" default:"1000"` // nolint:lll

	// RetryAfterMs indicates how long to wait before retrying.
	RetryAfterMs int64 `json:"retry_after_ms" yaml:"retry_after_ms" toml:"retry_after_ms" xml:"retry_after_ms" env:"DLOCK_RETRY_AFTER_MS" default:"100"` // nolint:lll
}
