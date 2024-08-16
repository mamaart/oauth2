package authorizer

import (
	"github.com/mamaart/oauth2/internal/cookiemanager"
	"github.com/mamaart/oauth2/internal/ports"
	"github.com/mamaart/viewmodel"
)

type Authorizer struct {
	cookieManager  *cookiemanager.Manager
	userAuthorizer ports.UserAuthorizer
	vmFn           func(clientID string, err error) viewmodel.Root
}

func New(cm *cookiemanager.Manager, userAuthorizer ports.UserAuthorizer, vmFn func(clientID string, err error) viewmodel.Root) *Authorizer {
	return &Authorizer{cookieManager: cm, userAuthorizer: userAuthorizer, vmFn: vmFn}
}
