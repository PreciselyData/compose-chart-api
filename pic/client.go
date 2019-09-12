package pic

// Client is responsible for creating the chart Builder.
type Client interface {
	NewBuilder(c *Config) Builder
}

// SetClient specifies the implementation.
func SetClient(c Client, o Options) {
	client = c
	options = o
	initLogger(o)
}

var client Client
