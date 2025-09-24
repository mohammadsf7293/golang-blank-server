# Blank Go Project

A Go project template with HTTP server, MySQL database, and sqlc integration.

## Features

- Standard Go HTTP router
- MySQL database with Docker Compose setup
- SQL code generation using sqlc
- Environment-based configuration
- Basic user management API endpoints

## Prerequisites

- Go 1.21 or later
- Docker and Docker Compose
- sqlc (installed via `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`)

## Project Structure

```
.
├── cmd/
│   └── api/            # Application entrypoint
├── internal/
│   ├── config/         # Configuration package
│   └── db/            # Database and generated sqlc code
├── db/
│   ├── migrations/     # SQL migration files
│   └── queries/        # SQL queries for sqlc
├── docker-compose.yml  # Docker Compose configuration
└── sqlc.yaml          # sqlc configuration
```

## Getting Started

1. Clone the repository:
   ```bash
   git clone https://github.com/mohammadsf7293/blank-go-project.git
   cd blank-go-project
   ```

2. Start the MySQL database:
   ```bash
   docker-compose up -d
   ```

3. Generate database code (if you make changes to queries):
   ```bash
   sqlc generate
   ```

4. Run the application:
   ```bash
   go run cmd/api/main.go
   ```

## API Endpoints

### Health Check
```
GET /api/health
```

### Users
```
GET /api/users?limit=10&offset=0  # List users
POST /api/users                   # Create user
```

Example POST request:
```json
{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
}
```

## Environment Variables

The application can be configured using environment variables:

- `DB_HOST` - Database host (default: "localhost")
- `DB_PORT` - Database port (default: "3306")
- `DB_USER` - Database user (default: "app_user")
- `DB_PASSWORD` - Database password (default: "app_pass")
- `DB_NAME` - Database name (default: "blank_project")
- `SERVER_PORT` - HTTP server port (default: "8080")

## Development

### Adding New Database Migrations

1. Create a new SQL file in `db/migrations/`
2. Add your SQL statements
3. Restart the database container to apply migrations

### Adding New API Endpoints

1. Add new SQL queries in `db/queries/`
2. Generate new code with `sqlc generate`
3. Add new handlers in `cmd/api/main.go`

## Notes

- This is a basic template and should be enhanced with proper error handling, middleware, and security measures for production use
- Password hashing is not implemented in this template
- Consider adding proper logging, metrics, and monitoring for production