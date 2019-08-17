package main

import (
	"github.com/rivo/tview"
)

type ModalDialog struct {
	*tview.Modal
}

func ShowModalDialog(
	text string,
	buttons []string,
	doneFunc func(buttonIndex int, buttonLabel string)) {
	v := &ModalDialog{}
	v.Modal = tview.NewModal().
		SetText(text).
		AddButtons(buttons).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			doneFunc(buttonIndex, buttonLabel)
			app.Pop()
		})
	app.Push(v)
}

func (v *ModalDialog) Content() tview.Primitive {
	return v.Modal
}

func (v *ModalDialog) OnActivate() {
}
