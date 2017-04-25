package main

import (
	"github.com/andlabs/ui"
)

func main() {
	err := ui.Main(func() {
		name := ui.NewEntry()
		button := ui.NewButton("Greet")
		greeting := ui.NewLabel("")
		box := ui.NewVerticalBox()
		label := ui.NewLabel("Enter your name:")
		//label.SetText("")
		box.Append(label, false)
		box.Append(name, false)
		box.Append(button, false)
		box.Append(greeting, false)
		window := ui.NewWindow("Hello", 600, 400, false)
		window.SetChild(box)
		button.OnClicked(func(*ui.Button) {
			greeting.SetText("Hello, " + name.Text() + "!")
		})
		window.OnClosing(func(*ui.Window) bool {
			ui.Quit()
			return true
		})
		window.Show()
	})
	if err != nil {
		panic(err)
	}
}