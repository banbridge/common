package logs

import (
	"context"
	"os"
	"testing"

	"log/slog"
)

func TestGenLogID(t *testing.T) {

	logId := GenLogID()

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

func TestStdLogger(t *testing.T) {

	ctx := context.Background()
	ctx = context.WithValue(ctx, CTXKeyLogID, GenLogID())
	ctx = CtxAddKVs(ctx, "key-add", "value")
	logID := GenLogID()
	CtxInfo(ctx, "test id %s", logID)
	logID = GenLogID()
	CtxInfo(ctx, "test id %s", logID)
}
