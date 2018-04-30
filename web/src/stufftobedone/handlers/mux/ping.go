package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"stufftobedone/templates"
)

type ping struct {
	Message string `json:"message"`
}

type pong struct {
	Message string `json:"message"`
}

/*Ping ... */
func Ping(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//message := vars["message"]
	//message := "message"
	query := r.URL.Query()
	message := query.Get("message")
	if message == "" {
		message = "empty"
	}

	data := ping{Message: message}
	templates.JSON(w, data)
}

/*Pang ... */
func Pang(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	message := query.Get("message")
	if message == "" {
		message = "Hello"
	}
	c := getAppContext(r)
	message = message + " " + c.User.Email
	data := ping{Message: message}
	templates.JSON(w, data)
}

/*Pong ... */
func Pong(w http.ResponseWriter, r *http.Request) {
	// get the data
	var pongJSON pong
	err := json.NewDecoder(r.Body).Decode(&pongJSON)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}

	data := pong{Message: pongJSON.Message}
	templates.JSON(w, data)
}
