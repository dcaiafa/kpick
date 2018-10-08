// kpick is a graphical terminal app to change the current kubectl config
// context.
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

var app *tview.Application

func die(err interface{}) {
	if app != nil {
		app.Stop()
	}
	log.Fatal(err)
}

func getCurrentContext() string {
	cmd := exec.Command("kubectl", "config", "current-context")
	res, err := cmd.CombinedOutput()
	if err != nil {
		die(err)
	}
	return strings.TrimSpace(string(res))
}

func listContexts() []string {
	cmd := exec.Command("kubectl", "config", "get-contexts", "-o", "name")
	res, err := cmd.CombinedOutput()
	if err != nil {
		die(err)
	}
	var lines []string
	scanner := bufio.NewScanner(bytes.NewReader(res))
	for scanner.Scan() {
		l := scanner.Text()
		if l == "" {
			continue
		}
		lines = append(lines, l)
	}
	return lines
}

func useContext(c string) {
	cmd := exec.Command("kubectl", "config", "use-context", c)
	err := cmd.Run()
	if err != nil {
		die(err)
	}
}

type contextList struct {
	*tview.List
}

func newContextList() *contextList {
	l := &contextList{}
	l.List = tview.NewList()
	l.List.ShowSecondaryText(false)
	l.List.SetInputCapture(l.onInput)
	l.List.SetSelectedFunc(l.onSelect)
	l.Refresh()
	return l
}

func (l *contextList) Refresh() {
	current := getCurrentContext()
	contexts := listContexts()

	l.Clear()
	for ndx, context := range contexts {
		var key rune
		if ndx < 10 {
			key = '0' + rune(ndx)
		}
		l.AddItem(context, "", key, nil)
		if context == current {
			l.SetCurrentItem(ndx)
		}
	}
}

func (l *contextList) onInput(e *tcell.EventKey) *tcell.EventKey {
	if e.Key() == tcell.KeyRune {
		switch e.Rune() {
		case 'j':
			return tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
		case 'k':
			return tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone)
		case 'q':
			app.Stop()
		}
	}
	return e
}

func (l *contextList) onSelect(i int, context string, _ string, _ rune) {
	useContext(context)
	app.Stop()
	fmt.Println("Switched to context", context)
}

func main() {
	app = tview.NewApplication()

	list := newContextList()

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView().
			SetText("j:down k:up return:select q:quit").
			SetTextColor(tcell.ColorYellow), 2, 1, false).
		AddItem(list, 0, 1, true)

	if err := app.SetRoot(flex, true).SetFocus(flex).Run(); err != nil {
		panic(err)
	}
}
