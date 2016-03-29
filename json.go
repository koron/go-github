package github

import "encoding/json"

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
