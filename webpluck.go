package main

import (
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"gopkg.in/xmlpath.v2"
)

type targetList struct {
	TargetList []dataLocation `yaml:"targetList"`
}

type dataLocation struct {
	Name    string `yaml:"name"`
	BaseUrl string `yaml:"baseUrl"`
	Xpath   string `yaml:"xpath"`
	Regex   string `yaml:"regex"`
}

/* Params:
   url   - URL of the page to be scraped
   xpath - the xpath of the target text in the html document
   regex - the regex to be applied to the extracted text found at xpath to extract final value
   returns:
   a string value that is the final plucked value
*/
func ExtractTextFromUrl(
	url string,
	xpath string,
	regex string) string {

	text := ""

	//logIt("Fetch from URL: "+url, 1)
	parsedHtml := fetchUrl(url) // returns a xmlpath.Node object
	path := xmlpath.MustCompile(xpath)
	value, ok := path.String(parsedHtml)
	if ok {
		if regex != "" {
			// try applying regex
			regexMatch := regexp.MustCompile(regex)
			text = regexMatch.FindStringSubmatch(string(value))[1]
		} else {
			text = value // no regex, the xpath element is the value
		}
	}
	return strings.TrimSpace(text)
}

// does a HTTP GET and returns the HTML body for that URL
func fetchUrl(url string) *xmlpath.Node {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("User-Agent", UserAgent)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	html, _ := ioutil.ReadAll(resp.Body)
	htmlStr := string(html)

	parsedHtml, err := xmlpath.ParseHTML(strings.NewReader(htmlStr))
	if err != nil {
		panic(err)
	}
	return parsedHtml
}
