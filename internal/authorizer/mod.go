package authorizer

import (
	"github.com/mamaart/oauth2/internal/cookiemanager"
	"github.com/mamaart/oauth2/internal/ports"
)

type Authorizer struct {
	cookieManager  *cookiemanager.Manager
	userAuthorizer ports.UserAuthorizer
}

func New(cm *cookiemanager.Manager, userAuthorizer ports.UserAuthorizer) *Authorizer {
	return &Authorizer{cookieManager: cm, userAuthorizer: userAuthorizer}
}
