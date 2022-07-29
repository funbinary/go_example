package bfile

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"

	"github.com/funbinary/go_example/pkg/errors"
)

var (
	// 默认读取缓冲区大小
	DefaultReadBuffer = 1024
)

// GetContents
//  @Description: 读取文件的内容，并将其转为字符串
//  @param path 文件路径
//  @return string 文件内容
//
func GetContents(path string) string {
	return string(GetBytes(path))
}

// GetBytes
//  @Description: 读取文件的内容
//  @param path 文件路径
//  @return []byte 文件内容
//
func GetBytes(path string) []byte {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil
	}
	return data
}

// SetContents
//  @Description: 往指定路径文件添加字符串内容。如果文件不存在将会递归的形式自动创建，如果文件有内容会被清空
//  @param path 文件所在的路径
//  @param content 要写入文件的内容
//  @return error
//
func SetContents(path string, content string) error {
	return setContents(path, []byte(content), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, DefaultPermOpen)
}

// SetBytes
//  @Description: 以字节形式写入文件。如果文件不存在将会递归的形式自动创建，如果文件有内容会被清空
//  @param path 文件所在的路径
//  @param content 字节数组
//  @return error
//
func SetBytes(path string, content []byte) error {
	return setContents(path, content, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, DefaultPermOpen)
}

// AppendBytes
//  @Description: 以字节形式写入文件。如果文件不存在将会递归的形式自动创建
//  @param path 文件所在的路径
//  @param content 字节数组
//  @return error
//
func AppendBytes(path string, content []byte) error {
	return setContents(path, content, os.O_WRONLY|os.O_CREATE|os.O_APPEND, DefaultPermOpen)
}

// AppendContents
//  @Description: 以字符串形式写入文件，使用追加的方式。如果文件不存在将会递归的形式自动创建
//  @param path 文件所在的路径
//  @param content 字节数组
//  @return error
//
func AppendContents(path string, content string) error {
	return setContents(path, []byte(content), os.O_WRONLY|os.O_CREATE|os.O_APPEND, DefaultPermOpen)
}

func setContents(path string, data []byte, flag int, perm os.FileMode) error {
	// 递归创建目录
	dir := Dir(path)
	if !Exists(dir) {
		if err := Mkdir(dir); err != nil {
			return err
		}
	}
	// 打开文件
	f, err := OpenFile(path, flag, perm)
	if err != nil {
		return err
	}
	defer f.Close()
	// 写入数据
	var n int
	if n, err = f.Write(data); err != nil {
		err = errors.Wrapf(err, `Write data to file "%s" failed`, path)
		return err
	} else if n < len(data) {
		return io.ErrShortWrite
	}
	return nil
}

// GetCharFromOffset
//  @Description: 从某个偏移量开始，获取文件中指定字符所在下标
//  @param reader 文件接口
//  @param char 要查找的字符
//  @param start 偏移量
//  @return int64 如果找到返回序号，如果找不到则返回-1
//
func GetCharFromOffset(reader io.ReaderAt, char byte, start int64) int64 {
	buffer := make([]byte, DefaultReadBuffer)
	offset := start
	for {
		if n, err := reader.ReadAt(buffer, offset); n > 0 {
			for i := 0; i < n; i++ {
				if buffer[i] == char {
					return int64(i) + offset
				}
			}
			offset += int64(n)
		} else if err != nil {
			break
		}
	}
	return -1
}

// GetCharOffsetFromByPath
//  @Description: 从某个偏移量开始，获取文件中指定字符所在下标
//  @param path 文件路径
//  @param char 要查找的字符
//  @param start 偏移量
//  @return int64 如果找到返回序号，如果找不到则返回-1
//
func GetCharOffsetFromByPath(path string, char byte, start int64) int64 {
	if f, err := OpenFile(path, os.O_RDONLY, DefaultPermOpen); err == nil {
		defer f.Close()
		return GetCharFromOffset(f, char, start)
	}
	return -1
}

// GetBytesByRange
//  @Description: 读取指定区间的文件字节内容，左闭右开
//  @param reader 文件接口
//  @param start 起始偏移量
//  @param end  终止偏移量
//  @return []byte
//
func GetBytesByRange(reader io.ReaderAt, start int64, end int64) []byte {
	buffer := make([]byte, end-start)
	if _, err := reader.ReadAt(buffer, start); err != nil {
		return nil
	}
	return buffer
}

// GetBytesByRangesByPath
//  @Description: 读取指定区间的文件字节内容，左闭右开
//  @param path 文件路径
//  @param start 起始偏移量
//  @param end  终止偏移量
//  @return []byte
//
func GetBytesByRangesByPath(path string, start int64, end int64) []byte {
	if f, err := OpenFile(path, os.O_RDONLY, DefaultPermOpen); err == nil {
		defer f.Close()
		return GetBytesByRange(f, start, end)
	}
	return nil
}

// GetBytesTilChar
//  @Description: 从文件中截取start到char字符位置的文件内容，以字节形式返回
//  @param reader 文件接口
//  @param start 起始位置
//  @param char 字符
//  @return []byte
//  @return int64 返回的字节数组数量
//
func GetBytesTilChar(reader io.ReaderAt, start int64, char byte) ([]byte, int64) {
	if offset := GetCharFromOffset(reader, char, start); offset != -1 {
		return GetBytesByRange(reader, start, offset+1), offset
	}
	return nil, -1
}

// GetBytesTilCharByPath
//  @Description: 从文件中截取start到char字符位置的文件内容，以字节形式返回
//  @param path 文件路径
//  @param start 起始位置
//  @param char 字符
//  @return []byte
//  @return int64 返回的字节数组数量
//
func GetBytesTilCharByPath(path string, start int64, char byte) ([]byte, int64) {
	if f, err := OpenFile(path, os.O_RDONLY, DefaultPermOpen); err == nil {
		defer f.Close()
		return GetBytesTilChar(f, start, char)
	}
	return nil, -1
}

// ReadLines
//  @Description: 以字符串的形式逐行读取文件内容
//  @param file 文件名
//  @param callback
//  @return error
//
func ReadLines(file string, callback func(text string) error) error {
	f, err := Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if err = callback(scanner.Text()); err != nil {
			return err
		}
	}
	return nil
}

// ReadLinesBytes
//  @Description:以字节形式逐行读取文件内容
//  @param file
//  @param callback
//  @return error
//
func ReadLinesBytes(file string, callback func(bytes []byte) error) error {
	f, err := Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if err = callback(scanner.Bytes()); err != nil {
			return err
		}
	}
	return nil
}

// Truncate
//  @Description: 将文件裁剪为指定大小
//  @param path 文件路径
//  @param size 截取的文件大小
//  @return err
//
func Truncate(path string, size int) (err error) {
	err = os.Truncate(path, int64(size))
	if err != nil {
		err = errors.Wrapf(err, `os.Truncate failed for file "%s", size "%d"`, path, size)
	}
	return
}
