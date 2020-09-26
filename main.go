package main

import (
	"fmt"
	"time"

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

	query := prompt.Input("> ", completer)

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

	// Navigate to the simple playground interface.
	if err := wd.Get("http://google.com/search?tbm=isch&q=" + query); err != nil {
		panic(err)
	}

	i := 0
	for true {
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

			time.Sleep(3 * time.Second)

			fmt.Println(elem2.GetAttribute("src"))
		}
	}
}
