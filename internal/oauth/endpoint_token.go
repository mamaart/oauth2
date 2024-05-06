package oauth

import (
	"encoding/json"
	"net/http"

	"github.com/mamaart/oauth2/internal/claims"
	"golang.org/x/oauth2"
)

// TOKEN ENDPOINT
// used by the client to obtain an access token
// no cookie sessions here
func (o *OAuth) Token(w http.ResponseWriter, r *http.Request) {
	clientID := r.FormValue("client_id")
	client, err := o.clientDB.Client(clientID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	clientSecret := r.FormValue("client_secret")
	if clientSecret != client.Secret {
		http.Error(w, "oauth client not authorized", http.StatusUnauthorized)
		return
	}

	code := r.FormValue("code")
	if code == "" {
		http.Error(w, "no authorization code provided", http.StatusNotAcceptable)
		return
	}

	codeVerifier := r.FormValue("code_verifier")
	if codeVerifier == "" {
		http.Error(w, "no code verifier provided", http.StatusNotAcceptable)
		return
	}

	codeChallenge, err := o.clientDB.CheckAuthorizationCode(clientID, code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if oauth2.S256ChallengeFromVerifier(codeVerifier) != codeChallenge {
		http.Error(w, "code veifier does not match code challenge", http.StatusUnauthorized)
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
	if err := json.NewEncoder(w).Encode(returnData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
