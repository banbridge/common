package logs_

import (
	"context"
	"os"
	"path"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/banbridge/common/pkg/logs"
)

type ZapLogger struct {
	zapLogger *zap.Logger
}

func NewZapLogger(opts ...Options) *ZapLogger {
	opt := defaultOptions()
	for _, o := range opts {
		o(&opt)
	}
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(opt.Level)); err != nil {
		zapLevel = zapcore.InfoLevel
	}

	coreList := make([]zapcore.Core, 0)

	zapOptions := []zap.Option{
		zap.AddCaller(),
		zap.AddCallerSkip(2),
		zap.AddStacktrace(zapcore.ErrorLevel),
	}

	// ======================输出界面============================
	consoleConfig := zap.NewDevelopmentEncoderConfig()
	consoleConfig.EncodeTime = timeEncoder
	encoder := zapcore.NewConsoleEncoder(consoleConfig)
	coreList = append(coreList, zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapLevel))

	// ======================输出文件============================
	// 将文件名设置为日期
	logFileName := time.Now().Format("2006-01-02_15") + ".log"
	fileName := path.Join(opt.LogDir, logFileName)
	// 提供压缩和删除
	lumberjackLogger := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    200,  // 一个文件最大可达 20M。
		MaxBackups: 50,   // 最多同时保存 5 个文件。
		MaxAge:     10,   // 一个文件最多可以保存 10 天。
		Compress:   true, // 用 gzip 压缩。
	}
	defer lumberjackLogger.Close()

	coreList = append(coreList, zapcore.NewCore(encoder, zapcore.AddSync(lumberjackLogger), zapLevel))
	core := zapcore.NewTee(coreList...)
	l := zap.New(
		core,
		zapOptions...,
	)

	defer l.Sync()

	logger := &ZapLogger{
		zapLogger: l.Named(opt.Name),
	}
	zap.RedirectStdLog(l)
	return logger
}

func (l *ZapLogger) Debugf(format string, v ...interface{}) {
	l.zapLogger.Sugar().Debugf(format, v...)
}

func (l *ZapLogger) Debug(msg string, fields ...Field) {
	l.zapLogger.Debug(msg, fields...)
}

func (l *ZapLogger) Warnf(format string, v ...interface{}) {
	l.zapLogger.Sugar().Warnf(format, v...)
}

func (l *ZapLogger) Warn(msg string, fields ...Field) {
	l.zapLogger.Warn(msg, fields...)
}

func (l *ZapLogger) Infof(format string, v ...interface{}) {
	l.zapLogger.Sugar().Infof(format, v...)
}

func (l *ZapLogger) Info(msg string, fields ...Field) {
	l.zapLogger.Info(msg, fields...)
}

func (l *ZapLogger) Errorf(format string, v ...interface{}) {
	l.zapLogger.Sugar().Errorf(format, v...)
}

func (l *ZapLogger) Error(msg string, fields ...Field) {
	l.zapLogger.Error(msg, fields...)
}

func (l *ZapLogger) Fatalf(format string, v ...interface{}) {
	l.zapLogger.Sugar().Fatalf(format, v...)
}

func (l *ZapLogger) Fatal(msg string, fields ...Field) {
	l.zapLogger.Fatal(msg, fields...)
}

func (l *ZapLogger) Panicf(format string, v ...interface{}) {
	l.zapLogger.Sugar().Panicf(format, v...)
}

func (l *ZapLogger) Panic(msg string, fields ...Field) {
	l.zapLogger.Panic(msg, fields...)
}

func (l *ZapLogger) CtxDebug(ctx context.Context, format string, v ...interface{}) {
	l.withCtx(ctx).zapLogger.Sugar().Debugf(format, v...)
}

func (l *ZapLogger) CtxWarn(ctx context.Context, format string, v ...interface{}) {
	l.withCtx(ctx).zapLogger.Sugar().Warnf(format, v...)
}

func (l *ZapLogger) CtxInfo(ctx context.Context, format string, v ...interface{}) {
	l.withCtx(ctx).zapLogger.Sugar().Infof(format, v...)
}

func (l *ZapLogger) CtxError(ctx context.Context, format string, v ...interface{}) {
	l.withCtx(ctx).zapLogger.Sugar().Errorf(format, v...)
}

func (l *ZapLogger) CtxFatal(ctx context.Context, format string, v ...interface{}) {
	l.withCtx(ctx).zapLogger.Sugar().Fatalf(format, v...)
}

func (l *ZapLogger) CtxPanic(ctx context.Context, format string, v ...interface{}) {
	l.withCtx(ctx).zapLogger.Sugar().Panicf(format, v...)
}

func (l *ZapLogger) withCtx(ctx context.Context) *ZapLogger {
	lg := l.clone()
	logId, _ := logs.CtxLogID(ctx)

	fields := make([]zap.Field, 0)
	if logId != "" {
		fields = append(fields, zap.String("LogId", logId))
	}

	lg.zapLogger = lg.zapLogger.With(fields...)
	return lg
}

func (l *ZapLogger) clone() *ZapLogger {
	copyL := *l
	return &copyL
}

func (l *ZapLogger) Flush() {
	_ = l.zapLogger.Sync()
}
