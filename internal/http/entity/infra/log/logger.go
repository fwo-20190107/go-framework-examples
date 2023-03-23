//go:generate stringer -type LogLevel logger.go
package log

import (
	"context"
)

type LogLevel int

const (
	Info LogLevel = iota + 1
	Warn
	Error
	Fatal
	Debug
)

type Logger interface {
	Log(ctx context.Context, msg string, level LogLevel) error
	Info(ctx context.Context, msg string) error
	Warn(ctx context.Context, msg string) error
	Err(ctx context.Context, msg string) error
	Debug(ctx context.Context, msg string) error
	Fatal(ctx context.Context, msg string) error
	Send(ctx context.Context) error
}
