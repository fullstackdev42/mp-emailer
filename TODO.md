# TODO

## Campaigns

### Observations and Suggestions

1. **Error Handling**
   - Standardize the errors returned. For example, instead of always wrapping errors with `fmt.Errorf`, consider using custom error types or constants for common errors like `ErrCampaignNotFound`.

2. **Closing Rows**
   - Close rows properly. A `defer` statement for closing rows should ideally follow immediately after checking for errors from `db.Query`.

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