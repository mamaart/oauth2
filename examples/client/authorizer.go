package main

import (
	"github.com/mamaart/oauth2"
	"github.com/mamaart/oauth2/internal/models"
)

type User struct {
	Password string
	Info     models.UserInfo
}

type UserAuthorizer struct{ users map[string]User }

func NewUA() *UserAuthorizer {
	return &UserAuthorizer{
		users: map[string]User{
			"test": {
				Password: "test",
				Info: models.UserInfo{
					Sub:               "000",
					Name:              "Martin Maartensson",
					GivenName:         "Martin",
					FamilyName:        "Maartensson",
					PreferredUsername: "test",
					Email:             "martinmaartensson@gmail.com",
					Picture:           "none",
				},
			},
		},
	}
}

func (ua *UserAuthorizer) Login(username, password string) error {
	user, ok := ua.users[username]
	if !ok || user.Password != password {
		return oauth2.ErrUnauthorized
	}
	return nil
}

func (ua *UserAuthorizer) UserInfo(username string) (models.UserInfo, error) {
	user, ok := ua.users[username]
	if !ok {
		return user.Info, oauth2.ErrNotFound
	}
	return user.Info, nil
}
