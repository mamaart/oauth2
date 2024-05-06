package models

import (
	"bytes"
	"net/url"
	"strings"
)

type Client struct {
	Id          string
	Secret      string
	RedirectUrl string
	Scopes      []string
}

type OAuthParams struct {
	RedirectURI  string `json:"redirect_uri"`
	State        string `json:"state"`
	Scope        string `json:"scope"`
	ClientID     string `json:"client_id"`
	ResponseType string `json:"response_type"`

	CodeChallengeMethod string `json:"code_challenge_method"`
	CodeChallenge       string `json:"code_challenge"`
}

func (o OAuthParams) IsEmpty() bool {
	vs := []string{o.State, o.ClientID, o.RedirectURI, o.Scope, o.ResponseType}
	return strings.Join(vs, "") == ""
}

func (o OAuthParams) URL(code string) string {
	var buf bytes.Buffer
	buf.WriteString(o.RedirectURI)
	if strings.Contains(o.RedirectURI, "?") {
		buf.WriteByte('&')
	} else {
		buf.WriteByte('?')
	}
	buf.WriteString(url.Values{"code": {code}, "state": {o.State}}.Encode())
	return buf.String()
}
