package main

import (
	"conc/customLog"
	"conc/imgParser"
	"conc/utils"
	"fmt"
	"sync"
	"time"
)

func main() {
	customLog.LogInit("./logs/app.log")

	urlSlice := []string{
		"https://stub.com/",
		"https://stub1.com/",
	}

	startTime := time.Now()
	defer func() {
		durations := utils.Duration(startTime)
		fmt.Println(durations)
	}()

	var wg sync.WaitGroup
	defer wg.Wait()
	parser := imgParser.ImgParser{}
	parser.Init()
	for _, url := range urlSlice {
		wg.Add(1)
		go func() {
			parser.GetImg(url)
			wg.Done()
		}()
	}
}
