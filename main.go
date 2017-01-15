package main

import (
	"help"
	"html/template"
	"io"
	"log"
	"net/http"
	"path"
)

var items []help.Item
var articleTmpl *template.Template

func setHTTPContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}

func getItem(path string) *help.Item {
	for i := range items {
		item := &items[i]
		if item.Path == path {
			return item
		}
	}
	return nil
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	slug := r.URL.Path[1:]
	if slug == "" {
		setHTTPContentType(w)
		// TODO: render mainpage
	} else if slug[len(slug)-4:] == ".css" {
		file := path.Base(slug)
		name := file[:len(file)-4]
		http.ServeFile(w, r, path.Join("views", name, name+".css"))
	} else {
		setHTTPContentType(w)

		item := getItem("help/" + slug + ".md")
		if item == nil {
			// TODO: status code
			io.WriteString(w, "404: not found")
		} else {
			articleTmpl.Execute(w, item)
		}
	}
}

func main() {
	var err error
	items, err = help.LoadItems("entries")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("loaded %d help items\n", len(items))

	articleTmpl, err = template.ParseFiles("views/article/article.html")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", mainPage)
	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
