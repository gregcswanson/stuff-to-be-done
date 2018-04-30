package templates

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"appengine"
)

//Compile templates on start
var templates = template.Must(template.ParseFiles(
	"templates/includes/header2.html",
	"templates/includes/footer2.html",
	"templates/layouts/404.html",
	"templates/layouts/app.html",
	"templates/layouts/error2.html",
	"templates/layouts/message.html"))

type page struct {
	IsProduction bool
	Model        interface{}
}

type errorData struct {
	Title   string
	Message string
}

type errorDataJSON struct {
	Error string `json:"error"`
}

/*HTML ... */
func HTML(w http.ResponseWriter, templateName string, data interface{}) {
	page := page{IsProduction: true, Model: data}
	page.IsProduction = !appengine.IsDevAppServer()
	templates.ExecuteTemplate(w, templateName, page)
}

/*HTMLError ... */
func HTMLError(w http.ResponseWriter, title string, message string) {
	page := page{IsProduction: true, Model: errorData{Title: title, Message: message}}
	page.IsProduction = !appengine.IsDevAppServer()
	templates.ExecuteTemplate(w, "error2.html", page)
}

/*HTMLMessage ... */
func HTMLMessage(w http.ResponseWriter, title string, message string) {
	page := page{IsProduction: true, Model: errorData{Title: title, Message: message}}
	page.IsProduction = !appengine.IsDevAppServer()
	templates.ExecuteTemplate(w, "message.html", page)
}

/*Blank .. */
func Blank(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Accept", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
}

/*JSON ... */
func JSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Accept", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	log.Println(data)
	jsonData, err := json.Marshal(data)
	log.Println(jsonData)
	if err != nil {
		panic(err)
	}
	w.Write(jsonData)
}

/*JSONErrorMessage .. */
func JSONErrorMessage(w http.ResponseWriter, e string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Accept", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.WriteHeader(http.StatusInternalServerError)
	jsonData, err := json.Marshal(&errorDataJSON{Error: e})
	if err != nil {
		panic(err)
	}
	w.Write(jsonData)
}

func JSONErrorExpiredToken(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Accept", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.WriteHeader(http.StatusUnauthorized)
	jsonData, err := json.Marshal(&errorDataJSON{Error: "The token has expired"})
	if err != nil {
		panic(err)
	}
	w.Write(jsonData)
}

/*JSONError .. */
func JSONError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Accept", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.WriteHeader(http.StatusInternalServerError)
	jsonData, err := json.Marshal(&errorDataJSON{Error: err.Error()})
	if err != nil {
		panic(err)
	}
	w.Write(jsonData)
}
