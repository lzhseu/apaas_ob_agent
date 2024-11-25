package logger

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"strings"

	"github.com/fatih/color"
)

type PrettyHandlerOptions struct {
	SlogOpts slog.HandlerOptions
}

type PrettyHandler struct {
	slog.Handler
	l *log.Logger
}

func (h *PrettyHandler) Handle(ctx context.Context, r slog.Record) error {
	level := r.Level.String() + ":"

	switch r.Level {
	case slog.LevelDebug:
		level = color.MagentaString(level)
	case slog.LevelInfo:
		level = color.BlueString(level)
	case slog.LevelWarn:
		level = color.YellowString(level)
	case slog.LevelError:
		level = color.RedString(level)
	}

	sb := strings.Builder{}
	if r.NumAttrs() > 0 {
		sb.WriteString("{ ")
		r.Attrs(func(a slog.Attr) bool {
			sb.WriteString(fmt.Sprintf("%s:%v", a.Key, a.Value.Any()))
			return true
		})
		sb.WriteString(" }")
	}

	timeStr := r.Time.Format("[15:05:05.000]")
	msg := color.CyanString(r.Message)

	h.l.Println(timeStr, level, msg, color.WhiteString(sb.String()))

	return nil
}

func NewPrettyHandler(out io.Writer, opts *PrettyHandlerOptions) *PrettyHandler {
	color.NoColor = false
	h := &PrettyHandler{
		Handler: slog.NewJSONHandler(out, &opts.SlogOpts),
		l:       log.New(out, "", 0),
	}
	return h
}
