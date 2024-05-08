package userinfo

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/mamaart/jwtengine/issuer"
	"github.com/mamaart/jwtengine/validator"
	"github.com/mamaart/oauth2/internal/ports"
)

type s struct {
	db        ports.UserDB
	validator *validator.Validator
}

func New(db ports.UserDB, v *validator.Validator) *s {
	return &s{db: db, validator: v}
}

func (s *s) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	bearerToken := r.Header.Get("Authorization")

	if bearerToken == "" {
		writeErr(w, ErrInvalidRequest, "missing bearer token", http.StatusUnauthorized)
		return
	}

	token, err := s.validator.GetToken(strings.TrimPrefix(bearerToken, "Bearer "))
	if err != nil || !token.Valid {
		writeErr(w, ErrInvalidToken, err.Error(), http.StatusUnauthorized)
		return
	}

	var info struct {
		UserID   string   `json:"user"`
		ClientID string   `json:"client"`
		Scope    []string `json:"scope"`
	}

	j, err := json.Marshal(issuer.ExtractClaims(token))
	if err != nil {
		log.Println(err)
	}
	if err := json.Unmarshal(j, &info); err != nil {
		log.Println(err)
	}

	if len(info.Scope) == 0 {
		writeErr(w, ErrInvalidToken, "missing scope in token", http.StatusUnauthorized)
		return
	}
	if info.UserID == "" {
		writeErr(w, ErrInvalidToken, "missing username in token", http.StatusUnauthorized)
		return
	}
	if info.ClientID == "" {
		writeErr(w, ErrInvalidToken, "missing clientID in token", http.StatusUnauthorized)
		return
	}

	user, err := s.db.UserInfo(info.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
