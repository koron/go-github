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
	cache  *cache
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

func (c *Client) cacheDo(id string, proc cacheProc) (interface{}, error) {
	c.l.Lock()
	if c.cache == nil {
		c.cache = newCache()
	}
	c.l.Unlock()
	return c.cache.do(id, proc)
}
