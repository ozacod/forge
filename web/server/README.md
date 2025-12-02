# Forge Web Server

Minimal Go server for the Forge web interface using Gin.

## Features

- Version information API
- Static file serving for frontend
- CORS support

## Building

```bash
cd web/server
go build ./cmd/server
```

## Running

```bash
./server
```

Or set environment variables:

```bash
PORT=8000 ./server
```

## API Endpoints

- `GET /api` - API root
- `GET /api/version` - Version information

## Structure

```
web/server/
├── cmd/
│   └── server/
│       └── main.go          # Main server entry point
├── pkg/
│   └── server/
│       └── server.go         # Server setup
├── static/                   # Frontend build output (generated)
└── go.mod                    # Go module definition
```

## Dependencies

- `github.com/gin-gonic/gin` - HTTP web framework
- `github.com/gin-contrib/cors` - CORS middleware
