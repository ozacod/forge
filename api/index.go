package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ozacod/forge/web/server/pkg/server"
)

var app *gin.Engine

func init() {
	var err error
	app, err = server.SetupServer()
	if err != nil {
		panic(err)
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	app.ServeHTTP(w, r)
}
