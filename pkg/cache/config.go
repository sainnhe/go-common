package cache

// Config defines the cache config model.
type Config struct {
	// Host is the host of valkey server.
	Host string `json:"host" yaml:"host" env:"CacheHost" default:"localhost"`
	// Port is the port of valkey server.
	Port int `json:"port" yaml:"port" env:"CachePort" default:"6379"`
	// Username is the username.
	Username string `json:"username" yaml:"username" env:"CacheUsername"`
	// Password is the password.
	Password string `json:"password" yaml:"password" env:"CachePassword"`
}
