# Forge - C++ Project Generator Makefile

.PHONY: all build-client build-frontend build-server install clean setup-frontend run-server run-frontend run-go stop-server stop-frontend stop help

# Default target
all: build-client

# Build the Go CLI client (statically linked)
build-client:
	@echo "🔨 Building cpx client..."
	cd cpx && \
		CGO_ENABLED=0 go build -ldflags="-s -w" -o ../bin/cpx .
	@echo "✅ Built: bin/cpx"

# Build for all platforms
build-all: build-client
	@echo "🔨 Building for all platforms..."
	@mkdir -p bin
	cd cpx && \
		GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o ../bin/cpx-linux-amd64 . && \
		GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -o ../bin/cpx-linux-arm64 . && \
		GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o ../bin/cpx-darwin-amd64 . && \
		GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-s -w" -o ../bin/cpx-darwin-arm64 . && \
		GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o ../bin/cpx-windows-amd64.exe .
	@echo "✅ Built binaries for all platforms in bin/"

# Install the client to /usr/local/bin
install: build-client
	@echo "📦 Installing cpx to /usr/local/bin..."
	sudo cp bin/cpx /usr/local/bin/
	@echo "✅ Installed! Run 'cpx --help' to get started"

# Setup the frontend (install npm dependencies)
setup-frontend:
	@echo "📦 Setting up frontend..."
	cd web/frontend && npm install
	@echo "✅ Frontend setup complete"

# Build frontend for production (outputs to web/server/static)
build-frontend:
	@echo "🔨 Building frontend..."
	cd web/frontend && npm run build
	@rm -rf web/server/static
	@mv web/frontend/dist web/server/static
	@echo "✅ Frontend built to web/server/static"

# Build the server
build-server:
	@echo "🔨 Building server..."
	cd web/server && go build -o server ./cmd/server
	@echo "✅ Built: web/server/server"

# Run the server (serves API + static frontend)
run-server: build-server
	@echo "🚀 Starting cpx server on http://localhost:8000..."
	cd web/server && \
		PORT=8000 ./server

# Run the frontend in dev mode
run-frontend:
	@echo "🚀 Starting frontend dev server on http://localhost:5173..."
	cd web/frontend && npm run dev

# Build frontend and run server (production mode)
run-go: build-frontend build-server
	@echo "🚀 Starting cpx server with bundled frontend on http://localhost:8000..."
	cd web/server && \
		PORT=8000 ./server

# Stop the server (kills process on port 8000)
stop-server:
	@echo "🛑 Stopping server on port 8000..."
	@-lsof -ti:8000 | xargs kill -9 2>/dev/null || true
	@echo "✅ Server stopped"

# Stop the frontend (kills process on port 5173)
stop-frontend:
	@echo "🛑 Stopping frontend on port 5173..."
	@-lsof -ti:5173 | xargs kill -9 2>/dev/null || true
	@echo "✅ Frontend stopped"

# Stop both server and frontend
stop: stop-server stop-frontend
	@echo "✅ All services stopped"

# Clean build artifacts
clean:
	rm -rf bin/
	rm -rf cpx/cpx
	rm -rf web/server/server
	@echo "✅ Cleaned build artifacts"

# Download Go dependencies
deps:
	cd cpx && go mod tidy
	cd web/server && go mod tidy
	cd api && go mod tidy

# Help
help:
	@echo "Forge - C++ Project Generator"
	@echo ""
	@echo "Usage:"
	@echo "  make build-client      Build the Go CLI client"
	@echo "  make build-frontend    Build frontend (to web/server/static)"
	@echo "  make build-server      Build the backend server"
	@echo "  make build-all         Build for all platforms (Linux, macOS, Windows)"
	@echo "  make install           Install cpx to /usr/local/bin"
	@echo "  make setup-frontend    Install frontend npm dependencies"
	@echo "  make run-go            Build frontend & run server (production)"
	@echo "  make run-server        Start the server only"
	@echo "  make run-frontend      Start the React dev server"
	@echo "  make stop-server       Stop the server on port 8000"
	@echo "  make stop-frontend     Stop the React frontend"
	@echo "  make stop              Stop both server and frontend"
	@echo "  make clean             Remove build artifacts"
	@echo "  make deps              Download Go dependencies"
	@echo ""
	@echo "Quick Start (Development):"
	@echo "  1. make setup-frontend"
	@echo "  2. make run-server     (terminal 1)"
	@echo "  3. make run-frontend    (terminal 2)"
	@echo ""
	@echo "Quick Start (Production):"
	@echo "  1. make setup-frontend"
	@echo "  2. make run-go          (builds frontend & serves at :8000)"
