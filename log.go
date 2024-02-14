package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"
)

var lock = &sync.Mutex{}

type LogHandler struct{}

var loggerInstance *slog.Logger

var globalLogLevel = slog.LevelDebug

var logLevelMap = map[int]slog.Level{
	0: slog.LevelError,
	1: slog.LevelWarn,
	2: slog.LevelInfo,
	3: slog.LevelDebug,
}

func GetLogger() *slog.Logger {
	lock.Lock()
	defer lock.Unlock()
	if loggerInstance == nil {
		loggerInstance = slog.New(&LogHandler{})
		slog.SetDefault(loggerInstance)
	}

	return loggerInstance
}

func SetLogLevel(logLevel int) {
	if level, ok := logLevelMap[logLevel]; ok {
		globalLogLevel = level
	} else {
		panic(fmt.Sprintf("invalid log level: %v", logLevel))
	}
}

func (l *LogHandler) Enabled(context context.Context, level slog.Level) bool {
	return level >= globalLogLevel
}

func (l *LogHandler) Handle(context context.Context, record slog.Record) error {
	message := record.Message

	record.Attrs(func(attr slog.Attr) bool {
		message += fmt.Sprintf(" %v", attr)
		return true
	})

	timestamp := record.Time.Format(time.RFC3339)

	if record.Level == slog.LevelError {
		fmt.Fprintf(os.Stderr, "%v [%v] %v\n", timestamp, record.Level, message)
	} else {
		fmt.Fprintf(os.Stdout, "%v [%v] %v\n", timestamp, record.Level, message)
	}

	return nil
}

func (l *LogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	panic("unimplemented")
}

func (l *LogHandler) WithGroup(name string) slog.Handler {
	panic("unimplemented")
}
