package config

// Config holds the application's configuration values.
type Config struct {
	AppDebug                    bool          `env:"APP_DEBUG" envDefault:"false"`
	AppEnv                      Environment   `env:"APP_ENV" envDefault:"development"`
	AppHost                     string        `env:"APP_HOST" envDefault:"localhost"`
	AppPort                     int           `env:"APP_PORT" envDefault:"8080"`
	DBUser                      string        `env:"DB_USER,required"`
	DBPassword                  string        `env:"DB_PASSWORD,required"`
	DBHost                      string        `env:"DB_HOST,required"`
	DBPort                      int           `env:"DB_PORT" envDefault:"3306"`
	DBName                      string        `env:"DB_NAME,required"`
	EmailProvider               EmailProvider `env:"EMAIL_PROVIDER" envDefault:"smtp"`
	JWTExpiry                   string        `env:"JWT_EXPIRY" envDefault:"24h"`
	JWTSecret                   string        `env:"JWT_SECRET,required"`
	LogFile                     string        `env:"LOG_FILE" envDefault:"storage/logs/app.log"`
	LogLevel                    string        `env:"LOG_LEVEL" envDefault:"info"`
	MailgunAPIKey               string        `env:"MAILGUN_API_KEY" envDefault:""`
	MailgunDomain               string        `env:"MAILGUN_DOMAIN" envDefault:""`
	MigrationsPath              string        `env:"MIGRATIONS_PATH" envDefault:"database/migrations"`
	RepresentativeLookupBaseURL string        `env:"REPRESENTATIVE_LOOKUP_BASE_URL" envDefault:"https://represent.opennorth.ca"`
	SessionName                 string        `env:"SESSION_NAME" envDefault:"mp_emailer_session"`
	SessionSecret               string        `env:"SESSION_SECRET,required"`
	SMTPFrom                    string        `env:"SMTP_FROM" envDefault:"noreply@localhost"`
	SMTPHost                    string        `env:"SMTP_HOST" envDefault:"localhost"`
	SMTPPassword                string        `env:"SMTP_PASSWORD" envDefault:""`
	SMTPPort                    int           `env:"SMTP_PORT" envDefault:"587"`
	SMTPUsername                string        `env:"SMTP_USERNAME" envDefault:""`
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
