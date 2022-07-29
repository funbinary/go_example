package bfile

import (
	"bufio"
	"log"
	"os"
)

//
// IsFile
//  @Description: 判断路径是否是文件
//  @param path
//  @return bool 当路径为文件时返回true，否则返回false
//
func IsFile(path string) bool {
	s, err := Stat(path)
	if err != nil {
		return false
	}
	return !s.IsDir()
}

//
// IsDir
//  @Description: 判断路径是否是目录
//  @param path
//  @return bool 当路径为目录时返回true，否则返回false
//
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// IsBinary
//  @Description: 判断文件是否是二进制类型
//  @param path 文件路径
//  @return bool true - 文件为二进制类型
//
func IsBinary(path string) bool {

	file, err := os.Open(path)
	if err != nil {
		log.Printf("\033[31merror : IO error - \033[0m%s", err)
		return false
	}
	defer file.Close()

	r := bufio.NewReader(file)
	buf := make([]byte, 1024)
	n, err := r.Read(buf)
	if n <= 0 {
		return false
	}
	var white_byte int = 0
	for i := 0; i < n; i++ {
		if (buf[i] >= 0x20 && buf[i] <= 0xff) ||
			buf[i] == 9 ||
			buf[i] == 10 ||
			buf[i] == 13 {
			white_byte++
		} else if buf[i] <= 6 || (buf[i] >= 14 && buf[i] <= 31) {
			return true
		}
	}

	if white_byte >= 1 {
		return false
	}
	return true
}
