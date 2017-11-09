package domain

type AppUser struct {
	ID           string `datastore:"-"`
	IsLoggedIn   bool   `datastore:"-"`
	Email        string
	Name         string
	Nickname     string
	LoginUrl     string `datastore:"-"`
	LogoutUrl    string `datastore:"-"`
	IsProduction bool   `datastore:"-"`
}
