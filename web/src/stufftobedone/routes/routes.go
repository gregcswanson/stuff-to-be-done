package routes

import (
	"log"
	"net/http"
	"stufftobedone/handlers/mux"
	"stufftobedone/middleware"
	"stufftobedone/templates"

	//gorillahandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	//"github.com/rs/cors"
	//"appengine"
	//"appengine/datastore"
	//"errors"
	//"log"
	//"time"
)

/*SetupRoutes defines all the available routes and handlers for the application */
func SetupRoutes() error {
	r := mux.NewRouter()

	// redirect to the google login page, redirect to authenticated
	r.HandleFunc("/login", handlers.Login).Methods("GET")
	r.HandleFunc("/logout", handlers.Logout).Methods("GET")
	r.HandleFunc("/accesstoken", handlers.AccessToken).Methods("GET")
	// api authentication requests
	r.HandleFunc("/api/v1/request_token", handlers.RequestToken).Methods("POST")
	r.HandleFunc("/api/v1/refresh_token", handlers.RefreshToken).Methods("POST")

	// api test requests
	r.HandleFunc("/api/v1/ping", handlers.Ping).Methods("GET")
	r.HandleFunc("/api/v1/pong", handlers.Pong).Methods("POST")
	r.HandleFunc("/api/v1/pang", middleware.Auth(handlers.Pang)).Methods("GET")

	// profile
	r.HandleFunc("/api/v1/profile", middleware.Auth(handlers.Profile)).Methods("GET")

	// lists
	r.HandleFunc("/api/v1/active", middleware.Auth(handlers.ActiveTasks)).Methods("GET")

	// daytask
	r.HandleFunc("/api/v1/daytasks", middleware.Auth(handlers.PostDayTasks)).Methods("POST")
	r.HandleFunc("/api/v1/daytasks", middleware.Auth(handlers.PutDayTasks)).Methods("PUT")

	// tasks

	/*
		r.HandleFunc("/api/1/gettoken", authenticated(handlerGetToken)).Methods("GET")
		r.HandleFunc("/api/1/checktoken", handlerCheckToken).Methods("GET")
		r.HandleFunc("/api/1/token", TokenHandler).Methods("GET")
		r.HandleFunc("/api/1/tokenread", TokenReadHander).Methods("GET")

		// projects
		r.HandleFunc("/api/1/projects", useCaseMiddleware(ProjectsHandler)).Methods("GET")
		r.HandleFunc("/api/1/projects/{id}", useCaseMiddleware(ProjectHandler)).Methods("GET")
		r.HandleFunc("/api/1/projects", useCaseMiddleware(CreateProjectHandler)).Methods("POST")
		r.HandleFunc("/api/1/projects/{id}", useCaseMiddleware(UpdateProjectHandler)).Methods("PUT")
		r.HandleFunc("/api/1/projects/{id}", DeleteListHandler).Methods("DELETE")

		// project lines
		r.HandleFunc("/api/1/projects/{project_id}/lines",  useCaseMiddleware(CreateProjectLineHandler)).Methods("POST")
		r.HandleFunc("/api/1/projects/{project_id}/lines/{id}",  useCaseMiddleware(UpdateProjectLineHandler)).Methods("PUT")
		r.HandleFunc("/api/1/projects/{project_id}/lines/{id}", useCaseMiddleware(DeleteProjectLineHandler)).Methods("DELETE")

		// day items
		r.HandleFunc("/api/1/day/{day_id}", useCaseMiddleware(DayHandler)).Methods("GET")
		//r.HandleFunc("/api/1/dayitem/{id}", useCaseMiddleware(ProjectHandler)).Methods("GET")
		r.HandleFunc("/api/1/dayitem", useCaseMiddleware(CreateDayItemHandler)).Methods("POST")
		//r.HandleFunc("/api/1/dayitem/{id}", useCaseMiddleware(UpdateProjectHandler)).Methods("PUT")
		//r.HandleFunc("/api/1/dayitem/{id}", DeleteListHandler).Methods("DELETE")

		// lists
		r.HandleFunc("/api/1/lists", ListsHandler).Methods("GET")
		r.HandleFunc("/api/1/lists/{id}", ListHandler).Methods("GET")
		r.HandleFunc("/api/1/lists", CreateListHandler).Methods("POST")
		r.HandleFunc("/api/1/lists/{id}", UpdateListHandler).Methods("PUT")
		r.HandleFunc("/api/1/lists/{id}", DeleteListHandler).Methods("DELETE")

		// items
		r.HandleFunc("/api/1/items", ItemsHandler).Methods("GET")
		r.HandleFunc("/api/1/items/{id}", ItemHandler).Methods("GET")

		r.HandleFunc("/", useCaseRequest(indexHander)).Methods("GET")
		r.HandleFunc("/day/{day_id}", useCaseRequest(dayHandler)).Methods("GET")
		r.HandleFunc("/day/{day_id}", useCaseRequest(dayPostHandler)).Methods("POST")
		r.HandleFunc("/day/{day_id}/item/{item_id}", useCaseRequest(dayItemHandler)).Methods("GET")
		r.HandleFunc("/day/{day_id}/item/{item_id}", useCaseRequest(dayItemPostHandler)).Methods("POST")
		r.HandleFunc("/day/{day_id}/item/{item_id}/toggle", useCaseRequest(togglePostHandler)).Methods("POST")
		r.HandleFunc("/day/{day_id}/item/{item_id}/delete", useCaseRequest(dayItemDeleteHandler)).Methods("GET")

		r.HandleFunc("/month/{month_id}", useCaseRequest(monthHandler)).Methods("GET")
		r.HandleFunc("/month/{month_id}/items", useCaseRequest(monthItemsHandler)).Methods("GET")
		r.HandleFunc("/month/{month_id}/items", useCaseRequest(monthItemsPostHandler)).Methods("POST")
		r.HandleFunc("/month/{month_id}/overview", useCaseRequest(monthOverviewHandler)).Methods("GET")
		r.HandleFunc("/month/{month_id}/item/{item_id}", useCaseRequest(monthItemHandler)).Methods("GET")
		r.HandleFunc("/month/{month_id}/item/{item_id}", useCaseRequest(monthItemPostHandler)).Methods("POST")
		r.HandleFunc("/month/{month_id}/item/{item_id}/toggle", useCaseRequest(monthItemTogglePostHandler)).Methods("POST")
		r.HandleFunc("/month/{month_id}/item/{item_id}/delete", useCaseRequest(monthItemDeleteHandler)).Methods("GET")

		r.HandleFunc("/projects", useCaseRequest(projectslistHandler)).Methods("GET")
		r.HandleFunc("/projects/upsert", useCaseRequest(projectUpsertHandler)).Methods("GET")
		r.HandleFunc("/projects/upsert", useCaseRequest(projectUpsertPostHandler)).Methods("POST")
		r.HandleFunc("/projects/upsert/{id}", useCaseRequest(projectUpsertHandler)).Methods("GET")
		r.HandleFunc("/projects/upsert/{id}", useCaseRequest(projectUpsertPostHandler)).Methods("POST")
		r.HandleFunc("/projects/delete/{id}", useCaseRequest(projectDeleteHandler)).Methods("GET")
		r.HandleFunc("/project/{id}", useCaseRequest(projectHandler)).Methods("GET")
		r.HandleFunc("/project/{id}", useCaseRequest(projectPostHandler)).Methods("POST")
		r.HandleFunc("/project/{project_id}/item/{item_id}", useCaseRequest(projectItemHandler)).Methods("GET")
		r.HandleFunc("/project/{project_id}/item/{item_id}", useCaseRequest(projectItemPostHandler)).Methods("POST")
		r.HandleFunc("/project/{project_id}/item/{item_id}/toggle", useCaseRequest(projectItemTogglePostHandler)).Methods("POST")
		r.HandleFunc("/project/{project_id}/item/{item_id}/delete", useCaseRequest(projectItemDeleteHandler)).Methods("GET")
	*/

	//r.HandleFunc("/", handlers.App).Methods("GET")
	//r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	//r.PathPrefix("/app").Handler(http.FileServer(http.Dir("./../apps/polymer2app/")))
	//r.PathPrefix
	r.HandleFunc("/{path:.*}", handlers.App)

	//originsOK := gorillahandlers.AllowedOrigins([]string{"*"})
	//headersOK := gorillahandlers.AllowedHeaders([]string{"Access-Control-Allow-Origin", "Accept", "Accept-Language", "Content-Language", "Origin","X-Requested-With", "Content-Type"})
	//methodsOK := gorillahandlers.AllowedMethods([]string{"GET", "HEAD", "PUT", "DELETE", "POST", "OPTIONS"})

	// TO DO - support CORS Preflight with middleware
	http.Handle("/", createCORSHandler()(r))

	// authentication

	// api

	// application

	return nil
}

type cors struct {
	http http.Handler
}

type blank struct {
	Message string `json:"message"`
}

func createCORSHandler() func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		ch := &cors{}
		ch.http = h
		return ch
	}
}

func (ch *cors) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method == "OPTIONS" {
		log.Println("CORS OPTIONS")
		method := r.Header.Get("Access-Control-Request-Method")
		w.Header().Set("Access-Control-Max-Age", "600")
		log.Println(method)
		templates.Blank(w)
		return
	}
	ch.http.ServeHTTP(w, r)
}
