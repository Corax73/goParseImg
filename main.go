package main

import (
	"conc/customLog"
	"conc/imgParser"
	"conc/parserGui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	parserApp := app.New()
	window := parserApp.NewWindow("Image parser")

	btnExit := widget.NewButton("Exit", func() {
		parserApp.Quit()
	})

	parser := imgParser.ImgParser{}
	parser.Init()

	parserGui := parserGui.ParserGui{
		Parser:  &parser,
		Input:   widget.NewEntry(),
		Display: widget.NewLabel("Duration: "),
	}
	parserGui.SendBtn = parserGui.SendBtnHandler()

	content := container.NewGridWithColumns(
		1,
		container.NewGridWithRows(
			4,
			parserGui.Input,
			parserGui.SendBtn,
			parserGui.Display,
			btnExit,
		),
	)

	window.SetContent(content)
	window.CenterOnScreen()
	window.Resize(fyne.NewSize(800, 600))
	window.ShowAndRun()

	customLog.LogInit("./logs/app.log")
}
