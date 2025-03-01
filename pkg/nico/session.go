package nico

import (
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/shinosaki/namagent/internal/config"
	"golang.org/x/net/http2"
)

func NewSession(config *config.Config) *http.Client {
	userSession := config.Auth.Nico.UserSession

	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Panicln("failed to create cookiejar", err)
	}

	if userSession != "" {
		origin, _ := url.Parse("https://nicovideo.jp/")
		jar.SetCookies(origin, []*http.Cookie{
			{
				Name:   "user_session",
				Value:  userSession,
				Path:   "/",
				Domain: ".nicovideo.jp",
			},
		})
	}

	client := &http.Client{
		Transport: &http2.Transport{},
		Jar:       jar,
	}

	return client
}
