package cookiemanager

import (
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/mamaart/oauth2/internal/models"
	"github.com/mamaart/oauth2/internal/ports"
	"github.com/wader/gormstore/v2"
	"golang.org/x/oauth2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Manager struct {
	store      sessions.Store
	cookieName string
}

func New(name string) *Manager {
	db, err := gorm.Open(sqlite.Open("./cookies.sql"))
	if err != nil {
		panic(err)
	}
	return &Manager{
		store:      gormstore.New(db, []byte("hello")),
		cookieName: name,
	}
}

type session struct {
	inner *sessions.Session
	w     http.ResponseWriter
	r     *http.Request
}

func (m *Manager) Session(r *http.Request, w http.ResponseWriter) (*session, error) {
	s, err := m.store.New(r, m.cookieName)
	if err != nil {
		return nil, err
	}
	log.Println("session with id:", s.ID, "created")
	return &session{inner: s, w: w, r: r}, nil
}

func (s *session) ID() string {
	return s.inner.ID
}

func (s *session) CheckCodeVerifier(verifier string) error {
	if verifier == "" {
		return ports.ErrEmptyVerifier
	}
	c1, ok := s.inner.Values["code_challenge"].(string)
	if !ok {
		return errors.New("missing code challenge from authorization")
	}
	c2 := oauth2.S256ChallengeFromVerifier(verifier)

	if c1 != c2 {
		return errors.New("code veifier does not match code challenge")
	}
	return nil
}

func (s *session) IsAuthorized() (bool, string) {
	id, ok := s.inner.Values["userID"].(string)
	if !ok {
		return false, ""
	}
	return id != "", id
}

func (s *session) SetAuthorized(username string) {
	s.inner.Values["userID"] = username
	s.save()
}

func (s *session) StoreParams(params models.OAuthParams) {
	s.inner.Values["client_id"] = params.ClientID
	s.inner.Values["redirect_uri"] = params.RedirectURI
	s.inner.Values["scope"] = params.Scope
	s.inner.Values["state"] = params.State
	s.inner.Values["response_type"] = params.ResponseType
	s.inner.Values["code_challenge_method"] = params.CodeChallengeMethod
	s.inner.Values["code_challenge"] = params.CodeChallenge
	s.save()
}

func (s *session) save() {
	if err := s.inner.Save(s.r, s.w); err != nil {
		log.Fatal(err)
	}
}

func (s *session) GetParams() (params models.OAuthParams) {
	if x, ok := s.inner.Values["client_id"]; ok {
		params.ClientID = x.(string)
	}
	if x, ok := s.inner.Values["redirect_uri"]; ok {
		params.RedirectURI = x.(string)
	}
	if x, ok := s.inner.Values["scope"]; ok {
		params.Scope = x.(string)
	}
	if x, ok := s.inner.Values["state"]; ok {
		params.State = x.(string)
	}
	if x, ok := s.inner.Values["response_type"]; ok {
		params.ResponseType = x.(string)
	}
	if x, ok := s.inner.Values["code_challenge_method"]; ok {
		params.CodeChallengeMethod = x.(string)
	}
	if x, ok := s.inner.Values["code_challenge"]; ok {
		params.CodeChallenge = x.(string)
	}

	return params
}
