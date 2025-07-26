# Overlap Avalara

A Go-based REST API service that checks for overlapping time ranges. This application provides a simple HTTP endpoint to determine whether two time ranges overlap with each other.

## Features

- **Time Range Overlap Detection**: Check if two time ranges overlap
- **Multi-Environment Configuration**: Support for local, non-prod, and production environments
- **Structured Logging**: Built-in logging with daily log rotation
- **Error Handling**: Comprehensive error responses with proper HTTP status codes
- **CI/CD Ready**: GitHub Actions workflow for automated testing and building

##  Prerequisites

- Go 1.24.2 or higher
- Make (for using Makefile commands)

##  Project Structure

```
overlap-avalara/
├── Makefile                    # Build automation
├── README.md                   # Project documentation
├── bin/                        # Compiled binaries
├── cmd/
│   └── main.go                # Application entry point
├── config/                     # Configuration files
│   ├── config.go              # Configuration logic
│   ├── local/
│   │   └── server.yml         # Local environment config
│   ├── nonprod/
│   │   └── server.yml         # Non-production config
│   └── prod/
│       └── server.yml         # Production config
├── constants/
│   └── error.go               # Error constants
├── data/
│   └── data_range.go          # Data models
├── internal/
│   ├── api/                   # HTTP handlers
│   │   ├── overlap.go
│   │   └── register.go
│   ├── fx.go                  # Dependency injection
│   └── overlap/
│       └── overlap_service.go # Business logic
├── logger/                    # Logging utilities
├── logs/                      # Log files
├── pkg/                       # Shared packages
│   ├── customerror/
│   ├── error/
│   ├── http/
│   └── response/
├── server/                    # HTTP server setup
└── toolkit/                   # Utility functions
```

## 🛠️ Installation & Setup

### Clone the Repository
```bash
git clone <repository-url>
cd overlap-avalara
```

### Install Dependencies
```bash
go mod tidy
```

### Build the Application
```bash
make build
```

## 🚀 Running the Application

### Using Make (Recommended)
```bash
# Run with local configuration
make run

# Or specify a different environment
make run CONFIG=config/nonprod
```

### Using Go Command Directly
```bash
# Local environment
go run cmd/main.go -config config/local

# Non-production environment
go run cmd/main.go -config config/nonprod

# Production environment
go run cmd/main.go -config config/prod
```

The application will start on `http://localhost:8080` by default (configurable via environment-specific YAML files).

## API Endpoints

### POST /api/v1/overlap-check

Checks if two time ranges overlap.

#### Request Body
```json
{
  "range1": {
    "start": "2025-07-01T10:00:00Z",
    "end": "2025-07-01T12:00:00Z"
  },
  "range2": {
    "start": "2025-07-01T11:00:00Z",
    "end": "2025-07-01T13:00:00Z"
  }
}
```

#### Response
```json
{
  "is_success": true,
  "status_code": 200,
  "data": false
}
```

##  API Testing Examples

### 1. Overlapping Ranges (Expected: `overlap: true`)
```bash
curl -s -X POST http://localhost:8080/api/v1/overlap-check \
  -H "Content-Type: application/json" \
  -d '{
    "range1": {
      "start": "2025-07-01T10:00:00Z",
      "end": "2025-07-01T12:00:00Z"
    },
    "range2": {
      "start": "2025-07-01T11:00:00Z",
      "end": "2025-07-01T13:00:00Z"
    }
  }' | jq
```

**Expected Output:**
```json
{
  "is_success": true,
  "status_code": 200,
  "data": false
}
```

### 2. Non-overlapping Ranges (Expected: `overlap: false`)
```bash
curl -s -X POST http://localhost:8080/api/v1/overlap-check \
  -H "Content-Type: application/json" \
  -d '{
    "range1": {
      "start": "2025-07-01T10:00:00Z",
      "end": "2025-07-01T11:00:00Z"
    },
    "range2": {
      "start": "2025-07-01T11:01:00Z",
      "end": "2025-07-01T12:00:00Z"
    }
  }' | jq
```

**Expected Output:**
```json
{
  "is_success": true,
  "status_code": 200,
  "data": false
}
```

### 3. Invalid Time Format (Expected: Error message)
```bash
curl -s -X POST http://localhost:8080/api/v1/overlap-check \
  -H "Content-Type: application/json" \
  -d '{
    "range1": {
      "start": "foo",
      "end": "bar"
    },
    "range2": {
      "start": "",
      "end": ""
    }
  }' | jq
```

**Expected Output:**
```json
{
  "is_success": false,
  "status_code": 400,
  "error": {
    "message": "Invalid Request"
  }
}
```

##  Development Commands

### Available Make Targets
```bash
# Build the application
make build

# Run the application (local config)
make run

# Run tests with coverage
make test

# Format code
make fmt

# Clean build artifacts
make clean

# Complete CI pipeline (clean, format, test, build)
make ci
```

### Running Tests
```bash
# Run all tests with verbose output and coverage
make test

# Or use go command directly
go test ./... -v -cover
```

### Code Formatting
```bash
# Format all Go files
make fmt
```

##  Build & Deployment

### Local Build
```bash
make build
```
This creates the binary at `bin/overlap-avalara`.

### CI/CD Pipeline
The project includes a GitHub Actions workflow (`.github/workflows/ci.yml`) that:

1.  Checks out the code
2.  Sets up Go 1.24.2
3.  Installs dependencies
4.  Runs tests
5.  Builds the application
6.  Uploads the binary as an artifact
7.  Cleans up build artifacts

## 🌍 Environment Configuration

The application supports three environments:

- **Local** (`config/local/server.yml`): Development environment
- **Non-Prod** (`config/nonprod/server.yml`): Testing/staging environment  
- **Production** (`config/prod/server.yml`): Production environment

Configuration files should contain server settings like:
```yaml
server:
  port: 8080
  host: localhost
```

##  Logging

- Logs are stored in the `logs/` directory
- Daily log rotation (organized by date: `logs/2025-07-26/`)
- Structured logging for better observability

##  Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests (`make test`)
5. Format code (`make fmt`)
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

##  License

[Add your license information here]

##  Troubleshooting

### Common Issues

1. **Port already in use**: Change the port in your environment's `server.yml` file
2. **Binary not found**: Run `make build` to compile the application
3. **Tests failing**: Ensure all dependencies are installed with `go mod tidy`

### Getting Help

If you encounter any issues:
1. Check the logs in the `logs/` directory
2. Verify your configuration files
3. Ensure all dependencies are up to date
4. Run `make ci` to verify everything works end-to-end

---

**Built with  using Go and modern development practices**