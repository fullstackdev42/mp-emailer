# CLI Implementation Tasks

## Configuration Priority
1. Command line flags (highest)
2. Environment variables
3. Configuration file (Viper)
4. Default values (lowest)

## Core Flags
- [ ] `--config` - Path to configuration file
- [ ] `--db-dsn` - Database connection string
- [ ] `--port` - Server port
- [ ] `--env` - Environment (dev/prod)
- [ ] `--log-level` - Logging level

## Implementation Notes
- Use Cobra for CLI framework
- Implement Viper for config management
- All flags should have:
  - Corresponding env var (e.g., `MP_DB_DSN`)
  - Config file key (e.g., `database.dsn`)
  - Sensible default value
  - Help text explaining all configuration methods

## Example Usage
```shell
# Using flags
mp-emailer --port=8080 --env=dev

# Using env vars
export MP_DB_DSN="user:pass@tcp(localhost:3306)/db"
mp-emailer

# Using config file
mp-emailer --config=/etc/mp-emailer/config.yaml
```

## Testing Requirements
- [ ] Test configuration priority order
- [ ] Test default values
- [ ] Test environment variable loading
- [ ] Test config file parsing
- [ ] Test flag validation 