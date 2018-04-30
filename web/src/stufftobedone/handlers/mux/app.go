package handlers

import (
	//"errors"
	//"log"
	"net/http"
	"stufftobedone/templates"
)

type appPage struct {
}

/*App ... */
func App(w http.ResponseWriter, r *http.Request) {
	data := appPage{}
	templates.HTML(w, "app.html", data)
}
