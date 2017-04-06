package domain

type Registration struct {
	ID string `datastore:"-"`
	Email string
}