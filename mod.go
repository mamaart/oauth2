package oauth2

import (
	"errors"
	"log"
	"net/http"

	"github.com/mamaart/jwtengine/issuer"
	"github.com/mamaart/jwtengine/validator"
	"github.com/mamaart/oauth2/internal/authorizer"
	"github.com/mamaart/oauth2/internal/claims"
	"github.com/mamaart/oauth2/internal/cookiemanager"
	"github.com/mamaart/oauth2/internal/oauth"
	"github.com/mamaart/oauth2/internal/oauth/token"
	"github.com/mamaart/oauth2/internal/ports"
	"github.com/mamaart/oauth2/internal/userinfo"
)

type OAuthServer struct {
	mux *http.ServeMux

	cookieManager *cookiemanager.Manager
}

type Opts struct {
	ClientDB       ports.ClientDB
	UserAuthorizer ports.UserAuthorizer
	UserDB         ports.UserDB
	Issuer         *issuer.Issuer[*claims.OAuthClaims]
	// UserValidatorPublicKey crypto.PublicKey
}

var (
	ErrMissingClientDB       = errors.New("missing client db")
	ErrMissingUserAuthorizer = errors.New("missing user authorizer")
	ErrMissingUserDB         = errors.New("missing user db")
	ErrUnauthorized          = ports.ErrUnauthorized
	ErrNotFound              = errors.New("not found")
)

func New(opts Opts) (*OAuthServer, error) {
	if opts.ClientDB == nil {
		return nil, ErrMissingClientDB
	}

	if opts.UserAuthorizer == nil {
		return nil, ErrMissingUserAuthorizer
	}

	if opts.UserDB == nil {
		return nil, ErrMissingUserDB
	}

	if opts.Issuer == nil {
		issuer, err := issuer.NewIssuer[*claims.OAuthClaims](&claims.RefreshValidator{})
		if err != nil {
			panic(err)
		}
		opts.Issuer = issuer
	}

	validator, err := validator.NewValidator(opts.Issuer.PublicKeyRAW())
	if err != nil {
		panic(err)
	}

	var (
		cm       = cookiemanager.New("cookieman")
		mux      = http.NewServeMux()
		userAuth = authorizer.New(cm, opts.UserAuthorizer)
	)

	mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./static/"))))

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
	log.Println(r.URL.String())
	log.Println("PEER:", r.RemoteAddr)
	s.mux.ServeHTTP(w, r)
}
