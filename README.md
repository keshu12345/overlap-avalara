# Overlap Avalara

A Go-based REST API service that checks for overlapping time ranges. This application provides a simple HTTP endpoint to determine whether two time ranges overlap with each other.

## Features

- **Time Range Overlap Detection**: Check if two time ranges overlap
- **Multi-Environment Configuration**: Support for local, non-prod, and production environments
- **Structured Logging**: Built-in logging with daily log rotation
- **Error Handling**: Comprehensive error responses with proper HTTP status codes
- **Docker Support**: Containerized deployment with Docker and Docker Compose
- **CI/CD Pipeline**: GitHub Actions workflow for automated testing, building, and Docker image publishing
- **Jenkins Integration**: Automated builds triggered by GitHub webhooks
- **Development Tools**: ngrok integration for local development and webhook testing

## Prerequisites

- Go 1.24.2 or higher
- Make (for using Makefile commands)
- Docker and Docker Compose (for containerized deployment)
- ngrok (for local development and webhook testing)
- Jenkins (for CI/CD pipeline)

## Project Structure

```
overlap-avalara/
├── Makefile                    # Build automation
├── README.md                   # Project documentation
├── Dockerfile                  # Docker configuration
├── docker-compose.yml          # Docker Compose setup
├── .github/
│   └── workflows/
│       └── ci.yml             # GitHub Actions CI/CD pipeline
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

## Installation & Setup

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

##  Running the Application

### Method 1: Using Make (Recommended)
```bash
# Run with local configuration
make run

# Or specify a different environment
make run CONFIG=config/nonprod
```

### Method 2: Using Go Command Directly
```bash
# Local environment
go run cmd/main.go -config config/local

# Non-production environment
go run cmd/main.go -config config/nonprod

# Production environment
go run cmd/main.go -config config/prod
```

### Method 3: Using Docker
```bash
# Build and run with Docker
docker build -t overlap-avalara .
docker run -p 8081:8081 overlap-avalara

# Or use the published image from Docker Hub
docker run -p 8081:8081 210423/overlap-avalara:latest
```

### Method 4: Using Docker Compose
```bash
# Start the application and any dependencies
docker-compose up

# Run in detached mode
docker-compose up -d

# Stop the application
docker-compose down
```

The application will start on `http://localhost:8081` by default (configurable via environment-specific YAML files).

## Docker Configuration

### Dockerfile
The application includes a multi-stage Dockerfile for optimized builds:

```dockerfile
FROM golang:1.24.2-alpine AS builder

RUN apk add --no-cache git

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 \
    go build -trimpath -o overlap-avalara ./cmd/main.go

FROM alpine:latest

# Install certs (for outbound HTTPS calls, if any)
RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /src/overlap-avalara .

COPY config ./config

EXPOSE 8081

ENTRYPOINT ["./overlap-avalara", "-config", "./config/local"]
```

### Docker Compose
```yaml
version: '3.9'

services:
  overlap-avalara:
    image: 210423/overlap-avalara:latest
    container_name: overlap-avalara
    ports:
      - '8081:8081'
    volumes:
      - ./config:/app/config:ro
    restart: unless-stopped
```

### Docker Hub Image
The application is available on Docker Hub:
- **Image**: `210423/overlap-avalara:latest`
- **Pull Command**: `docker pull 210423/overlap-avalara:latest`

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

## API Testing Examples

### 1. Overlapping Ranges (Expected: `overlap: true`)
```bash
curl -s -X POST http://localhost:8081/api/v1/overlap-check \
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

### 2. Non-overlapping Ranges (Expected: `overlap: false`)
```bash
curl -s -X POST http://localhost:8081/api/v1/overlap-check \
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

### 3. Invalid Time Format (Expected: Error message)
```bash
curl -s -X POST http://localhost:8081/api/v1/overlap-check \
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

## Development Commands

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

#### Actual Makefile Configuration
```makefile
APP_NAME := overlap-avalara
OUT_DIR := bin
BIN := $(OUT_DIR)/$(APP_NAME)
CONFIG := config/local

build:
	@echo "🔨 Building $(APP_NAME)..."
	@mkdir -p $(OUT_DIR)
	@go build -o $(BIN) ./cmd/main.go

run:
	@echo "Running $(APP_NAME) with config=$(CONFIG)"
	@go run ./cmd/main.go -config $(CONFIG)

test:
	@echo "Running tests..."
	@go test ./... -v -cover

fmt:
	@echo "Formatting code..."
	@gofmt -s -w .

clean:
	@echo "Cleaning artifacts..."
	@rm -rf $(OUT_DIR)

ci: clean fmt test build
	@echo "CI pipeline complete"
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

## CI/CD Pipeline

### GitHub Actions Workflow

The project includes a comprehensive GitHub Actions workflow (`.github/workflows/ci.yml`) that:

1. **Code Quality Checks**
   - Checkout code
   - Setup Go 1.24.2
   - Format code validation
   - Run tests with coverage

2. **Build Process**
   - Build application binary
   - Build Docker image
   - Tag with commit SHA and 'latest'

3. **Docker Hub Deployment**
   - Login to Docker Hub
   - Push images with multiple tags
   - Update latest tag

4. **Artifact Management**
   - Upload build artifacts
   - Clean up temporary files

#### Workflow Triggers
- **Push** to `main` branch
- **Pull Request** to `main` branch
- **Manual trigger** via GitHub UI

#### Environment Variables Required
```bash
# Add these secrets in GitHub repository settings
DOCKER_HUB_USERNAME=210423
DOCKER_HUB_ACCESS_TOKEN=your_docker_hub_token
```

### Sample GitHub Actions Configuration
```yaml
name: CI

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  test-and-build:
    runs-on: ubuntu-latest

    steps:
      - name: Checkout code
        uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: 1.24.2

      - name: Install dependencies
        run: go mod tidy

      - name: Run tests
        run: make test

      - name: Build artifact
        run: make build

      - name: List bin contents
        run: ls -lh bin || echo " bin/ not found"

      - name: Upload binary
        uses: actions/upload-artifact@v4
        with:
          name: overlap-avalara
          path: bin/overlap-avalara

      - name:  Clean bin folder (explicit rm fallback)
        run: |
          echo "Cleaning bin manually..."
          rm -rf bin
          echo "bin cleaned"
```

## 🔗 Jenkins Integration & GitHub Webhooks

### Jenkins Setup

#### Required Jenkins Plugins
- GitHub Integration Plugin
- Docker Pipeline Plugin
- Generic Webhook Trigger Plugin

#### Jenkins Pipeline Configuration
Create a new Pipeline job with the following configuration:

```groovy
pipeline {
  agent any

  environment {
     PATH = "/usr/local/bin:/usr/local/go/bin:${env.PATH}"
  }

  stages {
    stage('Code') {
      steps {
        echo "→ Cloning code"
        git url: 'https://github.com/keshu12345/overlap-avalara', branch: 'main'
      }
    }


    stage('Docker Health Check') {
      steps {
        echo "→ Checking Docker status"
        script {
          def dockerRunning = false
          try {
            sh 'docker info > /dev/null 2>&1'
            dockerRunning = true
            echo "docker daemon is running"
          } catch (Exception e) {
            echo "Docker daemon is not running"
            echo "Please start Docker Desktop and wait for it to fully initialize"
            echo "You can verify Docker is ready by running: docker info"
            error "Docker daemon is not accessible. Pipeline stopped."
          }
        }
      }
    }
    stage('Build') {
      steps {
        echo "→ Building Docker image"
        sh 'which docker'
        sh 'docker version'
        
        // Login to Docker Hub to avoid 429 rate limit errors
        withCredentials([usernamePassword(credentialsId: 'credID', 
                                          passwordVariable: 'DOCKER_PASSWORD', 
                                          usernameVariable: 'DOCKER_USERNAME')]) {
          sh '''
            echo "→ Logging into Docker Hub"
            echo $DOCKER_PASSWORD | docker login -u $DOCKER_USERNAME --password-stdin
            
            echo "→ Pulling base images"
            docker pull golang:1.24.2-alpine
            docker pull alpine:latest
            
            echo "→ Building Docker image"
            docker build -t overlap-avalara:latest .
          '''
        }
      }
    }

  stage('Push Image') {
      steps {
        echo "→ Pushing Docker image to Docker Hub"
        withCredentials([usernamePassword(credentialsId: 'credID', 
                                          passwordVariable: 'DOCKER_PASSWORD', 
                                          usernameVariable: 'DOCKER_USERNAME')]) {
          sh '''
            echo "→ Logging into Docker Hub for push"
            echo $DOCKER_PASSWORD | docker login -u $DOCKER_USERNAME --password-stdin
            
            echo "→ Tagging image for Docker Hub"
            docker tag overlap-avalara:latest $DOCKER_USERNAME/overlap-avalara:latest
            
            echo "→ Pushing images to Docker Hub"
            docker image tag overlap-avalara overlap-avalara:latest
            docker push $DOCKER_USERNAME/overlap-avalara:latest
            
            echo " successfully pushed images:"
            echo "  - $DOCKER_USERNAME/overlap-avalara:latest"
            
            echo "→ Listing pushed images"
            docker images | grep overlap-avalara
          '''
        }
      }
    }
    
        stage('Test') {
          steps {
            echo "→ Running Go unit tests with coverage"
            sh 'go version'
            sh 'go mod download'
            sh 'go test -v -coverprofile=coverage.out ./...'
          }
          post {
            always {
              archiveArtifacts artifacts: 'coverage.out', fingerprint: true
            }
          }
        }

    stage('Deploy') {
      steps {
        echo "→ Deploying container"
        sh 'docker compose up -d --remove-orphans'
        echo "docker compose running and up in port :8081"
      }
    }
  }
}

```

### GitHub Webhook Configuration

#### Setup Steps

1. **Jenkins Configuration**
   - Go to your Jenkins job configuration
   - Enable "Generic Webhook Trigger"
   - Set token: `overlap-avalara-webhook-token`

2. **GitHub Repository Settings**
   - Navigate to Settings → Webhooks
   - Click "Add webhook"
   - Configure webhook:
     ```
     Payload URL: https://your-ngrok-url.ngrok.io/generic-webhook-trigger/invoke?token=overlap-avalara-webhook-token
     Content type: application/json
     Secret: (optional)
     Events: Just the push event
     Active: Checked
     ```

3. **Webhook Payload Processing**
   ```json
   Post content parameters:
   - Variable: ref
   - Expression: $.ref
   
   - Variable: repository_name
   - Expression: $.repository.name
   
   - Variable: pusher_name
   - Expression: $.pusher.name
   
   Optional filters:
   - Expression: $.ref
   - Text: refs/heads/main
   ```

##  ngrok Setup for Local Development

### Installation & Setup

#### Install ngrok
```bash
# Download and install from https://ngrok.com/download
# Or use package managers:

# macOS
brew install ngrok

# Linux
curl -s https://ngrok-agent.s3.amazonaws.com/ngrok.asc | sudo tee /etc/apt/trusted.gpg.d/ngrok.asc
echo "deb https://ngrok-agent.s3.amazonaws.com buster main" | sudo tee /etc/apt/sources.list.d/ngrok.list
sudo apt update && sudo apt install ngrok
```

#### Configure Authentication
```bash
# Get auth token from https://dashboard.ngrok.com/get-started/your-authtoken
ngrok authtoken YOUR_AUTH_TOKEN
```

### Usage for Development

#### Expose Local Application
```bash
# Start your application locally
make run

# In another terminal, expose it via ngrok
ngrok http 8081

# Note the generated HTTPS URL (e.g., https://abc123.ngrok.io)
```

#### Expose Jenkins for Webhooks
```bash
# Expose Jenkins (typically on port 8080)
ngrok http 8080

# Use the HTTPS URL for GitHub webhook configuration
# Example: https://5cba9608fb1d.ngrok-free.app/generic-webhook-trigger/invoke?token=your-token
```

#### Configuration File (Optional)
Create `~/.ngrok2/ngrok.yml`:
```yaml
version: "2"
authtoken: YOUR_AUTH_TOKEN
tunnels:
  app:
    proto: http
    addr: 8081
    subdomain: overlap-avalara  # requires paid plan
  jenkins:
    proto: http
    addr: 8080
    basic_auth:
      - "admin:password"
```

Start tunnels:
```bash
ngrok start app jenkins
```

### Development Workflow with ngrok

1. **Start local development**:
   ```bash
   make run  # Start application on localhost:8081
   ```

2. **Expose via ngrok**:
   ```bash
   ngrok http 8081  # Creates public tunnel
   ```

3. **Test webhook integration**:
   ```bash
   # Update GitHub webhook with ngrok URL
   # Make a commit to trigger the webhook
   # Check Jenkins job execution
   ```

4. **Monitor requests**:
   - Access ngrok web interface: http://localhost:4040
   - View request/response details
   - Replay requests for testing



##  Monitoring & Logging

### Application Logs
- **Location**: `logs/` directory
- **Rotation**: Daily log rotation (organized by date: `logs/2025-07-26/`)
- **Format**: Structured JSON logging in production
- **Levels**: DEBUG, INFO, WARN, ERROR

### Docker Container Logs
```bash
# View logs from Docker container
docker logs overlap-avalara

# Follow logs in real-time
docker logs -f overlap-avalara

# View logs from Docker Compose
docker-compose logs
docker-compose logs -f overlap-avalara
```

### Jenkins Build Logs
- Available in Jenkins UI: Job → Build History → Console Output
- Webhook trigger logs show GitHub payload processing
- Docker build and push logs for deployment tracking

### ngrok Request Monitoring
- Web interface: http://localhost:4040
- Request inspection and replay
- Webhook debugging and testing

##  Deployment Strategies

### Local Development
```bash
# Method 1: Direct Go execution
make run

# Method 2: Docker development
docker-compose up

# Method 3: With ngrok exposure
make run &
ngrok http 8080
```

### Staging/Testing
```bash
# Pull latest image
docker pull 210423/overlap-avalara:latest

# Deploy with Docker Compose
docker-compose -f docker-compose.staging.yml up -d

# Or use specific version
docker run -d -p 8081:8081 --name overlap-avalara-staging 210423/overlap-avalara:v1.2.3
```

### Production Deployment
```bash
# Using Docker Compose with production config
docker-compose -f docker-compose.prod.yml up -d

```

##  Troubleshooting

### Common Issues

#### Application Issues
1. **Port already in use**: Change port in `config/*/server.yml`
2. **Binary not found**: Run `make build`
3. **Tests failing**: Run `go mod tidy` and check dependencies

#### Docker Issues
1. **Image build fails**: Check Dockerfile syntax and dependencies
2. **Container won't start**: Verify exposed ports and configuration
3. **Permission denied**: Check file permissions and Docker daemon

#### Webhook Issues
1. **"Failed to connect to host"**: 
   - Restart ngrok and update webhook URL
   - Check firewall settings
   - Verify Jenkins is accessible

2. **"403 Forbidden"**:
   - Check Jenkins CSRF settings
   - Verify webhook token
   - Use correct endpoint URL

3. **Builds not triggering**:
   - Check GitHub webhook delivery status
   - Verify Jenkins job configuration
   - Review webhook payload filters

#### ngrok Issues
1. **Tunnel disconnected**: Free plan has time limits
2. **URL changed**: Update GitHub webhook with new URL
3. **Rate limiting**: Upgrade to paid plan for higher limits

### Debug Commands

```bash

# Test Docker image locally
docker run --rm -p 8081:8081 210423/overlap-avalara:latest

# Verify webhook endpoint
curl -X POST https://your-ngrok-url.ngrok.io/generic-webhook-trigger/invoke?token=your-token

# Check Jenkins logs
docker logs jenkins-container

# Monitor ngrok traffic
curl http://localhost:4040/api/tunnels
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests (`make test`)
5. Format code (`make fmt`)
6. Test Docker build (`make docker-build`)
7. Commit your changes (`git commit -m 'Add amazing feature'`)
8. Push to the branch (`git push origin feature/amazing-feature`)
9. Open a Pull Request

### Development Guidelines

- Follow Go best practices and conventions
- Write tests for new functionality
- Update documentation for API changes
- Ensure Docker image builds successfully
- Test webhook integration locally with ngrok
- Verify CI/CD pipeline passes

## Quick Start Checklist

### Development Setup
- [ ] Clone repository
- [ ] Install Go 1.24.2+
- [ ] Run `go mod tidy`
- [ ] Install Docker and Docker Compose
- [ ] Install and configure ngrok
- [ ] Set up Jenkins (optional)

### Local Development
- [ ] Start application: `make run`
- [ ] Test API endpoints
- [ ] Run tests: `make test`
- [ ] Build Docker image: `make docker-build`

### CI/CD Setup
- [ ] Configure GitHub secrets for Docker Hub
- [ ] Set up Jenkins job with webhook trigger
- [ ] Configure GitHub webhook with ngrok URL
- [ ] Test webhook delivery and build trigger

### Production Deployment
- [ ] Pull image: `docker pull 210423/overlap-avalara:latest`
- [ ] Deploy with Docker Compose
- [ ] Configure monitoring and logging
- [ ] Set up health checks and alerts

## 📄 License

[Add your license information here]

## Getting Help

If you encounter any issues:

1. **Check logs**: Application logs in `logs/` directory
2. **Review configuration**: Verify environment-specific YAML files
3. **Test components**: Use provided curl examples and debug commands
4. **CI/CD issues**: Check GitHub Actions logs and Jenkins console output
5. **Docker problems**: Use `docker logs` and `docker inspect` commands
6. **Webhook debugging**: Use ngrok web interface and GitHub webhook delivery logs

### Support Resources

- **Application Issues**: Check application logs and configuration
- **Docker Issues**: Docker documentation and container logs
- **CI/CD Problems**: GitHub Actions logs and Jenkins documentation
- **Webhook Integration**: ngrok documentation and GitHub webhook guides

---

**Built with using Go, Docker, GitHub Actions, Jenkins, and modern DevOps practices**

**Docker Hub**: `210423/overlap-avalara:latest`