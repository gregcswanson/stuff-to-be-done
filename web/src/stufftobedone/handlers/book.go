package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"stufftobedone/repositories"
	//"time"
)

func BookIndexsHandler(c *gin.Context) {
	// setup the required repositories
	user := GetAppUser(c)
	bookRepository := repositories.NewBookRepository(c.Request)
	// get or create the default book for this user
	defaultBook, err := bookRepository.GetDefault(user)
	if err != nil {
		ErrorPage(c, err)
	}
	// redirect to the current day
	c.Redirect(http.StatusFound, "/app/book/"+defaultBook.ID+"/live")
}

func BookLiveHander(c *gin.Context) {
	user := GetAppUser(c)
	bookID := c.Param("book")

	c.HTML(http.StatusOK, "book.html", gin.H{
		"title":        "Today",
		"logouturl":    user.LogoutUrl,
		"isLoggedIn":   true,
		"bookID":       bookID,
		"isProduction": user.IsProduction,
	})
}
