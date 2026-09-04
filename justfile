# List available recipes
default:
    @just --list

# Build backend binary
build:
    go build -o bin/metadata-mcp ./cmd/main.go

# Run linters
lint:
    golangci-lint run ./...

# Run go mod tidy
tidy:
    go mod tidy

# Update go mod dependencies
update-go-deps:
    go get -u -t ./...
    @just tidy

# Update all dependencies
update-deps: update-go-deps

# Run tests
test:
    go test -v ./...

# Run all tests with local environment variables if available
alltest:
    test -f .local/env.sh && source .local/env.sh && go test -v -count=1 ./...

# Format Go code using goimports
fmt:
    goimports -w -local "github.com/autoget-project/metadata-mcp" .

# Check Go code formatting using goimports
fmt-check:
    @test -z "$($(go env GOPATH)/bin/goimports -local github.com/autoget-project/metadata-mcp -l . 2>/dev/null || goimports -local github.com/autoget-project/metadata-mcp -l .)" || (echo "Unformatted Go files found:" && goimports -local github.com/autoget-project/metadata-mcp -l . && exit 1)

# Clean build artifacts
clean:
    rm -rf bin
