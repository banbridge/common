package logs

import (
	"log/slog"
	"os"
	"testing"
)

// TestGenLogID 是一个测试函数，用于测试 GenLogID 函数的功能
// 该函数属于 testing 包的测试用例，通过 t 参数来控制测试流程和报告测试结果
func TestGenLogID(t *testing.T) {

	// 调用 GenLogID 函数生成一个日志ID，并将其存储在 logId 变量中
	logId := GenLogID()

	// 使用 slog 包输出信息日志，打印生成的 logId
	// 这里将字符串 "log id is " 与 logId 拼接后输出
	slog.Info("log id is " + logId)

}

func TestJson(t *testing.T) {
	//logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	logger.Debug("Debug message")
	logger.Info("Info message")
	logger.Warn("Warning message")
	logger.Error("Error message")
}
