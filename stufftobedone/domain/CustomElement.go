package domain

type BookElement struct {
    Name string
    ElementName string
    Icon string
}

type CustomElement struct {
	ID string `datastore:"-"`
	UserID string // owner of the element
    IsPublic bool
	Name string
    ElementName string
    Icon string
    Body string
    LatestElementVersionID string
}

type CustomElementVersion struct {
	ID string `datastore:"-"`
    CustomElementID string
    Css string
    Html string
    Javascript string
    Icon string
}