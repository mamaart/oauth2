package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/oauth2"
)

func (s *someApp) root(w http.ResponseWriter, r *http.Request) {
	w.Write(
		[]byte(
			"root endpoint, here you dont have to be logged in to to access <a href=\"http://localhost:2002/secure\">secure</a>",
		),
	)
}

func (s *someApp) oauth(w http.ResponseWriter, r *http.Request) {
	if err := r.FormValue("error"); err != "" {
		description := r.FormValue("error_description")
		w.Write([]byte(fmt.Sprintf("ERROR: %s\nDESCRIPTION: %s\n", err, description)))
		return
	}

	session, err := s.store.New(r, s.cookieName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	verifier, ok := session.Values["verifier"].(string)
	if !ok || verifier == "" {
		http.Error(w, "no verifier code found", http.StatusUnauthorized)
		return
	}

	state := r.FormValue("state")
	if state != "xyz" {
		http.Error(w, "state invalid", http.StatusBadRequest)
		return
	}

	code := r.FormValue("code")
	if code == "" {
		http.Error(w, "code not found", http.StatusBadRequest)
		return
	}
	if _ /*token*/, err := s.config.Exchange(
		context.Background(),
		code,
		oauth2.VerifierOption(verifier),
	); err != nil {
		log.Println(fmt.Errorf("oauth Exchange failed: %w", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["authorized"] = "userIDMartiniBoy"
	session.Save(r, w)

	uri, ok := session.Values["redirect_uri"].(string)
	if ok {
		fmt.Println("INFO: URI IS:", uri)
		http.Redirect(w, r, uri, http.StatusFound)
	} else {
		http.Redirect(w, r, "/secure", http.StatusFound)
	}
}
