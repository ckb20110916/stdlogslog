package stdlogslog

import (
	"context"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/ckb20110916/lumberjacklogwriter"
	"github.com/ckb20110916/rotatelogswriter"
	"github.com/lmittmann/tint"
	"github.com/mattn/go-colorable"
)

const (
	LevelTrace = slog.LevelDebug - 4
	LevelDebug = slog.LevelDebug
	LevelInfo  = slog.LevelInfo
	LevelWarn  = slog.LevelWarn
	LevelError = slog.LevelError
	LevelFatal = slog.LevelError + 4
	LevelPanic = slog.LevelError + 8
)

const (
	defaultTimeFormat = "2006-01-02 15:04:05.000"
)

func replaceAttr(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.LevelKey {
		if level, ok := a.Value.Any().(slog.Level); ok {
			switch {
			case level == LevelTrace:
				return slog.String(slog.LevelKey, "TRC")
			case level == LevelFatal:
				return slog.String(slog.LevelKey, "FLT")
			case level == LevelPanic:
				return slog.String(slog.LevelKey, "PNC")
			}
		}
	}
	return a
}

func newHandler(f *os.File, logLevel slog.Level, logSource bool, logColour bool) slog.Handler {
	return tint.NewHandler(colorable.NewColorable(f), &tint.Options{
		NoColor:     !logColour,
		Level:       logLevel,
		AddSource:   logSource,
		ReplaceAttr: replaceAttr,
		TimeFormat:  defaultTimeFormat,
	})
}

var (
	Logger = slog.New(newHandler(os.Stdout, LevelDebug, false, true))
)

func EnableConsole(logLevel slog.Level, logSource bool, logColour bool) {
	Logger = slog.New(newHandler(os.Stdout, logLevel, logSource, logColour))
}

func EnableLogFile(logLevel slog.Level, logSource bool, folder, filename string, maxAge, rotateTime time.Duration) {
	logWriter := rotatelogswriter.New(folder, filename, maxAge, rotateTime)
	if logWriter == nil {
		EnableConsole(logLevel, logSource, false)
	} else {
		enableOutfile(logLevel, logSource, logWriter)
	}
}

func EnableLogFile2(logLevel slog.Level, logSource bool, folder, filename string, maxBackups, maxSize, maxAge int) {
	logWriter := lumberjacklogwriter.New(folder, filename, maxBackups, maxSize, maxAge)
	if logWriter == nil {
		EnableConsole(logLevel, logSource, false)
	} else {
		enableOutfile(logLevel, logSource, logWriter)
	}
}

func enableOutfile(logLevel slog.Level, logSource bool, logWriter io.Writer) {
	Logger = slog.New(
		tint.NewHandler(logWriter, &tint.Options{
			NoColor:     true,
			Level:       logLevel,
			AddSource:   logSource,
			ReplaceAttr: replaceAttr,
			TimeFormat:  defaultTimeFormat,
		}),
	)
}

func Trace(msg string, args ...any) {
	Logger.Log(context.Background(), LevelTrace, msg, args...)
}

func TraceContext(ctx context.Context, msg string, args ...any) {
	Logger.Log(ctx, LevelTrace, msg, args...)
}

func Debug(msg string, args ...any) {
	Logger.Debug(msg, args...)
}

func DebugContext(ctx context.Context, msg string, args ...any) {
	Logger.DebugContext(ctx, msg, args...)
}

func Info(msg string, args ...any) {
	Logger.Info(msg, args...)
}

func InfoContext(ctx context.Context, msg string, args ...any) {
	Logger.InfoContext(ctx, msg, args...)
}

func Warn(msg string, args ...any) {
	Logger.Warn(msg, args...)
}

func WarnContext(ctx context.Context, msg string, args ...any) {
	Logger.WarnContext(ctx, msg, args...)
}

func Error(msg string, args ...any) {
	Logger.Error(msg, args...)
}

func ErrorContext(ctx context.Context, msg string, args ...any) {
	Logger.ErrorContext(ctx, msg, args...)
}

func Fatal(msg string, args ...any) {
	Logger.Log(context.Background(), LevelFatal, msg, args...)
	os.Exit(1)
}

func FatalContext(ctx context.Context, msg string, args ...any) {
	Logger.Log(ctx, LevelFatal, msg, args...)
	os.Exit(1)
}

func Panic(msg string, args ...any) {
	Logger.Log(context.Background(), LevelPanic, msg, args...)
	panic(msg)
}

func PanicContext(ctx context.Context, msg string, args ...any) {
	Logger.Log(ctx, LevelPanic, msg, args...)
	panic(msg)
}
