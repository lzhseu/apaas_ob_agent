package logger

import (
	"log/slog"

	"github.com/fatih/color"
)

const (
	FormatText   = "text"
	FormatJSON   = "json"
	FormatPretty = "pretty" // 美观输出，一般用在 console
)

const (
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
)

var (
	levelFromString = map[string]slog.Level{
		LevelDebug: slog.LevelDebug,
		LevelInfo:  slog.LevelInfo,
		LevelWarn:  slog.LevelWarn,
		LevelError: slog.LevelError,
	}

	levelToColor = map[string]*color.Color{
		LevelDebug: color.New(color.FgCyan),
		LevelInfo:  color.New(color.FgBlue),
		LevelWarn:  color.New(color.FgYellow),
		LevelError: color.New(color.FgRed),
	}
)
