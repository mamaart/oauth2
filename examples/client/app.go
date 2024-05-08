package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/wader/gormstore/v2"
	"golang.org/x/oauth2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type someApp struct {
	mux        *http.ServeMux
	config     oauth2.Config
	store      sessions.Store
	cookieName string
}

func newSomeApp() *someApp {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"))
	if err != nil {
		panic(err)
	}
	app := &someApp{
		mux: http.NewServeMux(),
		config: oauth2.Config{
			ClientID:     "someapp",
			ClientSecret: "secret",
			Endpoint: oauth2.Endpoint{
				AuthURL:   "http://127.0.0.1:2000/authorize",
				TokenURL:  "http://127.0.0.1:2000/token",
				AuthStyle: oauth2.AuthStyleAutoDetect,
			},
			RedirectURL: "http://127.0.0.1:2002/oauth2",
			Scopes:      []string{},
		},
		store:      gormstore.New(db, []byte("hejsa")),
		cookieName: "someapp",
	}

	app.mux.HandleFunc("GET /{$}", app.root)
	app.mux.HandleFunc("GET /oauth2", app.oauth)
	secure := NewSecure(app.store, app.cookieName, app.config)
	app.mux.Handle("/secure/", http.StripPrefix("/secure", secure))

	return app
}

func (s *someApp) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.String())
	s.mux.ServeHTTP(w, r)
}
