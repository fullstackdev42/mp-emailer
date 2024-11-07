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
	EmailProvider               EmailProvider
	SMTPHost                    string
	SMTPPort                    string
	SMTPUsername                string
	SMTPPassword                string
	SMTPFrom                    string
}

// Log is used for logging configuration without sensitive fields
type Log struct {
	*Config
	JWTSecret     string `json:"-"`
	MailgunAPIKey string `json:"-"`
	SessionSecret string `json:"-"`
}

// RequiredEnvVars returns a map of required environment variables and their generation commands
func (c *Config) RequiredEnvVars() map[string]string {
	return map[string]string{
		"SESSION_SECRET": "echo 'SESSION_SECRET='$(openssl rand -base64 32) >> .env",
		"JWT_SECRET":     "echo 'JWT_SECRET='$(openssl rand -base64 32) >> .env",
		"DB_USER":        "echo 'DB_USER=your_database_user' >> .env",
		"DB_NAME":        "echo 'DB_NAME=your_database_name' >> .env",
		"DB_PASSWORD":    "echo 'DB_PASSWORD=your_database_password' >> .env",
	}
}

// EmailProvider represents the type of email service to use
type EmailProvider string

const (
	EmailProviderSMTP    EmailProvider = "smtp"
	EmailProviderMailgun EmailProvider = "mailgun"
)
