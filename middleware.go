package main

import (
	"net/http"

	"github.com/apex/log"
)

func (s *server) IsAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := s.store.Get(r, "session-name")
		if err != nil {
			log.WithError(err).Error("error getting session")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		profile := session.Values["profile"]
		if profile == nil {
			log.Error("no profile in session")
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}
