package balancer

import "net/url"

type Backend struct {
	URL *url.URL
}

func NewBackend(rawURL string) *Backend {
	parsed, _ := url.Parse(rawURL)
	return &Backend{URL: parsed}
}
