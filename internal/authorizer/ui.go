package authorizer

import (
	"net/http"

	"github.com/mamaart/viewmodel"
)

func (a *Authorizer) UI(w http.ResponseWriter, r *http.Request) {
	clientID := r.FormValue("client_id")
	if clientID == "" {
		http.Error(w, "missing client_id", http.StatusNotAcceptable)
		return
	}

	vm := newVM(clientID)

	if err := r.FormValue("error"); err != "" {
		vm.ErrorBox = ErrorBox{
			Title: "Login failed!",
			Errors: []struct{ Message string }{
				{Message: err},
			},
		}
	}

	session, err := a.cookieManager.Session(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if auth, _ := session.IsAuthorized(); !auth {
		viewmodel.New("Login", vm).Execute(w)
	} else {
		http.Redirect(w, r, "/authorize", http.StatusFound)
	}
}
