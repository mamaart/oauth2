package clientdb

import "gorm.io/gorm"

type Client struct {
	gorm.Model

	ClientID           string `gorm:"unique"`
	Secret             string
	RedirectUrl        string
	Scopes             []Scope `gorm:"many2many:client_scopes;"`
	AuthorizationCodes []AuthorizationCode
}

type Scope struct {
	gorm.Model

	Name string `gorm:"unique"`
}

type AuthorizationCode struct {
	gorm.Model

	ClientID      uint   `gorm:"uniqueIndex:client_code"`
	Code          string `gorm:"uniqueIndex:client_code"`
	CodeChallenge string
}
