package help

import (
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
	"time"
)

type infoBlock struct {
	Title  string
	Author string
	Date   time.Time
}

// Item relates to a file in the help/ folder
type Item struct {
	Name         string
	Content      string
	Path         string
	Author       string
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
		}
	}

	return res
}

func fileToItems(f file) ([]Item, error) {
	var resItem *Item
	res := make([]Item, len(f.Files)+1)

	if f.IsDir {
		for _, file := range f.Files {
			items, err := fileToItems(file)
			if err != nil {
				return res, err
			}
			for _, item := range items {
				res = append(res, item)
			}
		}
	} else {
		file, err := os.Open(f.Path)
		if err != nil {
			return res, err
		}
		defer file.Close()

		content, err := ioutil.ReadAll(file)
		if err != nil {
			return res, err
		}

		splitted := infoBlockRegexp.Split(string(content), 3)
		infoBlock := parseInfoBlock(splitted[1])

		resItem = &Item{
			Name: path.Base(f.Path),
			Path: f.Path,
		}

		resItem.Name = infoBlock.Title
		resItem.CreationDate = infoBlock.Date
		resItem.Author = infoBlock.Author
		resItem.Content = strings.TrimSpace(splitted[2])
	}

	if resItem != nil {
		res = append(res, *resItem)
	}
	return res, nil
}

// LoadItems returns all the help items available
func LoadItems(dir string) ([]Item, error) {
	res := make([]Item, 10)

	file, err := readDirRecursive(dir)
	if err != nil {
		return res, err
	}

	return fileToItems(file)
}
