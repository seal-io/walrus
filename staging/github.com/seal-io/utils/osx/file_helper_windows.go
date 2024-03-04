package osx

import (
	"os"
	"syscall"
	"time"
)

func fileTimes(fi os.FileInfo) (aTime, mTime, cTime time.Time) {
	s := fi.Sys().(*syscall.Win32FileAttributeData)
	aTime = time.Unix(0, s.LastAccessTime.Nanoseconds())
	mTime = time.Unix(0, s.LastWriteTime.Nanoseconds())
	cTime = time.Unix(0, s.CreationTime.Nanoseconds())

	return
}

func fileOwner(fi os.FileInfo) (uid, gid int) {
	return
}
