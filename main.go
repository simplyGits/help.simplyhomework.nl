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
	io.WriteString(w, "<!doctype html>")
	for _, item := range items {
		if item.Name != "" {
			io.WriteString(w, fmt.Sprintf("<h1>%s by %s</h1><hr>%s<br><br>", item.Name, item.Author, item.HTMLContent))
		}
	}
}

func main() {
	http.HandleFunc("/", mainPage)
	log.Fatal(http.ListenAndServe(":80", nil))
}
