package main

import (
	"fmt"

	"github.com/rivo/tview"
)

type RenameContextView struct {
	originalName string
	form         *tview.Form
}

func NewRenameContextView(contextName string) *RenameContextView {
	v := &RenameContextView{
		originalName: contextName,
	}
	v.form = tview.NewForm().
		AddInputField("Context", v.originalName, 60, nil, nil).
		AddButton("Rename", func() {
			newName := v.form.GetFormItem(0).(*tview.InputField).GetText()
			app.Pop()
			ShowModalDialog(
				fmt.Sprintf("Rename %q to %q?", v.originalName, newName),
				[]string{"No", "Yes"},
				func(_ int, label string) {
					if label == "Yes" {
						renameContext(v.originalName, newName)
					}
				})
		}).
		AddButton("Cancel", app.Pop).
		SetCancelFunc(app.Pop)

	return v
}

func (v *RenameContextView) Content() tview.Primitive {
	return v.form
}

func (v *RenameContextView) OnActivate() {
}
