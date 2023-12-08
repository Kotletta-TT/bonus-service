package logger

import (
	"context"
	"log/slog"
	"os"
	"sync"

	"github.com/Kotletta-TT/bonus-service/config"
)

var logger *slog.Logger
var once sync.Once

func Init(config *config.Config) {
	once.Do(func() {
		lvl := &slog.LevelVar{}
		switch config.LogLevel {
		case "INFO":
			lvl.Set(slog.LevelInfo)
		case "DEBUG":
			lvl.Set(slog.LevelDebug)
		case "ERROR":
			lvl.Set(slog.LevelError)
		default:
			panic("Log type must be INFO,DEBUG,ERROR")
		}
		opts := &slog.HandlerOptions{
			Level: lvl,
		}
		logger = slog.New(slog.NewJSONHandler(os.Stdout, opts))
	})
}

func Debug(msg string, args ...any) {
	logger.Debug(msg, args...)
}
func DebugContext(ctx context.Context, msg string, args ...any) {
	logger.DebugContext(ctx, msg, args...)
}

func Error(msg string, args ...any) {
	logger.Error(msg, args...)
}
func ErrorContext(ctx context.Context, msg string, args ...any) {
	logger.ErrorContext(ctx, msg, args...)
}
func Info(msg string, args ...any) {
	logger.Info(msg, args...)
}
func InfoContext(ctx context.Context, msg string, args ...any) {
	logger.InfoContext(ctx, msg, args...)
}
func Log(ctx context.Context, level slog.Level, msg string, args ...any) {
	logger.Log(ctx, level, msg, args...)
}
func LogAttrs(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) {
	logger.LogAttrs(ctx, level, msg, attrs...)
}
func Warn(msg string, args ...any) {
	logger.Warn(msg, args...)
}
func WarnContext(ctx context.Context, msg string, args ...any) {
	logger.WarnContext(ctx, msg, args...)
}
