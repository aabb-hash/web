package util

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

var htmls map[string]string = make(map[string]string)

func InitHtml() {
	files, err := os.ReadDir("./html")
	if err != nil {
		os.Mkdir("html", os.ModePerm)
		panic("no html")
	}

	for _, file := range files {
		key := strings.ReplaceAll(file.Name(), ".html", "")

		data, err := os.ReadFile("html/" + file.Name())
		if err != nil {
			panic("unable to read html/" + file.Name())
		}

		htmls[key] = string(data)
	}
}

func GetHeader(page string, title string, w http.ResponseWriter) string {
	html := htmls["header"]

	html = strings.ReplaceAll(html, "{name}", page)
	html = strings.ReplaceAll(html, "{title}", title)

	if w != nil {
		fmt.Fprint(w, html)
	}

	return html
}

func GetHtmlContent(name string, replace map[string]string, w http.ResponseWriter) string {
	html := htmls[name]

	for old, new := range replace {
		html = strings.ReplaceAll(html, old, new)
	}

	if w != nil {
		fmt.Fprint(w, html)
	}

	return html
}

func Replace(html string, old string, new string) string {
	return strings.ReplaceAll(html, old, new)
}
