package authorizer

import (
	"fmt"
	"net/http"
)

func (a *Authorizer) UI(w http.ResponseWriter, r *http.Request) {
	clientID := r.FormValue("client_id")
	if clientID == "" {
		http.Error(w, "missing client_id", http.StatusNotAcceptable)
		return
	}

	var err error
	if er := r.FormValue("error"); er != "" {
		err = fmt.Errorf("login failed: %s", er)
	}

	vm := a.vmFn(clientID, err)

	session, err := a.cookieManager.Session(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if auth, _ := session.IsAuthorized(); !auth {
		vm.Execute(w)
	} else {
		http.Redirect(w, r, "/authorize", http.StatusFound)
	}
}
