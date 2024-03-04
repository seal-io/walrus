package vcs

import (
	"io"
	"io/fs"

	"github.com/go-git/go-billy/v5"
	"github.com/seal-io/utils/osx"
	"github.com/seal-io/utils/pools/bytespool"
)

var _ fs.ReadFileFS = (*ClonedFilesystem)(nil)

// ClonedFilesystem ports the fs.FS interface to the billy.Filesystem interface.
type ClonedFilesystem struct {
	billy.Filesystem
}

func (d ClonedFilesystem) Open(name string) (fs.File, error) {
	f, err := d.Filesystem.Open(name)
	if err != nil {
		return nil, err
	}

	return ClonedFile{File: f, dir: d.Filesystem}, nil
}

func (d ClonedFilesystem) ReadDir(name string) ([]fs.DirEntry, error) {
	ds, err := d.Filesystem.ReadDir(name)
	if err != nil {
		return nil, err
	}

	r := make([]fs.DirEntry, len(ds))
	for i := range ds {
		r[i] = fs.FileInfoToDirEntry(ds[i])
	}
	return r, nil
}

func (d ClonedFilesystem) ReadFile(name string) ([]byte, error) {
	f, err := d.Filesystem.Open(name)
	if err != nil {
		return nil, err
	}
	defer osx.Close(f)

	buf := bytespool.GetBuffer()
	_, err = io.Copy(buf, f)
	return buf.Bytes(), err
}

var _ fs.File = (*ClonedFile)(nil)

type ClonedFile struct {
	billy.File

	dir billy.Filesystem
}

func (f ClonedFile) Stat() (fs.FileInfo, error) {
	return f.dir.Stat(f.Name())
}
