package repositories

import (
	"appengine"
	"appengine/memcache"
	"net/http"
	"strconv"
	"strings"
)

type AppCache struct {
	request *http.Request
}

func NewAppCache(request *http.Request) *AppCache {
	r := new(AppCache)
	r.request = request
	return r
}

func (a *AppCache) Prepare(subject string, key string) string {
	keys := []string{subject, key}
	cacheKey := strings.Join(keys, "-")
	return cacheKey
}

func (a *AppCache) Get(subject string, key string) (string, bool) {
	globalContext := appengine.NewContext(a.request)
	cacheKey := a.Prepare(subject, key)
	if item, err := memcache.Get(globalContext, cacheKey); err != memcache.ErrCacheMiss && err == nil {
		return string(item.Value), true
	}
	return "", false
}

func (a *AppCache) GetInt(subject string, key string) (int, bool) {
	value, found := a.Get(subject, key)
	if !found {
		return 0, found
	}
	stringValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, false
	}
	return stringValue, found
}

func (a *AppCache) Set(subject string, key string, value string) error {
	a.Clear(subject, key)
	globalContext := appengine.NewContext(a.request)
	cacheKey := a.Prepare(subject, key)
	item := &memcache.Item{
		Key:   cacheKey,
		Value: []byte(value),
	}
	memcache.Add(globalContext, item)
	return nil
}

func (a *AppCache) SetInt(subject string, key string, value int) error {
	return a.Set(subject, key, strconv.Itoa(value))
}

func (a *AppCache) Clear(subject string, key string) {
	globalContext := appengine.NewContext(a.request)
	cacheKey := a.Prepare(subject, key)
	if _, err := memcache.Get(globalContext, cacheKey); err == memcache.ErrCacheMiss {
		//
	} else if err != nil {
		//
	} else {
		memcache.Delete(globalContext, cacheKey)
	}
}
