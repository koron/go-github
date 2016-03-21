package github

import (
	"fmt"
	"time"
)

// Release represents a release on Github.
type Release struct {
	Name        string
	Draft       bool
	PreRelease  bool      `json:"prerelease"`
	CreatedAt   time.Time `json:"created_at"`
	PublishedAt time.Time `json:"published_at"`
	Assets      []Asset
}

// Latest gets latest release info.
func Latest(owner, repo string) (*Release, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest",
		owner, repo)
	rel := new(Release)
	err := jsonGet(url, rel)
	if err != nil {
		return nil, err
	}
	return rel, nil
}
