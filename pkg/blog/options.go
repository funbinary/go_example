package blog

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
)

const (
	flagLevel             = "log.level"
	flagDisableCaller     = "log.disable-caller"
	flagDisableStacktrace = "log.disable-stacktrace"
	flagFormat            = "log.format"
	flagEnableColor       = "log.enable-color"
	flagOutputPaths       = "log.output-paths"
	flagErrorOutputPaths  = "log.error-output-paths"
	flagDevelopment       = "log.development"
	flagName              = "log.name"

	consoleFormat    = "console"
	golbalTimeFormat = "2006-01-02"
	jsonFormat       = "json"
)

// Options contains configuration items related to log.
type Options struct {
	// ErrorOutputPaths  []string `json:"error-output-paths" mapstructure:"error-output-paths"`
	OutputPath        string `json:"output-path"       mapstructure:"output-path"`
	Level             string `json:"level"              mapstructure:"level"`
	Format            string `json:"format"             mapstructure:"format"`
	TimeFormat        string `json:"time-format"        mapstructure:"time-format"`
	DisableCaller     bool   `json:"disable-caller"     mapstructure:"disable-caller"`
	DisableStacktrace bool   `json:"disable-stacktrace" mapstructure:"disable-stacktrace"`
	EnableColor       bool   `json:"enable-color"       mapstructure:"enable-color"`
	Development       bool   `json:"development"        mapstructure:"development"`
	ModuleName        string `json:"module-name"        mapstructure:"module-name"`
	MaxSize           int    `json:"max-size"           mapstructure:"max-size"`
	Name              string `json:"name"               mapstructure:"name"`
	// OutputPaths       []string `json:"output-paths"       mapstructure:"output-paths"`
	// ErrorOutputPaths  []string `json:"error-output-paths" mapstructure:"error-output-paths"`
}

// 创建一个默认的Options
func NewOptions() *Options {
	return &Options{
		Level:             zapcore.InfoLevel.String(),
		DisableCaller:     false,
		DisableStacktrace: false,
		Format:            consoleFormat,
		TimeFormat:        golbalTimeFormat,
		EnableColor:       false,
		Development:       false,
		// OutputPaths:       []string{"stdout"},
		// ErrorOutputPaths:  []string{"stderr"},
	}
}

// Validate validate the options fields.
func (o *Options) Validate() []error {
	var errs []error

	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(o.Level)); err != nil {
		errs = append(errs, err)
	}

	format := strings.ToLower(o.Format)
	if format != consoleFormat && format != jsonFormat {
		errs = append(errs, fmt.Errorf("not a valid log format: %q", o.Format))
	}

	return errs
}

func (o *Options) String() string {
	data, _ := json.Marshal(o)

	return string(data)
}

// Build constructs a global zap logger from the Config and Options.
func (o *Options) Build() error {
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(o.Level)); err != nil {
		zapLevel = zapcore.InfoLevel
	}
	encodeLevel := zapcore.CapitalLevelEncoder
	if o.Format == consoleFormat && o.EnableColor {
		encodeLevel = zapcore.CapitalColorLevelEncoder
	}

	zc := &zap.Config{
		Level:             zap.NewAtomicLevelAt(zapLevel),
		Development:       o.Development,
		DisableCaller:     o.DisableCaller,
		DisableStacktrace: o.DisableStacktrace,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: o.Format,
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "message",
			LevelKey:       "level",
			TimeKey:        "timestamp",
			NameKey:        "logger",
			CallerKey:      "caller",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    encodeLevel,
			EncodeTime:     timeEncoder,
			EncodeDuration: milliSecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeName:     zapcore.FullNameEncoder,
		},
		// OutputPaths:      {o.OutputPath,},
		// ErrorOutputPaths: o.ErrorOutputPaths,
	}
	logger, err := zc.Build(zap.AddStacktrace(zapcore.PanicLevel))
	if err != nil {
		return err
	}
	zap.RedirectStdLog(logger.Named(o.Name))
	zap.ReplaceGlobals(logger)

	return nil
}
