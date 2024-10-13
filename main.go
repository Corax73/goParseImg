package main

import (
	"conc/customLog"
	"conc/imgParser"
	"conc/parserGui"
	"conc/utils"
	"log"
	"strconv"

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

	delayStr := utils.ConcatSlice([]string{"Enter delay, now installed: ", strconv.Itoa(parserGui.Parser.Delay), " seconds"})
	parserGui.DelayEntry.SetPlaceHolder(delayStr)
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
			btnExit,
		),
	)

	window.SetContent(content)
	window.CenterOnScreen()
	window.Resize(fyne.NewSize(800, 600))
	window.ShowAndRun()
}
