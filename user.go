package github

import "net/http"

// User represents
type User struct {
	Login string
	ID    int64
}

var (
	// Username used for OAuth when not empty.
	Username string

	// Token used for OAuth when not empty.
	Token string
)

// Login logins as a user.
func Login(username, token string) (*User, error) {
	g := func(url string) (*http.Response, error) {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		req.SetBasicAuth(username, token)
		return http.DefaultClient.Do(req)
	}
	url := "https://api.github.com/user"
	user := new(User)
	err := jsonGet0(url, user, httpGetter(g).toJSONGetter())
	if err != nil {
		return nil, err
	}
	return user, nil
}
