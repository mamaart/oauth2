package claims

import (
	"fmt"
)

type OAuthClaims struct {
	user   string
	scope  string
	client string
}

func (c *OAuthClaims) AccessClaimsAsMap() map[string]interface{} {
	return map[string]interface{}{
		"user":   c.user,
		"scope":  c.scope,
		"client": c.client,
	}
}

func (c *OAuthClaims) RefreshClaimsAsMap() map[string]interface{} {
	return c.AccessClaimsAsMap()
}

type RefreshValidator struct {
}

func (v *RefreshValidator) Validate(m map[string]interface{}) (*OAuthClaims, error) {
	user, ok := m["user"].(string)
	if !ok {
		return nil, fmt.Errorf("no user in claims")
	}
	scope, ok := m["scope"].(string)
	if !ok {
		return nil, fmt.Errorf("no scope in claims")
	}
	client, ok := m["client"].(string)
	if !ok {
		return nil, fmt.Errorf("no client in claims")
	}
	return &OAuthClaims{
		user:   user,
		scope:  scope,
		client: client,
	}, nil
}
