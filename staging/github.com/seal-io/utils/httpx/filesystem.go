package httpx

import (
	"io/fs"
	"net/http"
	"os"
	"time"
)

// FS converts fsys to a http.FileSystem implementation,
// for use with http.FileServer and http.NewFileTransport.
// The files provided by fsys must implement io.Seeker.
func FS(fsys fs.FS, opts ...*FSOption) http.FileSystem {
	var o *FSOption
	if len(opts) > 0 && opts[0] != nil {
		o = opts[0]
	} else {
		o = FSOptions()
	}

	return filesystem{
		FSOption:   o,
		FileSystem: http.FS(fsys),
	}
}

type filesystem struct {
	*FSOption
	http.FileSystem
}

func (fs filesystem) Open(name string) (http.File, error) {
	f, err := fs.FileSystem.Open(name)
	if err != nil {
		return nil, err
	}

	return file{FSOption: fs.FSOption, File: f}, nil
}

type file struct {
	*FSOption
	http.File
}

func (f file) Readdir(count int) ([]os.FileInfo, error) {
	if f.listable {
		return f.File.Readdir(count)
	}

	return nil, nil
}

func (f file) Stat() (fs.FileInfo, error) {
	i, err := f.File.Stat()
	if err != nil {
		return nil, err
	}
	return fileInfo{FSOption: f.FSOption, FileInfo: i}, nil
}

type fileInfo struct {
	*FSOption
	fs.FileInfo
}

var epoch = time.Now()

func (i fileInfo) ModTime() time.Time {
	if i.embedded {
		return epoch
	}
	return i.FileInfo.ModTime()
}
