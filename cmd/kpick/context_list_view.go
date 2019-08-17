package main

import (
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type ContextListView struct {
	list *tview.List
	flex *tview.Flex
}

func NewContextListView() *ContextListView {
	v := &ContextListView{}
	v.list = tview.NewList()
	v.list.ShowSecondaryText(false)
	v.list.SetInputCapture(v.onInput)
	v.list.SetSelectedFunc(v.onSelect)

	v.flex = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView().
			SetText("j:down k:up return:select R:rename D:delete q:quit").
			SetTextColor(tcell.ColorYellow), 2, 1, false).
		AddItem(v.list, 0, 1, true)

	return v
}

func (v *ContextListView) Content() tview.Primitive {
	return v.flex
}

func (l *ContextListView) OnActivate() {
	config := getConfig()

	l.list.Clear()
	for ndx, context := range config.Contexts {
		var key rune
		if ndx < 10 {
			key = '0' + rune(ndx)
		}
		l.list.AddItem(context.Name, "", key, nil)
		if context.Name == config.CurrentContext {
			l.list.SetCurrentItem(ndx)
		}
	}
}

func (l *ContextListView) onInput(e *tcell.EventKey) *tcell.EventKey {
	if e.Key() == tcell.KeyRune {
		switch e.Rune() {
		case 'j':
			return tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
		case 'k':
			return tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone)
		case 'D':
			l.deleteContext()
		case 'R':
			contextName, _ := l.list.GetItemText(l.list.GetCurrentItem())
			app.Push(NewRenameContextView(contextName))
		case 'q':
			app.Pop()
		}
	}
	return e
}

func (l *ContextListView) onSelect(i int, context string, _ string, _ rune) {
	useContext(context)
	app.Stop()
	fmt.Println("Switched to context", context)
}

func (l *ContextListView) deleteContext() {
	contextName, _ := l.list.GetItemText(l.list.GetCurrentItem())
	ShowModalDialog(
		fmt.Sprintf("Delete %q?", contextName),
		[]string{"No", "Yes"},
		func(_ int, label string) {
			if label == "Yes" {
				deleteContext(contextName)
			}
		})
}
