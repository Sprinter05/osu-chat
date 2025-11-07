package ui

import (
	"net/http"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"github.com/Sprinter05/osu-chat/api"
)

type GUI struct {
	Client *http.Client
	token  *api.Token
}

func (g GUI) Run() {
	a := app.New()
	w := a.NewWindow("Hello World")

	w.SetContent(widget.NewLabel("Hello World!"))
	w.ShowAndRun()
}
