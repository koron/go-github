package github

import (
	"net/http"
	"sync"
)

// Client packs access context/infomation to github.
type Client struct {
	DisableCache bool

	Username string
	Token    string

	Logger Logger

	l      sync.Mutex
	client *http.Client
	cache  map[string]*cacheEntry
}

var (
	DefaultClient = &Client{}
)

func (c *Client) logf(format string, v ...interface{}) {
	if c.Logger == nil {
		return
	}
	c.Logger.Printf(format, v...)
}

func cacheKey(method, url string) string {
	return method + ";" + url
}

func (c *Client) getCache(method, url string, proc cacheProc) (ce *cacheEntry, created bool) {
	k := cacheKey(method, url)
	c.l.Lock()
	ce, ok := c.cache[k]
	if !ok || ce.expired() {
		ce = newCacheEntry(method, url, proc)
		c.cache[k] = ce
		created = true
	}
	c.l.Unlock()
	return ce, created
}
