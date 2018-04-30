package usecases

import (
	"net/http"
	"stufftobedone/domain"

	"stufftobedone/repositories"
)

/*BookUseCases ... */
type BookUseCases struct {
	BookRepository *repositories.BookRepository
	Request        *http.Request
}

/*NewBookUseCases ... */
func NewBookUseCases(r *http.Request) *BookUseCases {
	u := new(BookUseCases)
	u.Request = r
	u.BookRepository = repositories.NewBookRepository(r)
	return u
}

/*GetDefault ... for now pass through, refactor later */
func (u *BookUseCases) GetDefault(user domain.AppUser) (domain.Book, error) {
	return u.BookRepository.GetDefault(user)
}
