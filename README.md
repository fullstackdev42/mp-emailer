# MP Emailer

MP Emailer is a web application that facilitates communication between constituents and their Members of Parliament (MPs). It allows users to create email campaigns and easily contact their representatives on important issues.

## Features

- User registration and authentication
- Create and manage email campaigns
- Find MPs based on postal codes
- Send emails to MPs
- Responsive web design

## Technologies Used

- Go 1.23
- Echo web framework
- MariaDB
- Docker
- Taskfile for task automation
- Tailwind CSS

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

7. Access the application at `http://localhost:8080`

## Development

- To run tests: `task test`
- To lint the code: `task lint`
- To build the application: `task build`
- To run the application with live reloading: `task dev`
- To watch and compile frontend assets: `task frontend-dev`

## Database Management

- To create a new migration: `migrate create -ext sql -dir migrations -seq <migration_name>`
- To apply all migrations: `task migrate-up`
- To rollback one migration: `task migrate-down-one`
- To view current migration version: `task migrate-version`

## Email Testing

The development environment includes Mailpit for email testing. Access the Mailpit interface at `http://localhost:8025`.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the [MIT License](LICENSE).
