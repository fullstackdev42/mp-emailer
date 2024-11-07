# TODO

## main.go

### Middleware Registration: 

The middleware for setting the logger and session store could be combined into a single middleware for efficiency, or ideally, these should be set up in a way that they're part of the initial Echo setup or through a custom middleware group.

### Config Handling: 

You're passing config.Config around quite a bit. Consider if this configuration could be injected into structs once rather than passed around, especially for route handlers.

### Error Handling in Middleware: 

The rate limiter middleware doesn't handle the error case. You might want to add error handling there or ensure the middleware knows how to handle rate limit exceedances.

### Session Store: 

There's no error handling when setting the session store in the middleware. This should be checked to ensure the session store can actually be accessed:

```go
e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        if sessionStore != nil {
            c.Set("store", sessionStore)
        } else {
            return errors.New("session store not available")
        }
        return next(c)
    }
})
```

### Server Start and Stop: 

The server starts in a goroutine, which is good for non-blocking start. However, you might want to add some form of timeout for graceful shutdown to prevent indefinite waits.

```go
OnStop: func(ctx context.Context) error {
    shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()
    return e.Shutdown(shutdownCtx)
},
```

### Type Assertions: 

When using context.Context to retrieve values set by middleware, remember that value retrieval can return nil. In your handlers, ensure you handle this scenario properly.


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

### Observations and Suggestions

3. **Time Parsing**
   - Using `shared.ParseDateTime`, which is fine as long as potential parsing errors are handled appropriately. However, ensure that `ParseDateTime` properly handles all possible date/time formats you might encounter.

## Testing Plan

### Priority Areas

#### 1. Core Business Logic
- [ ] Campaign utils and handlers
- [ ] User authentication and management
- [ ] Email sending functionality

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
├── campaign/
│   ├── handler_test.go
│   ├── repository_test.go
│   ├── service_test.go
│   └── lookup_service_test.go
├── user/
│   ├── handler_test.go
│   ├── repository_test.go
│   ├── auth_test.go
│   └── service_test.go
├── email/
│   ├── template_test.go
│   └── service_test.go
├── config/
│   └── config_test.go
└── shared/
    ├── template_test.go
    └── session_test.go
```

### Package Testing Requirements

#### Campaign Package
- [ ] Expand existing tests in `campaign/utils_test.go`
- [ ] Test campaign handlers
- [ ] Test campaign repository methods
- [ ] Test campaign service layer
- [ ] Test representative lookup service

#### User Package
- [ ] User authentication
- [ ] User registration
- [ ] Password hashing/validation
- [ ] User repository methods

#### Email Package
- [ ] Email template rendering
- [ ] Email sending failures
- [ ] Rate limiting
- [ ] Email validation

#### Database Package
- [ ] Database connection handling
- [ ] Query methods
- [ ] Transaction handling
- [ ] Error scenarios

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
- [ ] Use table-driven tests for comprehensive coverage
- [ ] Include edge cases and boundary conditions
- [ ] Test both valid and invalid inputs

#### 2. Mocking
- [ ] Use testify/mock for external dependencies
- [ ] Create mock implementations of interfaces
- [ ] Test interaction between components

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