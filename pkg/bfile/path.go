package bfile

import (
	"path/filepath"
	"strings"
)

// Join
//  @Description: 将多个字符串路径通过/进行拼接，windows为\
//  @param paths: 路径字符串
//  @return string 拼接后的字符串路径
//
func Join(paths ...string) string {
	var s string
	for _, path := range paths {
		if s != "" {
			s += Separator
		}
		s += strings.TrimRight(path, Separator)
	}
	return s
}

// Basename
//  @Description: 返回路径的最后一个元素，包含文件扩展名,如果path为空则返回.
//  @param path 路径
//  @return string 返回路径的最后一个元素，包含文件扩展名
//
func Basename(path string) string {
	return filepath.Base(path)
}

// Name
//  @Description: 返回当前程序名，不包含扩展名
//  @param path 路径
//  @return string 返回当前程序名，不包含扩展名
//
func Name(path string) string {
	base := filepath.Base(path)
	if i := strings.LastIndexByte(base, '.'); i != -1 {
		return base[:i]
	}
	return base
}

// Ext
//  @Description: 获取给定路径的扩展名，包含.
//  @param path 路径
//  @return string给定路径的扩展名，包含.
//
func Ext(path string) string {
	ext := filepath.Ext(path)
	if p := strings.IndexByte(ext, '?'); p != -1 {
		ext = ext[0:p]
	}
	return ext
}

// ExtName
//  @Description: 获取给定路径的扩展名，不包含.
//  @param path 路径
//  @return string 给定路径的扩展名
//
func ExtName(path string) string {
	return strings.TrimLeft(Ext(path), ".")
}

// Abs
//  @Description: 使用filepath的ABS API，返回路径的绝对路径
//  @param path:  路径
//  @return string 路径的绝对路径
//
func Abs(path string) string {
	p, _ := filepath.Abs(path)
	return p
}

// RealPath
//  @Description: 返回绝对路径
//  @param path 文件或者目录的路径
//  @return string 如果文件/目录存在则返回绝对路径，否则返回空字符串
//
func RealPath(path string) string {
	p, err := filepath.Abs(path)
	if err != nil {
		return ""
	}
	if !Exists(p) {
		return ""
	}
	return p
}

// Dir
//  @Description:
//  @param path 路径
//  @return string 给定路径的目录部分
//
func Dir(path string) string {
	if path == "." {
		return filepath.Dir(RealPath(path))
	}
	return filepath.Dir(path)
}
