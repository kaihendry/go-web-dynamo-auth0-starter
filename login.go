package main

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/apex/log"
)

func (s *server) login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("login")

		state, err := generateRandomState()
		if err != nil {
			log.WithError(err).Error("failed to generate random state")
			http.Error(w, "failed to generate random state", http.StatusInternalServerError)
			return
		}

		session, err := s.store.Get(r, "session-name")
		if err != nil {
			log.WithError(err).Error("failed to get session")
			return
		}

		session.Values["state"] = state

		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, s.auth.AuthCodeURL(state), http.StatusTemporaryRedirect)
	}
}

func generateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	state := base64.StdEncoding.EncodeToString(b)

	return state, nil
}
