package sources

// SourceConfig defines a source
type SourceConfig struct {
	Type     string
	Host     string
	Port     string
	Database string
	Username string
	Password string
}

// Sources contains all available sources
var Sources map[string]Source = make(map[string]Source)
