package main

import (
	"conc/customLog"
	"conc/imgParser"
	"conc/parserGui"
	"conc/utils"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	_, err := utils.CreateDir("./logs")
	if err != nil {
		log.Fatalln(err.Error())
	}
	customLog.LogInit("./logs/app.log")

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
		Display: widget.NewEntry(),
		DelayEntry:    widget.NewEntry(),
	}

	parserGui.DelayEntry.SetPlaceHolder(parserGui.GetDelayPlaceholder())
	parserGui.ClearWindowBtn = parserGui.ClearWindowBtnHandler()
	parserGui.SendBtn = parserGui.SendBtnHandler()
	parserGui.ScrollContainer = parserGui.GetScrollDisplay()

	content := container.NewGridWithColumns(
		1,
		container.NewGridWithRows(
			4,
			parserGui.Input,
			container.NewGridWithColumns(
				2,
				parserGui.DelayEntry,
				parserGui.SendBtn,
			),
			parserGui.ScrollContainer,
			container.NewGridWithColumns(
				2,
				parserGui.ClearWindowBtn,
				btnExit,
			),
		),
	)

	window.SetContent(content)
	window.CenterOnScreen()
	window.Resize(fyne.NewSize(800, 600))
	window.ShowAndRun()
}
