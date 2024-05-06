package authorizer

import (
	"fmt"
	"net/http"
)

func (a *Authorizer) Login(w http.ResponseWriter, r *http.Request) {
	clientID := r.FormValue("client_id")
	if clientID == "" {
		http.Error(w, "invalid client_id", http.StatusNotAcceptable)
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

	fmt.Println(username, password)

	if username == "test" && password == "test" {
		session.SetAuthorized(username)
	}

	http.Redirect(w, r, fmt.Sprintf("/auth?client_id=%s", clientID), http.StatusFound)
}
