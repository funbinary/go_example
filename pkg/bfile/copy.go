package bfile

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/funbinary/go_example/pkg/errors"
)

// Copy
//  @Description: 拷贝文件和目录
//  @param src
//  @param dst
//  @return error
//
func Copy(src string, dst string) error {
	if src == "" {
		return errors.New("source path cannot be empty")
	}
	if dst == "" {
		return errors.New("destination path cannot be empty")
	}
	if IsFile(src) {
		return CopyFile(src, dst)
	}
	return CopyDir(src, dst)
}

// CopyFile
//  @Description: 拷贝文件
//  @param src
//  @param dst
//  @return err
//
func CopyFile(src, dst string) (err error) {
	if src == "" {
		return errors.New("source file cannot be empty")
	}
	if dst == "" {
		return errors.New("destination file cannot be empty")
	}
	// If src and dst are the same path, it does nothing.
	if src == dst {
		return nil
	}
	in, err := Open(src)
	if err != nil {
		return
	}
	defer func() {
		if e := in.Close(); e != nil {
			err = errors.Wrapf(e, `file close failed for "%s"`, src)
		}
	}()
	out, err := Create(dst)
	if err != nil {
		return
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = errors.Wrapf(e, `file close failed for "%s"`, dst)
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		err = errors.Wrapf(err, `io.Copy failed from "%s" to "%s"`, src, dst)
		return
	}
	if err = out.Sync(); err != nil {
		err = errors.Wrapf(err, `file sync failed for file "%s"`, dst)
		return
	}
	err = Chmod(dst, DefaultPermCopy)
	if err != nil {
		return
	}
	return
}

// CopyDir
//  @Description: 拷贝目录
//  @param src
//  @param dst
//  @return err
//
func CopyDir(src string, dst string) (err error) {
	if src == "" {
		return errors.New("source directory cannot be empty")
	}
	if dst == "" {
		return errors.New("destination directory cannot be empty")
	}
	// If src and dst are the same path, it does nothing.
	if src == dst {
		return
	}
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)
	si, err := Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return errors.New("source is not a directory")
	}
	if !Exists(dst) {
		if err = os.MkdirAll(dst, DefaultPermCopy); err != nil {
			err = errors.Wrapf(err, `create directory failed for path "%s", perm "%s"`, dst, DefaultPermCopy)
			return
		}
	}
	entries, err := ioutil.ReadDir(src)
	if err != nil {
		err = errors.Wrapf(err, `read directory failed for path "%s"`, src)
		return
	}
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())
		if entry.IsDir() {
			if err = CopyDir(srcPath, dstPath); err != nil {
				return
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}
			if err = CopyFile(srcPath, dstPath); err != nil {
				return
			}
		}
	}
	return
}
