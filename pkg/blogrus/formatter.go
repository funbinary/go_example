package blogrus

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"path"
	"sort"
	"strings"
)

const (
	DefaultTimestampFormat = "2006-1-2 15:04:05.000"
)

type BFormatter struct {
	TimestampFormat string // 时间序列格式 default: DefaultTimestampFormat = "2006-1-2 15:04:05.000"
	HideKeys        bool   // WriteField字段是否隐藏其key值
	NoColor         bool   // 是否禁用颜色输出
	ShowFullLevel   bool   // 是否显示日志级别所有内容例如[WARNING]，否则[W]
	TrimMessages    bool   // 是否移除末尾的空行

	buffer *bytes.Buffer
}

func NewDefaultBFormatter() *BFormatter {
	return &BFormatter{
		TimestampFormat: DefaultTimestampFormat,
		HideKeys:        false,
		NoColor:         false,
		ShowFullLevel:   true,
		TrimMessages:    false,
	}
}

func (self *BFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	if self.buffer == nil {
		self.buffer = &bytes.Buffer{}
	} else {
		self.buffer.Reset()
	}
	if self.TimestampFormat == "" {
		self.TimestampFormat = DefaultTimestampFormat
	}
	levelColor := getColorByLevel(entry.Level)
	self.wrapColor(levelColor, func() {
		self.WriteCommon(entry)
		self.WriteFields(entry)
		self.WriteMessages(entry)
	})
	self.buffer.WriteByte('\n')
	return self.buffer.Bytes(), nil
}

func (self *BFormatter) WriteCommon(entry *logrus.Entry) {
	// 日志级别
	self.WriteLevel(entry)
	// 日期
	self.WriteTime(entry)
	// 堆栈信息
	self.WriteCaller(entry)
}

func (self *BFormatter) WriteFields(entry *logrus.Entry) {
	// write fields
	if len(entry.Data) != 0 {
		fields := make([]string, 0, len(entry.Data))
		for field := range entry.Data {
			fields = append(fields, field)
		}

		sort.Strings(fields)

		for _, field := range fields {
			if self.HideKeys {
				fmt.Fprintf(self.buffer, "[%v] ", entry.Data[field])
			} else {
				fmt.Fprintf(self.buffer, "[%s:%v] ", field, entry.Data[field])
			}
		}
	}
}

func (self *BFormatter) WriteMessages(entry *logrus.Entry) {
	if self.TrimMessages {
		self.buffer.WriteString(strings.TrimSpace(entry.Message)) //去除末尾的空格
	} else {
		self.buffer.WriteString(entry.Message)
	}
}

func (self *BFormatter) WriteLevel(entry *logrus.Entry) {
	self.wrapFrame(func() {
		level := strings.ToUpper(entry.Level.String())
		if self.ShowFullLevel {
			self.buffer.WriteString(level)
		} else {
			self.buffer.WriteString(level[:1])
		}
	})
}

func (self *BFormatter) WriteTime(entry *logrus.Entry) {
	self.wrapFrame(func() {
		self.buffer.WriteString(entry.Time.Format(self.TimestampFormat))
	})
}

func (self *BFormatter) WriteCaller(entry *logrus.Entry) {
	if entry.HasCaller() {
		functionName := entry.Caller.Function
		fileName := path.Base(entry.Caller.File)
		line := entry.Caller.Line
		fmt.Fprintf(self.buffer, "[%s:%d] [%s] ", fileName, line, functionName)
	}
}

func (self *BFormatter) wrapColor(levelColor int, run func()) {
	if !self.NoColor {
		fmt.Fprintf(self.buffer, "\x1b[%dm", levelColor)
	}
	run()
	if !self.NoColor {
		self.buffer.WriteString("\x1b[0m")
	}
}

func (self *BFormatter) wrapColorAndFrame(levelColor int, run func()) {
	if !self.NoColor {
		fmt.Fprintf(self.buffer, "\x1b[%dm", levelColor)
	}
	self.buffer.WriteString("[")
	run()
	if !self.NoColor {
		self.buffer.WriteString("\x1b[0m")
	}
	self.buffer.WriteString("] ")
}

func (self *BFormatter) wrapFrame(run func()) {
	self.buffer.WriteString("[")
	run()
	self.buffer.WriteString("] ")
}
