package main

import (
	"errors"
	"fmt"
	"sync"

	"github.com/mamaart/oauth2/internal/models"
)

type ClientDB struct {
	clients sync.Map
}

type cli struct {
	client    *models.Client
	codeStore sync.Map
}

func New() *ClientDB {
	c := &ClientDB{}

	c.clients.Store("someapp", &cli{
		client: &models.Client{
			Id:          "someapp",
			Secret:      "secret",
			RedirectUrl: "http://127.0.0.1:2002/oauth2",
			Scopes:      []string{},
		},
	})

	return c
}

func (c *ClientDB) Client(id string) (*models.Client, error) {
	fmt.Println(id)

	if c, ok := c.clients.Load(id); ok {
		fmt.Println("hej")
		if c, ok := c.(*cli); ok {
			return c.client, nil
		}
	}
	fmt.Println("not approved:", id)

	return nil, errors.New("client not found")
}

func (c *ClientDB) SetAuthorizationCode(clientID, code string) {
	if c, ok := c.clients.Load(clientID); ok {
		if c, ok := c.(*cli); ok {
			c.codeStore.Store(code, true)
		}
	}
}

func (c *ClientDB) CheckAuthorizationCode(clientID, code string) bool {
	if code == "" {
		return false
	}

	if c, ok := c.clients.Load(clientID); ok {
		if c, ok := c.(*cli); ok {
			_, ok := c.codeStore.Load(code)
			return ok
		}
	}

	return false
}
