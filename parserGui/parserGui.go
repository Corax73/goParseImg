package parserGui

import (
	"conc/imgParser"
	"conc/utils"
	"fmt"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2/widget"
)

type ParserGui struct {
	Parser                  *imgParser.ImgParser
	Input                   *widget.Entry
	Display                 *widget.Label
	SendBtn, ClearResultBtn *widget.Button
}

func (parserGui *ParserGui) SendBtnHandler() *widget.Button {
	return widget.NewButton("Send", func() {
		if parserGui.Input.Text != "" {
			urlSlice := strings.Split(parserGui.Input.Text, ",")
			fmt.Println(urlSlice)
			startTime := time.Now()
			defer func() {
				durations := utils.Duration(startTime)
				parserGui.Display.SetText(durations)
			}()

			var wg sync.WaitGroup
			defer wg.Wait()
			for _, url := range urlSlice {
				wg.Add(1)
				go func() {
					parserGui.Parser.GetImg(url)
					wg.Done()
				}()
			}
		}
	})
}
