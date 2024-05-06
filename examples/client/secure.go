package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
)

type Secure struct {
	mux        *http.ServeMux
	store      sessions.Store
	cookieName string
	config     oauth2.Config
}

func NewSecure(store sessions.Store, cookieName string, config oauth2.Config) *Secure {
	s := &Secure{
		mux:        http.NewServeMux(),
		store:      store,
		cookieName: cookieName,
		config:     config,
	}
	s.mux.HandleFunc("/", s.main)
	s.mux.HandleFunc("/other", s.other)
	return s
}

func (s *Secure) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Secure endpoing accessed")
	session, err := s.store.New(r, s.cookieName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userID, ok := session.Values["authorized"].(string)
	if !ok || userID == "" {

		verifier := oauth2.GenerateVerifier()

		session.Values["verifier"] = verifier
		session.Values["redirect_uri"] = fmt.Sprintf("/secure%s", r.URL.String())
		session.Save(r, w)

		u := s.config.AuthCodeURL(
			"xyz",
			oauth2.S256ChallengeOption(verifier),
		)
		http.Redirect(w, r, u, http.StatusFound)
		return
	}

	s.mux.ServeHTTP(w, r)
}

func (s *Secure) main(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("main secure endpoint accessed"))
}

func (s *Secure) other(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("other secure endpoint accessed"))
}
