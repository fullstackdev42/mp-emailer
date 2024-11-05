package config

// Environment represents the application runtime environment
type Environment string

const (
	EnvDevelopment Environment = "development"
	EnvStaging     Environment = "staging"
	EnvProduction  Environment = "production"
	EnvTesting     Environment = "testing"
)

// IsValidEnvironment checks if the current environment is valid
func (e Environment) IsValidEnvironment() bool {
	switch e {
	case EnvDevelopment, EnvStaging, EnvProduction, EnvTesting:
		return true
	default:
		return false
	}
}

// String representation of Environment
func (e Environment) String() string {
	return string(e)
}
