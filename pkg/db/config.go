package db

// Config is the config model for initializing a new database pool.
// The driver and DSN (Data Source Name) can be found in your SQL driver documentation, for example
// [github.com/go-sql-driver/mysql].
type Config struct {
	Driver string `json:"driver,omitempty" yaml:"driver" toml:"driver" xml:"driver"`
	DSN    string `json:"dsn,omitempty" yaml:"dsn" toml:"dsn" xml:"dsn"`
}
