package token

import (
	"net/http"

	"github.com/mamaart/jwtengine/issuer"
	"github.com/mamaart/oauth2/internal/claims"
	"github.com/mamaart/oauth2/internal/clienterrors"
	"github.com/mamaart/oauth2/internal/ports"
)

type s struct {
	clientDB          ports.ClientDB
	clientTokenIssuer *issuer.Issuer[*claims.OAuthClaims]
}

func New(clientDB ports.ClientDB, iss *issuer.Issuer[*claims.OAuthClaims]) *s {
	return &s{
		clientDB:          clientDB,
		clientTokenIssuer: iss,
	}
}

func (s *s) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.FormValue("grant_type") {
	case "authorization_code":
		client_id, client_secret, ok := r.BasicAuth()
		if !ok {
			handleErr(w, &AuthGrantError{
				HttpStatus:  http.StatusUnauthorized,
				ClientError: clienterrors.ErrInvalidClient,
				Description: "missing basic auth",
			})
		}
		if tokens, err := s.authGrantFlow(
			client_id,
			client_secret,
			r.FormValue("code"),
			r.FormValue("code_verifier"),
		); err != nil {
			handleErr(w, err)
		} else {
			writeTokens(w, tokens)
		}
	case "refresh_token":
		if tokens, err := s.clientTokenIssuer.Refresh(r.FormValue("refresh_token")); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		} else {
			writeTokens(w, tokens)
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
		clienterrors.Write(
			w,
			clienterrors.ErrUnsupportedGrantType,
			"server only supports authorization_code",
		)
	}
}
