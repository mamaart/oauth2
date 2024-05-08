package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mamaart/oauth2"
	"github.com/mamaart/oauth2/internal/clientdb"
	"github.com/mamaart/oauth2/internal/models"
)

func main() {
	// fmt.Println("Serving OAuth client on port :2002")
	// go http.ListenAndServe(":2002", newSomeApp())

	db := clientdb.New()
	if err := db.AddScope("matrix"); err != nil {
		panic(err)
	}
	db.AddClient(models.Client{
		Id:          "matrix",
		Secret:      "matrix",
		RedirectUrl: "https://localhost/_synapse/client/oidc/callback",
		Scopes:      []string{"matrix"},
	})

	ua := NewUA()

	o, err := oauth2.New(oauth2.Opts{
		ClientDB:       db,
		UserAuthorizer: ua,
		UserDB:         ua,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Serving OAuth server on port :2000")
	log.Fatal(http.ListenAndServeTLS(":2000", "./cert/server.crt", "./cert/server.key", o))
}
