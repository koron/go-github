package github

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func jsonGet(url string, v interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %d for %s", r.StatusCode, url)
	}
	defer r.Body.Close()
	d := json.NewDecoder(r.Body)
	if err := d.Decode(v); err != nil {
		return err
	}
	return nil
}
