//go:generate stringer -type LogLevel logger.go
package log

type LogLevel int

const (
	Info LogLevel = iota + 1
	Warn
	Error
	Fatal
	Debug
)

type Logger interface {
	Info(msg string)
	Warn(msg string)
	Err(msg string)
	Fatal(msg string)
	Debug(msg string)
}
