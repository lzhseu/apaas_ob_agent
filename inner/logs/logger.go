package logs

import (
	"fmt"
	"runtime/debug"

	"github.com/lzhseu/apaas_ob_agent/conf"
	"github.com/lzhseu/apaas_ob_agent/inner/metrics"
	"github.com/lzhseu/apaas_ob_agent/pkg/logger"
)

var (
	innerLogger *InnerLogger // 内部日志器，用于 agent 本身的日志记录
)

func MustInit() {
	innerLogger = NewInnerLogger()

	cfg := conf.GetConfig().InnerLogsCfg

	opts := make([]logger.OptionFunc, 0)
	if cfg.LogLevel != nil {
		opts = append(opts, logger.WithLevel(*cfg.LogLevel))
	}

	if cfg.Console != nil && cfg.Console.Enable {
		innerLogger.Add(logger.NewConsoleLogger(opts...))
	}

	if cfg.File != nil && cfg.File.Enable {
		size := 0
		if cfg.File.SegMaxSize != nil {
			size = *cfg.File.SegMaxSize
		}
		dur := logger.NoDur
		if cfg.File.SegDur != nil {
			dur = logger.SegDuration(*cfg.File.SegDur)
		}
		maxAge := 30
		if cfg.File.MaxAge != nil {
			maxAge = *cfg.File.MaxAge
		}
		innerLogger.Add(logger.NewFileLogger(cfg.File.FileName, dur, size, maxAge, opts...))
	}

	if cfg.Loki != nil && cfg.Loki.Enable {
		innerLogger.Add(logger.NewLokiLogger(fmt.Sprintf("%v://%v:%v", cfg.Loki.Schema, cfg.Loki.Host, cfg.Loki.Port), cfg.Loki.Labels, opts...))

	}

	// 兜底用 console
	if len(innerLogger.loggers) == 0 {
		innerLogger.Add(logger.NewConsoleLogger())
	}
}

type InnerLogger struct {
	loggers []*logger.Logger
}

func NewInnerLogger(logger ...*logger.Logger) *InnerLogger {
	return &InnerLogger{
		loggers: logger,
	}
}

func (l *InnerLogger) Add(logger *logger.Logger) *InnerLogger {
	l.loggers = append(l.loggers, logger)
	return l
}

func (l *InnerLogger) Debug(msg string, args ...any) {
	l.Log(logger.LevelDebug, msg, args...)
}

func (l *InnerLogger) Info(msg string, args ...any) {
	l.Log(logger.LevelInfo, msg, args...)
}

func (l *InnerLogger) Warn(msg string, args ...any) {
	l.Log(logger.LevelWarn, msg, args...)
}

func (l *InnerLogger) Error(msg string, args ...any) {
	l.Log(logger.LevelError, msg, args...)
}

func (l *InnerLogger) Log(level string, msg string, args ...any) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("🍀panic recover: %+v\n%s\n", r, debug.Stack())
				metrics.PanicTotal.WithLabelValues("inner_log").Add(1)
			}
		}()
		for _, logger := range l.loggers {
			logger.Log(level, msg, args...)
		}
	}()
}
