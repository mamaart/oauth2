package oauth

import (
	"sync"

	"github.com/mamaart/jwtengine/issuer"
	"github.com/mamaart/oauth2/internal/claims"
	"github.com/mamaart/oauth2/internal/cookiemanager"
	"github.com/mamaart/oauth2/internal/ports"
)

type OAuth struct {
	cookieManager *cookiemanager.Manager
	cookieName    string
	clientDB      ports.ClientDB

	clientTokenIssuer *issuer.Issuer[*claims.OAuthClaims]

	codeStore      sync.Map
	challengeStore sync.Map
}

func New(clientDB ports.ClientDB, cookieManager *cookiemanager.Manager) *OAuth {
	iss, err := issuer.NewIssuer[*claims.OAuthClaims](&claims.RefreshValidator{})
	if err != nil {
		panic(err)
	}
	return &OAuth{
		clientTokenIssuer: iss,

		cookieManager: cookieManager,
		clientDB:      clientDB,
	}
}
