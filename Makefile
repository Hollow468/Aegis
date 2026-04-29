BINARY := aegis
GO := go
NPM := npm
FRONTEND_DIR := frontend
EMBED_DIR := internal/web/dist/frontend

.PHONY: all build frontend clean dev test

all: frontend build

# Build frontend and embed into binary
build: frontend
	$(GO) build -o $(BINARY) ./cmd/apigateway

# Build frontend only
frontend:
	cd $(FRONTEND_DIR) && $(NPM) ci && $(NPM) run build

# Development: build frontend, run with hot reload
dev: frontend
	$(GO) run ./cmd/apigateway

# Go tests
test:
	$(GO) test ./...

# Clean build artifacts
clean:
	rm -f $(BINARY)
	rm -rf $(EMBED_DIR)/*
	rm -rf $(FRONTEND_DIR)/node_modules
	rm -rf $(FRONTEND_DIR)/dist

# Frontend dev server (standalone, proxies to backend at :8080)
frontend-dev:
	cd $(FRONTEND_DIR) && $(NPM) run dev

# Build for multiple platforms
build-all: frontend
	GOOS=linux GOARCH=amd64 $(GO) build -o $(BINARY)-linux-amd64 ./cmd/apigateway
	GOOS=darwin GOARCH=amd64 $(GO) build -o $(BINARY)-darwin-amd64 ./cmd/apigateway
	GOOS=darwin GOARCH=arm64 $(GO) build -o $(BINARY)-darwin-arm64 ./cmd/apigateway
	GOOS=windows GOARCH=amd64 $(GO) build -o $(BINARY)-windows-amd64.exe ./cmd/apigateway
