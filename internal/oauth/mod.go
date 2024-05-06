package oauth

import (
	"sync"

	"github.com/mamaart/oauth2/internal/cookiemanager"
	"github.com/mamaart/oauth2/internal/ports"
)

type OAuth struct {
	cookieManager *cookiemanager.Manager
	cookieName    string
	clientDB      ports.ClientDB

	codeStore      sync.Map
	challengeStore sync.Map
}

func New(clientDB ports.ClientDB, cookieManager *cookiemanager.Manager) *OAuth {
	return &OAuth{
		cookieManager: cookieManager,
		clientDB:      clientDB,
	}
}
