package github

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// Client packs access context/infomation to github.
type Client struct {
	DisableCache bool

	Username string
	Token    string

	Logger Logger

	client *http.Client
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

func (c *Client) newRequest(method, url string, body io.Reader) (*http.Request, error) {
	r, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	if c.Username != "" && c.Token != "" {
		r.SetBasicAuth(c.Username, c.Token)
	}
	return r, nil
}

func (c *Client) httpDo(req *http.Request) (*http.Response, error) {
	if c.client == nil {
		c.client = http.DefaultClient
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	c.logf("Response.Header=%#v", resp.Header)
	return resp, nil
}

func (c *Client) httpGet(url string) ([]byte, error) {
	req, err := c.newRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpDo(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d for %s",
			resp.StatusCode, url)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	c.logf("Response.Body=%s", string(b))
	return b, nil
}

func (c *Client) jsonGet(url string, v interface{}) error {
	b, err := c.httpGet(url)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, v); err != nil {
		return err
	}
	return nil
}
