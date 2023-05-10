package files

import (
	"os"
	"syscall"
	"time"
)

func fileTimes(fi os.FileInfo) (aTime, mTime, cTime time.Time) {
	s := fi.Sys().(*syscall.Stat_t)
	aTime = time.Unix(s.Atimespec.Unix())
	mTime = fi.ModTime()
	cTime = time.Unix(s.Ctimespec.Unix())
	return
}

func fileOwner(fi os.FileInfo) (uid, gid int) {
	s := fi.Sys().(*syscall.Stat_t)
	return int(s.Uid), int(s.Gid)
}
