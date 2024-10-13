package imgParser

import (
	"conc/customLog"
	"conc/utils"
	"errors"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type State struct {
	Delay                        int
	ImageDir, SrtAdded, StrError string
}

func (state *State) ResetState() {
	state.Delay = 0
}

type ImgParser struct {
	State
}

func (parser *ImgParser) Init() {
	var err error

	if _, err := os.Stat(".env"); errors.Is(err, os.ErrNotExist) {
		_, err := os.Create(".env")
		if err != nil {
			customLog.Logging(err)
		}
	}

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
	response, err := parser.SendRequest(url)
	if err == nil {
		defer response.Body.Close()
		currentTime := time.Now()
		dirName := utils.ConcatSlice([]string{parser.ImageDir, currentTime.Format("2006_01_2-15_04_05")})
		dirName, err = utils.CreateDir(dirName)
		if err == nil {
			doc, err := html.Parse(response.Body)
			if err != nil {
				parser.StrError = err.Error()
				customLog.Logging(err)
			}
			pathSlice := strings.Split(url, "/")
			pathSlice = pathSlice[:3]
			strDomain := utils.ConcatSlice([]string{pathSlice[0], "//", pathSlice[2]})
			parser.ProcessHtmlDoc(doc, "img", dirName, strDomain)
		} else {
			parser.StrError = err.Error()
			customLog.Logging(err)
		}
	} else {
		parser.StrError = err.Error()
		customLog.Logging(err)
	}
	utils.GCRunAndPrintMemory()
}

func (parser *ImgParser) ProcessHtmlDoc(n *html.Node, tagName string, dirName string, strDomain string) {
	if n.Data == tagName {
		parser.GetSrc(n, dirName, strDomain)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		parser.ProcessHtmlDoc(c, tagName, dirName, strDomain)
	}
}

func (parser *ImgParser) GetSrc(n *html.Node, dirName string, strDomain string) {
	for _, a := range n.Attr {
		if a.Key == "src" && strings.Contains(a.Val, "/") && !strings.Contains(a.Val, ".svg") && len(a.Val) > 5 {
			imgUrl := a.Val
			if !strings.Contains(a.Val, "http") {
				imgUrl = utils.ConcatSlice([]string{strDomain, imgUrl})
			}
			pathSlice := strings.Split(imgUrl, "?")
			pathSlice = pathSlice[:1]
			imgUrl = pathSlice[0]
			var err error
			response, err := parser.SendRequest(imgUrl)
			defer response.Body.Close()
			if err != nil {
				customLog.Logging(err)
			} else {
				var fileName string
				if !strings.Contains(imgUrl, ".jpg") && !strings.Contains(imgUrl, ".png") {
					pathSlice = strings.Split(pathSlice[0], "/")
					if len(pathSlice[len(pathSlice)-1]) > 0 {
						fileName = utils.ConcatSlice([]string{dirName, "/", pathSlice[len(pathSlice)-1], ".jpg"})
					}
				} else {
					pathSlice := strings.Split(a.Val, "/")
					fileName = utils.ConcatSlice([]string{dirName, "/", pathSlice[len(pathSlice)-1]})
				}
				if fileName != "" {
					file, err := os.Create(fileName)
					defer file.Close()

					_, err = io.Copy(file, response.Body)
					if err != nil {
						customLog.Logging(err)
					}
					parser.SrtAdded = utils.ConcatSlice([]string{parser.SrtAdded, utils.ConcatSlice([]string{"added: ", fileName}), "\n"})
				}
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		parser.GetSrc(c, dirName, strDomain)
	}
}
