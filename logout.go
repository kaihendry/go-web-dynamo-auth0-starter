package main

import (
	"net/http"
	"net/url"
	"os"

	"github.com/apex/log"
)

func (s *server) logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		session, err := s.store.Get(r, "session-name")
		if err != nil {
			log.WithError(err).Error("error getting session")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		profile := session.Values["profile"]
		log.WithField("profile", profile).Info("user")

		logoutUrl, err := url.Parse("https://" + os.Getenv("AUTH0_DOMAIN") + "/v2/logout")
		log.WithFields(log.Fields{
			"url":     logoutUrl,
			"profile": profile,
		}).Info("logouting user out")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		scheme := "http"
		if r.Header.Get("X-Forwarded-Proto") == "https" {
			scheme = "https"
		}

		log.WithFields(log.Fields{
			"scheme": scheme,
		}).Info("setting scheme")

		log.WithField("host", r.Host).Info("returning to")
		returnTo, err := url.Parse(scheme + "://" + r.Host)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		parameters := url.Values{}
		parameters.Add("returnTo", returnTo.String())
		parameters.Add("client_id", os.Getenv("AUTH0_CLIENT_ID"))
		logoutUrl.RawQuery = parameters.Encode()

		http.Redirect(w, r, logoutUrl.String(), http.StatusTemporaryRedirect)

	}
}
