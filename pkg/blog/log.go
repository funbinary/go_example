package blog

import (
	"go.uber.org/zap/zapcore"
	"sync"

	"go.uber.org/zap"
)

var (
	std = New(NewOptions())
	mu  sync.Mutex
)

type zapLogger struct {
	zapLogger *zap.Logger
}

//初始化全局log对象
func Init(opts *Options) {
	mu.Lock()
	defer mu.Unlock()
	WrapOption(opts)
}

func New(opts *Options) *zapLogger {
	if opts == nil {
		opts = NewOptions()
	}

	// 调试级别
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(opts.Level)); err != nil {
		zapLevel = zapcore.InfoLevel
	}

	//是否开启颜色输出
	encodeLevel := zapcore.CapitalLevelEncoder
	if opts.Format == consoleFormat && opts.EnableColor {
		encodeLevel = zapcore.CapitalColorLevelEncoder
	}

	encoderConfig := zapcore.EncoderConfig{
		MessageKey: "message", //日志内容的key
		// FunctionKey:    "func",//日志函数key
		NameKey:        "logger",     //日志模块名称的key
		LevelKey:       "level",      //日志级别的Key
		TimeKey:        "timestamp",  //日志时间的key
		CallerKey:      "caller",     //调用栈的key eg."caller":"example/example.go:81"
		StacktraceKey:  "stacktrace", //
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    encodeLevel, //打印的编码
		EncodeTime:     timeEncoder, //时间的编码格式
		EncodeDuration: milliSecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	loggerConfig := &zap.Config{
		Level:             zap.NewAtomicLevelAt(zapLevel), //调试级别
		Development:       opts.Development,               //设置开发者模式，更容易接受堆栈追踪
		DisableCaller:     opts.DisableCaller,             // 禁用行号、函数调用信息。
		DisableStacktrace: opts.DisableStacktrace,         // 禁用自动堆栈跟踪捕获。
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		}, //采样配置,单位是没秒钟，限制日志在每秒钟的输出数量，防止CPU和IO被过度占用
		Encoding:         opts.Format,        //指定日志编码器, 目前仅支持两种编码器：console和json，默认为json。
		EncoderConfig:    encoderConfig,      //编码器配置
		OutputPaths:      []string{"stdout"}, //标准日志输出
		ErrorOutputPaths: []string{"stderr"}, //错误输出
	}

	l, _ := loggerConfig.Build(zap.AddCaller(), zap.AddStacktrace(zapcore.PanicLevel), zap.AddCallerSkip(1))

	logger := &zapLogger{
		zapLogger: l.Named(opts.Name),
	}

	return logger
}

func WrapOption(opts *Options) {
	if opts == nil {
		opts = NewOptions()
	}
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(opts.Level)); err != nil {
		zapLevel = zapcore.InfoLevel
	}

	//是否开启颜色输出
	encodeLevel := zapcore.CapitalLevelEncoder
	if opts.Format == consoleFormat && opts.EnableColor {
		encodeLevel = zapcore.CapitalColorLevelEncoder
	}

	encoderConfig := zapcore.EncoderConfig{
		MessageKey: "message", //日志内容的key
		// FunctionKey:    "func",//日志函数key
		NameKey:        "logger",     //日志模块名称的key
		LevelKey:       "level",      //日志级别的Key
		TimeKey:        "timestamp",  //日志时间的key
		CallerKey:      "caller",     //调用栈的key eg."caller":"example/example.go:81"
		StacktraceKey:  "stacktrace", //
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    encodeLevel, //打印的编码
		EncodeTime:     timeEncoder, //时间的编码格式
		EncodeDuration: milliSecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	w := zapcore.AddSync(NewWriter(opts.OutputPath, opts.TimeFormat, opts.ModuleName, opts.MaxSize))
	std.zapLogger = std.zapLogger.WithOptions(zap.WrapCore(func(zapcore.Core) zapcore.Core {
		return zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			w,
			zap.NewAtomicLevelAt(zapLevel))
	}))
}

func StdLogger() *zapLogger {
	return std
}

func SetStdLog(l *zapLogger) {
	mu.Lock()
	defer mu.Unlock()
	std = l
}

func Flush() { std.Flush() }

func (l *zapLogger) Flush() {
	_ = l.zapLogger.Sync()
}

func Debug(msg string, fields ...Field) {
	std.zapLogger.Debug(msg, fields...)
}

func Debugf(format string, v ...interface{}) {
	std.zapLogger.Sugar().Debugf(format, v...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	std.zapLogger.Sugar().Debugw(msg, keysAndValues...)
}

func Info(msg string, fields ...Field) {
	std.zapLogger.Info(msg, fields...)
}

func Infof(format string, v ...interface{}) {
	std.zapLogger.Sugar().Infof(format, v...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	std.zapLogger.Sugar().Infow(msg, keysAndValues...)
}

func Warn(msg string, fields ...Field) {
	std.zapLogger.Warn(msg, fields...)
}

func Warnf(format string, v ...interface{}) {
	std.zapLogger.Sugar().Warnf(format, v...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	std.zapLogger.Sugar().Warnw(msg, keysAndValues...)
}

func Error(msg string, fields ...Field) {
	std.zapLogger.Error(msg, fields...)
}

func Errorf(format string, v ...interface{}) {
	std.zapLogger.Sugar().Errorf(format, v...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	std.zapLogger.Sugar().Errorw(msg, keysAndValues...)
}

func Panic(msg string, fields ...Field) {
	std.zapLogger.Panic(msg, fields...)
}

func Panicf(format string, v ...interface{}) {
	std.zapLogger.Sugar().Panicf(format, v...)
}

func Panicw(msg string, keysAndValues ...interface{}) {
	std.zapLogger.Sugar().Panicw(msg, keysAndValues...)
}

func Fatal(msg string, fields ...Field) {
	std.zapLogger.Fatal(msg, fields...)
}

func Fatalf(format string, v ...interface{}) {
	std.zapLogger.Sugar().Fatalf(format, v...)
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	std.zapLogger.Sugar().Fatalw(msg, keysAndValues...)
}

func (l *zapLogger) Debug(msg string, fields ...Field) {
	l.zapLogger.Debug(msg, fields...)
}

func (l *zapLogger) Debugf(format string, v ...interface{}) {
	l.zapLogger.Sugar().Debugf(format, v...)
}

func (l *zapLogger) Debugw(msg string, keysAndValues ...interface{}) {
	l.zapLogger.Sugar().Debugw(msg, keysAndValues...)
}

func (l *zapLogger) Info(msg string, fields ...Field) {
	l.zapLogger.Info(msg, fields...)
}

func (l *zapLogger) Infof(format string, v ...interface{}) {
	l.zapLogger.Sugar().Infof(format, v...)
}

func (l *zapLogger) Infow(msg string, keysAndValues ...interface{}) {
	l.zapLogger.Sugar().Infow(msg, keysAndValues...)
}

func (l *zapLogger) Warn(msg string, fields ...Field) {
	l.zapLogger.Warn(msg, fields...)
}

func (l *zapLogger) Warnf(format string, v ...interface{}) {
	l.zapLogger.Sugar().Warnf(format, v...)
}

func (l *zapLogger) Warnw(msg string, keysAndValues ...interface{}) {
	l.zapLogger.Sugar().Warnw(msg, keysAndValues...)
}

func (l *zapLogger) Error(msg string, fields ...Field) {
	l.zapLogger.Error(msg, fields...)
}

func (l *zapLogger) Errorf(format string, v ...interface{}) {
	l.zapLogger.Sugar().Errorf(format, v...)
}

func (l *zapLogger) Errorw(msg string, keysAndValues ...interface{}) {
	l.zapLogger.Sugar().Errorw(msg, keysAndValues...)
}

func (l *zapLogger) Panic(msg string, fields ...Field) {
	l.zapLogger.Panic(msg, fields...)
}

func (l *zapLogger) Panicf(format string, v ...interface{}) {
	l.zapLogger.Sugar().Panicf(format, v...)
}

func (l *zapLogger) Panicw(msg string, keysAndValues ...interface{}) {
	l.zapLogger.Sugar().Panicw(msg, keysAndValues...)
}

func (l *zapLogger) Fatal(msg string, fields ...Field) {
	l.zapLogger.Fatal(msg, fields...)
}

func (l *zapLogger) Fatalf(format string, v ...interface{}) {
	l.zapLogger.Sugar().Fatalf(format, v...)
}

func (l *zapLogger) Fatalw(msg string, keysAndValues ...interface{}) {
	l.zapLogger.Sugar().Fatalw(msg, keysAndValues...)
}

// 添加一个固定字段到std logger中
func WithValues(keysAndValues ...interface{}) { std.WithValues(keysAndValues...) }

// 添加一个固定字段到logger中
func (l *zapLogger) WithValues(keysAndValues ...interface{}) {
	newLogger := l.zapLogger.With(handleFields(l.zapLogger, keysAndValues)...)
	l.zapLogger = newLogger
}

func handleFields(l *zap.Logger, args []interface{}, additional ...zap.Field) []zap.Field {
	if len(args) == 0 {
		return additional
	}

	fields := make([]zap.Field, 0, len(args)/2+len(additional))
	for i := 0; i < len(args); {
		if _, ok := args[i].(zap.Field); ok {
			l.DPanic("strongly-typed Zap Field passed to logr", zap.Any("zap field", args[i]))
			break
		}

		if i == len(args)-1 {
			l.DPanic("odd number of arguments passed as key-value pairs for logging", zap.Any("ignored key", args[i]))
			break
		}

		//确保key是个字符串
		key, val := args[i], args[i+1]
		keyStr, isString := key.(string)
		if !isString {
			l.DPanic(
				"non-string key argument passed to logging, ignoring all later arguments",
				zap.Any("invalid key", key),
			)
			break
		}

		fields = append(fields, zap.Any(keyStr, val))
		i += 2
	}

	return append(fields, additional...)
}

type Logger interface {
	Debug(msg string, fields ...Field)
	Debugf(format string, v ...interface{})
	Debugw(msg string, keysAndValues ...interface{})

	Warn(msg string, fields ...Field)
	Warnf(format string, v ...interface{})
	Warnw(msg string, keysAndValues ...interface{})

	Error(msg string, fields ...Field)
	Errorf(format string, v ...interface{})
	Errorw(msg string, keysAndValues ...interface{})

	Panic(msg string, fields ...Field)
	Panicf(format string, v ...interface{})
	Panicw(msg string, keysAndValues ...interface{})

	Fatal(msg string, fields ...Field)
	Fatalf(format string, v ...interface{})
	Fatalw(msg string, keysAndValues ...interface{})

	// WithName adds a new element to the logger's name.
	// Successive calls with WithName continue to append
	// suffixes to the logger's name.  It's strongly recommended
	// that name segments contain only letters, digits, and hyphens
	// (see the package documentation for more information).
	WithName(name string) Logger

	// Flush calls the underlying Core's Sync method, flushing any buffered
	// log entries. Applications should take care to call Sync before exiting.
	Flush()
}

//给模块命名
func WithName(s string) Logger {
	newLogger := std.zapLogger.Named(s)
	std.zapLogger = newLogger
	newLogger.Debug("with")
	return &zapLogger{
		zapLogger: newLogger,
	}
}

func (l *zapLogger) WithName(name string) Logger {
	obj := l.zapLogger.Named(name)
	obj.Debug("with")
	l.zapLogger = obj
	return &zapLogger{
		zapLogger: obj,
	}
}
