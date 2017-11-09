package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"stufftobedone/repositories"
)

func ElementHandler(c *gin.Context) {
	// get the element latest version - there needs to be a global version element key to invalidate this
	// when an element is updated
	name := c.Param("name")
	log.Println("requesting element", name)
	r := repositories.NewCustomElementRepository(c.Request)
	customElement, err := r.GetByName(name)
	if err != nil {
		ErrorPage(c, err)
	}
	c.String(http.StatusOK, customElement.Body)
}

func BookElementsHandler(c *gin.Context) {
	// get a list of available elements for this book
	bookID := c.Param("bookId")
	log.Println("getting list of elements available for the book", bookID)
	r := repositories.NewCustomElementRepository(c.Request)
	bookElements, err := r.FindBookElements(bookID)
	if err != nil {
		c.JSON(500, gin.H{
			"message": err,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"data": bookElements,
		})
	}
}
