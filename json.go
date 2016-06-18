package github

import (
	"encoding/json"
	"time"
)

func (c *Client) jsonGet(url string, pivot time.Time, v interface{}) error {
	b, err := c.httpGet(url, pivot)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, v); err != nil {
		return err
	}
	return nil
}
