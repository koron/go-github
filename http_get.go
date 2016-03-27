package github

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func httpGet(url string) ([]byte, error) {
	if DisableCache {
		return httpRawGet(url)
	}
	ce, created := getCacheEntry(url)
	// GET via cacheEntry
	if created {
		return ce.get()
	}
	ce.cond.L.Lock()
	for ce.processing() {
		ce.cond.Wait()
	}
	ce.cond.L.Unlock()
	return ce.data, ce.err
}

func httpRawGet(url string) ([]byte, error) {
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d for %s", r.StatusCode, url)
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}
