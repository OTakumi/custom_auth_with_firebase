# Backend Server

Go-based backend API implementing passwordless OTP authentication with Firebase.

## Tech Stack

- Go 1.25.5
- Gin 1.11 (HTTP framework)
- Firebase Admin SDK v4.18 (Auth + Firestore)

## Project Structure

```
server/
├── cmd/api/main.go              # Entry point
├── internal/
│   ├── config/                  # Environment config
│   ├── domain/                  # Entities, VOs, interfaces
│   ├── usecase/                 # Business logic
│   ├── infrastructure/          # Firebase, Firestore, email
│   └── interface/               # Handlers, middleware, router
├── tests/                       # Integration tests
├── Dockerfile
├── Makefile
└── go.mod
```

## Architecture

Clean Architecture with 4-layer separation:

- **Domain**: Entities, value objects, repository interfaces
- **Use Case**: Business logic (OTP service, Auth service)
- **Infrastructure**: Firebase, Firestore, email sender
- **Interface**: HTTP handlers, middleware, router

## Security Features

- **Timing Attack Prevention**: Constant-time OTP comparison
- **Email Enumeration Prevention**: Generic error messages
- **Brute Force Prevention**: 3 attempts + rate limiting (5 req/min)
- **OTP Security**: Secure random generation, 5-minute expiration
- **IP Privacy**: SHA-256 hashing (GDPR compliant)
- **CORS**: Environment-based origin whitelist

## Quick Start

### 1. Start Firebase Emulator

From project root:

```bash
docker compose up -d firebase
```

### 2. Run Server

**Local:**

```bash
cd server
GOOGLE_CLOUD_PROJECT=demo-project \
FIRESTORE_EMULATOR_HOST=localhost:8080 \
FIREBASE_AUTH_EMULATOR_HOST=localhost:9099 \
go run ./cmd/api/main.go
```

**Docker:**

```bash
docker compose up -d app
```

**API Server**: <http://localhost:8000>

## Development

### Make Commands

```bash
make build    # Lint + format + build
make test     # Run tests
make lint     # Lint only
make format   # Format only
make clean    # Clean build artifacts
```

### Environment Variables

**Development (Emulator):**

```bash
GOOGLE_CLOUD_PROJECT=demo-project
FIRESTORE_EMULATOR_HOST=localhost:8080
FIREBASE_AUTH_EMULATOR_HOST=localhost:9099
PORT=8000                                    # Optional, default: 8000
ENV=development
```

**Production:**

```bash
PORT=8000
ENV=production
ALLOWED_ORIGINS=https://yourdomain.com       # Comma-separated
RATE_LIMIT_REQUESTS_PER_MINUTE=5            # Optional, default: 5
```

## API Endpoints

### `POST /auth/otp`

Request OTP for email address.

**Request:**

```json
{"email": "user@example.com"}
```

**Response (200):**

```json
{"message": "OTP sent successfully."}
```

**Dev Mode:** OTP printed to console

```
OTP Code (Development): 123456
```

### `POST /auth/verify`

Verify OTP and get custom token.

**Request:**

```json
{"email": "user@example.com", "otp": "123456"}
```

**Response (200):**

```json
{"token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."}
```

### `GET /health`

Health check endpoint.

## Testing

```bash
# Run all tests (requires emulator)
FIRESTORE_EMULATOR_HOST=localhost:8080 go test -v ./tests/...

# Coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

## Troubleshooting

**Emulator connection issues:**

```bash
docker compose restart firebase
docker compose logs firebase
```

**Build errors:**

```bash
go mod tidy
make clean && make build
```

See [main README](../README.md) for more details.
