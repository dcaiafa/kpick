// kpick is a graphical terminal app to change the current kubectl config
// context.
package main

import (
	"log"
)

var app *App

func die(err interface{}) {
	if app != nil {
		app.Stop()
	}
	log.Fatal(err)
}

func main() {
	app = NewApp()

	app.Push(NewContextListView())
	app.Run()
}
