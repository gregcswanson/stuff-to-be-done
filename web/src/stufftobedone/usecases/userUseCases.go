package usecases

import (
	"net/http"
	"stufftobedone/domain"

	"appengine"
	"appengine/user"

	"stufftobedone/repositories"
)

/*UserUseCases ... */
type UserUseCases struct {
	UserRepository *repositories.UserRepository
	Request        *http.Request
}

/*NewUserUseCases ... */
func NewUserUseCases(r *http.Request) *UserUseCases {
	u := new(UserUseCases)
	u.Request = r
	u.UserRepository = repositories.NewUserRepository(r)
	return u
}

/*FindOrCreateByAppEngine ... */
func (r *UserUseCases) FindOrCreateByAppEngine() (domain.AppUser, error) {
	appengineContext := appengine.NewContext(r.Request)
	appengineUser := user.Current(appengineContext)

	// get or create the app user
	appUser := domain.AppUser{}
	appUser.ID = ""
	appUser.Name = appengineUser.ID
	appUser.Nickname = ""
	appUser.Email = appengineUser.Email
	appUser.IsLoggedIn = true
	url, err := user.LogoutURL(appengineContext, "/")
	if err == nil {
		appUser.LogoutUrl = url
	}
	appUser, err = r.UserRepository.Get(appUser)

	return appUser, err
}

/*FindOrCreateByEmail ... */
func (r *UserUseCases) FindOrCreateByEmail(email string) (domain.AppUser, error) {

	// get or create the app user
	appUser := domain.AppUser{}
	appUser.ID = ""
	appUser.Name = email
	appUser.Nickname = ""
	appUser.Email = email
	appUser.IsLoggedIn = true
	foundUser, err := r.UserRepository.Get(appUser)

	return foundUser, err
}
