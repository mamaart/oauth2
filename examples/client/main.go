package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"

	"github.com/mamaart/oauth2"
	"github.com/mamaart/oauth2/examples/client/internal/clientdb"
	"github.com/mamaart/viewmodel"
)

//go:embed all:templates
var templates embed.FS

type vm struct {
	ClientID string
	Error    string
}

func (vm *vm) Data() viewmodel.VM { return nil }
func (vm *vm) FS() fs.FS          { return templates }

func main() {
	// fmt.Println("Serving OAuth client on port :2002")
	// go http.ListenAndServe(":2002", newSomeApp())

	db := clientdb.New()
	if err := db.AddScope("matrix"); err != nil {
		panic(err)
	}
	db.AddClient(oauth2.Client{
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
		Viewmodel: func(clientID string, err error) viewmodel.Root {
			if err != nil {
				return viewmodel.New("Login", templates, &vm{
					ClientID: clientID,
					Error:    err.Error(),
				})
			}
			return viewmodel.New("Login", templates, &vm{ClientID: clientID})
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Serving OAuth server on port :2000")
	log.Fatal(http.ListenAndServeTLS(":2000", "./cert/server.crt", "./cert/server.key", o))
}
