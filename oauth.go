package oauth2

import (
	"crypto/sha256"
	"errors"
	"log"
	"net/http"

	"github.com/mamaart/oauth2/internal/authorizer"
	"github.com/mamaart/oauth2/internal/cookiemanager"
	"github.com/mamaart/oauth2/internal/oauth"
	"github.com/mamaart/oauth2/internal/ports"
	"github.com/mamaart/oauth2/pkg/uuid"
)

// var hashKey = []byte("FF51A553-72FC-478B-9AEF-93D6F506DE91")
var hashKey = sha256.Sum256([]byte(uuid.New()))

type OAuthServer struct {
	mux *http.ServeMux

	cookieManager *cookiemanager.Manager
}

type Opts struct {
	ClientDB       ports.ClientDB
	UserAuthorizer ports.UserAuthorizer
	// UserValidatorPublicKey crypto.PublicKey
}

var (
	ErrMissingClientDB       = errors.New("missing client db")
	ErrMissingUserAuthorizer = errors.New("missing user authorizer")
	ErrUnauthorized          = ports.ErrUnauthorized
)

func New(opts Opts) (*OAuthServer, error) {
	if opts.ClientDB == nil {
		return nil, ErrMissingClientDB
	}

	if opts.UserAuthorizer == nil {
		return nil, ErrMissingUserAuthorizer
	}

	var (
		cm       = cookiemanager.New("cookieman")
		mux      = http.NewServeMux()
		oauth    = oauth.New(opts.ClientDB, cm)
		userAuth = authorizer.New(cm, opts.UserAuthorizer)
	)

	mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./static/"))))

	mux.HandleFunc("GET /authorize", oauth.Authorize)
	mux.HandleFunc("POST /token", oauth.Token)

	mux.HandleFunc("GET /auth", userAuth.UI)
	mux.HandleFunc("POST /auth", userAuth.Login)

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
