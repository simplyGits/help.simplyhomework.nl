package main

import (
	"fmt"
	"help"
	"io"
	"log"
	"net/http"
)

var items, _ = help.LoadItems("help")

func mainPage(w http.ResponseWriter, req *http.Request) {
	for _, item := range items {
		if item.Name != "" {
			io.WriteString(w, fmt.Sprintf("%s\n---------\n\n%s", item.Name, item.Content))
		}
	}
}

func main() {
	http.HandleFunc("/", mainPage)
	log.Fatal(http.ListenAndServe(":80", nil))
}
