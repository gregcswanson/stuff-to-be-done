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
    //    apiV1book.GET("/day/:day", handlers.DayHandler)
    apiV1book.GET("/elements", handlers.BookElementsHandler)
    // later actions
    apiV1book.GET("/later", handlers.BookLaterHandler)
    apiV1book.POST("/later", handlers.BookLaterPostHandler)
    apiV1book.PUT("/later/:taskId/restore", handlers.BookLaterRestoreHandler)
    apiV1book.PUT("/later/:taskId/dotoday", handlers.BookLaterDoTodayHandler)
    apiV1book.PUT("/later/:taskId", handlers.BookLaterPutHandler)
    apiV1book.DELETE("/later/:taskId", handlers.BookLaterDeleteHandler)
    // completed actions
    //apiV1book.GET("/completed", handlers.ApiCompletedHandler) // remove
    // trash handlers
    apiV1book.GET("/trash", handlers.ApiTrashGetHandler) // get trash items

    // day handlers
    apiV1book.GET("/day/:dayAsString", handlers.ApiDayHander)
    apiV1book.POST("/day/:dayAsString", handlers.ApiDayNewElementHandler)
    apiV1book.PUT("/day/:taskId", handlers.ApiDayPutHandler)
    //apiV1book.PUT("/day/:taskId/dolater", handlers.ApiDayDoLaterPutHandler)
    apiV1book.DELETE("/task/:taskId", handlers.ApiDayDeleteHandler)
    // move
    apiV1book.PUT("/move/:taskId", handlers.ApiDayDoOnDatePutHander)
    // previous
    apiV1book.GET("/previous/:dayAsString", handlers.ApiPreviousHandler)
    apiV1book.PUT("/previous/:taskId", handlers.ApiDayPutHandler)

    apiV1book.PUT("/dolater/:taskId", handlers.ApiDayDoLaterPutHandler)
    
  }

  r.GET("/book", handlers.BookIndexsHandler) // will redirect to the default book for the user
  //r.GET("/live/:book", handlers.BookLiveHander) // will replace the sub pages
  //r.GET("/trash/:book", handlers.TrashHandler)
  //r.GET("/later/:book", handlers.LaterHandler) // will replace the sub pages
  book := r.Group("book/:book")
  {
    book.GET("/live", handlers.BookLiveHander) // will replace the sub pages
    book.GET("/trash", handlers.TrashHandler)
    book.GET("/later", handlers.LaterHandler)
    //book.GET("/history", handlers.LaterHandler)
      // add a middleware to setup the book context
      
      //book.GET("/:day/later", handlers.LaterHandler)
      //book.GET("/:day/completed", handlers.CompletedHandler)
      //book.GET("/:day/day", handlers.DayHandler)
      //book.GET("/:day/previous", handlers.PreviousHandler)
  }

  appProfile := r.Group("profile")
  {
    appProfile.GET("/", handlers.ProfileIndexsHandler)
  }

  r.GET("/login", handlers.LoginHandler)
  r.GET("/logout", handlers.LogoutHandler)
  r.GET("/about", handlers.AboutHandler)
  r.GET("/error", handlers.TestErrorPage)
  r.GET("/admin/registrations", handlers.RegistrationsHandler)
  r.GET("/element/*name", handlers.ElementHandler)
  r.GET("/", handlers.IndexHandler)

  http.Handle("/", r)
}
