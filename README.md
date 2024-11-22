# MP Emailer

MP Emailer is a web application that facilitates communication between constituents and their Members of Parliament (MPs). It allows users to create email campaigns and easily contact their representatives on important issues.

## Features

- User registration and authentication with JWT and session support
- Create and manage email campaigns with rich text editing
- Find MPs based on postal codes using the OpenNorth Represent API
- Send emails via SMTP or Mailgun
- Responsive web design with Tailwind CSS
- Real-time frontend development with Browser-Sync
- Comprehensive logging with Zap
- Database migrations with Goose
- Dependency injection with Uber's fx

## Technologies Used

### Backend
- Go 1.23
- Echo web framework
- MariaDB with GORM
- Uber fx for dependency injection
- Zap for logging
- JWT and session-based authentication
- Mailgun/SMTP for email delivery

### Frontend
- Tailwind CSS
- Quill rich text editor
- Browser-Sync for live reload
- Responsive design

### Development & Testing
- Docker and Docker Compose
- Taskfile for task automation
- Goose for database migrations
- Testify for testing
- Air for live reload
- Mailpit for local email testing
- GolangCI-Lint for code quality

## Prerequisites

- Docker and Docker Compose
- Go 1.23 or later
- Task (https://taskfile.dev)
- Node.js and npm (for Tailwind CSS)

## Getting Started

1. Clone the repository:
   ```
   git clone https://github.com/jonesrussell/mp-emailer.git
   cd mp-emailer
   ```

2. Set up the environment variables:
   Copy the `.env.example` file to `.env` and update the values as needed.

3. Start the development environment:
   ```
   task dev
   ```

4. Run database migrations:
   ```
   task migrate-up
   ```

5. Install frontend dependencies:
   ```
   task frontend-install
   ```

6. Build the frontend assets:
   ```
   task frontend-build
   ```

7. Access the application:
   - Main application: http://localhost:8080
   - Browser-Sync UI: http://localhost:3000
   - Mailpit interface: http://localhost:8025

## Development

### Common Tasks
- Run tests: `task test`
- Lint code: `task lint`
- Build application: `task build`
- Run with live reload: `task dev`
- Watch frontend assets: `task frontend-dev`

### Database Management
- Create migration: `task migrate:create -- <migration_name>`
- Apply migrations: `task migrate:up`
- Rollback last migration: `task migrate:down`
- Reset all migrations: `task migrate:reset` (Warning: This will drop all tables and data!)

Example:
```bash
# Create a new migration
task migrate:create -- add_status_to_users

# Apply all pending migrations
task migrate:up

# Rollback the most recent migration
task migrate:down

# Reset all migrations (interactive prompt)
task migrate:reset
```

### Email Testing
The development environment includes Mailpit for email testing. Access the Mailpit interface at `http://localhost:8025`.

## Configuration

The application uses a layered configuration approach:
1. Default values from `config.yaml`
2. Environment variables
3. `.env` file overrides

Key configuration areas:
- Database settings
- Email provider (SMTP/Mailgun)
- Authentication (JWT/Sessions)
- Feature flags
- Logging
- Rate limiting

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the [MIT License](LICENSE).
