package blogrus

import "github.com/sirupsen/logrus"

const (
	colorBlack  = 30 //黑色
	colorRed    = 31 //红色
	colorGreen  = 32 //绿色
	colorYellow = 33 //黄色
	colorBlue   = 34 //蓝色
	colorPurple = 35 //紫红色
	colorGBlue  = 36 //青蓝色
	colorGray   = 37 //灰色
)

func getColorByLevel(level logrus.Level) int {
	switch level {
	case logrus.DebugLevel:
		return colorBlack
	case logrus.InfoLevel:
		return colorGray
	case logrus.WarnLevel, logrus.ErrorLevel:
		return colorRed
	case logrus.FatalLevel, logrus.PanicLevel:
		return colorPurple
	default:
		return colorBlack
	}
}
