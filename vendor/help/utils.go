package help

import (
	"os"
	"path/filepath"
)

type file struct {
	Path string
}

// readDirRecursive reads the given `dir` recursively
func readDirRecursive(dir string) ([]file, error) {
	var res []file

	f := func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			f := file{
				Path: path,
			}
			res = append(res, f)
		}

		return nil
	}

	return res, filepath.Walk(dir, f)
}
