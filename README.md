# acme-gateway

acme-gateway is the entry point to the Acme platform. Built in Go for performance and simplicity, it authenticates and authorizes incoming requests by verifying JWTs from acme-auth-service, then proxies them to internal services such as acme-tickets-service.

## Features

- **JWT Authentication**: Validates JWTs issued by acme-auth-service
- **Claims Extraction**: Extracts and forwards `tenant_id`, `roles`, and `scopes` from JWT claims
- **CORS Support**: Configurable Cross-Origin Resource Sharing for web applications
- **Rate Limiting**: Token bucket-based rate limiting to protect backend services
- **Request Tracing**: Unique request IDs for tracking requests across services
- **Unified Error Responses**: Standardized error format matching acme-contracts schemas
- **Reverse Proxy**: Routes requests to appropriate backend services
- **Route Mappings**: Simple prefix-based routing (e.g., `/api/tickets` → tickets service)

## Installation

```bash
# Clone the repository
git clone https://github.com/RobertLukenbillIV/acme-gateway.git
cd acme-gateway

# Build the binary
go build -o acme-gateway .
```

## Configuration

Configure the gateway using environment variables. See `.env.example` for available options:

- `PORT`: Server port (default: 8080)
- `JWT_SECRET`: Secret key for validating JWT tokens
- `RATE_LIMIT_PER_SECOND`: Maximum requests per second (default: 100)
- `TICKETS_SERVICE_URL`: URL of the acme-tickets-service backend

Example:
```bash
export PORT=8080
export JWT_SECRET=my-secret-key
export RATE_LIMIT_PER_SECOND=100
export TICKETS_SERVICE_URL=http://localhost:8081
```

## Usage

```bash
# Run the gateway
./acme-gateway
```

The gateway will start on the configured port and listen for incoming requests.

## API Routes

| Gateway Route     | Backend Service        | Description                    |
|-------------------|------------------------|--------------------------------|
| `/api/tickets/*`  | acme-tickets-service   | Tickets management endpoints  |

## Request Flow

1. **Request ID Generation**: Each request receives a unique `X-Request-ID` header
2. **CORS Handling**: CORS headers are added to responses, OPTIONS requests are handled
3. **Rate Limiting**: Requests are checked against the rate limit
4. **JWT Validation**: Authorization header is validated and JWT claims are extracted
5. **Proxy**: Request is forwarded to the appropriate backend service with added headers:
   - `X-Tenant-ID`: Tenant identifier from JWT
   - `X-Roles`: Comma-separated list of user roles
   - `X-Scopes`: Comma-separated list of user scopes

## Error Response Format

All errors follow a unified format matching acme-contracts schemas:

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable error message",
    "status": 401
  },
  "trace_id": "request-id-here"
}
```

Error codes include:
- `UNAUTHORIZED`: Missing or invalid authorization
- `INVALID_TOKEN`: JWT validation failed
- `RATE_LIMIT_EXCEEDED`: Too many requests
- `NOT_FOUND`: Route not found
- `INTERNAL_ERROR`: Server error

## Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...
```

### Project Structure

```
acme-gateway/
├── main.go                          # Application entry point
├── internal/
│   ├── config/                      # Configuration management
│   │   └── config.go
│   ├── errors/                      # Unified error handling
│   │   └── errors.go
│   ├── middleware/                  # HTTP middleware
│   │   ├── chain.go                 # Middleware chaining
│   │   ├── cors.go                  # CORS middleware
│   │   ├── error_handler.go         # Error recovery middleware
│   │   ├── jwt.go                   # JWT validation
│   │   ├── rate_limit.go            # Rate limiting
│   │   └── request_id.go            # Request ID generation
│   └── proxy/                       # Reverse proxy and routing
│       └── router.go
├── go.mod                           # Go module definition
└── README.md                        # This file
```

## Example Request

```bash
# Request with JWT token
curl -H "Authorization: Bearer eyJhbGc..." \
     http://localhost:8080/api/tickets/123

# Response headers include:
# X-Request-ID: unique-request-id
# Access-Control-Allow-Origin: *
```

## License

MIT
