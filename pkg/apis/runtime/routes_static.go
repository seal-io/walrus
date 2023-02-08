package runtime

import (
	"io/fs"
	"net/http"
	"os"
	"time"
)

type StaticHttpFileSystem struct {
	http.FileSystem

	Listable bool
	Embedded bool
}

func (fs StaticHttpFileSystem) Open(name string) (http.File, error) {
	f, err := fs.FileSystem.Open(name)
	if err != nil {
		return nil, err
	}
	return StaticHttpFile{File: f, Listable: fs.Listable, Embedded: fs.Embedded}, nil
}

type StaticHttpFile struct {
	http.File

	Listable bool
	Embedded bool
}

func (f StaticHttpFile) Readdir(count int) ([]os.FileInfo, error) {
	if f.Listable {
		return f.File.Readdir(count)
	}
	return nil, nil
}

func (f StaticHttpFile) Stat() (fs.FileInfo, error) {
	var i, err = f.File.Stat()
	if err != nil {
		return nil, err
	}
	if f.Embedded {
		return embeddedFileInfo{FileInfo: i}, nil
	}
	return i, nil
}

type embeddedFileInfo struct {
	fs.FileInfo
}

var embeddedAt = time.Now()

func (i embeddedFileInfo) ModTime() time.Time {
	return embeddedAt
}
