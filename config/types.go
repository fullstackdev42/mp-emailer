package config

// Config holds the application's configuration values.
type Config struct {
	AppDebug                    bool
	AppEnv                      Environment
	AppPort                     string
	DBHost                      string
	DBName                      string
	DBPassword                  string
	DBPort                      string
	DBUser                      string
	JWTExpiry                   string
	JWTSecret                   string
	LogFile                     string
	LogLevel                    string
	MailgunAPIKey               string
	MailgunDomain               string
	MailpitHost                 string
	MailpitPort                 string
	MigrationsPath              string
	RepresentativeLookupBaseURL string
	SessionName                 string
	SessionSecret               string
}

// Log is used for logging configuration without sensitive fields
type Log struct {
	*Config
	JWTSecret     string `json:"-"`
	MailgunAPIKey string `json:"-"`
	SessionSecret string `json:"-"`
}
