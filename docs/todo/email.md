# Email Service Tasks

## Core Service
### Current Implementation
- [x] Basic Email Sending
  - [x] Test successful sending
  - [x] Test sending failures
  - [ ] Test retry mechanism
  # Rate limiting moved to security.md

## Template System
### Current Implementation
- [ ] Template Preprocessing
  - [ ] Cache compiled templates
  - [ ] Handle template rendering errors
  - [ ] Implement template versioning

### Testing Requirements
- [ ] Template Rendering Tests
  # See validation.md for template validation requirements

### Rate Limiting
# See security.md for rate limiting implementation

## Provider Integration
### Mailgun Provider
- [x] Basic Integration
  - [x] Test API interaction
  - [x] Test error handling
  - [ ] Test rate limiting
  - [ ] Test webhook handling

### Mailpit Provider (Development)
- [x] Basic Integration
  - [x] Test local sending
  - [x] Test SMTP interaction
  - [ ] Test debugging features

## Implementation Notes
- Each decorator should be independently configurable
- Consider using builder pattern for decorator chain setup
- Implement proper context handling for timeouts/cancellation
- Add appropriate test coverage for each decorator
- Document failure modes and recovery strategies
