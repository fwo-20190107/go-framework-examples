package model

import (
	"context"
)

type LogLevel int

const (
	LogLevelInfo LogLevel = iota + 1
	LogLevelWarn
	LogLevelErr
	LogLevelFatal
	LogLevelDebug
)

type Logging interface {
	Log(ctx context.Context, msg string, level LogLevel) error
	Info(ctx context.Context, msg string) error
	Warn(ctx context.Context, msg string) error
	Err(ctx context.Context, msg string) error
	Debug(ctx context.Context, msg string) error
	Fatal(ctx context.Context, msg string) error
	Send(ctx context.Context) error
}

var Logger Logging

func (l LogLevel) String() (s string) {
	switch l {
	case LogLevelInfo:
		s = "INFO"
	case LogLevelWarn:
		s = "WARNING"
	case LogLevelErr:
		s = "ERROR"
	case LogLevelFatal:
		s = "FATAL"
	case LogLevelDebug:
		s = "DEBUG"
	default:
		s = "NONE"
	}
	return
}
