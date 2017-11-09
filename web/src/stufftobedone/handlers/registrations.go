package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"stufftobedone/domain"
	"stufftobedone/repositories"
)

func RegistrationsHandler(c *gin.Context) {
	// setup the required repositories and data
	r := repositories.NewRegistrationRepository(c.Request)
	user := GetAppUser(c)
	// get the current registrations
	data, err := r.Find()
	if err != nil {
		// show the custom error page
		ErrorPage(c, err)
		return
	}
	log.Println(len(data))
	// return the result
	c.HTML(http.StatusOK, "registrations.html", gin.H{
		"title":        "Registrations",
		"isLoggedIn":   user.IsLoggedIn,
		"hasData":      len(data) > 0,
		"data":         data,
		"isProduction": user.IsProduction,
	})
}

func RegisterHandler(c *gin.Context) {
	registrations := repositories.NewRegistrationRepository(c.Request)
	log.Println("register handler")
	registration := domain.Registration{}
	registration.Email = c.PostForm("email")

	registrations.Store(registration)

	c.JSON(http.StatusOK, gin.H{
		"status":  "posted",
		"message": registration.Email,
	})

}
