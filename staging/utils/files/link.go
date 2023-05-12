package files

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type LinkOptions func(*linkOptions)

// LinkEvenIfNotFound links the src,
// and create the src with given perm if not found.
func LinkEvenIfNotFound(isFile bool, perm os.FileMode) LinkOptions {
	return func(o *linkOptions) {
		o.create = true
		o.createFile = isFile

		if perm != 0 {
			o.createPerm = &perm
		}
	}
}

// LinkInReplace deletes dst if found.
func LinkInReplace() LinkOptions {
	return func(o *linkOptions) {
		o.replace = true
	}
}

// LinkInReal reads the real path from src.
func LinkInReal() LinkOptions {
	return func(o *linkOptions) {
		o.real = true
	}
}

// LinkWithTimes preserves the times of src.
func LinkWithTimes() LinkOptions {
	return func(o *linkOptions) {
		o.preserveTimes = true
	}
}

// LinkWithOwner preserves the owner of src.
func LinkWithOwner() LinkOptions {
	return func(o *linkOptions) {
		o.preserveOwner = true
	}
}

type linkOptions struct {
	create        bool
	createFile    bool
	createPerm    *os.FileMode
	replace       bool
	real          bool
	preserveTimes bool
	preserveOwner bool
}

func Link(src, dst string, opts ...LinkOptions) error {
	var o linkOptions

	for i := range opts {
		if opts[i] == nil {
			continue
		}

		opts[i](&o)
	}

	srcInfo, err := os.Lstat(src)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) || !o.create {
			return err
		}
		m := os.FileMode(0o666)

		if o.createPerm != nil {
			m = *o.createPerm
		}

		if !o.createFile {
			if err = os.MkdirAll(src, m); err != nil {
				return fmt.Errorf("cannot create src dir: %w", err)
			}
		} else {
			dstDir := filepath.Dir(dst)
			if !Exists(dstDir) {
				if err = os.MkdirAll(dstDir, 0o666); err != nil {
					return fmt.Errorf("cannot create dir of dst: %w", err)
				}
			}
			if _, err = os.Create(src); err != nil {
				return fmt.Errorf("cannot create src file: %w", err)
			}
			if err = os.Chmod(src, m); err != nil {
				return fmt.Errorf("cannot change src file mod: %w", err)
			}
		}
		srcInfo, _ = os.Stat(src)
	}

	if _, err := os.Lstat(dst); err == nil {
		if !o.replace {
			return errors.New("dst is not empty")
		}

		if err = os.Remove(dst); err != nil {
			return fmt.Errorf("cannot clean dst: %w", err)
		}
	} else {
		dstDir := filepath.Dir(dst)
		if !Exists(dstDir) {
			if err = os.MkdirAll(dstDir, 0o666); err != nil {
				return fmt.Errorf("cannot create dir of dst: %w", err)
			}
		}
	}

	if o.real {
		for {
			if srcInfo.Mode()&os.ModeSymlink == 0 {
				break
			}

			if src, err = os.Readlink(src); err != nil {
				return fmt.Errorf("failed to read origin of link src: %w", err)
			}

			if srcInfo, err = os.Lstat(src); err != nil {
				return fmt.Errorf("failed to stat origin of link src: %w", err)
			}
		}
	}

	if err = os.Symlink(src, dst); err != nil {
		return fmt.Errorf("cannot link dst: %w", err)
	}

	if o.preserveTimes {
		aTime, mTime, _ := fileTimes(srcInfo)
		if err = os.Chtimes(dst, aTime, mTime); err != nil {
			return fmt.Errorf("cannot preserve times: %w", err)
		}
	}

	if o.preserveOwner {
		uid, gid := fileOwner(srcInfo)
		if err = os.Lchown(dst, uid, gid); err != nil {
			return fmt.Errorf("cannot preserve owner: %w", err)
		}
	}

	return nil
}
