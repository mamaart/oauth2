package ports

import (
	"errors"

	"github.com/mamaart/oauth2/internal/models"
)

type UserAuthorizer interface {
	Login(username, password string) error
}

type ClientDB interface {
	Client(id string) (*models.Client, error)
	SetAuthorizationCode(clientID, code, codeChallenge string) error
	CheckAuthorizationCode(clientID, code string) (string, error)
	AddClient(models.Client) error
}

type Session interface {
	ID() string
	IsAuthorized() (bool, string)
	SetAuthorized(username string)
	StoreParams(params models.OAuthParams)
	GetParams() (params models.OAuthParams)
}

var (
	ErrClientNotFound     = errors.New("client not found")
	ErrCodeNotFound       = errors.New("code not found")
	ErrInvalidScope       = errors.New("invalid scope")
	ErrUnexpected         = errors.New("unexpected error")
	ErrAlreadyExist       = errors.New("already exist")
	ErrEmptyClientID      = errors.New("empty client id")
	ErrEmptyCodeChallenge = errors.New("empty code challenge")
	ErrEmptyVerifier      = errors.New("empty verifier")
	ErrEmptyCode          = errors.New("empty code")
	ErrCodeExpired        = errors.New("code expired")
	ErrUnauthorized       = errors.New("unauthorized")
)
