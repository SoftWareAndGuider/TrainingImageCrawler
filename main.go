package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/tebeka/selenium"
)

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func main() {
	fmt.Println("Crwaler v1")

	const (
		// These paths will be different on your system.
		seleniumPath    = "bin/selenium-server.jar"
		geckoDriverPath = "bin/geckodriver"
		port            = 8080
	)

	querys := strings.Split(prompt.Input("> ", completer), ",")

	opts := []selenium.ServiceOption{
		selenium.GeckoDriver(geckoDriverPath), // Specify the path to GeckoDriver in order to use Firefox
	}
	// selenium.SetDebug(true)
	service, err := selenium.NewSeleniumService(seleniumPath, port, opts...)
	if err != nil {
		panic(err) // panic is used only as an example and is not otherwise recommended.
	}
	defer service.Stop()

	// Connect to the WebDriver instance running locally.
	caps := selenium.Capabilities{"browserName": "firefox"}
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		panic(err)
	}
	defer wd.Quit()

	for _, query := range querys {

		os.MkdirAll("./downloads/"+query, os.ModeDir)
		// Navigate to the simple playground interface.
		if err := wd.Get("http://google.com/search?tbm=isch&q=" + query); err != nil {
			panic(err)
		}

		i := -1
		str := ""
		for true {
			if i > 130 {
				break
			}

			elems, err := wd.FindElements(selenium.ByCSSSelector, ".rg_i")
			if err != nil {
				panic(err)
			}

			for i2, elem := range elems {
				if i >= i2 {
					continue
				} else {
					i = i2
				}

				elem.Click()
				elem2, err := wd.FindElement(selenium.ByCSSSelector, ".n3VNCb")
				if err != nil {
					panic(err)
				}

				src, err := elem2.GetAttribute("src")
				if err != nil {
					panic(err)
				}

				if strings.HasPrefix(src, "https://encrypted-tbn0.gstatic.com/") {
					if str == src {
						continue
					} else {
						str = src
					}

					filepath := fmt.Sprintf("./downloads/%s/%d", query, i2)
					fmt.Printf("saved \"%s\" to \"%s\"\n", src, filepath)
					downloadFile(filepath, src)
				}
			}
		}
	}
}

func downloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	contype := strings.Split(resp.Header.Get("content-type"), "/")[1]
	filepath = filepath + "." + contype

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
