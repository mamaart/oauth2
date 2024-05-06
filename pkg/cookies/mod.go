package cookies

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/mamaart/oauth2/pkg/uuid"
	"gorm.io/gorm"
)

// Cookie holds the session id
// Session id points to the session in postgres
// Schema is uniform for all sessions

type (
	Session interface {
		Accepted() bool
		Accept() error
		Values() (*SessionValues, error)
	}

	SessionValues struct {
		Username string
	}

	store struct {
		db *gorm.DB
	}

	session struct {
		store *store

		accepted bool
		id       string
	}

	manager struct {
		name  string
		store *store
	}
)

func New() *manager {
	return &manager{
		name:  "my_cookie_name",
		store: newStore(),
	}
}

func (m *manager) newSession() (*session, error) {
	id := uuid.New()
	if err := m.store.addSession(id); err != nil {
		return nil, fmt.Errorf("could not add session to store: %w", err)
	}
	vals, err := m.store.getSession(id)
	if err != nil {
		return nil, fmt.Errorf("could not get session info from store: %w", err)
	}
	return &session{id: id, store: m.store, accepted: vals.Accepted}, nil
}

func (m *manager) getSession(r *http.Request) (*session, error) {
	c, err := r.Cookie(m.name)
	if err == nil {
		if s, err := m.store.getSession(c.Value); err == nil {
			return &session{id: s.SessionID, store: m.store, accepted: s.Accepted}, nil
		} else {
			log.Println("ERROR:", err)
		}
	}
	return m.newSession()
}

func (m *manager) Session(r *http.Request, w http.ResponseWriter) (*session, error) {
	s, err := m.getSession(r)
	if err != nil {
		return nil, fmt.Errorf("could not get session from manager: %w", err)
	}

	// always update when used (to extend expiry)
	http.SetCookie(w, &http.Cookie{
		Name:     m.name,
		Value:    s.id,
		Domain:   "localhost",
		Expires:  time.Now().Add(time.Hour * 24 * 7), // one week from now
		SameSite: http.SameSiteStrictMode,
	})

	return s, nil
}

func (s *session) Accept() error {
	if err := s.store.acceptCookies(s.id); err != nil {
		return fmt.Errorf("could not accept cookies in store: %w ", err)
	}

	if vals, err := s.store.getSession(s.id); err == nil {
		s.accepted = vals.Accepted
	} else {
		return fmt.Errorf("failed to get session from store: %w", err)
	}

	return nil
}

func (s *session) Accepted() bool {
	return s.accepted
}

func (s *session) Values() (*SessionValues, error) {
	vals, err := s.store.getSession(s.id)
	if err != nil {
		return nil, fmt.Errorf("failed to get session values from store: %w", err)
	}
	return &SessionValues{
		Username: vals.Username,
	}, nil
}
