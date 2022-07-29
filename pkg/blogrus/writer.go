package blogrus

import (
	"fmt"
	"os"
	"time"
)

type MultiWriter struct {
	fileWriter *FileWriter
}

func NewDefaultMultiWriter(filename string, maxSize int64, maxLines int) *MultiWriter {
	fd, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		fmt.Errorf("打开日志文件%s失败", filename)
		return nil
	}
	ffinfo, err := fd.Stat()
	if err != nil {
		fmt.Errorf("获取日志文件%s状态失败", filename)
		fd.Close()
		return nil
	}

	fileWriter := &FileWriter{
		FileName:       filename,
		MaxSize:        maxSize,
		MaxLines:       maxLines,
		Daily:          true,
		curFd:          fd,
		curSize:        ffinfo.Size(),
		lastUpdateDate: time.Now().Day(),
	}
	mul := &MultiWriter{
		fileWriter: fileWriter,
	}
	return mul
}

func (self *MultiWriter) Write(p []byte) (n int, err error) {
	os.Stdout.Write(p)
	return self.fileWriter.Write(p)
}

type FileWriter struct {
	FileName       string
	MaxSize        int64
	MaxLines       int
	Daily          bool
	curFd          *os.File
	curLines       int
	curSize        int64
	lastUpdateDate int
}

func NewDefaultFileWriter(filename string, maxSize int64, maxLines int) *FileWriter {
	fd, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		fmt.Errorf("打开日志文件%s失败", filename)
		return nil
	}
	ffinfo, err := fd.Stat()
	if err != nil {
		fmt.Errorf("获取日志文件%s状态失败", filename)
		fd.Close()
		return nil
	}

	fileWriter := &FileWriter{
		FileName:       filename,
		MaxSize:        maxSize,
		MaxLines:       maxLines,
		Daily:          true,
		curFd:          fd,
		curSize:        ffinfo.Size(),
		lastUpdateDate: time.Now().Day(),
	}
	return fileWriter
}

func (self *FileWriter) Write(p []byte) (n int, err error) {
	self.docheckAndRotate(len(p))
	return self.curFd.Write(p)
}

func (self *FileWriter) docheckAndRotate(size int) {
	n := self.curSize + int64(size)
	if (self.MaxLines > 0 && self.curLines >= self.MaxLines) ||
		(self.MaxSize > 0 && n >= self.MaxSize) ||
		(self.Daily && time.Now().Day() != self.lastUpdateDate) {
		if err := self.doRotate(); err != nil {
			fmt.Fprintf(os.Stderr, "FileLogWriter(%q): %s\n", self.FileName, err)
			return
		}
	}
	self.curLines++
	self.curSize += int64(size)
}

func (self *FileWriter) doRotate() error {
	_, err := os.Lstat(self.FileName)
	if err == nil { // file exists
		// Find the next available number
		num := 1
		fname := ""
		for ; err == nil && num <= 999; num++ {
			fname = self.FileName + fmt.Sprintf(".%s.%03d", time.Now().Format("2006-01-02"), num)
			_, err = os.Lstat(fname)
		}
		// return error if the last file checked still existed
		if err == nil {
			return fmt.Errorf("Rotate: Cannot find free log number to rename %s", self.FileName)
		}

		self.curFd.Close()

		//重命名
		err = os.Rename(self.FileName, fname)
		if err != nil {
			return fmt.Errorf("Rotate: %s", err)
		}

		// re-start logger
		self.curFd, err = os.OpenFile(self.FileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
		self.curLines = 0
		self.curSize = 0
		if err != nil {
			return fmt.Errorf("Rotate StartLogger: %s", err)
		}
	}
	return nil
}
