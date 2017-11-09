package domain

type Book struct {
	ID             string `datastore:"-"`
	SubscriptionID string // link to the payments preocessor for multi user books
	Name           string
}
