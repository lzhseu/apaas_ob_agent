package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
)

type Logger struct {
	slogger *slog.Logger
	cfg     *Config
}

type Config struct {
	out       io.Writer // 输出
	format    string    // 输出格式：text，json
	level     string    // 日志级别
	addSource bool      // 是否需要打印日志文件信息
}

type OptionFunc func(cfg *Config)

func NewLogger(out io.Writer, opts ...OptionFunc) *Logger {
	cfg := defaultConfig(out)
	for _, opt := range opts {
		opt(cfg)
	}

	var handler slog.Handler
	switch cfg.format {
	case FormatJSON:
		handler = slog.NewJSONHandler(cfg.out,
			&slog.HandlerOptions{
				AddSource: cfg.addSource,
				Level:     levelFromString[cfg.level],
			},
		)
	case FormatPretty:
		handler = NewPrettyHandler(cfg.out,
			&PrettyHandlerOptions{SlogOpts: slog.HandlerOptions{
				AddSource: cfg.addSource,
				Level:     levelFromString[cfg.level],
			}},
		)
	default:
		handler = slog.NewTextHandler(cfg.out,
			&slog.HandlerOptions{
				AddSource: cfg.addSource,
				Level:     levelFromString[cfg.level],
			},
		)
	}

	slogger := slog.New(handler)
	return &Logger{
		slogger: slogger,
		cfg:     cfg,
	}
}

func NewConsoleLogger(opts ...OptionFunc) *Logger {
	options := []OptionFunc{
		WithFormat(FormatPretty),
	}
	options = append(options, opts...)
	return NewLogger(os.Stdout, options...)
}

// NewFileLogger 创建文件日志记录器
// params：
//   - filename: 要写入日志的文件名
//   - dur: 日志文件分割的时间间隔
//   - size: 日志文件分割的大小阈值，当希望通过文件大小切割日志时，设置 [dur = NoDur]; if size is 0, it defaults to 100 megabytes
func NewFileLogger(filename string, dur SegDuration, size, maxAge int, opts ...OptionFunc) *Logger {
	if dur == NoDur {
		return NewLogger(NewProductionFileWriterSegBySize(filename, size, maxAge), opts...)
	}
	return NewLogger(NewProductionFileWriterSegByTime(filename, dur, maxAge), opts...)
}

// NewLokiLogger 创建 Loki 日志记录器，通过 HTTP 接口写入 Loki 服务，详见：https://grafana.com/docs/loki/latest/reference/loki-http-api/#ingest-logs
// params：
//   - baseUrl: loki 服务地址
//   - labels: loki 日志标签
func NewLokiLogger(baseUrl string, labels map[string]string, opts ...OptionFunc) *Logger {
	return NewLogger(&LokiWriter{
		BaseURL: baseUrl,
		Labels:  labels,
	}, opts...)
}

func (l *Logger) Debug(msg string, args ...any) {
	l.Log(LevelDebug, msg, args...)
}

func (l *Logger) Info(msg string, args ...any) {
	l.Log(LevelInfo, msg, args...)
}

func (l *Logger) Warn(msg string, args ...any) {
	l.Log(LevelWarn, msg, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	l.Log(LevelError, msg, args...)
}

func (l *Logger) Log(level string, msg string, args ...any) {
	l.log(context.Background(), level, msg, args...)
}

func (l *Logger) log(ctx context.Context, level string, msg string, args ...any) {
	l.slogger.Log(ctx, levelFromString[level], msg, args...)
}

func (l *Logger) GetFormat() string {
	return l.cfg.format
}

func defaultConfig(out io.Writer) *Config {
	return &Config{
		out:       out,
		format:    FormatText,
		level:     LevelInfo,
		addSource: false,
	}
}
