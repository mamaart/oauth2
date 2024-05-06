package authorizer

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/mamaart/oauth2/internal/ports"
)

func (a *Authorizer) Login(w http.ResponseWriter, r *http.Request) {
	clientID := r.FormValue("client_id")
	if clientID == "" {
		http.Error(w, "missing client id", http.StatusBadRequest)
		return
	}

	session, err := a.cookieManager.Session(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Form == nil {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	username := r.Form.Get("username")
	password := r.Form.Get("password")

	if username == "" || password == "" {
		http.Redirect(w, r, fmt.Sprintf("/auth?%s", url.Values{
			"client_id": {clientID},
			"error":     {"missing username or password"},
		}.Encode()), http.StatusFound)
		return
	}

	if err := a.userAuthorizer.Login(username, password); err != nil {
		if errors.Is(err, ports.ErrUnauthorized) {
			http.Redirect(w, r, fmt.Sprintf("/auth?%s", url.Values{
				"client_id": {clientID},
				"error":     {"wrong username or password"},
			}.Encode()), http.StatusFound)
		}
		http.Redirect(w, r, fmt.Sprintf("/auth?%s", url.Values{
			"client_id": {clientID},
			"error":     {err.Error()},
		}.Encode()), http.StatusFound)
	} else {
		session.SetAuthorized(username)
		http.Redirect(w, r, fmt.Sprintf("/auth?client_id=%s", clientID), http.StatusFound)
	}
}
