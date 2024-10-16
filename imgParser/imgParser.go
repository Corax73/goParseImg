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
	Delay, CountAdded            int
	ImageDir, StrAdded, StrError string
}

// ResetState resets State structure parameters to default.
func (state *State) ResetState() {
	state.Delay = 0
	state.CountAdded = 0
	state.ImageDir = "./images/"
	state.StrAdded, state.StrError = "", ""
}

type HtmlDataToParse struct {
	HtmlDoc                      *html.Node
	TagName, DirName, DomainName string
}

type ImgParser struct {
	State
}

// Init checks and creates if there is no environment file, gets the Delay value.
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
}

// SendRequest sends Get requests using the passed string.
func (parser *ImgParser) SendRequest(url string) (*http.Response, error) {
	var response *http.Response
	if parser.Delay > 0 {
		time.Sleep(time.Duration(parser.Delay) * time.Second)
	}

	url = strings.Trim(url, " ")

	client := &http.Client{
		Transport: &http.Transport{},
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		customLog.Logging(err)
	} else {
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
		response, err = client.Do(req)
		if err != nil {
			customLog.Logging(err)
		}
	}
	return response, err
}

// GetHtmlFromUrl receives images from the address of the passed string from the IMG tag.
func (parser *ImgParser) GetHtmlFromUrl(url string, ch chan<- *HtmlDataToParse) {
	var err error
	if !strings.Contains(url, "http") {
		url = utils.ConcatSlice([]string{"http://", url})
	}
	response, err := parser.SendRequest(url)
	if err == nil {
		defer response.Body.Close()
		dirName, err := parser.createTimestampAsDirForFiles()
		if err == nil {
			doc, err := html.Parse(response.Body)
			if err != nil {
				parser.StrError = err.Error()
				customLog.Logging(err)
			}
			pathSlice := strings.Split(url, "/")
			pathSlice = pathSlice[:3]
			strDomain := utils.ConcatSlice([]string{pathSlice[0], "//", pathSlice[2]})
			ch <- &HtmlDataToParse{
				doc,
				"img",
				dirName,
				strDomain,
			}
		} else {
			parser.StrError = err.Error()
			customLog.Logging(err)
		}
	} else {
		var doc *html.Node
		ch <- &HtmlDataToParse{
			doc,
			"img",
			"",
			"",
		}
		parser.StrError = err.Error()
		customLog.Logging(err)
	}
	utils.GCRunAndPrintMemory()
}

// ProcessHtmlDoc processes html nodes.
func (parser *ImgParser) ProcessHtmlDoc(ch <-chan *HtmlDataToParse) {
	htmlData := <-ch
	if htmlData.HtmlDoc != nil {
		if htmlData.HtmlDoc.Data == htmlData.TagName {
			parser.GetSrc(htmlData)
		}
		for c := htmlData.HtmlDoc.FirstChild; c != nil; c = c.NextSibling {
			ch := make(chan *HtmlDataToParse, 1)
			defer close(ch)
			ch <- &HtmlDataToParse{
				c,
				htmlData.TagName,
				htmlData.DirName,
				htmlData.DomainName,
			}
			parser.ProcessHtmlDoc(ch)
		}
	}
}

// GetSrc gets the image by value in the SRC attribute.
func (parser *ImgParser) GetSrc(htmlData *HtmlDataToParse) {
	for _, a := range htmlData.HtmlDoc.Attr {
		if a.Key == "src" && strings.Contains(a.Val, "/") && !strings.Contains(a.Val, ".svg") && len(a.Val) > 5 {
			imgUrl := a.Val
			if !strings.Contains(a.Val, "http") && !strings.Contains(a.Val, "//") {
				imgUrl = utils.ConcatSlice([]string{htmlData.DomainName, imgUrl})
			} else if !strings.Contains(a.Val, "http") && strings.Contains(a.Val, "//") {
				imgUrl = utils.ConcatSlice([]string{"http:", imgUrl})
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
						fileName = utils.ConcatSlice([]string{htmlData.DirName, "/", pathSlice[len(pathSlice)-1], ".jpg"})
					}
				} else {
					pathSlice = strings.Split(a.Val, "/")
					if strings.Contains(pathSlice[len(pathSlice)-1], "?") {
						pathSlice = strings.Split(pathSlice[len(pathSlice)-1], "?")
						fileName = pathSlice[0]
					} else {
						fileName = pathSlice[len(pathSlice)-1]
					}
					fileName = utils.ConcatSlice([]string{htmlData.DirName, "/", fileName})
				}
				if fileName != "" {
					file, err := os.Create(fileName)
					defer file.Close()
					_, err = io.Copy(file, response.Body)
					if err != nil {
						customLog.Logging(err)
					} else {
						parser.CountAdded++
						parser.StrAdded = utils.ConcatSlice([]string{parser.StrAdded, utils.ConcatSlice([]string{"added: ", fileName}), "\n"})
					}
				}
			}
		}
	}

	for c := htmlData.HtmlDoc.FirstChild; c != nil; c = c.NextSibling {
		parser.GetSrc(htmlData)
	}
}

// createTimestampAsDirForFiles creates a directory for files with a name - timestamp.
func (parser *ImgParser) createTimestampAsDirForFiles() (string, error) {
	currentTime := time.Now()
	dirName := utils.ConcatSlice([]string{parser.ImageDir, currentTime.Format("2006_01_2-15_04_05")})
	return utils.CreateDir(dirName)
}
