package blog

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const (
	defaltTimeFormat = "2006-01-02"
)

var (
	//函数别名
	currentTime = time.Now
	// 兆
	megabyte = 1024 * 1024
)

type Writer struct {

	//文件路径
	filePath string
	//时间格式,默认"2006-01-02"
	timeFormat string
	//模块名称
	moduleName string
	//文件的最大值，0不转轮
	maxSize int

	//要写入的文件名称: 文件路径 + 日期 + "/kds/log/2021-08-03/bysc_ftp_1.log"
	fileName string
	//当前文件大小
	size int64
	// 文件指针
	file *os.File
	//当前写入的文件编号
	currentNum int
	//当前写入的日期
	currentTime time.Time

	// mu sync.Mutex
}

// 创建Writer
// path:文件的目录
// format: 日期的格式化，请使用go起点日期2006-01-02 15:04:05
// modulename: 模块名称
func NewWriter(path string, format string, n string, size int) *Writer {
	if _, err := time.Parse(format, format); err != nil {
		format = defaltTimeFormat
	}

	if path == "" {
		path, _ = os.Getwd()
	}

	w := &Writer{
		filePath:    path,
		timeFormat:  format,
		moduleName:  n,
		maxSize:     size,
		currentNum:  1,
		currentTime: currentTime(),
	}

	if err := w.openExistingOrNew(); err != nil {
		panic(err)
	}
	return w
}

func (w *Writer) Write(p []byte) (n int, err error) {

	// w.mu.Lock()
	// defer w.mu.Unlock()

	writeLen := int64(len(p))

	// 1. 判断日期
	now, _ := time.Parse(w.timeFormat, currentTime().Format(w.timeFormat))
	if w.currentTime.IsZero() ||
		now.Year() != w.currentTime.Year() ||
		now.Month() != w.currentTime.Month() ||
		now.Day() != w.currentTime.Day() {
		w.currentTime = now
		w.currentNum = 1
		w.rotate()
	}

	if w.file == nil {
		if err = w.openExistingOrNew(); err != nil {
			return 0, err
		}
	}

	if int64(w.maxSize) > 0 && w.size+writeLen > w.max() {
		w.currentNum++
		if err := w.rotate(); err != nil {
			return 0, err
		}
	}
	// info, _ := w.file.Stat()
	n, err = w.file.Write(p)
	os.Stdout.Write(p)
	w.size += int64(n)

	return n, err
}

func (w *Writer) max() int64 {
	return int64(w.maxSize) * int64(megabyte)
}

func (w *Writer) close() error {
	if w.file == nil {
		return nil
	}
	err := w.file.Close()
	w.file = nil
	return err
}

func (w *Writer) rotate() error {
	if err := w.close(); err != nil {
		return err
	}
	if err := w.openNew(); err != nil {
		return err
	}
	return nil
}

func (w *Writer) openNew() error {
	w.setFilename()

	os.MkdirAll(w.dir(), 0744)
	f, err := os.OpenFile(w.fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("打开日志文件失败:%s", err)
	}
	w.file = f
	w.size = 0
	return nil
}

func (w *Writer) openExistingOrNew() error {
	//1. 判断文件是否存在
	w.setFilename()
	info, err := os.Stat(w.fileName)
	if os.IsNotExist(err) {
		return w.openNew()
	}
	if err != nil {
		return fmt.Errorf("获取日志信息失败:%s", err)
	}

	//2. 打开文件
	f, err := os.OpenFile(w.fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return w.openNew()
	}
	w.file = f
	w.size = info.Size()
	return nil
}

func (w *Writer) setFilename() error {
	if _, err := time.Parse(w.timeFormat, w.timeFormat); w.filePath == "" || err != nil {
		panic("日志目录未设置")
	}
	if w.moduleName == "" {
		w.moduleName = "writer"
	}
	//1. 拼接路径 文件路径 + 日期 + "/kds/log/2021-08-03/bysc_ftp_1.log"
	path := w.filePath + "/" + w.currentTime.Format(w.timeFormat)
	//2. 拼接文件名
	name := w.moduleName + "_" + strconv.Itoa(w.currentNum) + ".log"
	w.fileName = path + "/" + name
	return nil
}

//返回当前写入的日志文件目录
func (w *Writer) dir() string {
	if w.fileName == "" {
		w.setFilename()
	}
	return filepath.Dir(w.fileName)
}

//返回当前写入的日志文件名
// func (w *Writer) logFileName() string {
// 	if w.fileName == "" {
// 		w.setFilename()
// 	}

// 	return filepath.Base(w.fileName)
// }
