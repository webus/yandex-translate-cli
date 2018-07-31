package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/atotto/clipboard"
)

func showHelp() {
	fmt.Println("Simple command line Yandex Translate client")
	fmt.Println("Usage of yandex-translate-cli:")
	flag.PrintDefaults()
	fmt.Println()
	fmt.Println("ex: yandex-translate-cli -lang=en привет")
	fmt.Println("    yandex-translate-cli -lang=ru morning")
	fmt.Println("    yandex-translate-cli -lang=ru \"How are you?\"")
}

func main() {
	curLang := flag.String("lang", "en", "language of translation")
	showMeHelp := flag.Bool("help", false, "show help")
	showMeHelpShort := flag.Bool("h", false, "show short help")
	flag.Parse()

	if *showMeHelp || *showMeHelpShort {
		showHelp()
		os.Exit(0)
	}

	if len(os.Args) < 3 {
		fmt.Println("There are not enough arguments!")
		flag.PrintDefaults()
		os.Exit(1)
	}

	data := os.Args[2]
	data = url.QueryEscape(data)

	var lang string

	if strings.ToLower(*curLang) == "ru" {
		lang = "en-ru"
	} else {
		lang = "ru-" + *curLang
	}

	yandexKey := os.Getenv("YANDEX_TRANSLATE_API_KEY")
	yandexAPI := "https://translate.yandex.net/api/v1.5/tr.json/translate"
	yandexAPIURL := yandexAPI + "?key=" + yandexKey + "&text=" + data + "&lang=" + lang

	resp, err := http.Get(yandexAPIURL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var f map[string]interface{}
	err = json.Unmarshal(body, &f)

	dataStr := f["text"].([]interface{})[0].(string)
	fmt.Println(dataStr)
	clipboard.WriteAll(dataStr)
}
