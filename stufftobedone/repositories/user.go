package repositories

import (
	"net/http"
	"appengine"
	"appengine/datastore"
    "stufftobedone/domain"
    "log"
)

type UserRepository struct {
  request    *http.Request
}

func NewUserRepository(request *http.Request) *UserRepository {
	r := new(UserRepository)
	r.request = request
	return r
}

func (r *UserRepository) Get(appUser domain.AppUser)(domain.AppUser, error) {
    // get or create, a partial user is passed in from the authentication
    // use the email address to find the user
    var appUsers []domain.AppUser
    globalContext := appengine.NewContext(r.request)
    q := datastore.NewQuery("AppUsers").Filter("Email =",appUser.Email)
    keys, err := q.GetAll(globalContext, &appUsers)
    
	if err != nil {    
		log.Println(err)
		return appUser, err
	} else if len(appUsers) > 0 {    
		log.Println("Found the user")
        appUser.ID = keys[0].Encode()
		return appUser, err
	}
    // the user profile was not found, save it to the datastore
    return r.Store(appUser)
}

func (r *UserRepository) Store(appUser domain.AppUser) (domain.AppUser, error) {
	// upsert operation
    globalContext := appengine.NewContext(r.request)
    err := datastore.RunInTransaction(globalContext, func(c appengine.Context) error {
        
        if appUser.ID != "" {
			// update
			key , err := datastore.DecodeKey(appUser.ID)
			if err != nil {
                log.Println(err)
				return err
			}
			_, err = datastore.Put(c, key, &appUser)
			if err != nil {
                log.Println(err)
				return err
			}
            log.Println("App user updated")
		} else {
			// new
			key, err := datastore.Put(c, datastore.NewIncompleteKey(c, "AppUsers", nil), &appUser)
			if err != nil {
                log.Println(err)
				return err
			} else {
				appUser.ID = key.Encode()
                log.Println("App user created")
			}
		}
        
        return nil
    }, nil)
    if err != nil {
        log.Println(err)
        return appUser, err
    }
  
	return appUser, nil
}

