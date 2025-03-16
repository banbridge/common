package handler

import (
	"bytes"
	"context"
	"io"
	"sync"

	"github.com/fatih/color"
	"log/slog"

	"github.com/banbridge/common/pkg/consts"
)

type ConsoleHandler struct {
	opts ConsoleOptions
	// TODO: state for WithGroup and WithAttrs
	mu  *sync.Mutex
	out io.Writer
}

type ConsoleOptions struct {
	// Level reports the minimum level to log.
	// Levels with lower levels are discarded.
	// If nil, the Handler uses [slog.LevelInfo].
	Level slog.Leveler
}

func NewConsoleHandler(w io.Writer, opts *ConsoleOptions) *ConsoleHandler {
	h := &ConsoleHandler{
		mu:  &sync.Mutex{},
		out: w,
	}
	if opts != nil {
		h.opts = *opts
	}
	if h.opts.Level == nil {
		h.opts.Level = slog.LevelInfo
	}
	return h
}

func (c *ConsoleHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= c.opts.Level.Level()
}

var colorFuncMap = map[slog.Level]func(format string, a ...any) string{
	slog.LevelDebug: color.MagentaString,
	slog.LevelInfo:  color.BlueString,
	slog.LevelWarn:  color.YellowString,
	slog.LevelError: color.RedString,
}

func (c *ConsoleHandler) Handle(ctx context.Context, r slog.Record) error {
	buf := new(bytes.Buffer)

	timeStr := r.Time.Format(consts.TimeWithMS)
	buf.WriteString(timeStr)
	buf.WriteString(" ")

	level := r.Level.String()
	if fn, ok := colorFuncMap[r.Level]; ok {
		level = fn(level)
	}

	if r.PC != 0 {
		buf.WriteString(getCallInfo(r.PC))
		buf.WriteString(" ")
	}

	buf.WriteString("[" + level + "]")
	buf.WriteString(" ")

	msg := color.CyanString(r.Message)
	buf.WriteString(slog.String(slog.MessageKey, msg).String())
	buf.WriteString(" ")

	r.Attrs(func(attr slog.Attr) bool {
		buf.WriteString(attr.String() + " ")
		return true
	})

	buf.WriteString("\n")

	c.mu.Lock()
	defer c.mu.Unlock()
	_, err := c.out.Write(buf.Bytes())
	return err
}

func (c *ConsoleHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	//TODO implement me
	panic("implement me")
}

func (c *ConsoleHandler) WithGroup(name string) slog.Handler {
	//TODO implement me
	panic("implement me")
}

var _ slog.Handler = &ConsoleHandler{}
