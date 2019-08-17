package main

import (
	"github.com/rivo/tview"
)

type View interface {
	Content() tview.Primitive
	OnActivate()
}

type App struct {
	app       *tview.Application
	viewStack []View
}

func NewApp() *App {
	c := &App{
		app: tview.NewApplication(),
	}
	return c
}

func (c *App) Push(v View) {
	c.viewStack = append(c.viewStack, v)
	c.app.SetRoot(v.Content(), true)
	v.OnActivate()
}

func (c *App) Pop() {
	c.viewStack = c.viewStack[:len(c.viewStack)-1]
	if len(c.viewStack) == 0 {
		c.Stop()
		return
	}
	v := c.viewStack[len(c.viewStack)-1]
	c.app.SetRoot(v.Content(), true)
	v.OnActivate()
}

func (c *App) Run() error {
	return c.app.Run()
}

func (c *App) Stop() {
	c.app.Stop()
}
