package github

import "encoding/json"

type jsonGetter func(url string) ([]byte, error)

// jsonGet0 is common logic for jsonGet.
func jsonGet0(url string, v interface{}, getter jsonGetter) error {
	b, err := getter(url)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, v); err != nil {
		return err
	}
	return nil
}

func jsonGet(url string, v interface{}) error {
	return jsonGet0(url, v, httpGet)
}
