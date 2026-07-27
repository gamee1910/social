package main

import (
	"net/http"

	"github.com/gamee1910/social/internal/httpx"
)

func (app *application) healthCheckHandler(response http.ResponseWriter, request *http.Request) {
	data := map[string]string{
		"status":  "ok",
		"env":     app.config.env,
		"version": version,
	}
	if err := httpx.WriteJSON(response, http.StatusOK, data); err != nil {
		httpx.InternalServerError(response, request, err)
	}
}
