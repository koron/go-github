package github

import (
	"net/http"
	"sync"
	"time"
)

var (
	// DisableCache disables HTTP response cache.
	DisableCache = false

	// CacheTTL set time to live of cached responses.
	CacheTTL = time.Minute * 5
)

type cacheProc func(url string) (http.Header, []byte, error)

type cacheEntry struct {
	method string
	url    string
	proc   cacheProc

	time   time.Time
	err    error
	header http.Header
	data   []byte

	cond *sync.Cond
}

func newCacheEntry(method, url string, proc cacheProc) *cacheEntry {
	return &cacheEntry{
		method: method,
		url:    url,
		proc:   proc,
		time:   time.Now(),
		cond:   sync.NewCond(new(sync.Mutex)),
	}
}

func (ce *cacheEntry) get() (http.Header, []byte, error) {
	ce.header, ce.data, ce.err = ce.proc(ce.url)
	ce.cond.L.Lock()
	ce.cond.Broadcast()
	ce.cond.L.Unlock()
	return ce.header, ce.data, ce.err
}

func (ce *cacheEntry) processing() bool {
	return ce.header == nil && ce.data == nil && ce.err == nil
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
		ce = newCacheEntry("", url, nil)
		cacheMap[url] = ce
		created = true
	}
	cacheLock.Unlock()
	return ce, created
}
