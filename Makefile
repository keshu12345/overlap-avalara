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
