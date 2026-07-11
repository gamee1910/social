package main

import (
	"log"
	"net/http"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "ok",
		"env":     app.config.env,
		"version": version,
	}
	if err := writeJSON(w, http.StatusOK, data); err != nil {
		err := responseJSONError(w, http.StatusInternalServerError, err.Error())
		if err != nil {
			log.Panic(err)
		}
	}
}
