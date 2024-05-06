package cookies

import (
	"errors"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type sessionValues struct {
	gorm.Model

	Accepted  bool
	SessionID string `gorm:"unique"`
	Username  string

	Code      string
	Challenge string
}

func newStore() *store {
	db, err := gorm.Open(sqlite.Open("session.db"))
	if err != nil {
		panic(err)
	}

	if err := db.AutoMigrate(&sessionValues{}); err != nil {
		panic(err)
	}

	return &store{
		db: db,
	}
}

func (s *store) acceptCookies(id string) error {
	ses, err := s.getSession(id)
	if err != nil {
		return err
	}
	ses.Accepted = true
	if res := s.db.Save(ses); res.Error != nil {
		return fmt.Errorf("failed to save changes to db %w ", res.Error)
	}
	return nil
}

func (s *store) addSession(id string) error {
	if res := s.db.Create(&sessionValues{SessionID: id}); res.Error != nil {
		return fmt.Errorf("failed to create session in db: %w ", res.Error)
	}
	return nil
}

func (s *store) getSession(id string) (*sessionValues, error) {
	var ses sessionValues
	if res := s.db.First(&ses, "session_id = ?", id); res.Error != nil {
		return nil, errors.Join(ErrSessionNotFound, res.Error)
	}
	return &ses, nil
}
