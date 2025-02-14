package utils

import (
	"net/http"

	"golang.org/x/net/http2"
)

func NewHttp2Client() *http.Client {
	return &http.Client{
		Transport: &http2.Transport{},
	}
}
