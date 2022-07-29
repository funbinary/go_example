package bfile

import (
	"github.com/funbinary/go_example/pkg/errors"

	"os"
)

// Pwd
//  @Description: 获取当前的工作目录
//  @return string
//
func Pwd() string {
	path, err := os.Getwd()
	if err != nil {
		return ""
	}
	return path
}

// Chmod
//  @Description: 修改权限
//  @param path
//  @param mode
//  @return err
//
func Chmod(path string, mode os.FileMode) (err error) {
	err = os.Chmod(path, mode)
	if err != nil {
		err = errors.Wrapf(err, `os.Chmod failed with path "%s" and mode "%s"`, path, mode)
	}
	return
}

// Chdir
//  @Description: 更改当前的工作路径
//  @param dir
//  @return err
//
func Chdir(dir string) (err error) {
	err = os.Chdir(dir)
	if err != nil {
		err = errors.Wrapf(err, `os.Chdir failed with dir "%s"`, dir)
	}
	return
}
