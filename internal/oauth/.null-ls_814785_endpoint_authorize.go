package oauth

import (
	"fmt"
	"net/http"

	"github.com/mamaart/oauth2/internal/models"
)

// AUTHORIZATION ENDPOINT
// Used to interact with the resource owner
// and obtain an authorization_code.
//
// - Client has to be valid
//
// If user is not authenticated:
//   - store the auth params
//   - redirect to login screen
//
// Otherwise:
//   - validate code challenge
//   - generate authorization code
//   - redirect to client_redirect_uri
func (o *OAuth) Authorize(w http.ResponseWriter, r *http.Request) {
	session, err := o.cookieManager.Session(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	params := models.OAuthParams{
		ClientID:     r.FormValue("client_id"),
		RedirectURI:  r.FormValue("redirect_uri"),
		Scope:        r.FormValue("scope"),
		State:        r.FormValue("state"),
		ResponseType: r.FormValue("response_type"),
	}

	if _, err := o.clientDB.GetClient(params.ClientID); err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	if !session.IsAuthorized() {
		session.StoreParams(params)
		// Redirect to user login endpoint
		http.Redirect(w, r, fmt.Sprintf("/auth?client_id=%s", params.ClientID), http.StatusFound)
		return
	}

	if params.IsEmpty() {
		params = session.GetParams()
		if params.IsEmpty() {
			http.Error(w, "no params in session", http.StatusNotAcceptable)
			return
		}
	}

	if err := session.CheckCodeChallenge(
		r.FormValue("code_challenge"),
		r.FormValue("code_challenge_method"),
	); err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	authorizationCode := "hej" // TODO: generate ?
	session.SaveAuthCode(authorizationCode)
	http.Redirect(w, r, params.URL(authorizationCode), http.StatusFound)
}
