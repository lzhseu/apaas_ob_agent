package logger

import (
	"io"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"gopkg.in/natefinch/lumberjack.v2"
)

type SegDuration string

const (
	HourDur  SegDuration = "hour"
	DayDur   SegDuration = "day"
	WeekDur  SegDuration = "week"
	MonthDur SegDuration = "month"
	NoDur    SegDuration = "no"
)

// FileSegConfig 日志文件切割配置
type FileSegConfig struct {
	// 公共配置
	Filename  string // 日志文件名
	MaxAge    int    // 日志文件最大保留天数，默认保留30天
	LocalTime bool   // 是否使用本地时间，默认为true

	// 按时间切割
	RotationTime time.Duration // 按多长时间切割一次

	// 按文件大小切割
	MaxSize    int  // 日志文件最大大小(MB)，默认100M
	MaxBackups int  // 保留日志文件的最大数量，会受到 MaxAge 的影响
	Compress   bool // 是否对日志文件进行压缩归档
}

func NewProductionFileWriterSegByTime(filename string, dur SegDuration, maxAge int) io.Writer {
	var rotationTime time.Duration
	switch dur {
	case HourDur:
		rotationTime = time.Hour
	case DayDur:
		rotationTime = time.Hour * 24
	case WeekDur:
		rotationTime = time.Hour * 24 * 7
	case MonthDur:
		rotationTime = time.Hour * 24 * 30
	default:
		rotationTime = time.Hour * 24
	}
	return NewFileWriterSegByTime(DefaultFileSegConfig(filename, &rotationTime, nil, &maxAge))
}

func NewProductionFileWriterSegBySize(filename string, size, maxAge int) io.Writer {
	return NewFileWriterSegBySize(DefaultFileSegConfig(filename, nil, &size, &maxAge))
}

func DefaultFileSegConfig(filename string, dur *time.Duration, size *int, maxAge *int) *FileSegConfig {
	var (
		rotationTime = time.Hour * 24
		maxSize      = 100
		age          = 30
	)

	if dur != nil {
		rotationTime = *dur
	}
	if size != nil {
		maxSize = *size
	}
	if maxAge != nil {
		age = *maxAge
	}

	return &FileSegConfig{
		Filename:     filename,
		MaxAge:       age,
		LocalTime:    true,
		RotationTime: rotationTime,
		MaxSize:      maxSize,
		MaxBackups:   100,
		Compress:     false,
	}
}

func NewFileWriterSegByTime(cfg *FileSegConfig) io.Writer {
	opts := []rotatelogs.Option{
		rotatelogs.WithMaxAge(time.Duration(cfg.MaxAge) * time.Hour * 24),
		rotatelogs.WithRotationTime(cfg.RotationTime),
		rotatelogs.WithLinkName(cfg.Filename),
	}
	if !cfg.LocalTime {
		opts = append(opts, rotatelogs.WithClock(rotatelogs.UTC))
	}
	l, _ := rotatelogs.New(
		cfg.Filename+".%Y-%m-%d_%H_%M.log",
		opts...,
	)
	return l
}

func NewFileWriterSegBySize(cfg *FileSegConfig) io.Writer {
	return &lumberjack.Logger{
		Filename:   cfg.Filename,
		MaxSize:    cfg.MaxSize,
		MaxAge:     cfg.MaxAge,
		MaxBackups: cfg.MaxBackups,
		LocalTime:  cfg.LocalTime,
		Compress:   cfg.Compress,
	}
}
