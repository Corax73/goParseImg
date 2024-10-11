package imgParser

import (
	"conc/customLog"
	"conc/utils"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type State struct {
	Delay    int
	ImageDir string
}

func (state *State) ResetState() {
	state.Delay = 0
}

type ImgParser struct {
	State
}

func (parser *ImgParser) Init() {
	var err error
	envDelay := utils.GetEnvByKey("DELAY")
	if envDelay != "" {
		parser.Delay, err = strconv.Atoi(envDelay)
		if err != nil {
			customLog.Logging(err)
		}
	}
	parser.ImageDir = "./images/"
}

func (parser *ImgParser) SendRequest(url string) (*http.Response, error) {
	if parser.Delay > 0 {
		time.Sleep(time.Duration(parser.Delay) * time.Second)
	}
	response, err := http.Get(url)
	if err != nil {
		customLog.Logging(err)
	}
	return response, err
}

func (parser *ImgParser) GetImg(url string) {
	var err error
	currentTime := time.Now()
	dirName := utils.ConcatSlice([]string{parser.ImageDir, currentTime.Format("2006_01_2-15_04_05")})
	dirName, err = utils.CreateDir(dirName)
	if err == nil {
		response, err := parser.SendRequest(url)
		defer response.Body.Close()
		doc, err := html.Parse(response.Body)
		if err != nil {
			customLog.Logging(err)
		}
		parser.ProcessHtmlDoc(doc, "img", dirName)
	} else {
		customLog.Logging(err)
	}
}

func (parser *ImgParser) ProcessHtmlDoc(n *html.Node, tagName string, dirName string) {
	if n.Data == tagName {
		parser.GetSrc(n, dirName)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		parser.ProcessHtmlDoc(c, tagName, dirName)
	}
}

func (parser *ImgParser) GetSrc(n *html.Node, dirName string) {
	for _, a := range n.Attr {
		if a.Key == "src" && strings.Contains(a.Val, ".jpg") {
			var err error
			response, err := parser.SendRequest(a.Val)
			defer response.Body.Close()
			if err != nil {
				customLog.Logging(err)
			} else {
				pathSlice := strings.Split(a.Val, "/")
				fileName := utils.ConcatSlice([]string{dirName, "/", pathSlice[len(pathSlice)-1]})
				file, err := os.Create(fileName)
				defer file.Close()

				_, err = io.Copy(file, response.Body)
				if err != nil {
					customLog.Logging(err)
				}
				fmt.Println(utils.ConcatSlice([]string{"added: ", fileName}))
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		parser.GetSrc(c, dirName)
	}
}
