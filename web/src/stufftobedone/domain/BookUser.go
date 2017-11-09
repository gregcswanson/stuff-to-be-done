package domain

type BookUser struct {
	ID        string `datastore:"-"`
	BookID    string
	UserID    string
	Email     string
	IsDefault bool
}
