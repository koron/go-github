package github

import "encoding/json"

func jsonGet(url string, v interface{}) error {
	b, err := httpGet(url)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, v); err != nil {
		return err
	}
	return nil
}
