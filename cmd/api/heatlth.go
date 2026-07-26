package main

import (
	"net/http"
)

func (app *application) healthCheckHandler(response http.ResponseWriter, request *http.Request) {
	data := map[string]string{
		"status":  "ok",
		"env":     app.config.env,
		"version": version,
	}
	if err := writeJSON(response, http.StatusOK, data); err != nil {
		app.internalServerError(response, request, err)
	}
}
