package main

import (
	"context"
	"net/http"

	"github.com/apex/log"
)

func (s *server) callback() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("callback")

		session, err := s.store.Get(r, "session-name")
		if err != nil {
			log.WithError(err).Error("error getting session")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		state := session.Values["state"]

		// get state from http get state parameter
		httpState := r.URL.Query().Get("state")

		if state != httpState {
			log.WithFields(log.Fields{
				"state":     state,
				"httpState": httpState,
			}).Error("state mismatch")
			http.Error(w, "state mismatch", http.StatusBadRequest)
			return
		}

		code := r.URL.Query().Get("code")
		token, err := s.auth.Exchange(context.Background(), code)
		if err != nil {
			log.WithError(err).Error("error exchanging code")
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		idToken, err := s.auth.VerifyIDToken(context.Background(), token)
		if err != nil {
			log.WithError(err).Error("failed to verify id token")
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		var profile map[string]interface{}
		if err := idToken.Claims(&profile); err != nil {
			log.WithError(err).Error("failed to parse id token claims")
			return
		}

		session.Values["access_token"] = token.AccessToken
		session.Values["profile"] = profile

		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Info("callback done")
		// redirect to /user
		http.Redirect(w, r, "/user", http.StatusFound)

	}
}
