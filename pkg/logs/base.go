package logs

import "context"

type Logger interface {
	Debug(msg string, args ...any)
	Warn(msg string, args ...any)
	Info(msg string, args ...any)
	Error(msg string, args ...any)
	//Fatal(msg string, args ...any)
	//Panic(msg string, args ...any)

	CtxDebug(ctx context.Context, msg string, args ...any)
	CtxWarn(ctx context.Context, msg string, args ...any)
	CtxInfo(ctx context.Context, msg string, args ...any)
	CtxError(ctx context.Context, msg string, args ...any)
	//CtxFatal(ctx context.Context, msg string, args ...any)
	//CtxPanic(ctx context.Context, msg string, args ...any)
}

var (
	stdLog        = NewLogger()
	DefaultLogger = NewLogger(WithCallDepth(1))
)

func Info(msg string, args ...any) {
	stdLog.Info(msg, args...)
}

func Debug(msg string, args ...any) {
	stdLog.Debug(msg, args...)
}

func Warn(msg string, args ...any) {
	stdLog.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	stdLog.Error(msg, args...)
}

func CtxInfo(ctx context.Context, msg string, args ...any) {
	stdLog.CtxInfo(ctx, msg, args...)
}

func CtxDebug(ctx context.Context, msg string, args ...any) {
	stdLog.CtxDebug(ctx, msg, args...)
}

func CtxWarn(ctx context.Context, msg string, args ...any) {
	stdLog.CtxWarn(ctx, msg, args...)
}

func CtxError(ctx context.Context, msg string, args ...any) {
	stdLog.CtxError(ctx, msg, args...)
}
