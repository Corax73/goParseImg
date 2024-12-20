package main

import (
	"conc/customLog"
	"conc/customTheme"
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
	parserApp.Settings().SetTheme(customTheme.NewCustomTheme())
	window := parserApp.NewWindow("Image parser")

	btnExit := widget.NewButton("Exit", func() {
		parserApp.Quit()
	})

	parser := imgParser.ImgParser{}
	parser.Init()

	parserGui := parserGui.ParserGui{
		Parser:         &parser,
		Input:          widget.NewEntry(),
		Display:        widget.NewEntry(),
		DelayEntry:     widget.NewEntry(),
		TagEntry:       widget.NewEntry(),
		AttributeEntry: widget.NewEntry(),
		DisplayTotal:   widget.NewLabel("Total added: "),
	}

	parserGui.Input.SetPlaceHolder("Enter the URL of the html resource for the request (enter multiple values ​​separated by commas)")
	parserGui.DelayEntry.SetPlaceHolder(parserGui.GetDelayPlaceholder())
	parserGui.TagEntry.SetPlaceHolder("Enter a tag storing images, default `img`")
	parserGui.AttributeEntry.SetPlaceHolder("Enter an attribute that stores a link to images, default `src`")
	parserGui.ClearWindowBtn = parserGui.ClearWindowBtnHandler()
	parserGui.SendBtn = parserGui.SendBtnHandler()
	parserGui.ScrollContainer = parserGui.GetScrollDisplay()

	content := container.NewGridWithColumns(
		1,
		container.NewGridWithRows(
			4,
			parserGui.Input,
			container.NewGridWithColumns(
				4,
				parserGui.DelayEntry,
				parserGui.TagEntry,
				parserGui.AttributeEntry,
				parserGui.SendBtn,
			),
			parserGui.ScrollContainer,
			container.NewGridWithColumns(
				3,
				parserGui.ClearWindowBtn,
				parserGui.DisplayTotal,
				btnExit,
			),
		),
	)

	window.SetContent(content)
	window.CenterOnScreen()
	window.Resize(fyne.NewSize(800, 600))
	window.ShowAndRun()
}
