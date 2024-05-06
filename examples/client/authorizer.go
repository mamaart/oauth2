package main

import "github.com/mamaart/oauth2"

type UserAuthorizer struct{}

func (ua *UserAuthorizer) Login(username, password string) error {
	if username == "test" && password == "test" {
		return nil
	}
	return oauth2.ErrUnauthorized
}
