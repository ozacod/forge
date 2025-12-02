package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const (
	Version    = "0.0.44"
	CLIVersion = "0.0.44"
)

func main() {
	r, err := SetupServer()
	if err != nil {
		fmt.Printf("Failed to setup server: %v\n", err)
		os.Exit(1)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	fmt.Printf("Forge server starting on port %s\n", port)
	if err := r.Run(":" + port); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
		os.Exit(1)
	}
}

// SetupServer initializes the Gin engine
func SetupServer() (*gin.Engine, error) {

	// Setup Gin router
	r := gin.Default()

	// CORS middleware
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"*"}
	r.Use(cors.New(config))

	// API routes
	api := r.Group("/api")
	{
		api.GET("", apiRoot)
		api.GET("/version", getVersion)
	}

	// Static file serving
	staticDir := "static"
	if envDir := os.Getenv("FORGE_STATIC_DIR"); envDir != "" {
		staticDir = envDir
	}

	hasStatic := false
	if _, err := os.Stat(staticDir); err == nil {
		if _, err := os.Stat(filepath.Join(staticDir, "index.html")); err == nil {
			hasStatic = true
			// Serve static assets
			r.Static("/assets", filepath.Join(staticDir, "assets"))
			r.StaticFile("/cpx.svg", filepath.Join(staticDir, "cpx.svg"))

			// Serve index.html for root
			r.GET("/", func(c *gin.Context) {
				c.File(filepath.Join(staticDir, "index.html"))
			})

			// Catch-all for SPA routes
			r.NoRoute(func(c *gin.Context) {
				path := c.Request.URL.Path
				if strings.HasPrefix(path, "/api") {
					c.JSON(http.StatusNotFound, gin.H{"detail": "Not found"})
					return
				}
				c.File(filepath.Join(staticDir, "index.html"))
			})
		}
	}

	// Fallback root if no static files
	if !hasStatic {
		r.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message":     "Forge API - C++ Project Generator",
				"version":     Version,
				"cli_version": CLIVersion,
				"docs":        "/docs",
				"frontend":    "Not built. Run 'make build-frontend-go' to build the UI.",
			})
		})
	}

	return r, nil
}
func apiRoot(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message":     "Forge API - C++ Project Generator",
		"version":     Version,
		"cli_version": CLIVersion,
		"docs":        "/docs",
	})
}

func getVersion(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version":     Version,
		"cli_version": CLIVersion,
		"name":        "cpx",
		"description": "C++ Project Generator - Like Cargo for Rust, but for C++!",
	})
}
