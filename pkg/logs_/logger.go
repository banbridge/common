package logs_

import (
	"context"
	"log"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Debugf(msg string, args ...interface{})
	Debug(msg string, fields ...Field)
	Warnf(msg string, args ...interface{})
	Warn(msg string, fields ...Field)
	Infof(msg string, args ...interface{})
	Info(msg string, fields ...Field)
	Errorf(msg string, args ...interface{})
	Error(msg string, fields ...Field)
	Fatalf(msg string, args ...interface{})
	Fatal(msg string, fields ...Field)
	Panicf(msg string, args ...interface{})
	Panic(msg string, fields ...Field)

	CtxDebug(ctx context.Context, msg string, args ...interface{})
	CtxWarn(ctx context.Context, msg string, args ...interface{})
	CtxInfo(ctx context.Context, msg string, args ...interface{})
	CtxError(ctx context.Context, msg string, args ...interface{})
	CtxFatal(ctx context.Context, msg string, args ...interface{})
	CtxPanic(ctx context.Context, msg string, args ...interface{})

	Flush()
}

var _ Logger = &ZapLogger{}

var (
	std = NewZapLogger()
	mu  sync.Mutex
)

func Init(opts ...Options) {
	mu.Lock()
	defer mu.Unlock()
	std = NewZapLogger(opts...)
}

func StdLogger() *ZapLogger {
	return std
}

func StdErrLogger() *log.Logger {
	if std == nil {
		return nil
	}
	if l, err := zap.NewStdLogAt(std.zapLogger, zapcore.ErrorLevel); err == nil {
		return l
	}
	return nil
}

func SugaredLogger() *zap.SugaredLogger {
	return std.zapLogger.Sugar()
}

func Debugf(msg string, args ...interface{}) {
	std.Debugf(msg, args...)
}

func Debug(msg string, fields ...Field) {
	std.Debug(msg, fields...)
}

func CtxDebug(ctx context.Context, msg string, args ...interface{}) {
	std.CtxDebug(ctx, msg, args...)
}

func Warnf(msg string, args ...interface{}) {
	std.Warnf(msg, args...)
}

func Warn(msg string, fields ...Field) {
	std.Warn(msg, fields...)
}

func CtxWarn(ctx context.Context, msg string, args ...interface{}) {
	std.CtxWarn(ctx, msg, args...)
}

func Infof(msg string, args ...interface{}) {
	std.Infof(msg, args...)
}

func Info(msg string, fields ...Field) {
	std.Info(msg, fields...)
}

func CtxInfo(ctx context.Context, msg string, args ...interface{}) {
	std.CtxInfo(ctx, msg, args...)
}

func Errorf(msg string, args ...interface{}) {
	std.Errorf(msg, args...)
}

func Error(msg string, fields ...Field) {
	std.Error(msg, fields...)
}

func CtxError(ctx context.Context, msg string, args ...interface{}) {
	std.CtxError(ctx, msg, args...)
}

func Fatalf(msg string, args ...interface{}) {
	std.Fatalf(msg, args...)
}

func Fatal(msg string, fields ...Field) {
	std.Fatal(msg, fields...)
}

func CtxFatal(ctx context.Context, msg string, args ...interface{}) {
	std.CtxFatal(ctx, msg, args...)
}

func Panicf(msg string, args ...interface{}) {
	std.Panicf(msg, args...)
}

func Panic(msg string, fields ...Field) {
	std.Panic(msg, fields...)
}

func CtxPanic(ctx context.Context, msg string, args ...interface{}) {
	std.CtxPanic(ctx, msg, args...)
}

func Flush() { std.Flush() }
