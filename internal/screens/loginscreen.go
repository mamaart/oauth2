package screens

import (
	"fmt"
	"os"
)

type LoginScreen struct {
	title     string
	href      string
	wrongPass bool
}

func NewLoginScreen(title, href string, wrongPass bool) LoginScreen {
	return LoginScreen{
		title:     title,
		href:      href,
		wrongPass: wrongPass,
	}
}

func (l LoginScreen) Html() []byte {
	data, err := os.ReadFile("./html/login.html")
	if err != nil {
		panic(err)
	}
	show := "none"
	if l.wrongPass {
		show = "block"
	}
	return []byte(fmt.Sprintf(string(data), l.title, show, l.title, l.href))
}
