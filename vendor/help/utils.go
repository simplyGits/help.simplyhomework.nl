package help

import (
	"io/ioutil"
	"path"
)

type file struct {
	Path  string
	IsDir bool
	Files []file
}

// readDirRecursive reads the given `dir` recursively
func readDirRecursive(dir string) (file, error) {
	res := file{
		Path:  dir + "/",
		IsDir: true,
	}

	infos, err := ioutil.ReadDir(dir)
	if err != nil {
		return res, err
	}

	for _, info := range infos {
		f := file{
			Path: path.Join(dir, info.Name()),
		}

		if info.IsDir() {
			f, err = readDirRecursive(f.Path)
			if err != nil {
				return f, err
			}
		} else {
			f.IsDir = false
		}

		res.Files = append(res.Files, f)
	}

	return res, nil
}
