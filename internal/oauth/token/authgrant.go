package token

import (
	"net/http"

	"github.com/mamaart/jwtengine"
	"github.com/mamaart/oauth2/internal/claims"
	"github.com/mamaart/oauth2/internal/clienterrors"
)

type AuthGrantError struct {
	HttpStatus  int
	ClientError clienterrors.Error
	Description string
}

func (err *AuthGrantError) Error() string {
	return string(err.ClientError) + ": " + err.Description
}

func (s *s) authGrantFlow(
	clientID, clientSecret, code, codeVerifier string,
) (*jwtengine.Tokens, error) {
	client, err := s.clientDB.Client(clientID)
	if err != nil {
		return nil, &AuthGrantError{
			HttpStatus:  http.StatusUnauthorized,
			ClientError: clienterrors.ErrInvalidClient,
			Description: err.Error(),
		}
	}

	if clientSecret != client.Secret {
		return nil, &AuthGrantError{
			HttpStatus:  http.StatusBadRequest,
			ClientError: clienterrors.ErrUnauthorizedClient,
			Description: "",
		}
	}

	if code == "" {
		return nil, &AuthGrantError{
			HttpStatus:  http.StatusBadRequest,
			ClientError: clienterrors.ErrInvalidGrant,
			Description: "missing authcode",
		}
	}

	//if codeVerifier == "" {
	//	return nil, &AuthGrantError{
	//		HttpStatus:  http.StatusBadRequest,
	//		ClientError: clienterrors.ErrInvalidGrant,
	//		Description: "missing code verifier",
	//	}
	//}

	_, userID, err := s.clientDB.CheckAuthorizationCode(clientID, code)
	if err != nil {
		return nil, &AuthGrantError{
			HttpStatus:  http.StatusBadRequest,
			ClientError: clienterrors.ErrInvalidGrant,
			Description: err.Error(),
		}
	}

	//if oauth2.S256ChallengeFromVerifier(codeVerifier) != codeChallenge {
	//	return nil, &AuthGrantError{
	//		HttpStatus:  http.StatusBadRequest,
	//		ClientError: clienterrors.ErrInvalidGrant,
	//		Description: "code veifier does not match code challenge",
	//	}
	//}

	return s.clientTokenIssuer.IssueTokens(&claims.OAuthClaims{
		User:   userID,
		Client: clientID,
		Scope:  client.Scopes,
	})
}
