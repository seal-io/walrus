package files

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

type CopyOptions func(*copyOptions)

// CopyIfFound copies the src if found.
func CopyIfFound() CopyOptions {
	return func(o *copyOptions) {
		o.ignore = true
	}
}

// CopyInReplace deletes dst if found.
func CopyInReplace() CopyOptions {
	return func(o *copyOptions) {
		o.replace = true
	}
}

// CopyInShadow creates a symlink to the dst from the src.
func CopyInShadow() CopyOptions {
	return func(o *copyOptions) {
		o.shadow = true
	}
}

// CopyWithTimes preserves the times of all entries from src.
func CopyWithTimes() CopyOptions {
	return func(o *copyOptions) {
		o.preserveTimes = true
	}
}

// CopyWithOwner preserves the times of all entries from src.
func CopyWithOwner() CopyOptions {
	return func(o *copyOptions) {
		o.preserveOwner = true
	}
}

// CopyWithPerm specifies the perm of the dst,
// if the dst is directory,
// this only change the dst without its entries.
func CopyWithPerm(perm os.FileMode) CopyOptions {
	return func(o *copyOptions) {
		o.perm = &perm
	}
}

// CopyWithBuffer indicates to use buffer during copying file.
func CopyWithBuffer() CopyOptions {
	return func(o *copyOptions) {
		o.buffer = true
	}
}

// CopyWithModifier configures a modifier for change the content during copying file,
// it only works without buffer.
func CopyWithModifier(m func([]byte) ([]byte, error)) CopyOptions {
	return func(o *copyOptions) {
		o.modify = m
	}
}

type copyOptions struct {
	ignore        bool
	replace       bool
	shadow        bool
	preserveTimes bool
	preserveOwner bool
	perm          *os.FileMode
	buffer        bool
	modify        func([]byte) ([]byte, error)
}

func Copy(src, dst string, opts ...CopyOptions) error {
	var o copyOptions
	for i := range opts {
		if opts[i] == nil {
			continue
		}
		opts[i](&o)
	}

	var srcInfo, err = os.Lstat(src)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) && o.ignore {
			return nil
		}
		return err
	}

	var cp func(string, string, os.FileInfo, copyOptions) error
	for {
		switch m := srcInfo.Mode(); {
		default:
			return fmt.Errorf("unsupported src mode: %v", m)
		case m&os.ModeSymlink != 0:
			if src, err = os.Readlink(src); err != nil {
				return fmt.Errorf("failed to read origin of link src: %w", err)
			}
			if srcInfo, err = os.Lstat(src); err != nil {
				return fmt.Errorf("failed to stat origin of link src: %w", err)
			}
			continue
		case m.IsDir():
			cp = copyFromDir
		case m.IsRegular():
			cp = copyFromFile
		}
		break
	}

	var dstInfo = srcInfo
	if s, err := os.Lstat(dst); err == nil {
		if o.replace {
			if err = os.RemoveAll(dst); err != nil {
				return fmt.Errorf("cannot clean dst: %w", err)
			}
		}
		dstInfo = s
	}
	if o.shadow {
		if err = os.Symlink(src, dst); err != nil {
			return fmt.Errorf("cannot shadow dst: %w", err)
		}
		return nil
	}
	defer func() {
		var m = dstInfo.Mode()
		if o.perm != nil {
			m = *o.perm
		}
		if err = os.Chmod(dst, m); err != nil {
			err = fmt.Errorf("cannot change perm: %w", err)
		}
	}()

	if err = cp(src, dst, srcInfo, o); err != nil {
		return fmt.Errorf("cannot copy: %w", err)
	}

	if o.preserveTimes {
		var aTime, mTime, _ = fileTimes(srcInfo)
		if err = os.Chtimes(dst, aTime, mTime); err != nil {
			return fmt.Errorf("cannot preserve times: %w", err)
		}
	}
	if o.preserveOwner {
		var uid, gid = fileOwner(srcInfo)
		if err = os.Lchown(dst, uid, gid); err != nil {
			return fmt.Errorf("cannot preserve owner: %w", err)
		}
	}
	return nil
}

func copyFromDir(src, dst string, srcInfo os.FileInfo, o copyOptions) error {
	// Switch directory permission for copying.
	if err := os.Mkdir(dst, 0666); err != nil {
		if !os.IsExist(err) {
			return err
		}
		if err = os.Chmod(dst, 0666); err != nil {
			return err
		}
	}

	return filepath.Walk(src, func(srcSub string, srcSubInfo fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		sub, err := filepath.Rel(src, srcSub)
		if err != nil {
			return err
		}
		var dstSub = filepath.Join(dst, sub)

		switch m := srcSubInfo.Mode(); {
		default:
			// Unsupported src mode.
			return nil
		case m&os.ModeSymlink != 0:
			if srcSub, err = os.Readlink(srcSub); err != nil {
				return err
			}
			if err = os.Symlink(srcSub, dstSub); err != nil {
				return err
			}
		case m.IsDir():
			if err = os.Mkdir(dstSub, srcSubInfo.Mode()); err != nil && !os.IsExist(err) {
				return err
			}
		case m.IsRegular():
			if err = copyFromFile(srcSub, dstSub, srcSubInfo, o); err != nil {
				return err
			}
		}

		if o.preserveTimes {
			var aTime, mTime, _ = fileTimes(srcInfo)
			if err = os.Chtimes(dstSub, aTime, mTime); err != nil {
				return err
			}
		}
		if o.preserveOwner {
			var uid, gid = fileOwner(srcInfo)
			if err = os.Lchown(dstSub, uid, gid); err != nil {
				return err
			}
		}
		if srcSubInfo.Mode()&os.ModeSymlink == 0 {
			if err = os.Chmod(dstSub, srcSubInfo.Mode()); err != nil {
				return err
			}
		}

		return nil
	})
}

func copyFromFile(src, dst string, scrInfo os.FileInfo, o copyOptions) error {
	var dstDir = filepath.Dir(dst)
	if !Exists(dstDir) {
		if err := os.MkdirAll(dstDir, 0666); err != nil {
			return err
		}
	}

	if o.buffer {
		// Copy with buffer.
		var srcFile, err = os.Open(src)
		if err != nil {
			return err
		}
		defer func() { _ = srcFile.Close() }()
		dstFile, err := os.Create(dst)
		if err != nil {
			return err
		}
		defer func() { _ = dstFile.Close() }()
		_, err = io.Copy(srcFile, dstFile)
		return err
	}

	var content, err = os.ReadFile(src)
	if err != nil {
		return err
	}
	if o.modify != nil {
		content, err = o.modify(content)
		if err != nil {
			return err
		}
	}
	return os.WriteFile(dst, content, 0600)
}
