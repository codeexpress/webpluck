package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/codeexpress/webpluck/logger"
	"github.com/codeexpress/webpluck/webpluck"
	"gopkg.in/yaml.v2"
)

var (
	//argument flags
	filePtr       *string
	outputTextPtr *bool
	serverModePtr *int
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

func main() {
	initFlags()
	logger.InitLogger()

	serverMode := isFlagPassed("p")

	if serverMode {
		serveApi()
	} else { // things to pluck passed via a yaml file
		pluckFromFile()
	}
}

/* Server Mode
Listens on a port and answers online queries of type:
http://localhost:8080?baseUrl="example.com"&xpath="/html/body"&regex=""
*/
func serveApi() {
	logger.LogIt("Started HTTP server on localhost: "+strconv.Itoa(*serverModePtr), true)

	http.HandleFunc("/", handleHttp)
	fmt.Println(http.ListenAndServe(":"+strconv.Itoa(*serverModePtr), nil))
}

func handleHttp(w http.ResponseWriter, req *http.Request) {
	results := make(map[string]string)
	req.ParseForm()
	baseUrl := req.Form.Get("baseUrl")
	xpath := req.Form.Get("xpath")
	regex := req.Form.Get("regex")

	results["baseUrl"] = baseUrl
	results["xpath"] = xpath
	results["regex"] = regex

	logger.LogIt(getIp(req) + "  " + req.Header.Get("User-Agent") + " Request: ")
	logger.LogIt(results)
	defer func() { // in case of panic
		if err := recover(); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			fmt.Fprintf(w, "Webpluck encountered an error. Make sure that the baseUrl is a valid URL and xpath and regex are valid\n")
			fmt.Fprintf(w, "Error description:\n%s\n", err)
			logger.LogIt(err)
		}
	}()
	text := webpluck.ExtractTextFromUrl(baseUrl, xpath, regex)
	results["pluckedData"] = text
	jsonString, err := json.MarshalIndent(results, "", "  ")
	check(err)
	fmt.Fprintf(w, string(jsonString))
	logger.LogIt("Answer: " + text)
}

func pluckFromFile() {
	data, err := ioutil.ReadFile(*filePtr)
	check(err)
	var list targetList

	err = yaml.Unmarshal(data, &list)
	check(err)

	results := make(map[string]string)

	for _, t := range list.TargetList {
		text := webpluck.ExtractTextFromUrl(t.BaseUrl, t.Xpath, t.Regex)
		results[t.Name] = text
		if *outputTextPtr { // if output to text (t) flag is set
			fmt.Println(t.Name + ": " + text)
		}
	}

	logger.LogIt("Webpluck invoked. Reading from file: " + *filePtr)
	logger.LogIt(results)

	if !*outputTextPtr { // default case is to print in JSON
		jsonString, err := json.MarshalIndent(results, "", "  ")
		check(err)
		fmt.Println(string(jsonString))
	}
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

/* Parse command line flage and initiaize the global flag variables */
func initFlags() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "---------------------------------------------------------")
		fmt.Fprintln(os.Stderr, "Usage:   $ ./webpluck [-f filename.yml]")
		fmt.Fprintln(os.Stderr, "Example: $ ./webpluck -f ./extract_list.yml")
		fmt.Fprintln(os.Stderr, "---------------------------------------------------------\nFlags:")
		flag.PrintDefaults()
	}

	filePtr = flag.String("f", "./targets.yml",
		"`File name (yml)` with the list of targets to pluck/extract")
	outputTextPtr = flag.Bool("t", false,
		"Outputs the results in `text format` instead of JSON (applicable only in non server mode)")
	serverModePtr = flag.Int("p", 0,
		"`Port number` to serve webpluck as in HTTP API")
	flag.Parse()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Get IP address of the incoming HTTP request based on forwarded-for
// header (present in case of proxy). If not, use the remote address
func getIp(req *http.Request) string {
	forwardedIp := req.Header.Get("X-Forwarded-For")
	if forwardedIp != "" {
		return forwardedIp
	}

	addr := req.RemoteAddr
	ip, _, _ := net.SplitHostPort(addr)
	return ip
}
