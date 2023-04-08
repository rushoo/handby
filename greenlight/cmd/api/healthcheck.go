package main

import (
	"net/http"
)

/*
/v1/healthcheck 	healthcheckHandler 		Show application information
*/
const version = "1.0.0"

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	env := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     version,
		},
	}
	//通过marshal结构体/map来实现json响应
	err := app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
