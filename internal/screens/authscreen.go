package screens

import (
	"fmt"
	"os"
)

type AuthScreen struct {
	title string
	href  string
}

func NewAuthScreen(title, href string) AuthScreen {
	return AuthScreen{
		title: title,
		href:  href,
	}
}

func (a AuthScreen) Html() []byte {
	data, err := os.ReadFile("./html/auth.html")
	if err != nil {
		panic(err)
	}
	return []byte(fmt.Sprintf(string(data), a.title, a.title, a.href))
}
