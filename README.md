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

## Prerequisites

- Docker and Docker Compose
- Go 1.23 or later
- Task (https://taskfile.dev)

## Getting Started

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/mp-emailer.git
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

5. Access the application at `http://localhost:8080`

## Development

- To run tests: `task test`
- To lint the code: `task lint`
- To build the application: `task build`

## Deployment

(Add instructions for deploying the application to production)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the [MIT License](LICENSE).
