package github

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	ErrNotModified = errors.New("not modified")
)

func (c *Client) newRequest(method, url string, pivot time.Time, body io.Reader) (*http.Request, error) {
	r, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	if c.Token != "" {
		if c.Username != "" {
			r.SetBasicAuth(c.Username, c.Token)
		} else {
			r.Header.Set("Authorization", "token "+c.Token)
		}
	}
	if !pivot.IsZero() {
		t := pivot.UTC().Format(http.TimeFormat)
		r.Header.Set("If-Modified-Since", t)
	}
	return r, nil
}

func (c *Client) httpClient() *http.Client {
	c.l.Lock()
	if c.client == nil {
		c.client = http.DefaultClient
	}
	c.l.Unlock()
	return c.client
}

func (c *Client) httpDo(req *http.Request) (*http.Response, error) {
	resp, err := c.httpClient().Do(req)
	if err != nil {
		return nil, err
	}
	c.logf("Response.Header=%#v", resp.Header)
	return resp, nil
}

func (c *Client) httpGet0(url string, pivot time.Time) ([]byte, error) {
	req, err := c.newRequest("GET", url, pivot, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpDo(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		c.logf("Response.Body=%s", string(b))
		return b, nil
	case http.StatusNotModified:
		return nil, ErrNotModified
	default:
		return nil, fmt.Errorf("unexpected status %d for %s",
			resp.StatusCode, url)
	}
}

func (c *Client) httpGet(url string, pivot time.Time) ([]byte, error) {
	if c.DisableCache {
		return c.httpGet0(url, pivot)
	}
	v, err := c.cacheDo(url, func(id string) (interface{}, error) {
		return c.httpGet0(url, pivot)
	})
	if err != nil {
		return nil, err
	}
	return v.([]byte), err
}
