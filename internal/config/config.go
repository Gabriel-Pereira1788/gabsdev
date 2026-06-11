package config

import (
	"os"
	"strings"
	"sync"
)

var (
	baseURLOnce sync.Once
	baseURL     string
)

func BaseURL() string {
	baseURLOnce.Do(func() {
		baseURL = os.Getenv("BASE_URL")
		if baseURL == "" {
			baseURL = "http://localhost:8080"
		}
		baseURL = strings.TrimRight(baseURL, "/")
	})
	return baseURL
}
