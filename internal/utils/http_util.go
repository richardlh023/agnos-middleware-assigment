package utils

import (
	"net/http"
	"time"
)

type HTTPClientUtil struct {
	Client *http.Client
}

func NewHTTPClientUtil(timeout time.Duration) *HTTPClientUtil {
	return &HTTPClientUtil{
		Client: &http.Client{
			Timeout: timeout,
		},
	}
}

func DefaultHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 10 * time.Second,
	}
}
