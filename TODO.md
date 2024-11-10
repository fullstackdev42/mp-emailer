# TODO

## main.go

### Middleware Registration
- [x] Combine logger and session store middleware for efficiency
  - [x] Create unified middleware function
  - [x] Add proper error handling
  - [x] Add proper logging
  - [x] Consider moving to initial Echo setup

### Config Handling
- [ ] Reduce config.Config passing
  - [ ] Implement dependency injection for route handlers
  - [ ] Use struct injection instead of parameter passing
  - [ ] Consider using context values for config where appropriate

### Error Handling in Middleware
- [x] Add error handling to rate limiter middleware
- [ ] Implement proper error responses
  - [ ] Add structured error types
  - [ ] Add error logging
  - [ ] Add user-friendly error messages
- [ ] Add metrics collection for rate limiting

### Session Store
- [x] Add error handling for nil session store
- [ ] Improve session management
  - [ ] Add session validation
  - [ ] Add session cleanup
  - [ ] Add session security headers
  - [ ] Add session encryption
  - [ ] Add session timeout handling

### Server Start and Stop
- [x] Implement non-blocking server start
- [x] Add graceful shutdown with timeout
- [ ] Enhance shutdown process
  - [ ] Add connection draining
  - [ ] Add cleanup procedures
  - [ ] Add shutdown status logging
  - [ ] Add health check endpoint

### Type Assertions
- [ ] Improve context value handling
  - [ ] Add type assertion error handling
  - [ ] Add default values
  - [ ] Add validation
  - [ ] Add logging for missing values
  - [ ] Consider using strongly typed context values

### Implementation References
- Server startup (see main.go:startServer)
- Middleware registration (see main.go:registerMiddlewares)
- Route registration (see main.go:registerRoutes)
- Session store configuration (see main.go:registerMiddlewares)

### Testing Requirements
- [ ] Test middleware chain
- [ ] Test configuration loading
- [ ] Test server startup/shutdown
- [ ] Test error scenarios
- [ ] Test rate limiting
- [ ] Test session handling

## shared/app.go

### Session Secret Management:

Hardcoding the session secret in newSessionStore might not be ideal. Consider loading this from configuration or environment variables:

```go
func newSessionStore(cfg *config.Config) sessions.Store {
    return sessions.NewCookieStore([]byte(cfg.SessionSecret))
}
```

### Database Connection Attempts:

The retry logic in connectToDB could benefit from exponential backoff instead of a fixed retry interval. This makes it more resilient to network issues:

```go
for retries := 5; retries > 0; retries-- {
    db, err := database.NewDB(dsn, logger)
    if err == nil {
        return db, nil
    }
    sleepDuration := time.Duration(5 * (6 - retries)) * time.Second // Exponential backoff
    logger.Warn("Failed to connect to database, retrying...", "error", err, "retry in", sleepDuration)
    time.Sleep(sleepDuration)
}
```

### Template Loading:

The ProvideTemplates function doesn't handle errors in a very descriptive way. Perhaps adding more context about which template caused the error would be helpful:

```go
tmpl, err := tmpl.ParseFiles(templates...)
if err != nil {
    return nil, fmt.Errorf("failed to parse one or more templates: %w", err)
}
```

### Environment Specific Configuration:

If email.NewMailpitEmailService is for development or testing, consider making this environment-dependent or configurable.

### Testing Considerations:
While there's no test code here, the design does facilitate testing by providing functions for creating instances. However, the use of real database connections and file system operations in ProvideTemplates might complicate testing. Mocking these would be necessary.

### Redundancy in Logger Creation:

There's a possibility of creating multiple loggers if config.Load is called multiple times. Ensure this doesn't happen or that there's a mechanism to reuse or ensure a singleton logger.

### Additional Observations:

The code includes a custom logging decorator for the database, which is good for tracking database operations. Ensure that this decorator does not overly affect performance, especially in high-throughput scenarios.

The use of ** in the filepath pattern for template glob might not be supported in all file systems or Go versions. It's generally safer to use * for recursive matches if supported, or ensure your Go version supports it.
Make sure config.SessionSecret is secure and not hardcoded or exposed in any way in production.

## Email Service Decorators

### Current Implementation Needed
- [ ] Add logging before and after email sending
- [ ] Log any errors that occur
- [ ] Preserve the original email service interface
- [ ] Follow decorator pattern for clean separation of concerns

### Future Enhancements
- [ ] Rate limiting decorator
  - Implement token bucket or sliding window algorithm
  - Configure limits per email domain/recipient
  - Handle rate limit exceeded scenarios

- [ ] Retry decorator with backoff
  - Implement exponential backoff strategy
  - Configure max retry attempts
  - Handle permanent vs temporary failures

- [ ] Metrics collection decorator
  - Track success/failure rates
  - Measure response times
  - Monitor rate limit usage
  - Export Prometheus metrics

- [ ] Email validation decorator
  - Validate email format
  - Check for disposable email domains
  - Implement DNS MX record validation
  - Handle validation failures gracefully

- [ ] Template preprocessing decorator
  - Cache compiled templates
  - Validate template variables
  - Handle template rendering errors
  - Implement template versioning

### Implementation Notes
- Each decorator should be independently configurable
- Consider using builder pattern for decorator chain setup
- Implement proper context handling for timeouts/cancellation
- Add appropriate test coverage for each decorator
- Document failure modes and recovery strategies

## Campaigns

## Testing Plan

### Priority Areas

#### 1. Core Business Logic
- [x] Campaign utils and handlers
- [-] User authentication and management
  - [x] Basic login functionality
  - [ ] User registration
  - [ ] Password reset
  - [ ] Account management
  - [ ] Session handling
- [x] Email sending functionality

#### 2. Infrastructure
- [ ] Database operations
- [ ] Configuration loading
- [ ] Template rendering

#### 3. Integration
- [ ] API endpoints
- [ ] Form handling
- [ ] Session management

### Test Files to Create

```plaintext
tests/
├── campaign/ ✓
│   ├── ✓ handler_test.go
│   ├── ✓ repository_test.go
│   ├── ✓ service_test.go
│   └── ✓ utils_test.go
├── email/ ✓
│   ├── ✓ mailgun_test.go
│   └── ✓ mailpit_test.go
├── database/ ✓
│   ├── ✓ migrations_test.go
│   └── ✓ factories/
│       ├── ✓ campaign_factory.go
│       └── ✓ user_factory.go
├── mocks/ ✓
│   ├── ✓ campaign/
│   ├── ✓ database/
│   ├── ✓ email/
│   ├── ✓ shared/
│   └── ✓ user/
├── user/ ✓
│   └── ✓ service_test.go
├── config/
│   └── config_test.go
└── shared/
    ├── template_test.go
    └── session_test.go
```

### Package Testing Requirements

#### Campaign Package
- [x] Expand existing tests in `campaign/utils_test.go`
- [x] Test campaign handlers
- [x] Test campaign repository methods
- [x] Test campaign service layer
- [ ] Test representative lookup service

#### User Package
- [ ] User authentication
- [ ] User registration
- [ ] Password hashing/validation
- [ ] User repository methods

#### Email Package
- [x] Email sending failures (mailgun_test.go, mailpit_test.go)
- [ ] Email template rendering
- [ ] Rate limiting
- [ ] Email validation

#### Database Package
- [x] Database connection handling
- [x] Migration testing
- [ ] Query methods
- [ ] Transaction handling
- [ ] Error scenarios
- [x] Seeding functionality
- [x] Factory implementations

#### Config Package
- [ ] Environment variable loading
- [ ] Default values
- [ ] Configuration validation
- [ ] Error handling

#### Shared Package
- [ ] Template rendering
- [ ] Custom functions
- [ ] Session management
- [ ] Error handling

### Testing Guidelines

#### 1. Table-Driven Tests
- [x] Use table-driven tests for comprehensive coverage (see migrations_test.go)
- [x] Include edge cases and boundary conditions
- [x] Test both valid and invalid inputs

#### 2. Mocking
- [x] Use testify/mock for external dependencies (see mocks/)
- [x] Create mock implementations of interfaces
- [x] Test interaction between components

#### 3. Error Handling
- [ ] Test error conditions
- [ ] Verify error messages
- [ ] Test error propagation

#### 4. Integration Testing
- [ ] Test critical user flows
- [ ] Test API endpoints
- [ ] Test database interactions

#### 5. Test Coverage
- [ ] Aim for 80% code coverage in critical packages
- [ ] Use `go test -cover` to measure coverage
- [ ] Identify and test edge cases

#### 6. Best Practices
- [ ] Write clear test descriptions
- [ ] Use test helpers for common operations
- [ ] Keep tests maintainable and readable
- [ ] Follow Go testing conventions

### Notes
- Use `testify/assert` for assertions
- Avoid global state in tests
- Use test fixtures where appropriate
- Document complex test scenarios
- Consider adding benchmarks for performance-critical code

### Resources
- Go testing documentation
- Testify documentation
- Echo framework testing guide
- Go testing best practices

### CLI Flags

#### Configuration Priority
1. Command line flags (highest)
2. Environment variables
3. Configuration file (Viper)
4. Default values (lowest)

#### Core Flags
- [ ] `--config` - Path to configuration file
- [ ] `--db-dsn` - Database connection string
- [ ] `--port` - Server port
- [ ] `--env` - Environment (dev/prod)
- [ ] `--log-level` - Logging level

#### Implementation Notes
- Use Cobra for CLI framework
- Implement Viper for config management
- All flags should have:
  - Corresponding env var (e.g., `MP_DB_DSN`)
  - Config file key (e.g., `database.dsn`)
  - Sensible default value
  - Help text explaining all configuration methods

#### Example Usage
```shell
# Using flags
mp-emailer --port=8080 --env=dev

# Using env vars
export MP_DB_DSN="user:pass@tcp(localhost:3306)/db"
mp-emailer

# Using config file
mp-emailer --config=/etc/mp-emailer/config.yaml
```

#### Testing Requirements
- [ ] Test configuration priority order
- [ ] Test default values
- [ ] Test environment variable loading
- [ ] Test config file parsing
- [ ] Test flag validation

## API Authentication

### JWT Implementation
- [ ] Move JWT handling to API layer
  - [ ] Create JWT middleware for API routes
  - [ ] Implement token generation in API handlers
  - [ ] Configure JWT secret and expiry via environment variables
  - [ ] Add refresh token functionality
  - [ ] Implement token revocation

### Testing Requirements
- [ ] Test JWT token generation
- [ ] Test token validation
- [ ] Test token expiry
- [ ] Test invalid token scenarios
- [ ] Test refresh token flow

## Testing
- [x] Implement user service tests
  - [x] Test user login functionality
  - [ ] Test user registration
  - [ ] Test password validation
  - [ ] Test edge cases (empty username/password)
  - [ ] Test password hashing
- [ ] Implement repository tests
- [ ] Implement API handler tests

#### 2. Infrastructure

##### Database Package
- [ ] Connection Management
  - [ ] Test retry mechanism with exponential backoff
  - [ ] Test connection timeouts
  - [ ] Test connection failures
  - [ ] Test successful connections

- [ ] Migration System
  - [x] Test migration execution
  - [x] Test migration failures
  - [x] Test migration rollbacks
  - [ ] Test migration version tracking

- [ ] Seeding System
  - [x] Test user seeder
  - [x] Test campaign seeder
  - [ ] Test data relationships
  - [ ] Test seeding failures
  - [ ] Test data validation

- [ ] Factory System
  - [x] Test user factory generation
  - [x] Test campaign factory generation
  - [ ] Test factory relationships
  - [ ] Test custom factory attributes
  - [ ] Test factory validation rules

##### Config Package
- [ ] Environment Variables
  - [ ] Test required env vars validation
  - [ ] Test default values
  - [ ] Test sensitive data handling
  - [ ] Test config overrides

- [ ] JWT Configuration
  - [ ] Test secret key management
  - [ ] Test token expiration settings
  - [ ] Test token validation rules

- [ ] Email Provider Configuration
  - [ ] Test SMTP settings
  - [ ] Test Mailgun settings
  - [ ] Test provider switching
  - [ ] Test credentials validation

##### Shared Package
- [ ] Template System
  - [ ] Test template loading
  - [ ] Test template parsing
  - [ ] Test custom functions
  - [ ] Test error handling
  - [ ] Test template caching

- [ ] Session Management
  - [ ] Test session store initialization
  - [ ] Test session data persistence
  - [ ] Test session expiration
  - [ ] Test session security

- [ ] JWT Implementation
  - [ ] Test token generation
  - [ ] Test token validation
  - [ ] Test claims handling
  - [ ] Test error scenarios

#### Implementation References
- Database connection retry logic (see shared/app.go:79-93)
- Session store configuration (see shared/app.go:96-98)
- Template rendering system (see shared/app.go:101-124)
- JWT token handling (see shared/jwt.go)
- Configuration management (see config/types.go)

#### Testing Guidelines
- Use table-driven tests for configuration scenarios
- Implement mocks for external dependencies
- Test both success and failure paths
- Ensure proper cleanup in teardown
- Test configuration validation rules
- Verify security-sensitive operations

## API Layer

### Handler Implementation
- [x] Campaign endpoints
  - [x] GET /api/campaigns
  - [x] GET /api/campaign/:id
  - [x] POST /api/campaign
  - [x] PUT /api/campaign/:id
  - [x] DELETE /api/campaign/:id
- [x] User endpoints
  - [x] POST /api/user/register
  - [x] POST /api/user/login
  - [x] GET /api/user/:username

### Authentication & Authorization
- [x] JWT middleware implementation (see api/middleware.go)
- [x] Protected route group setup
- [ ] Rate limiting
- [ ] Request validation
- [ ] Response validation
- [ ] Error handling standardization

### Testing Requirements
- [ ] Handler Tests
  - [ ] Test campaign endpoints
  - [ ] Test user endpoints
  - [ ] Test authentication flows
  - [ ] Test error scenarios
  - [ ] Test input validation
  - [ ] Test response formats

- [ ] Middleware Tests
  - [ ] Test JWT validation
  - [ ] Test authorization failures
  - [ ] Test token expiration
  - [ ] Test invalid tokens
  - [ ] Test missing tokens

### Error Handling
- [ ] Implement consistent error responses
- [ ] Add error codes
- [ ] Add error messages
- [ ] Add validation errors
- [ ] Add logging for errors

### Security
- [ ] Input sanitization
- [ ] CORS configuration
- [ ] Rate limiting
- [ ] Request size limits
- [ ] Security headers

### Documentation
- [ ] API documentation
- [ ] Swagger/OpenAPI specs
- [ ] Authentication documentation
- [ ] Error codes documentation
- [ ] Example requests/responses

### Monitoring & Logging
- [ ] Request logging
- [ ] Error logging
- [ ] Performance metrics
- [ ] Health checks
- [ ] Audit logging

### Code References
- Handler implementation (see api/handler.go)