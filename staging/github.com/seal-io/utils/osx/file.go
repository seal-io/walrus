package osx

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Exists checks if the given path exists.
func Exists(path string, checks ...func(os.FileInfo) bool) bool {
	stat, err := os.Lstat(path)
	if err != nil {
		return false
	}

	for i := range checks {
		if checks[i] == nil {
			continue
		}

		if !checks[i](stat) {
			return false
		}
	}

	return true
}

// ExistsDir checks if the given path exists and is a directory.
func ExistsDir(path string) bool {
	return Exists(path, func(stat os.FileInfo) bool {
		return stat.Mode().IsDir()
	})
}

// ExistsLink checks if the given path exists and is a symbolic link.
func ExistsLink(path string) bool {
	return Exists(path, func(stat os.FileInfo) bool {
		return stat.Mode()&os.ModeSymlink != 0
	})
}

// ExistsFile checks if the given path exists and is a regular file.
func ExistsFile(path string) bool {
	return Exists(path, func(stat os.FileInfo) bool {
		return stat.Mode().IsRegular()
	})
}

// ExistsSocket checks if the given path exists and is a socket.
func ExistsSocket(path string) bool {
	return Exists(path, func(stat os.FileInfo) bool {
		return stat.Mode()&os.ModeSocket != 0
	})
}

// ExistsDevice checks if the given path exists and is a device.
func ExistsDevice(path string) bool {
	return Exists(path, func(stat os.FileInfo) bool {
		return stat.Mode()&os.ModeDevice != 0
	})
}

// TempFile creates a temporary file with the given pattern.
func TempFile(pattern string) string {
	f, err := os.CreateTemp("", pattern)
	if err != nil {
		panic(fmt.Errorf("create temp file: %w", err))
	}

	defer func() { _ = f.Close() }()

	return f.Name()
}

// TempDir creates a temporary directory with the given pattern.
func TempDir(pattern string) string {
	n, err := os.MkdirTemp("", pattern)
	if err != nil {
		panic(fmt.Errorf("create temp dir: %w", err))
	}

	return n
}

// SubTempDir is different to TempDir.
//
// TempDir creates a temporary directory randomly,
// but SubTempDir creates a subdirectory under the temporary directory with the given path.
func SubTempDir(path string) string {
	n := filepath.Join(os.TempDir(), filepath.Clean(path))
	err := os.MkdirAll(n, 0o700)
	if err != nil {
		panic(fmt.Errorf("create temp subdir: %w", err))
	}

	return n
}

// Close closes the given io.Closer without error.
func Close(c io.Closer) {
	if c == nil {
		return
	}
	_ = c.Close()
}

// IsEmptyDir checks if the given directory is empty.
func IsEmptyDir(dir string) bool {
	f, err := os.Open(dir)
	if err != nil {
		return false
	}
	defer Close(f)

	_, err = f.Readdir(1)
	return errors.Is(err, io.EOF)
}

// IsEmptyFile checks if the given file is empty.
func IsEmptyFile(file string) bool {
	s, err := os.Lstat(file)
	if err != nil {
		return false
	}
	if !s.Mode().IsRegular() {
		return false
	}
	return s.Size() == 0
}
