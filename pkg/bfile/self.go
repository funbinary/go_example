package bfile

import (
	"os"
	"os/exec"
	"path/filepath"
)

var (
	//  当前运行的程序的路径
	selfPath = ""
)

func init() {
	selfPath, _ = exec.LookPath(os.Args[0])
	if selfPath != "" {
		selfPath, _ = filepath.Abs(selfPath)
	}
	if selfPath == "" {
		selfPath, _ = filepath.Abs(os.Args[0])
	}
}

//
// SelfPath
//  @Description: 获取当前运行程序的路径，包含程序名
//  @return string 当前运行程序的路径
//
func SelfPath() string {
	return selfPath
}

// SelfName
//  @Description: 获取当前运行程序的名称
//  @return string 当前运行程序的名称
//
func SelfName() string {
	return Basename(SelfPath())
}

//
// SelfDir
//  @Description: 获取当前程序所在的目录
//  @return string 当前程序所在的目录
//
func SelfDir() string {
	return filepath.Dir(SelfPath())
}
