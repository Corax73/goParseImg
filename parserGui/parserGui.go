package parserGui

import (
	"conc/imgParser"
	"conc/utils"
	"strconv"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/joho/godotenv"
)

type ParserGui struct {
	Parser                     *imgParser.ImgParser
	Input, Display, DelayEntry *widget.Entry
	ScrollContainer            *container.Scroll
	SendBtn, ClearResultBtn    *widget.Button
}

func (parserGui *ParserGui) SendBtnHandler() *widget.Button {
	return widget.NewButton("Send", func() {
		if parserGui.Input.Text != "" {
			parserGui.Display.SetText("")
			parserGui.Parser.StrError = ""
			urlSlice := strings.Split(parserGui.Input.Text, ",")
			startTime := time.Now()
			defer func() {
				durations := utils.Duration(startTime)
				parserGui.Display.SetText(utils.ConcatSlice([]string{parserGui.Display.Text, "Duration: ", durations}))
			}()

			var wg sync.WaitGroup
			defer wg.Wait()
			for _, url := range urlSlice {
				wg.Add(1)
				go func() {
					parserGui.getDelay()
					parserGui.Parser.GetImg(url)
					if parserGui.Parser.StrError == "" {
						parserGui.Display.SetText(utils.ConcatSlice([]string{parserGui.Display.Text, parserGui.Parser.SrtAdded}))
					} else {
						parserGui.Display.SetText(utils.ConcatSlice([]string{parserGui.Parser.StrError}))
					}
					wg.Done()
				}()
			}
		}
	})
}

func (parserGui *ParserGui) GetScrollDisplay() *container.Scroll {
	return container.NewVScroll(container.NewGridWithRows(
		1,
		parserGui.Display,
	))
}

func (parserGui *ParserGui) getDelay() {
	if parserGui.DelayEntry.Text != "" {
		number, err := strconv.Atoi(parserGui.DelayEntry.Text)
		if err == nil {
			parserGui.Parser.Delay = number
			envMap := utils.GetConfFromEnvFile()
			envMap["DELAY"] = strconv.Itoa(number)
			godotenv.Write(envMap, ".env")
		}
	}
}
