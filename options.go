package oauth2

import (
	"github.com/mamaart/jwtengine/issuer"
	"github.com/mamaart/viewmodel"
)

type Opts struct {
	ClientDB       ClientDB
	UserAuthorizer UserAuthorizer
	UserDB         UserDB
	Issuer         *issuer.Issuer[*OAuthClaims]
	Viewmodel      func(clientID string, err error) viewmodel.Root
	// UserValidatorPublicKey crypto.PublicKey
}

func (opts Opts) validate() error {
	if opts.ClientDB == nil {
		return ErrMissingClientDB
	}

	if opts.UserAuthorizer == nil {
		return ErrMissingUserAuthorizer
	}

	if opts.UserDB == nil {
		return ErrMissingUserDB
	}

	if opts.Viewmodel == nil {
		return ErrMissingViewmodel
	}

	if opts.Issuer == nil {
		issuer, err := issuer.NewIssuer(&RefreshValidator{})
		if err != nil {
			panic(err)
		}
		opts.Issuer = issuer
	}
	return nil
}
