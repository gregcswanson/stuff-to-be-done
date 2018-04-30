package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"stufftobedone/domain"
	"stufftobedone/templates"
	"stufftobedone/usecases"
	"time"

	"github.com/dgrijalva/jwt-go"

	"appengine"
	"appengine/user"
)

/*Login ... not sure we need this function at all */
func Login(w http.ResponseWriter, r *http.Request) {
	appengineContext := appengine.NewContext(r)
	appengineUser := user.Current(appengineContext)
	if appengineUser != nil {
		url, err := user.LogoutURL(appengineContext, "/")
		if err == nil {
			w.Header().Set("Location", url)
			w.WriteHeader(http.StatusFound)
			return
		}
	}
	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusFound)
}

type accessTokenPage struct {
	Token string
}

type accessToken struct {
	IsValid bool
	Email   string
	ID      string
}

type requestTokenRequest struct {
	Token string `json:"token"`
}

/*AppContext ... */
type AppContext struct {
	Token        string
	UserUseCases *usecases.UserUseCases
	TaskUseCases *usecases.TaskUseCases
	BookUseCases *usecases.BookUseCases
	User         domain.AppUser
}

type requestTokenResponse struct {
	RequestToken string `json:"request_token"`
	RefreshToken string `json:"refresh_token"`
}

type CustomClaims struct {
	Username  string `json:"username"`
	Namespace string `json:"namespace"`
	TokenType string `json:"tokentype"`
	Expires   int64  `json:"exp"`
	jwt.StandardClaims
}

const AccessTokenKey = "Gtya59de12FPLe9"

/*AccessToken ... */
func AccessToken(w http.ResponseWriter, r *http.Request) {
	// the return URL may be sent (FUTURE)
	appengineContext := appengine.NewContext(r)
	appengineUser := user.Current(appengineContext)
	if appengineUser == nil {
		templates.HTMLError(w, "AccessToken", "Failed to authenticate for access token request")
		return
	}
	// get the return url from the request, otherwise return the current url
	query := r.URL.Query()
	returnURL := query.Get("redirectUrl")

	if returnURL == "" {
		returnURL = "/authorized"
	}
	log.Println(returnURL)

	//if !appengine.IsDevAppServer() {
	//	returnURL = "/authorized"
	//}

	//if token == "" {
	// get or create the app user
	userUseCases := usecases.NewUserUseCases(r)
	appUser, err := userUseCases.FindOrCreateByAppEngine()

	token, err := createToken(appUser.Email, appUser.ID, "access-token")
	if err != nil {
		templates.HTMLError(w, "AccessToken", "Failed to create access token")
		return
	}

	http.Redirect(w, r, returnURL+"?token="+token, http.StatusSeeOther)
}

/*RequestToken reads the token from the json body and returns a valid request token*/
func RequestToken(w http.ResponseWriter, r *http.Request) {
	// get the token from the body
	var requestData requestTokenRequest
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		templates.JSONError(w, err)
		return
	}

	// read and validate the token
	token, err := jwt.ParseWithClaims(requestData.Token, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(AccessTokenKey), nil
	})

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		now := time.Now().Unix()
		if now > claims.Expires {
			templates.JSONErrorExpiredToken(w)
			return
		}
		if claims.TokenType != "access-token" {
			templates.JSONErrorMessage(w, "The token is not valid")
			return
		}

		// create a request and refresh token
		requestToken, _ := createToken(claims.Username, claims.Namespace, "request-token")
		refreshToken, _ := createToken(claims.Username, claims.Namespace, "refresh-token")
		data := requestTokenResponse{RequestToken: requestToken, RefreshToken: refreshToken}
		templates.JSON(w, data)
	} else {
		templates.JSONError(w, err)
		return
	}

}

/*RefreshToken reads the token from the json body and returns a valid request token*/
func RefreshToken(w http.ResponseWriter, r *http.Request) {
	// get the token from the body
	var requestData requestTokenRequest
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		templates.JSONError(w, err)
		return
	}

	// read and validate the token
	token, err := jwt.ParseWithClaims(requestData.Token, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(AccessTokenKey), nil
	})

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		now := time.Now().Unix()
		if now > claims.Expires {
			templates.JSONErrorExpiredToken(w)
			return
		}
		if claims.TokenType != "refresh-token" {
			templates.JSONErrorMessage(w, "The refresh token is not valid")
			return
		}

		// create a request and refresh token
		requestToken, _ := createToken(claims.Username, claims.Namespace, "request-token")
		refreshToken, _ := createToken(claims.Username, claims.Namespace, "refresh-token")
		data := requestTokenResponse{RequestToken: requestToken, RefreshToken: refreshToken}
		templates.JSON(w, data)
	} else {
		templates.JSONError(w, err)
		return
	}

}

/*Logout ... */
func Logout(w http.ResponseWriter, r *http.Request) {
	appengineContext := appengine.NewContext(r)
	appengineUser := user.Current(appengineContext)
	if appengineUser != nil {
		url, err := user.LogoutURL(appengineContext, "/logout")
		if err == nil {
			w.Header().Set("Location", url)
			w.WriteHeader(http.StatusFound)
			return
		}
	}
	templates.HTMLMessage(w, "Logout", "Successful")
}

func validateRequestToken(w http.ResponseWriter, r *http.Request) {
	// middleware for api

	// get the token and validate

	// continue
}

func getAppContext(r *http.Request) AppContext {
	if value := r.Context().Value("appcontext"); value != nil {
		return value.(AppContext)
	}
	return AppContext{}
}

func createToken(username string, namespance string, tokenType string) (string, error) { //w http.ResponseWriter, r *http.Request) (string, error) {

	expires := time.Now().Add(time.Hour * 24 * 7).Unix()
	if tokenType == "access-token" {
		expires = time.Now().Add(time.Minute * 10).Unix()
	}
	if tokenType == "refresh-token" {
		expires = time.Now().Add(time.Hour * 24 * 30).Unix()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username":  username,
		"namespace": namespance,
		"tokentype": tokenType,
		"exp":       expires,
	})

	key := []byte(AccessTokenKey)

	tokenString, _ := token.SignedString(key)

	return tokenString, nil
}
