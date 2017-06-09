package middleware

import (
  "github.com/gin-gonic/gin"
  "appengine"
  "appengine/user"
  "stufftobedone/repositories"
  "stufftobedone/domain"
  "log"
)

func UserMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        appUser := domain.AppUser{}
        appUser.Name = "unknown"

        // get the google logged in user details
        appengineContext := appengine.NewContext(c.Request)
        appengineUser := user.Current(appengineContext)
        if appengineUser == nil {
            appUser.ID = ""
            appUser.Name = ""
            appUser.Nickname = ""
            appUser.Email = ""
            appUser.IsLoggedIn = false
            log.Println("***** Logged out")
            //url, err := user.LoginURL(appengineContext, "/")//c.Request.URL.String())
            //if err == nil {
                appUser.LoginUrl = "/app/profile"
            //}
        } else {
            appUser.ID = "" // appengineUser.ID
            appUser.Name =  appengineUser.ID
            appUser.Nickname = ""
            appUser.Email = appengineUser.Email
            appUser.IsLoggedIn = true
            log.Println("***** Logged in as " + appUser.Email)
            url, err := user.LogoutURL(appengineContext, "/") //repository.request.URL.String())
            if err == nil {
                appUser.LogoutUrl = url
            }
            // cache this
            r := repositories.NewUserRepository(c.Request)
            appUser, err = r.Get(appUser)
        }  
        
        c.Set("user", appUser)

        c.Next()
    }
}