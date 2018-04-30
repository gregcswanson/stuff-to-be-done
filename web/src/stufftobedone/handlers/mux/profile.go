package handlers

import (
	"net/http"
	"stufftobedone/templates"
)

type profile struct {
	Email  string `json:"email"`
	BookID string `json:"bookId"`
}

/*Profile ... */
func Profile(w http.ResponseWriter, r *http.Request) {
	c := getAppContext(r)

	// get or create the default book for this user
	defaultBook, err := c.BookUseCases.GetDefault(c.User)
	if err != nil {
		templates.JSONError(w, err)
	}
	data := profile{Email: c.User.Email, BookID: defaultBook.ID}

	templates.JSON(w, data)
}
