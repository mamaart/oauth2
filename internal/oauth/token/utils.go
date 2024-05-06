package token

import (
	"encoding/json"
	"net/http"

	"github.com/mamaart/jwtengine"
	"github.com/mamaart/oauth2/internal/clienterrors"
)

func writeTokens(w http.ResponseWriter, tokens *jwtengine.Tokens) {
	returnData := struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
	}{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    3600,
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(returnData)
}

func handleErr(w http.ResponseWriter, err error) {
	switch err.(type) {
	case *AuthGrantError:
		err := err.(*AuthGrantError)
		w.WriteHeader(err.HttpStatus)
		clienterrors.Write(w, err.ClientError, err.Description)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)

	}
	return
}
