package logs_

import (
	"time"

	"go.uber.org/zap/zapcore"
)

const (
	consoleFormat = "console"
	jsonFormat    = "json"
)

type Options func(o *options)

type options struct {
	OutputPaths        []string
	ErrorOutputPaths   []string
	Name               string
	Level              string
	Format             string
	EnableColor        bool
	MaxSizeInMegabytes int
	MaxBackups         int
	MaxAgeInDays       int
	LogDir             string
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func milliSecondsDurationEncoder(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendFloat64(float64(d) / float64(time.Millisecond))
}

func WithOutputPaths(paths []string) Options {
	return func(o *options) { o.OutputPaths = paths }
}

func WithErrorOutputPaths(paths []string) Options {
	return func(o *options) { o.ErrorOutputPaths = paths }
}

func WithName(name string) Options {
	return func(o *options) { o.Name = name }
}

func WithLevel(level string) Options {
	return func(o *options) { o.Level = level }
}

func WithFormat(format string) Options {
	return func(o *options) { o.Format = format }
}

func EnableColor(enable bool) Options {
	return func(o *options) { o.EnableColor = enable }
}

func WithSizeInMegaBytes(size int) Options {
	return func(o *options) { o.MaxSizeInMegabytes = size }
}

func WithBackups(size int) Options {
	return func(o *options) { o.MaxBackups = size }
}

func WithAgeInDays(size int) Options {
	return func(o *options) { o.MaxAgeInDays = size }
}

func WithLogDir(dir string) Options {
	return func(o *options) { o.LogDir = dir }
}

func defaultOptions() options {
	return options{
		Level:              zapcore.InfoLevel.String(),
		Format:             jsonFormat,
		OutputPaths:        []string{"stdout"},
		ErrorOutputPaths:   []string{"stderr"},
		EnableColor:        false,
		MaxSizeInMegabytes: 500,
		MaxBackups:         3,
		MaxAgeInDays:       28,
		LogDir:             "log",
	}
}
