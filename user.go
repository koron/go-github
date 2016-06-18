package github

import (
	"time"
)

// User represents
type User struct {
	Login string
	ID    int64
}

// CurrentUser obtains current login user.
func CurrentUser() (*User, error) {
	return DefaultClient.currentUser(time.Time{})
}

// currentUser obtains current login user.
func (c *Client) currentUser(pivot time.Time) (*User, error) {
	url := "https://api.github.com/user"
	v := new(User)
	err := c.jsonGet(url, pivot, v)
	if err != nil {
		return nil, err
	}
	return v, nil
}
