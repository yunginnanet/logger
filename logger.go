package logger

import (
	"io"
	"os"
	"runtime"
	"sync"

	"git.tcp.direct/kayos/zwrap"
	"github.com/rs/zerolog"
)

var (
	globalLogger *Log
	globalMu     sync.RWMutex
)

type Log struct {
	l       *zwrap.Logger
	writers []io.Writer

	mu sync.RWMutex
}

// NewQuietLogger creates a logger that writes to the given writers with no console writer added.
func NewQuietLogger(writers ...io.Writer) *Log {
	return newLogger(writers...)
}

// NewLogger creates a logger that writes to the given writers, as well as pretty prints to stdout.
func NewLogger(writers ...io.Writer) *Log {
	zcl := &zerolog.ConsoleWriter{
		Out:         os.Stdout,
		NoColor:     false,
		FormatLevel: zwrap.LogLevelFmt(runtime.GOOS == "windows"),
	}
	newW := make([]io.Writer, 0, len(writers))
	newW = append(newW, zcl)
	newW = append(newW, writers...)
	return newLogger(newW...)
}

func NewLoggerNoColor(writers ...io.Writer) *Log {
	zcl := &zerolog.ConsoleWriter{
		Out:         os.Stdout,
		NoColor:     false,
		FormatLevel: zwrap.LogLevelFmt(true),
	}
	newW := make([]io.Writer, 0, len(writers))
	newW = append(newW, zcl)
	newW = append(newW, writers...)
	return newLogger(newW...)
}

func newLogger(writers ...io.Writer) *Log {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	l := &Log{
		writers: make([]io.Writer, 0, 1),
	}

	l.writers = append(l.writers, writers...)
	w := zerolog.MultiLevelWriter(l.writers...)

	zl := zerolog.New(w).With().Timestamp().Logger()

	wrapped := zwrap.Wrap(zl)
	l.l = wrapped

	return l
}

// AddWriter adds a writer to the logger.
//
// Note: this may have unintended consequences if certain [zwrap.Logger] configuration values have been set via [Log.C].
// That said, the [zwrap.Logger] prefix value will be preserved.
func (l *Log) AddWriter(w io.Writer) {
	l.mu.Lock()
	l.writers = append(l.writers, w)
	zw := zerolog.MultiLevelWriter(l.writers...)
	zl := zerolog.New(zw).With().Timestamp().Logger()
	oldPrefix := l.l.Prefix()
	wrapped := zwrap.Wrap(zl).WithPrefix(oldPrefix)
	l.l = wrapped
	l.mu.Unlock()
}

// C returns a [zwrap.ZWrapLogger] which is a highly compattible interface to fit many other log intrefaces.
func (l *Log) C() zwrap.ZWrapLogger {
	l.mu.RLock()
	ll := l.l
	l.mu.RUnlock()
	return ll
}

// Z rerturns a pointer to the underlying [zerolog.Logger].
func (l *Log) Z() *zerolog.Logger {
	l.mu.RLock()
	ll := l.l.ZLogger()
	l.mu.RUnlock()
	return ll
}

func (l *Log) WithGlobalPackageAccess() *Log {
	globalMu.Lock()
	globalLogger = l
	globalMu.Unlock()
	return l
}

// Global acquires the assigned global logger.
//
// IMPORTANT: you MUST make your instance of [Log] globally accecible by calling [WithGlobalPackageAccess].
func Global() *Log {
	globalMu.RLock()
	l := globalLogger
	globalMu.RUnlock()
	return l
}
