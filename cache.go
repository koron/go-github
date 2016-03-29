package github

import (
	"sync"
	"time"
)

type cache struct {
	l       sync.Mutex
	entries map[string]*cacheEntry
}

func newCache() *cache {
	return &cache{
		entries: make(map[string]*cacheEntry),
	}
}

func (ca *cache) do(id string, proc cacheProc) (interface{}, error) {
	ca.l.Lock()
	ce, ok := ca.entries[id]
	if !ok || ce.expired() {
		ce = newCacheEntry(id, proc)
		ca.entries[id] = ce
		ca.l.Unlock()
		return ce.callProc()
	}
	ca.l.Unlock()
	return ce.waitProc()
}

var (
	// CacheTTL set time to live of cached responses.
	CacheTTL = time.Minute * 5
)

type cacheProc func(id string) (interface{}, error)

type cacheEntry struct {
	id   string
	proc cacheProc
	cond *sync.Cond
	time time.Time

	data interface{}
	err  error
}

func newCacheEntry(id string, proc cacheProc) *cacheEntry {
	return &cacheEntry{
		id:   id,
		proc: proc,
		cond: sync.NewCond(new(sync.Mutex)),
		time: time.Now(),
	}
}

func (ce *cacheEntry) callProc() (interface{}, error) {
	ce.data, ce.err = ce.proc(ce.id)
	ce.cond.L.Lock()
	ce.cond.Broadcast()
	ce.cond.L.Unlock()
	return ce.data, ce.err
}

func (ce *cacheEntry) waitProc() (interface{}, error) {
	ce.cond.L.Lock()
	for ce.processing() {
		ce.cond.Wait()
	}
	ce.cond.L.Unlock()
	return ce.data, ce.err
}

func (ce *cacheEntry) processing() bool {
	return ce.data == nil && ce.err == nil
}

func (ce *cacheEntry) expired() bool {
	return time.Since(ce.time) > CacheTTL
}
