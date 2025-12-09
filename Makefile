.PHONY: all build-client build-all install clean deps help

# Default target
all: build-client

# Build the Go CLI client (statically linked)
build-client:
	@echo "🔨 Building cpx client..."
	@mkdir -p bin
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

clean:
	rm -rf bin/
	rm -rf cpx/cpx
	@echo "✅ Cleaned build artifacts"

# Download Go dependencies
deps:
	cd cpx && go mod tidy

# Help
help:
	@echo "Cpx - C++ Project Generator"
	@echo ""
	@echo "Usage:"
	@echo "  make build-client   Build the Go CLI client"
	@echo "  make build-all      Build for all platforms (Linux, macOS, Windows)"
	@echo "  make install        Install cpx to /usr/local/bin"
	@echo "  make clean          Remove build artifacts"
	@echo "  make deps           Download Go dependencies"
