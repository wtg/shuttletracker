package api

import (
	"net/http"
	"os"
)

type staticFileSystem struct {
	fs http.FileSystem
}

// Open will only return Files that are not directories.
func (sfs staticFileSystem) Open(path string) (http.File, error) {
	f, err := sfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if s.IsDir() {
		return nil, os.ErrNotExist
	}

	return f, nil
}
