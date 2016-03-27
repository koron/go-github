package github

import (
	"sync"
	"time"
)

var (
	// DisableCache disables HTTP response cache.
	DisableCache = false

	// CacheTTL set time to live of cached responses.
	CacheTTL = time.Minute * 5
)

type cacheEntry struct {
	url  string
	time time.Time
	err  error
	data []byte

	cond *sync.Cond
}

func newCacheEntry(url string) *cacheEntry {
	return &cacheEntry{
		url:  url,
		time: time.Now(),
		cond: sync.NewCond(new(sync.Mutex)),
	}
}

func (ce *cacheEntry) get() ([]byte, error) {
	ce.data, ce.err = httpRawGet(ce.url)
	ce.cond.L.Lock()
	ce.cond.Broadcast()
	ce.cond.L.Unlock()
	return ce.data, ce.err
}

func (ce *cacheEntry) processing() bool {
	return ce.data == nil && ce.err == nil
}

func (ce *cacheEntry) expired() bool {
	return time.Since(ce.time) > CacheTTL
}

var (
	cacheMap  = make(map[string]*cacheEntry)
	cacheLock = new(sync.Mutex)
)

func getCacheEntry(url string) (ce *cacheEntry, created bool) {
	cacheLock.Lock()
	ce, ok := cacheMap[url]
	if !ok || ce.expired() {
		ce = newCacheEntry(url)
		cacheMap[url] = ce
		created = true
	}
	cacheLock.Unlock()
	return ce, created
}
