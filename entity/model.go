package entity

// Metric structure for URL analytics
type Metric struct {
	Visits int `json:"url_visits"`
}

// Config defines the configuration for TOML file
type Config struct {
	DB database `toml:"database"`
}

type database struct {
	Type     string
	Host     string
	Port     int
	Username string
	Password string
	Database string
}
