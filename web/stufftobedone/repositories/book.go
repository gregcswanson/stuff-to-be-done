///
// BOOK Repository
// Manage creating the default book for new users
// TO DO
//	- support adding multiple users to a book along with payment processing
//	- invitations to new users - sending emails and automatic linking
///

package repositories

import (
	"net/http"
	"appengine"
	"appengine/datastore"
	"stufftobedone/domain"
	"log"
)

type BookRepository struct {
  request    *http.Request
}

func NewBookRepository(request *http.Request) *BookRepository {
	r := new(BookRepository)
	r.request = request
	return r
}

func (r *BookRepository) GetDefault(user domain.AppUser)(domain.Book, error) {
	book := domain.Book{}
	var bookUsers []domain.BookUser
	
	// find the default book for this user
	globalContext := appengine.NewContext(r.request)
	q := datastore.NewQuery("BookUsers").Filter("UserID =", user.ID).Filter("IsDefault =", true)
	keys, err := q.GetAll(globalContext, &bookUsers)
    
	if err != nil {    
		log.Println(err)
		return book, err
	} else if len(bookUsers) > 0 {    
		log.Println("Found the default book")
		// we only care about the first default book
		book, err = r.GetBook(bookUsers[0].BookID)
		return book, err
	}
    
	// No default book was found to process invitations
	// select first book, set it as default, if none then
	q = datastore.NewQuery("BookUsers").Filter("UserID =", user.ID)
	keys, err = q.GetAll(globalContext, &bookUsers)
	if err != nil {   
		log.Println(err) 
		return book, err
	} else if len(keys) > 0 {    
		// select the first book as the defeault and store it
		log.Println("setting the first linked book as default and returning")
		bookUser := bookUsers[0]
		bookUser.ID = keys[0].Encode()
		bookUser.IsDefault = true
		r.StoreBookUser(bookUser)
		book, err = r.GetBook(bookUsers[0].BookID)
		return book, err
	}

	// create the default book
	log.Println("Creating the default book")
	return r.Create(user, "Stuff to do")
}

func (r *BookRepository) ProcessPendingInvitations(user domain.AppUser) error {
	// TO DO
	return nil
}

func (r *BookRepository) GetBook(bookID string) (domain.Book, error) {
	var book domain.Book

	// create the context
	globalContext := appengine.NewContext(r.request)

	// get the key
	key, err := datastore.DecodeKey(bookID)
	if err != nil {
		return book, err
	}
	// get the item from the datastore
	err = datastore.Get(globalContext, key, &book);
	book.ID = bookID

	return book, err
}

func (r *BookRepository) GetBookUsers(user domain.AppUser, bookID string) (domain.Book, error) {
	book := domain.Book{}
	// TO DO
	// process invitations
	// Check the user has access, if not throw error
	// return the book and set it as default if it is not
	return book, nil
}

func (r *BookRepository) StoreBookUser(bookUser domain.BookUser) (domain.BookUser, error) {
	// upsert operation
	globalContext := appengine.NewContext(r.request)
  
  	err := datastore.RunInTransaction(globalContext, func(c appengine.Context) error {
        
        if bookUser.ID != "" {
			// update
			key , err := datastore.DecodeKey(bookUser.ID)
			if err != nil {
				return err
			}
			_, err = datastore.Put(c, key, &bookUser)
			if err != nil {
				return err
			}
		} else {
			// new
			key, err := datastore.Put(c, datastore.NewIncompleteKey(c, "BookUsers", nil), &bookUser)
			if err != nil {
				return err
			} else {
				bookUser.ID = key.Encode()
			}
		}
        
        return nil
    }, nil)
    if err != nil {
        return bookUser, err
    }
  
	return bookUser, nil
}

func (r *BookRepository) StoreBook(book domain.Book) (domain.Book, error) {
	// upsert operation
	globalContext := appengine.NewContext(r.request)
  
  	err := datastore.RunInTransaction(globalContext, func(c appengine.Context) error {
        
        if book.ID != "" {
			// update
			key , err := datastore.DecodeKey(book.ID)
			if err != nil {
				return err
			}
			_, err = datastore.Put(c, key, &book)
			if err != nil {
				return err
			}
		} else {
			// new
			key, err := datastore.Put(c, datastore.NewIncompleteKey(c, "Books", nil), &book)
			if err != nil {
				return err
			} else {
				book.ID = key.Encode()
			}
		}
        
        return nil
    }, nil)
    if err != nil {
        return book, err
    }
  
	return book, nil
}

func (r *BookRepository) Create(user domain.AppUser, name string) (domain.Book, error) {
	// Create the book
	book := domain.Book{}
	book.Name = name
	createdBook, err := r.StoreBook(book)
	if err != nil {
		return book, nil
	}
	log.Println("Book created")
	// link the user
	bookUser := domain.BookUser{}
	bookUser.BookID = createdBook.ID
	bookUser.Email = user.Email
	bookUser.UserID = user.ID
	bookUser.IsDefault = false // true
	_, err = r.StoreBookUser(bookUser)

	log.Println("BookUser created")
	return createdBook, err
}

