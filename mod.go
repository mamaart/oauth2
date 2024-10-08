package oauth2

import (
	"errors"
	"net/http"

	"github.com/mamaart/jwtengine/validator"
	"github.com/mamaart/oauth2/internal/authorizer"
	"github.com/mamaart/oauth2/internal/claims"
	"github.com/mamaart/oauth2/internal/cookiemanager"
	"github.com/mamaart/oauth2/internal/models"
	"github.com/mamaart/oauth2/internal/oauth"
	"github.com/mamaart/oauth2/internal/oauth/token"
	"github.com/mamaart/oauth2/internal/ports"
	"github.com/mamaart/oauth2/internal/userinfo"
)

type OAuthServer struct {
	mux *http.ServeMux

	cookieManager *cookiemanager.Manager
}

type (
	ClientDB       = ports.ClientDB
	UserAuthorizer = ports.UserAuthorizer
	UserDB         = ports.UserDB

	Client      = models.Client
	UserInfo    = models.UserInfo
	OAuthParams = models.OAuthParams

	OAuthClaims      = claims.OAuthClaims
	RefreshValidator = claims.RefreshValidator
)

var (
	ErrMissingClientDB       = errors.New("missing client db")
	ErrMissingViewmodel      = errors.New("missing viewmodel")
	ErrMissingUserAuthorizer = errors.New("missing user authorizer")
	ErrMissingUserDB         = errors.New("missing user db")
	ErrUnauthorized          = ports.ErrUnauthorized
	ErrNotFound              = errors.New("not found")
)

func New(opts Opts) (*OAuthServer, error) {

	if err := opts.validate(); err != nil {
		return nil, err
	}

	validator, err := validator.NewValidator(opts.Issuer.PublicKeyRAW())
	if err != nil {
		panic(err)
	}

	cm := cookiemanager.New("cookieman")
	mux := http.NewServeMux()
	userAuth := authorizer.New(cm, opts.UserAuthorizer, opts.Viewmodel)

	mux.Handle("GET /authorize", oauth.New(opts.ClientDB, cm))
	mux.Handle("POST /token", token.New(opts.ClientDB, opts.Issuer))

	mux.HandleFunc("GET /auth", userAuth.UI)
	mux.HandleFunc("POST /auth", userAuth.Login)

	ui := userinfo.New(opts.UserDB, validator)
	mux.Handle("GET /userinfo", ui)
	mux.Handle("POST /userinfo", ui)

	return &OAuthServer{
		mux:           mux,
		cookieManager: cm,
	}, nil
}

func (s *OAuthServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
