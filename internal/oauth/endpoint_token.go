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
	grantType := r.FormValue("grant_type")
	if grantType != "authorization_code" {
		w.WriteHeader(http.StatusBadRequest)
		clienterrors.Write(
			w,
			clienterrors.ErrUnsupportedGrantType,
			"server only supports authorization_code",
		)
	}

	clientID := r.FormValue("client_id")
	client, err := o.clientDB.Client(clientID)
	if err != nil {
		// http.Error(w, err.Error(), http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
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
	if code == "" {
		w.WriteHeader(http.StatusBadRequest)
		clienterrors.Write(w, clienterrors.ErrInvalidGrant, "missing authcode")
		// http.Error(w, "no authorization code provided", http.StatusNotAcceptable)
		return
	}

	codeVerifier := r.FormValue("code_verifier")
	if codeVerifier == "" {
		// http.Error(w, "no code verifier provided", http.StatusNotAcceptable)
		w.WriteHeader(http.StatusBadRequest)
		clienterrors.Write(w, clienterrors.ErrInvalidGrant, "missing code verifier")
		return
	}

	codeChallenge, err := o.clientDB.CheckAuthorizationCode(clientID, code)
	if err != nil {
		// http.Error(w, err.Error(), http.StatusUnauthorized)
		w.WriteHeader(http.StatusBadRequest)
		clienterrors.Write(w, clienterrors.ErrInvalidGrant, err.Error())
		return
	}

	if oauth2.S256ChallengeFromVerifier(codeVerifier) != codeChallenge {
		// http.Error(w, "code veifier does not match code challenge", http.StatusUnauthorized)
		w.WriteHeader(http.StatusBadRequest)
		clienterrors.Write(
			w,
			clienterrors.ErrInvalidGrant,
			"code veifier does not match code challenge",
		)
		return
	}

	token, err := o.clientTokenIssuer.IssueTokens(&claims.OAuthClaims{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
