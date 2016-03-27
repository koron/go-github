package github

// User represents
type User struct {
	Login string
	ID    int64
}

var (
	// Username used for OAuth when not empty.
	Username string

	// Toen used for OAuth when not empty.
	Token    string
)

func Login(username, token string) (*User, error) {
	url := "https://api.github.com/user"
	user := new(User)
	err := jsonGet(url, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
