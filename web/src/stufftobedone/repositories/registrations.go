package repositories

import (
	"appengine"
	"appengine/datastore"
	"log"
	"net/http"
	"stufftobedone/domain"
)

type RegistrationRepository struct {
	request *http.Request
}

func NewRegistrationRepository(request *http.Request) *RegistrationRepository {
	r := new(RegistrationRepository)
	r.request = request
	return r
}

func (r *RegistrationRepository) Get(ID string) (domain.Registration, error) {
	var registration domain.Registration

	// create the namespace context
	c := appengine.NewContext(r.request)
	// get the key
	key, err := datastore.DecodeKey(ID)
	if err != nil {
		return registration, err
	}
	// retrieve the Registration
	err = datastore.Get(c, key, &registration)
	registration.ID = ID

	return registration, err
}

func (r *RegistrationRepository) Store(registration domain.Registration) (domain.Registration, error) {
	// upsert operation
	c := appengine.NewContext(r.request)

	err := datastore.RunInTransaction(c, func(c appengine.Context) error {

		key := datastore.NewKey(c, "Registrations", registration.ID, 0, nil)
		_, err := datastore.Put(c, key, &registration)
		if err != nil {
			return err
		}
		return nil
	}, nil)
	if err != nil {
		return registration, err
	}

	return registration, nil
}

func (repository *RegistrationRepository) Find() ([]domain.Registration, error) {
	var registrations []domain.Registration

	globalContext := appengine.NewContext(repository.request)
	q := datastore.NewQuery("Registrations").Order("Email")

	keys, err := q.GetAll(globalContext, &registrations)

	if err != nil {
		log.Println(err)
		return registrations, err
	} else {
		// loop through and add the keys as ID
		for i := 0; i < len(keys); i++ {
			registrations[i].ID = keys[i].Encode()
		}
	}

	return registrations, nil
}
