package oauth

import (
	"encoding/json"
	"net/http"

	"github.com/mamaart/oauth2/internal/claims"
	"github.com/mamaart/oauth2/internal/clienterrors"
	"golang.org/x/oauth2"
)

// TOKEN ENDPOINT
// used by the client to obtain an access token
// no cookie sessions here
func (o *OAuth) Token(w http.ResponseWriter, r *http.Request) {
	clientID := r.FormValue("client_id")
	client, err := o.clientDB.Client(clientID)
	if err != nil {
		// http.Error(w, err.Error(), http.StatusUnauthorized)
		w.WriteHeader(http.StatusBadRequest)
		clienterrors.Write(w, clienterrors.ErrInvalidClient, err.Error())
		return
	}

	clientSecret := r.FormValue("client_secret")
	if clientSecret != client.Secret {
		// http.Error(w, "oauth client not authorized", http.StatusUnauthorized)
		w.WriteHeader(http.StatusBadRequest)
		clienterrors.Write(w, clienterrors.ErrUnauthorizedClient, "")
		return
	}

	code := r.FormValue("code")
	if code != "" {
		w.WriteHeader(http.StatusBadRequest)
		clienterrors.Write(w, clienterrors.ErrInvalidRequest, "missing authcode")
		// http.Error(w, "no authorization code provided", http.StatusNotAcceptable)
		return
	}

	codeVerifier := r.FormValue("code_verifier")
	if codeVerifier == "" {
		// http.Error(w, "no code verifier provided", http.StatusNotAcceptable)
		w.WriteHeader(http.StatusBadRequest)
		clienterrors.Write(w, clienterrors.ErrInvalidRequest, "missing code verifier")
		return
	}

	codeChallenge, err := o.clientDB.CheckAuthorizationCode(clientID, code)
	if err != nil {
		// http.Error(w, err.Error(), http.StatusUnauthorized)
		w.WriteHeader(http.StatusBadRequest)
		clienterrors.Write(w, clienterrors.ErrInvalidRequest, err.Error())
		return
	}

	if oauth2.S256ChallengeFromVerifier(codeVerifier) != codeChallenge {
		// http.Error(w, "code veifier does not match code challenge", http.StatusUnauthorized)
		w.WriteHeader(http.StatusBadRequest)
		clienterrors.Write(
			w,
			clienterrors.ErrInvalidRequest,
			"code veifier does not match code challenge",
		)
		return
	}

	token, err := o.clientTokenIssuer.IssueTokens(&claims.OAuthClaims{})
	if err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		w.WriteHeader(http.StatusBadRequest)
		clienterrors.Write(w, clienterrors.ErrInvalidRequest, err.Error())
		return
	}

	returnData := struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
	}{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresIn:    3600,
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(returnData)
}
