package main

import (
	"fmt"

	"github.com/rivo/tview"
)

type EditContextView struct {
	context *Context
	form    *tview.Form
}

func NewEditContextView(contextName string) *EditContextView {
	v := &EditContextView{
		context: getConfig().GetContext(contextName),
	}

	v.form = tview.NewForm().
		AddInputField("Context", v.context.Name, 60, nil, nil).
		AddInputField("Namespace", v.context.Context.Namespace, 60, nil, nil).
		AddButton("Apply", v.apply).
		AddButton("Cancel", app.Pop).
		SetCancelFunc(app.Pop)

	return v
}

func (v *EditContextView) Content() tview.Primitive {
	return v.form
}

func (v *EditContextView) OnActivate() {
}

func (v *EditContextView) apply() {
	newName := v.form.GetFormItem(0).(*tview.InputField).GetText()
	newNamespace := v.form.GetFormItem(1).(*tview.InputField).GetText()

	app.Pop()

	if newName != v.context.Name {
		ShowModalDialog(
			fmt.Sprintf("Rename %q to %q?", v.context.Name, newName),
			[]string{"No", "Yes"},
			func(_ int, label string) {
				if label == "Yes" {
					renameContext(v.context.Name, newName)
				}
			})
	}
	if newNamespace != v.context.Context.Namespace {
		ShowModalDialog(
			fmt.Sprintf("Change namespace %q to %q?",
				v.context.Context.Namespace, newNamespace),
			[]string{"No", "Yes"},
			func(_ int, label string) {
				if label == "Yes" {
					changeNamespace(v.context.Name, newNamespace)
				}
			})

	}
}
