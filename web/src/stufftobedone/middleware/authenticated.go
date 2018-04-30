package middleware

import (
	"context"
	"net/http"
	handlers "stufftobedone/handlers/mux"
	"stufftobedone/templates"
	"stufftobedone/usecases"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

/*Chain ... */
func Auth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// get the token
		token := r.Header.Get("Authorization")
		if token == "" {
			templates.JSONErrorMessage(w, "Authorization required")
			return
		}

		// validate
		tokenClaims, err := jwt.ParseWithClaims(token, &handlers.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(handlers.AccessTokenKey), nil
		})
		if err != nil {
			templates.JSONError(w, err)
			return
		}

		// get the user
		if claims, ok := tokenClaims.Claims.(*handlers.CustomClaims); ok && tokenClaims.Valid {
			now := time.Now().Unix()
			if now > claims.Expires {
				templates.JSONErrorExpiredToken(w)
				return
			}
			if claims.TokenType != "request-token" {
				templates.JSONErrorMessage(w, "The token is not valid")
				return
			}
			userUseCases := usecases.NewUserUseCases(r)
			appUser, err := userUseCases.FindOrCreateByEmail(claims.Username)
			if err != nil {
				templates.JSONError(w, err)
				return
			}
			taskUseCases := usecases.NewTaskUseCases2(r, appUser)
			bookUseCases := usecases.NewBookUseCases(r)
			newContext := handlers.AppContext{UserUseCases: userUseCases, TaskUseCases: taskUseCases, BookUseCases: bookUseCases, User: appUser, Token: "authenticated !!"}

			ctx := context.WithValue(r.Context(), "appcontext", newContext)
			handler(w, r.WithContext(ctx))
		} else {
			templates.JSONErrorMessage(w, "Token is invalid")
			return
		}
	}
}
