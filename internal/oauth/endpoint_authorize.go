package oauth

import (
	"fmt"
	"net/http"

	"github.com/mamaart/oauth2/internal/models"
	"github.com/mamaart/oauth2/internal/redirecterrors"
	"golang.org/x/oauth2"
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
func (o *OAuth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

		CodeChallengeMethod: r.FormValue("code_challenge_method"),
		CodeChallenge:       r.FormValue("code_challenge"),
	}
	if params.IsEmpty() {
		params = session.GetParams()
		if params.IsEmpty() {
			http.Error(w, "no params in session", http.StatusNotAcceptable)
			return
		}
	}

	if params.ClientID == "" {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	client, err := o.clientDB.Client(params.ClientID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	if params.RedirectURI == "" || client.RedirectUrl != params.RedirectURI {
		http.Error(w, "redirect uri invalid", http.StatusNotAcceptable)
		return
	}

	// Any errors after this point should redirect to redirect_uri with error in url params

	authorized, _ := session.IsAuthorized()
	if !authorized {
		session.StoreParams(params)
		http.Redirect(w, r, fmt.Sprintf("/auth?client_id=%s", params.ClientID), http.StatusFound)
		return
	}

	if params.CodeChallengeMethod != "S256" {
		u := redirecterrors.URI(
			params.RedirectURI,
			redirecterrors.ErrInvalidRequest,
			"requires challengetype S256",
			params.State,
		)
		http.Redirect(w, r, u, http.StatusFound)
		return
	}

	if params.CodeChallenge == "" {
		u := redirecterrors.URI(
			params.RedirectURI,
			redirecterrors.ErrInvalidRequest,
			"empty code challenge",
			params.State,
		)
		http.Redirect(w, r, u, http.StatusFound)
		return
	}

	// TODO: code must live max 10 minutes and should remove all
	// tokens previously based on that code.
	authorizationCode := oauth2.S256ChallengeFromVerifier(oauth2.GenerateVerifier()) // random
	if err := o.clientDB.SetAuthorizationCode(params.ClientID, authorizationCode, params.CodeChallenge); err != nil {
		u := redirecterrors.URI(
			params.RedirectURI,
			redirecterrors.ErrServerError,
			err.Error(),
			params.State,
		)
		http.Redirect(w, r, u, http.StatusFound)
		return
	}

	http.Redirect(w, r, params.URL(authorizationCode), http.StatusFound)
}
