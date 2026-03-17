# Bill manager API

A robust and scalable REST API for managing payable bills with comprehensive audit logging. Built with Go, following Domain-Driven Design (DDD) and Clean Architecture principles.

## Overview

Bill manager API is a production-ready application that provides a complete CRUD system for managing bills with the following features:

- **Bill Management**: Create, read, update, and approve bills
- **Audit Logging**: Comprehensive audit trail tracking all operations with user information, IP address, and timestamps
- **JWT Authentication**: Secure endpoints with JWT token validation
- **Transaction Management**: Unit of Work pattern ensures data consistency across multiple repositories
- **Database Migrations**: Automated schema management with golang-migrate
- **Clean Architecture**: Clear separation of concerns with domain, application, infrastructure, and presentation layers

## Technology Stack

- **Language**: Go 1.25+
- **Web Framework**: Fiber v3 (High-performance HTTP server)
- **ORM**: GORM (with PostgreSQL driver)
- **Database**: PostgreSQL 16
- **Dependency Injection**: Uber FX
- **Configuration**: Viper
- **Authentication**: JWT (golang-jwt/jwt v5)
- **Database Migrations**: golang-migrate
- **Containerization**: Docker & Docker Compose

## Prerequisites

- Go 1.25 or higher
- Docker and Docker Compose
- PostgreSQL 16 (or use Docker Compose)

## Getting Started

### 1. Clone the Repository

```bash
git clone <repository-url>
cd Bill manager-api
```

### 2. Setup Environment Variables

Copy the example environment file and configure it:

```bash
cp .env.example .env
```

Default `.env` values:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=Bill manager_user
DB_PASSWORD=Bill manager_password
DB_NAME=Bill manager_db
DB_SSL_MODE=disable

SERVER_PORT=8080
SERVER_HOST=0.0.0.0

JWT_SECRET=your_secret_key_here_change_in_production
JWT_EXPIRATION=24h

ENVIRONMENT=development
```

### 3. Start PostgreSQL with Docker Compose

```bash
docker-compose up -d
```

This will start a PostgreSQL 16 container with the configured credentials.

### 4. Run Database Migrations

Migrations run automatically on application startup. The application will execute all pending migrations from the `migrations/` directory.

### 5. Build and Run the Application

```bash
export PATH=$PATH:/usr/local/go/bin
go run cmd/api/main.go
```

The server will start on `http://localhost:8080`

### 6. Verify the Application

Check the health endpoint:

```bash
curl http://localhost:8080/health
```

Expected response:

```json
{
  "status": "ok"
}
```

## API Endpoints

All endpoints except `/health` require JWT authentication via the `Authorization: Bearer <token>` header.

### Health Check

- **GET** `/health` - Health check endpoint (no authentication required)

### Bills Management

- **POST** `/api/v1/bills` - Create a new bill
- **GET** `/api/v1/bills` - List all bills (supports `?status=pending|approved|rejected|paid` filter)
- **GET** `/api/v1/bills/:id` - Get a specific bill
- **POST** `/api/v1/bills/:id/approve` - Approve a bill
- **GET** `/api/v1/bills/:id/audits` - Get audit logs for a bill

## API Examples

### 1. Create a Bill

```bash
curl -X POST http://localhost:8080/api/v1/bills \
  -H "Authorization: Bearer <jwt-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "description": "Office Supplies",
    "amount": 1500.50,
    "due_date": "2026-04-15"
  }'
```

Response:

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "description": "Office Supplies",
  "amount": 1500.50,
  "due_date": "2026-04-15T00:00:00Z",
  "status": "pending",
  "created_by": "123e4567-e89b-12d3-a456-426614174000",
  "created_at": "2026-03-15T10:30:00Z"
}
```

### 2. List Bills

```bash
curl -X GET http://localhost:8080/api/v1/bills \
  -H "Authorization: Bearer <jwt-token>"
```

### 3. Get a Specific Bill

```bash
curl -X GET http://localhost:8080/api/v1/bills/550e8400-e29b-41d4-a716-446655440000 \
  -H "Authorization: Bearer <jwt-token>"
```

### 4. Approve a Bill

```bash
curl -X POST http://localhost:8080/api/v1/bills/550e8400-e29b-41d4-a716-446655440000/approve \
  -H "Authorization: Bearer <jwt-token>"
```

### 5. Get Bill Audit Logs

```bash
curl -X GET http://localhost:8080/api/v1/bills/550e8400-e29b-41d4-a716-446655440000/audits \
  -H "Authorization: Bearer <jwt-token>"
```

Response:

```json
[
  {
    "id": "660e8400-e29b-41d4-a716-446655440111",
    "bill_id": "550e8400-e29b-41d4-a716-446655440000",
    "action": "created",
    "performed_by": "123e4567-e89b-12d3-a456-426614174000",
    "ip_address": "127.0.0.1",
    "user_agent": "curl/7.68.0",
    "created_at": "2026-03-15T10:30:00Z"
  },
  {
    "id": "660e8400-e29b-41d4-a716-446655440222",
    "bill_id": "550e8400-e29b-41d4-a716-446655440000",
    "action": "approved",
    "performed_by": "223e4567-e89b-12d3-a456-426614174111",
    "ip_address": "192.168.1.100",
    "user_agent": "Mozilla/5.0...",
    "created_at": "2026-03-15T11:45:00Z"
  }
]
```

## Database Schema

### Bills Table

```sql
CREATE TABLE bills (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    description VARCHAR(500) NOT NULL,
    amount DECIMAL(15, 2) NOT NULL,
    due_date DATE NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    created_by UUID NOT NULL,
    approved_by UUID,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

### Bill Audits Table

```sql
CREATE TABLE bill_audits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bill_id UUID NOT NULL REFERENCES bills(id) ON DELETE CASCADE,
    action VARCHAR(50) NOT NULL,
    performed_by UUID NOT NULL,
    ip_address VARCHAR(45),
    user_agent VARCHAR(500),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

## Bill Status Lifecycle

- **pending**: Initial status when a bill is created
- **approved**: Bill has been approved by an authorized user
- **rejected**: Bill has been rejected (not yet implemented in API)
- **paid**: Bill has been marked as paid (not yet implemented in API)

## Audit Actions

The audit system logs the following actions:

- **created**: Bill was created
- **approved**: Bill was approved
- **rejected**: Bill was rejected
- **paid**: Bill was marked as paid

## Development

### Running Tests

```bash
go test ./...
```

### Building the Application

```bash
go build -o Bill manager-api cmd/api/main.go
```

### Docker Build

```bash
docker build -t Bill manager-api .
docker run -p 8080:8080 --env-file .env Bill manager-api
```

## Project Structure Best Practices

This project follows Go best practices:

- **Domain Layer**: Contains pure business logic with no external dependencies
- **Application Layer**: Orchestrates use cases using domain entities
- **Infrastructure Layer**: Handles persistence, configuration, and external integrations
- **Presentation Layer**: Manages HTTP handlers, middleware, and DTOs
- **Dependency Injection**: Uber FX manages all dependencies and lifecycle
- **Configuration**: Viper loads configuration from environment variables and `.env` files
- **Error Handling**: Consistent error handling across all layers
- **Audit Trail**: Complete audit logging of all operations

## Security Considerations

- **JWT Validation**: All protected endpoints validate JWT tokens
- **User Context**: User ID is extracted from JWT and stored in context
- **Audit Logging**: All operations are logged with user ID, IP address, and timestamp
- **Database Transactions**: Unit of Work ensures atomic operations
- **SQL Injection Prevention**: GORM parameterized queries prevent SQL injection

## Troubleshooting

### Database Connection Error

Ensure PostgreSQL is running and credentials in `.env` are correct:

```bash
docker-compose ps
```

### Migration Errors

Check migration files in the `migrations/` directory and ensure they are properly formatted.

### JWT Token Errors

Ensure the JWT token is valid and includes a `user_id` claim:

```bash
# Example JWT payload
{
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "iat": 1234567890,
  "exp": 1234654290
}
```

## Contributing

1. Follow Go code style guidelines
2. Add docstrings to all exported functions
3. Write tests for new features
4. Update README.md for API changes

## License

MIT License - see LICENSE file for details

## Support

For issues and questions, please open an issue in the repository.
