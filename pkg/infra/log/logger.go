package log

import (
	"context"
	"examples/pkg/entity/infra/log"
	pkgLogger "examples/pkg/logger"
	"fmt"
	"io"
)

type logger struct {
	out io.Writer
	ch  chan string
}

func InitLogger(ctx context.Context, w io.Writer) {
	l := &logger{
		out: w,
		ch:  make(chan string, 10),
	}
	go l.start(ctx)
	pkgLogger.L = l
}

func (l *logger) log(msg string, level log.LogLevel) {
	l.ch <- fmt.Sprintf("[%s]%s", level, msg)
}

func (l *logger) Info(msg string) {
	l.log(msg, log.Info)
}

func (l *logger) Warn(msg string) {
	l.log(msg, log.Warn)
}

func (l *logger) Err(msg string) {
	l.log(msg, log.Error)
}

func (l *logger) Debug(msg string) {
	l.log(msg, log.Debug)
}

func (l *logger) Fatal(msg string) {
	l.log(msg, log.Fatal)
}

func (l *logger) start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-l.ch:
			fmt.Fprintln(l.out, msg)
		}
	}
}

var _ log.Logger = (*logger)(nil)
