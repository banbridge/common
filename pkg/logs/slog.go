package logs

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"time"

	"log/slog"

	"github.com/banbridge/common/pkg/logs/handler"
)

// StdLog 自定义日志记录器
type StdLog struct {
	inner *slog.Logger
	op    SlogOption
}

var _ Logger = &StdLog{}

func NewLogger() *StdLog {

	// 控制台输出（文本格式）
	consoleHandler := handler.NewConsoleHandler(os.Stdout, nil)

	// 文件输出（JSON格式）
	file := mustOpenLogFile("app.log")
	fileHandler := slog.NewJSONHandler(file, &slog.HandlerOptions{
		AddSource: true,
	})

	return &StdLog{
		inner: slog.New(handler.MultiHandler{consoleHandler, fileHandler}),
		op: SlogOption{
			CallDepth: 1,
		},
	}
}

func mustOpenLogFile(path string) io.Writer {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	return file
}

func (l *StdLog) Info(msg string, args ...any) {
	l.log(context.Background(), slog.LevelInfo, msg, args...)
}

func (l *StdLog) Warn(msg string, args ...any) {
	l.log(context.Background(), slog.LevelWarn, msg, args...)
}

func (l *StdLog) Error(msg string, args ...any) {
	l.log(context.Background(), slog.LevelError, msg, args...)
}

func (l *StdLog) Debug(msg string, args ...any) {
	l.log(context.Background(), slog.LevelDebug, msg, args...)
}

func (l *StdLog) CtxDebug(ctx context.Context, msg string, args ...any) {
	l.log(ctx, slog.LevelDebug, msg, args...)
}

func (l *StdLog) CtxWarn(ctx context.Context, msg string, args ...any) {
	l.log(ctx, slog.LevelWarn, msg, args...)
}

func (l *StdLog) CtxInfo(ctx context.Context, msg string, args ...any) {
	l.log(ctx, slog.LevelInfo, msg, args...)
}

func (l *StdLog) CtxError(ctx context.Context, msg string, args ...any) {
	l.log(ctx, slog.LevelError, msg, args...)
}

func (l *StdLog) log(ctx context.Context, level slog.Level, msg string, args ...any) {
	var pcs [1]uintptr
	runtime.Callers(3+l.op.CallDepth, pcs[:]) // 3 = log + Info/Warn/Error + runtime.Callers

	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}

	r := slog.NewRecord(time.Now(), level, msg, pcs[0])

	var attrs []slog.Attr

	logID := LogIDDefault(ctx)
	if logID != "-" {
		attrs = append(attrs, slog.String("log_id", logID))
	}

	kvs := GetAllKVs(ctx)

	for i := 0; i < len(kvs); i += 2 {
		str := kvs[i].(string)
		attrs = append(attrs, slog.Any(str, kvs[i+1]))
	}

	r.AddAttrs(attrs...)
	//r.Add(args...)
	_ = l.inner.Handler().Handle(context.Background(), r)
}
