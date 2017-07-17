package stufftobedone

import (
  "net/http"
  "github.com/gin-gonic/gin"
  "stufftobedone/handlers"
  "stufftobedone/middleware"
  //"log"
)

func init() { 
  r := gin.Default()
  r.LoadHTMLGlob("templates/**/*")
  r.Use(middleware.SharedDataMiddleWare())
  r.Use(middleware.UserMiddleware())
  r.NoRoute(handlers.FourZeroFourHandler)

  apiV1internal := r.Group("api/v1/internal")
  {
      // TO DO - set a custom shared variable using middleware so that the 404 returns an api error
      apiV1internal.POST("/register", handlers.RegisterHandler)
      apiV1internal.GET("/ping", handlers.PingHandler)
  }

  
  apiV1book := r.Group("api/v1/book/:bookId")
  {
    // add a middleware to setup the book context
    apiV1book.GET("/elements", handlers.BookElementsHandler)
    // later actions
    apiV1book.GET("/later", handlers.ApiLaterHandler)
    apiV1book.POST("/later", handlers.BookLaterPostHandler)
    apiV1book.GET("/latercount", handlers.BookLaterCountHandler)
    // trash handlers
    apiV1book.GET("/trash", handlers.ApiTrashGetHandler) // get trash items
    apiV1book.DELETE("/trash", handlers.ApiTrashEmptyHandler)

    // day handlers
    apiV1book.GET("/day/:dayAsString", handlers.ApiDayHander)
    apiV1book.POST("/day/:dayAsString", handlers.ApiDayNewElementHandler)
    apiV1book.PUT("/day/:taskId", handlers.ApiDayPutHandler)
    apiV1book.DELETE("/task/:taskId", handlers.ApiDayDeleteHandler)
    // move
    apiV1book.PUT("/move/:taskId", handlers.ApiDayDoOnDatePutHander)
    // previous
    apiV1book.GET("/previous/:dayAsString", handlers.ApiPreviousHandler)
    apiV1book.PUT("/previous/:taskId", handlers.ApiDayPutHandler)

    apiV1book.PUT("/dolater/:taskId", handlers.ApiDayDoLaterPutHandler)
    apiV1book.PUT("/comment", handlers.CommentPutHandler)
    
  }

  r.GET("app/book", handlers.BookIndexsHandler) // will redirect to the default book for the user
  book := r.Group("app/book/:book")
  {
    book.GET("/live", handlers.BookLiveHander) // will replace the sub pages
    book.GET("/trash", handlers.TrashHandler)
    book.GET("/later", handlers.LaterHandler)
    book.GET("/history", handlers.HistoryHandler)
  }

  appProfile := r.Group("profile")
  {
    appProfile.GET("/", handlers.ProfileIndexsHandler)
  }

  r.GET("/login", handlers.LoginHandler)
  r.GET("/logout", handlers.LogoutHandler)
  r.GET("/terms", handlers.TermsHandler)
  r.GET("/privacy", handlers.PrivacyHandler)
  r.GET("/about", handlers.AboutHandler)
  r.GET("/error", handlers.TestErrorPage)
  r.GET("/admin/registrations", handlers.RegistrationsHandler)
  r.GET("/element/*name", handlers.ElementHandler)
  r.GET("/", handlers.IndexHandler)

  http.Handle("/", r)
}
