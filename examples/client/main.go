package main

import (
	"net/http"

	"github.com/mamaart/oauth2"
	"github.com/mamaart/oauth2/internal/clientdb"
	"github.com/mamaart/oauth2/internal/models"
)

func main() {
	go http.ListenAndServe(":2002", newSomeApp())

	db := clientdb.New()

	db.AddClient(models.Client{
		Id:          "someapp",
		Secret:      "secret",
		RedirectUrl: "http://127.0.0.1:2002/oauth2",
		Scopes:      []string{},
	})

	o, err := oauth2.New(oauth2.Opts{
		ClientDB:       db,
		UserAuthorizer: &UserAuthorizer{},
	})
	if err != nil {
		panic(err)
	}
	http.ListenAndServe(":2000", o)
}
