package help

import (
	"html/template"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/russross/blackfriday"
)

type infoBlock struct {
	Title  string
	Author string
	Date   time.Time
	Tags   []string
}

// Item relates to a file in the help/ folder
type Item struct {
	Title        string
	Content      string
	HTMLContent  template.HTML
	Path         string
	Author       string
	Tags         []string
	CreationDate time.Time
}

var infoBlockRegexp = regexp.MustCompile("(?m)^-+$")
var loc, _ = time.LoadLocation("Europe/Amsterdam")

func parseInfoBlock(s string) infoBlock {
	var res infoBlock

	s = strings.TrimSpace(s)
	lines := strings.Split(s, "\n")
	for _, line := range lines {
		splitted := strings.Split(line, ":")
		key := strings.TrimSpace(splitted[0])
		val := strings.TrimSpace(splitted[1])

		switch key {
		case "title":
			res.Title = val
		case "author":
			res.Author = val
		case "date":
			date, err := time.ParseInLocation("2006-01-02", val, loc)
			if err == nil {
				res.Date = date
			}
		case "tags":
			tags := strings.Split(val, ",")
			for _, tag := range tags {
				res.Tags = append(res.Tags, strings.TrimSpace(tag))
			}
		}
	}

	return res
}

func filesToItems(files []file) ([]Item, error) {
	var res []Item

	for _, f := range files {
		file, err := os.Open(f.Path)
		if err != nil {
			return res, err
		}

		content, err := ioutil.ReadAll(file)
		file.Close()
		if err != nil {
			return res, err
		}

		splitted := infoBlockRegexp.Split(string(content), 3)

		infoBlock := parseInfoBlock(splitted[1])
		body := strings.TrimSpace(splitted[2])

		item := Item{
			Title:        infoBlock.Title,
			Path:         f.Path,
			CreationDate: infoBlock.Date,
			Author:       infoBlock.Author,
			Tags:         infoBlock.Tags,
			Content:      body,
			HTMLContent:  template.HTML(blackfriday.MarkdownCommon([]byte(body))),
		}

		res = append(res, item)
	}

	return res, nil
}

// LoadItems returns all the help items available in the given dir
func LoadItems(dir string) ([]Item, error) {
	files, err := readDirRecursive(dir)
	if err != nil {
		return make([]Item, 0), err
	}
	return filesToItems(files)
}
