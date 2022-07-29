package bfile

import (
	"os"
	"time"
)

//
// MTime
//  @Description: 获取文件的修改时间，返回时间类型
//  @param path
//  @return time.Time
//
func MTime(path string) time.Time {
	s, e := os.Stat(path)
	if e != nil {
		return time.Time{}
	}
	return s.ModTime()
}

//
// MTimestamp
//  @Description: 获取文件的修改时间，返回unix时间类型，即int64
//  @param path
//  @return int64
//
func MTimestamp(path string) int64 {
	mtime := MTime(path)
	if mtime.IsZero() {
		return -1
	}
	return mtime.Unix()
}

//
// MTimestampMilli
//  @Description: 获取文件的修改时间毫秒级别，返回unix时间类型，即int64
//  @param path
//  @return int64
//
func MTimestampMilli(path string) int64 {
	mtime := MTime(path)
	if mtime.IsZero() {
		return -1
	}
	return mtime.UnixNano() / 1000000
}
