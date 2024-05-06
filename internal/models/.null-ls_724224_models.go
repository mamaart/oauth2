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
}

func (o OAuthParams) IsEmpty() bool {
	return o.State == "" &&
		o.ClientID == "" &&
		o.RedirectURI == "" &&
		o.Scope == "" &&
		o.ResponseType == ""
}

func (o OAuthParams) URL(code string) string {
	var buf bytes.Buffer
	buf.WriteString(o.RedirectURI)
	v := map[string][]string{
		"code":  {code},
		"state": {o.State},
	}.(url.Values)
	v.Add("code", code)
	v.Add("state", o.State)
	if strings.Contains(o.RedirectURI, "?") {
		buf.WriteByte('&')
	} else {
		buf.WriteByte('?')
	}
	buf.WriteString(v.Encode())
	return buf.String()
}
