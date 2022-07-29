package bfile

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/funbinary/go_example/pkg/errors"
)

const (
	Separator = string(filepath.Separator)

	DefaultPermOpen = os.FileMode(0666)

	DefaultPermCopy = os.FileMode(0777)
)

const (
	// Exactly one of O_RDONLY, O_WRONLY, or O_RDWR must be specified.
	O_RDONLY int = os.O_RDONLY // open the file read-only.
	O_WRONLY int = os.O_WRONLY // open the file write-only.
	O_RDWR   int = os.O_RDWR   // open the file read-write.
	// The remaining values may be or'ed in to control behavior.
	O_APPEND int = os.O_APPEND // append data to the file when writing.
	O_CREATE int = os.O_CREATE // create a new file if none exists.
	O_EXCL   int = os.O_EXCL   // used with O_CREATE, file must not exist.
	O_SYNC   int = os.O_SYNC   // open for synchronous I/O.
	O_TRUNC  int = os.O_TRUNC  // truncate regular writable file when opened.
)

// Mkdir
//  @Description: 递归创建目录
//  @param path 路径，建议使用绝对路径，不推荐使用相对路径
//  @return err 如果创建失败会返回错误的原因，否则返回nil
//
func Mkdir(path string) (err error) {
	if err = os.MkdirAll(path, os.ModePerm); err != nil {
		err = errors.Wrapf(err, `os.MkdirAll failed for path "%s" with perm "%d"`, path, os.ModePerm)
		return err
	}
	return nil
}

// Create
//  @Description: 创建文件,如果文件所在路径不存在，则会自动创建文件夹及文件，其权限为0666（任何人都可读写，不可执行）如果创建的文件已存在则会清空该文件的内容。
//  @param path: 文件路径
//  @return *os.File 文件指针
//  @return error
//
func Create(path string) (*os.File, error) {
	dir := Dir(path)
	if !Exists(dir) {
		if err := Mkdir(dir); err != nil {
			return nil, err
		}
	}
	file, err := os.Create(path)
	if err != nil {
		err = errors.Wrapf(err, `os.Create failed for name "%s"`, path)
	}
	return file, err
}

// Remove
//  @Description: 删除给定的文件或文件夹
//  @param path 给定文件或文件夹的路径
//  @return err
//
func Remove(path string) (err error) {
	err = os.RemoveAll(path)
	if err != nil {
		err = errors.Wrapf(err, `os.RemoveAll failed for path "%s"`, path)
	}
	return err
}

// Exists
//  @Description: 判断路径是否存在
//  @param path 文件或者目录的路径
//  @return bool 如果文件或目录存在则返回true，否则返回false
//
func Exists(path string) bool {
	if stat, err := os.Stat(path); stat != nil && !os.IsNotExist(err) {
		return true
	}
	return false
}

// Open
//  @Description: 以只读的方式打开文件
//  @param path
//  @return *os.File
//  @return error
//
func Open(path string) (*os.File, error) {
	file, err := os.Open(path)
	if err != nil {
		err = errors.Wrapf(err, `os.Open failed for name "%s"`, path)
	}
	return file, err
}

// OpenFile
//  @Description: 以指定的flag和perm打开文件
//  @param path
//  @param flag O_RDONLY、O_WRONLY、O_RDWR、O_APPEND、O_CREATE、O_EXCL、O_SYNC、O_TRUNC
//  @param perm 0666、0777
//  @return *os.File
//  @return error
//
func OpenFile(path string, flag int, perm os.FileMode) (*os.File, error) {
	file, err := os.OpenFile(path, flag, perm)
	if err != nil {
		err = errors.Wrapf(err, `os.OpenFile failed with name "%s", flag "%d", perm "%d"`, path, flag, perm)
	}
	return file, err
}

// Stat
//  @Description: 获取文件信息
//  @param path
//  @return os.FileInfo
//  @return error
//
func Stat(path string) (os.FileInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		err = errors.Wrapf(err, `os.Stat failed for file "%s"`, path)
	}
	return info, err
}

// Move
//  @Description: 將文件src移动到dst
//  @param src
//  @param dst
//  @return err
//
func Move(src string, dst string) (err error) {
	err = os.Rename(src, dst)
	if err != nil {
		err = errors.Wrapf(err, `os.Rename failed from "%s" to "%s"`, src, dst)
	}
	return
}

// Rename
//  @Description: 将文件/文件夹src重命名为dst
//  @param src
//  @param dst
//  @return error
//
func Rename(src string, dst string) error {
	return Move(src, dst)
}

// IsEmpty
//  @Description: 判断文件或者目录是否为空
//  @param path
//  @return bool 如果文件或目录为空将会返回true，否则返回false
//
func IsEmpty(path string) bool {
	stat, err := Stat(path)
	if err != nil {
		return true
	}
	if stat.IsDir() {
		file, err := os.Open(path)
		if err != nil {
			return true
		}
		defer file.Close()
		names, err := file.Readdirnames(-1)
		if err != nil {
			return true
		}
		return len(names) == 0
	} else {
		return stat.Size() == 0
	}
}

// DirSubs
//  @Description: 获取文件夹下所有文件和文件夹，非递归
//  @param path
//  @return []string path下的文件和文件夹列表
//  @return error
//
func DirSubs(path string) ([]string, error) {
	f, err := Open(path)
	if err != nil {
		return nil, err
	}
	list, err := f.Readdirnames(-1)
	_ = f.Close()
	if err != nil {
		err = errors.Wrapf(err, `Read dir files failed from path "%s"`, path)
		return nil, err
	}
	return list, nil
}

// Glob
//  @Description: 模糊搜索给定路径下的文件列表，支持正则，第二个参数控制返回的结果是否带上绝对路径。
//  @param pattern
//  @param onlyNames 是否带上
//  @return []string
//  @return error
//
func Glob(pattern string, onlyNames ...bool) ([]string, error) {
	list, err := filepath.Glob(pattern)
	if err != nil {
		err = errors.Wrapf(err, `filepath.Glob failed for pattern "%s"`, pattern)
		return nil, err
	}
	if len(onlyNames) > 0 && onlyNames[0] && len(list) > 0 {
		array := make([]string, len(list))
		for k, v := range list {
			array[k] = Basename(v)
		}
		return array, nil
	}
	return list, nil
}

// IsReadable
//  @Description: 判断文件/目录是否可读
//  @param path
//  @return bool
//
func IsReadable(path string) bool {
	result := true
	file, err := os.OpenFile(path, os.O_RDONLY, DefaultPermOpen)
	if err != nil {
		result = false
	}
	file.Close()
	return result
}

// IsWritable
//  @Description: 判断文件或目录是否可写
//  @param path
//  @return bool
//
func IsWritable(path string) bool {
	result := true
	if IsDir(path) {
		// If it's a directory, create a temporary file to test whether it's writable.
		tmpFile := strings.TrimRight(path, Separator) + Separator + strconv.Itoa(int(time.Now().UnixNano()))
		if f, err := Create(tmpFile); err != nil || !Exists(tmpFile) {
			result = false
		} else {
			_ = f.Close()
			_ = Remove(tmpFile)
		}
	} else {
		// If it's a file, check if it can open it.
		file, err := os.OpenFile(path, os.O_WRONLY, DefaultPermOpen)
		if err != nil {
			result = false
		}
		_ = file.Close()
	}
	return result
}
