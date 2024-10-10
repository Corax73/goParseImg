package main

import (
	"conc/customLog"
	"conc/utils"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

func main() {
	customLog.LogInit("./logs/app.log")

	urlSlice := []string{
		"https://stub.com/",
	}

	defer duration(time.Now())

	var wg sync.WaitGroup
	defer wg.Wait()
	for _, url := range urlSlice {
		wg.Add(1)
		go func() {
			getImg(url)
			wg.Done()
		}()
	}
}

func sendRequest(url string, delayInSecond int) (*http.Response, error) {
	if delayInSecond > 0 {
		time.Sleep(time.Duration(delayInSecond) * time.Second)
	}
	response, err := http.Get(url)
	if err != nil {
		customLog.Logging(err)
	}
	return response, err
}

func getImg(url string) {
	var delay int
	var err error
	envDelay := getEnv("DELAY")
	if envDelay != "" {
		delay, err = strconv.Atoi(envDelay)
		if err != nil {
			customLog.Logging(err)
		}
	}
	fmt.Println(delay)
	response, err := sendRequest(url, delay)
	defer response.Body.Close()
	doc, err := html.Parse(response.Body)
	if err != nil {
		customLog.Logging(err)
	}
	processHtmlDoc(doc, "img")
}

func processHtmlDoc(n *html.Node, tagName string) {
	if n.Data == tagName {
		getSrc(n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		processHtmlDoc(c, tagName)
	}
}

func duration(start time.Time) {
	fmt.Printf("%v\n", time.Since(start))
}

func getSrc(n *html.Node) {
	for _, a := range n.Attr {
		if a.Key == "src" && strings.Contains(a.Val, ".jpg") {
			var delay int
			var err error
			envDelay := getEnv("DELAY")
			if envDelay != "" {
				delay, err = strconv.Atoi(envDelay)
				if err != nil {
					customLog.Logging(err)
				}
			}
			response, err := sendRequest(a.Val, delay)
			defer response.Body.Close()
			if err != nil {
				customLog.Logging(err)
			} else if response.Header.Get("Etag") != "" {
				dir := "./images/"
				fileName := concatSlice([]string{dir, strings.Trim(response.Header.Get("Etag"), "\""), ".jpg"})
				file, err := os.Create(fileName)
				defer file.Close()

				_, err = io.Copy(file, response.Body)
				if err != nil {
					customLog.Logging(err)
				}
				fmt.Println(concatSlice([]string{"added: ", fileName}))
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		getSrc(c)
	}
}

func concatSlice(strSlice []string) string {
	resp := ""
	if len(strSlice) > 0 {
		var strBuilder strings.Builder
		for _, val := range strSlice {
			strBuilder.WriteString(val)
		}
		resp = strBuilder.String()
		strBuilder.Reset()
	}
	return resp
}

func getEnv(key string) string {
	mapEnv := utils.GetConfFromEnvFile()
	val, ok := mapEnv[key]
	if ok {
		return val
	} else {
		return ""
	}
}
