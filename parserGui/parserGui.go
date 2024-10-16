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
	SendBtn, ClearWindowBtn    *widget.Button
	DisplayTotal               *widget.Label
}

func (parserGui *ParserGui) SendBtnHandler() *widget.Button {
	return widget.NewButton("Send", func() {
		if parserGui.Input.Text != "" {
			parserGui.Parser.ResetState()
			urlSlice := strings.Split(parserGui.Input.Text, ",")
			startTime := time.Now()
			defer func() {
				durations := utils.Duration(startTime)
				parserGui.Display.SetText(utils.ConcatSlice([]string{parserGui.Display.Text, "Duration: ", durations}))
			}()

			var wg sync.WaitGroup
			chanHtmlData := make(chan *imgParser.HtmlDataToParse, len(urlSlice))
			defer close(chanHtmlData)
			defer wg.Wait()
			for _, url := range urlSlice {
				wg.Add(1)
				go func() {
					parserGui.getDelay()
					parserGui.Parser.GetHtmlFromUrl(url, chanHtmlData)
					defer wg.Done()
				}()
			}

			check := true
			i := 0
			for check {
				if len(chanHtmlData) > 0 {
					wg.Add(1)
					go func(ch chan *imgParser.HtmlDataToParse) {
						defer wg.Done()
						parserGui.Display.SetText(utils.ConcatSlice([]string{"Wait...", "\n"}))
						parserGui.Parser.ProcessHtmlDoc(chanHtmlData)
						parserGui.ShowResp()
					}(chanHtmlData)
					i++
				}
				if i == len(urlSlice) {
					check = false
				}
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

func (parserGui *ParserGui) ClearWindowBtnHandler() *widget.Button {
	return widget.NewButton("Clearing window data", func() {
		parserGui.Input.SetText("")
		parserGui.Display.SetText("")
		parserGui.DelayEntry.SetText(parserGui.GetDelayPlaceholder())
		parserGui.DisplayTotal.SetText("Total added: ")
		parserGui.Parser.ResetState()
	})
}

func (parserGui *ParserGui) GetDelayPlaceholder() string {
	return utils.ConcatSlice([]string{"Enter delay, now installed: ", strconv.Itoa(parserGui.Parser.Delay), " seconds"})
}

func (parserGui *ParserGui) ShowResp() {
	if parserGui.Parser.StrError == "" {
		parserGui.Display.SetText(utils.ConcatSlice([]string{parserGui.Display.Text, parserGui.Parser.StrAdded}))
	} else {
		parserGui.Display.SetText(utils.ConcatSlice([]string{parserGui.Parser.StrError}))
	}
	parserGui.DisplayTotal.SetText(utils.ConcatSlice([]string{"Total added: ", strconv.Itoa(parserGui.Parser.CountAdded)}))
}
