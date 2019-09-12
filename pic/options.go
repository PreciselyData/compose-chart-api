package pic

// LogLevel specifies which information will be logged.
type LogLevel int

// LogLevel enumeration.
const (
	LogErrors LogLevel = iota
	LogInfo
)

// Options supplied by the client of the API.
type Options struct {
	LogLevel
	LogFileName string
}

var options Options

// LogInfo determines whether info level logging is enabled.
func (o Options) LogInfo() bool {
	return o.LogLevel == LogInfo
}
