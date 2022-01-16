package main

import (
	"html/template"
	"net/http"

	"github.com/apex/log"
)

func (s *server) user() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := s.store.Get(r, "session-name")
		if err != nil {
			log.WithError(err).Error("error getting session")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		profile := session.Values["profile"]
		log.WithField("profile", profile).Info("user")

		t, err := template.ParseFS(tmpl, "templates/*.html")
		if err != nil {
			log.WithError(err).Fatal("Failed to parse templates")
		}

		w.Header().Set("Content-Type", "text/html")
		err = t.ExecuteTemplate(w, "logged-in.html", struct {
			// TODO: Shouldn't be an interface
			Profile interface{}
		}{
			profile,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.WithError(err).Fatal("Failed to execute templates")
		}
	}
}
