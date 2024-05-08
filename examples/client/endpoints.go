package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/mamaart/oauth2/internal/models"
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
	tokens, err := s.config.Exchange(
		context.Background(),
		code,
		oauth2.VerifierOption(verifier),
	)
	if err != nil {
		log.Println(fmt.Errorf("oauth Exchange failed: %w", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequest(http.MethodGet, "http://localhost:2000/userinfo", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tokens.AccessToken))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Println(resp.Status)
		log.Println(resp.Header.Values("WWW-Authenticate"))
	}

	var userInfo models.UserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		panic(err)
	}

	session.Values["authorized"] = userInfo.PreferredUsername
	session.Save(r, w)

	uri, ok := session.Values["redirect_uri"].(string)
	if ok {
		fmt.Println("INFO: URI IS:", uri)
		http.Redirect(w, r, uri, http.StatusFound)
	} else {
		http.Redirect(w, r, "/secure", http.StatusFound)
	}
}
