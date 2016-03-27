package github

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type httpGetter func(url string) (*http.Response, error)

func (g httpGetter) toJSONGetter() jsonGetter {
	return func(url string) ([]byte, error) {
		return httpGet0(url, g)
	}
}

func httpGet0(url string, getter httpGetter) ([]byte, error) {
	if getter == nil {
		getter = http.Get
	}
	r, err := getter(url)
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
	//log.Printf("HEADER: %#v", r.Header)
	//log.Printf("BODY: %s", string(b))
	return b, nil
}

func httpGet(url string) ([]byte, error) {
	if DisableCache {
		return httpGet0(url, nil)
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
