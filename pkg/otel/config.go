package otel

// Config defines the config model for otel.
type Config struct {
	// Enable specifies whether to enable OpenTelemetry.
	Enable bool `json:"enable" yaml:"enable" toml:"enable" xml:"enable" env:"OTEL_ENABLE" default:"true"`

	// TimeoutMs specifies the timeout used in opentelemetry in milliseconds.
	TimeoutMs int `json:"timeout_ms" yaml:"timeout_ms" toml:"timeout_ms" xml:"timeout_ms" env:"OTEL_TIMEOUT_MS" default:"5000"` // nolint:lll

	// EnableGzip specifies whether to enable gzip compression.
	EnableGzip bool `json:"enable_gzip" yaml:"enable_gzip" toml:"enable_gzip" xml:"enable_gzip" env:"OTEL_ENABLE_GZIP" default:"true"` // nolint:lll

	// Headers specifies additional headers appended in each requests.
	Headers map[string]string `json:"headers" yaml:"headers" toml:"headers" xml:"headers" env:"OTEL_HEADERS" default:"{}"`

	// Attributes specifies the resource attributes.
	Attributes map[string]string `json:"attributes" yaml:"attributes" toml:"attributes" xml:"attributes" env:"OTEL_ATTRIBUTES" default:"{}"` // nolint:lll

	// Conn is the gRPC connection config.
	Conn ConnConfig `json:"conn" yaml:"conn" toml:"conn" xml:"conn"`

	// Batch is the batch config.
	Batch BatchConfig `json:"batch" yaml:"batch" toml:"batch" xml:"batch"`

	// Trace is the trace config
	Trace TraceConfig `json:"trace" yaml:"trace" toml:"trace" xml:"trace"`

	// Metric is the metric config
	Metric MetricConfig `json:"metric" yaml:"metric" toml:"metric" xml:"metric"`

	// Log is the log config.
	Log LogConfig `json:"log" yaml:"log" toml:"log" xml:"log"`
}

// ConnConfig defines the config model for gRPC connection.
type ConnConfig struct {
	// Host specifies the host of the OTLP gRPC server.
	Host string `json:"host" yaml:"host" toml:"host" xml:"host" env:"OTEL_CONN_HOST" default:"localhost"`

	// Port specifies the port of the OTLP gRPC server.
	Port int `json:"port" yaml:"port" toml:"port" xml:"port" env:"OTEL_CONN_PORT" default:"4317"`

	// EnableTLS specifies whether to enable TLS.
	EnableTLS bool `json:"enable_tls" yaml:"enable_tls" toml:"enable_tls" xml:"enable_tls" env:"OTEL_CONN_ENABLE_TLS" default:"false"` // nolint:lll
}

// BatchConfig defines the config model for batch processing.
type BatchConfig struct {
	// MaxSize is the max size of each batch. Set to 0 to disable batch processing.
	MaxSize int `json:"max_size" yaml:"max_size" toml:"max_size" xml:"max_size" env:"OTEL_BATCH_MAX_SIZE" default:"512"`

	// QueueSize is the size of waiting queue, which should be larger than BatchSize.
	QueueSize int `json:"queue_size" yaml:"queue_size" toml:"queue_size" xml:"queue_size" env:"OTEL_BATCH_QUEUE_SIZE" default:"2048"` // nolint:lll

	// MaxDelayMs is the maximum delay for constructing a batch in milliseconds.
	// Processor will forcefully sends available data if this delay is reached, even if the current batch size does not
	// reach BatchSize.
	MaxDelayMs int `json:"max_delay_ms" yaml:"max_delay_ms" toml:"max_delay_ms" xml:"max_delay_ms" env:"OTEL_BATCH_MAX_DELAY_MS" default:"3000"` // nolint:lll
}

// TraceConfig defines the config model for traces.
type TraceConfig struct {
	// Path is the path of the trace endpoint.
	Path string `json:"path" yaml:"path" toml:"path" xml:"path" env:"OTEL_TRACE_PATH" default:"/v1/traces"`

	// AlwaysSample specifies whether to sample every trace.
	// Be careful about using this sampler in a production application with significant traffic:
	// a new trace will be started and exported for every request.
	AlwaysSample bool `json:"always_sample" yaml:"always_sample" toml:"always_sample" xml:"always_sample" env:"OTEL_TRACE_ALWAYS_SAMPLE" default:"false"` // nolint:lll
}

// MetricConfig defines the config model for metrics.
type MetricConfig struct {
	// Path is the path of the metric endpoint.
	Path string `json:"path" yaml:"path" toml:"path" xml:"path" env:"OTEL_METRIC_PATH" default:"/v1/metrics"`

	// Temporality specifies the temporality selector to be used.
	// Possible values are: "default", "cumulative" or "delta"
	Temporality string `json:"temporality" yaml:"temporality" toml:"temporality" xml:"temporality" env:"OTEL_METRIC_TEMPORALITY" default:"default"` // nolint:lll

	// ReaderIntervalMs specifies the collecting interval of a periodic reader in milliseconds.
	ReaderIntervalMs int `json:"reader_interval_ms" yaml:"reader_interval_ms" toml:"reader_interval_ms" xml:"reader_interval_ms" env:"OTEL_METRIC_READER_INTERVAL_MS" default:"60000"` // nolint:lll
}

// LogConfig defines the config model for logs.
type LogConfig struct {
	// Path is the path of the log endpoint.
	Path string `json:"path" yaml:"path" toml:"path" xml:"path" env:"OTEL_LOG_PATH" default:"/v1/logs"`
}
