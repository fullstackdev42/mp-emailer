# Testing Overview

## Priority Areas

### 1. Core Business Logic
- [x] Campaign utils and handlers
- [-] User authentication and management
  - [x] Basic login functionality
  - [ ] User registration
  - [ ] Password reset
  - [ ] Account management
  - [ ] Session handling
- [x] Email sending functionality

### 2. Infrastructure
- [ ] Database operations
- [ ] Configuration loading
- [ ] Template rendering

### 3. Integration
- [ ] API endpoints
- [ ] Form handling
- [ ] Session management

## Directory Structure
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
└── user/ ✓
    └── ✓ service_test.go
```

## Testing Guidelines
- Use `testify/assert` for assertions
- Avoid global state in tests
- Use test fixtures where appropriate
- Document complex test scenarios
- Consider adding benchmarks for performance-critical code 