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
	return DefaultClient.release(owner, repo, "latest", time.Time{})
}

// LatestIfModifiedSince gets latest release info if modified since pivot.
// If no modification, it returns ErrNotModified.
func LatestIfModifiedSince(owner, repo string, pivot time.Time) (*Release, error) {
	return DefaultClient.release(owner, repo, "latest", pivot)
}

// release get a release information.
func (c *Client) release(owner, repo, relName string, pivot time.Time) (*Release, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/%s",
		owner, repo, relName)
	v := new(Release)
	err := c.jsonGet(url, pivot, v)
	if err != nil {
		return nil, err
	}
	return v, nil
}
