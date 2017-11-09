package repositories

import (
	//"appengine"
	//"appengine/datastore"
	//"log"
	"net/http"
	//"strconv"
	"stufftobedone/domain"
	//"time"
)

type SearchRepository struct {
	request *http.Request
	cache   *AppCache
}

func NewSearchRepository(request *http.Request) *SearchRepository {
	r := new(SearchRepository)
	r.request = request
	r.cache = NewAppCache(request)
	return r
}

func (r *SearchRepository) Search(bookID string, search string) ([]domain.Task, error) {
	tasks := []domain.Task{}

	return tasks, nil
}

func (r *SearchRepository) Add(task domain.Task) (error) {
	return nil
}

func (r *SearchRepository) Remove(taskID string) (error) {
	return nil
} 