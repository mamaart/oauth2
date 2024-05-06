package screens

import (
	"embed"
)

type Screens struct {
	fs embed.FS
}

func New() *Screens {
	return &Screens{}
}
