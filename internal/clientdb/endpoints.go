package clientdb

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/mamaart/oauth2/internal/models"
	"github.com/mamaart/oauth2/internal/ports"
	"gorm.io/gorm"
)

func (c *clientDB) Client(id string) (*models.Client, error) {
	client, err := c.getClient(id)
	if err != nil {
		return nil, err
	}

	scopes := make([]string, len(client.Scopes))
	for i := range client.Scopes {
		scopes[i] = client.Scopes[i].Name
	}

	return &models.Client{
		Id:          client.ClientID,
		Secret:      client.Secret,
		RedirectUrl: client.RedirectUrl,
		Scopes:      scopes,
	}, nil
}

func (c *clientDB) getClient(id string) (*Client, error) {
	if id == "" {
		return nil, ports.ErrEmptyClientID
	}
	var client Client
	if err := c.db.Preload("Scopes").First(&client, "client_id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ports.ErrClientNotFound
		}
		return nil, errors.Join(ports.ErrUnexpected, err)
	}
	return &client, nil
}

func (c *clientDB) SetAuthorizationCode(clientID, code, codeChallenge, userID string) error {
	if code == "" {
		return ports.ErrEmptyCode
	}

	if clientID == "" {
		return ports.ErrEmptyClientID
	}

	//if codeChallenge == "" {
	//	return ports.ErrEmptyCodeChallenge
	//}

	if userID == "" {
		return ports.ErrEmptyUserID
	}

	client, err := c.getClient(clientID)
	if err != nil {
		return ports.ErrClientNotFound
	}

	if err := c.db.Create(&AuthorizationCode{
		ClientID:      client.ID,
		Code:          code,
		CodeChallenge: codeChallenge,
		UserID:        userID,
	}).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ports.ErrAlreadyExist
		}
		return errors.Join(err, ports.ErrUnexpected)
	}
	return nil
}

func (c *clientDB) CheckAuthorizationCode(clientID, code string) (string, string, error) {
	if code == "" {
		return "", "", ports.ErrEmptyCode
	}

	client, err := c.getClient(clientID)
	if err != nil {
		return "", "", ports.ErrEmptyClientID
	}

	var authCode AuthorizationCode
	if err := c.db.First(&authCode, "client_id = ? AND code == ?", client.ID, code).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", "", ports.ErrCodeNotFound
		}
		return "", "", errors.Join(err, ports.ErrUnexpected)
	}

	// Maximum 10 minutes allowed in OAuth
	if time.Now().Sub(authCode.CreatedAt).Seconds() > 60 {
		log.Println("authorization code is expired after one minute")
		return "", "", ports.ErrCodeExpired
	}

	return authCode.CodeChallenge, authCode.UserID, nil
}

func (c *clientDB) AddScope(scope string) error {
	if err := c.db.Create(&Scope{
		Name: scope,
	}).Error; err != nil {
		return err
	}
	return nil
}

func (c *clientDB) AddClient(client models.Client) error {
	scopes := make([]Scope, len(client.Scopes))
	for i, e := range client.Scopes {
		if err := c.db.First(&scopes[i], "name = ?", e).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("%w: %s not found", ports.ErrInvalidScope, e)
			}
		}
	}

	if err := c.db.Create(&Client{
		ClientID:    client.Id,
		RedirectUrl: client.RedirectUrl,
		Secret:      client.Secret,
		Scopes:      scopes,
	}).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errors.Join(err, ports.ErrAlreadyExist)
		}
		return errors.Join(err, ports.ErrUnexpected)
	}
	return nil
}
