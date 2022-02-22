package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
)

func getLoggingLevel(level string) zap.AtomicLevel {
	atom := zap.NewAtomicLevel()
	switch strings.ToLower(level) {
	case "panic":
		atom.SetLevel(zapcore.PanicLevel)
		return atom
	case "fatal":
		atom.SetLevel(zapcore.FatalLevel)
		return atom
	case "error":
		atom.SetLevel(zapcore.ErrorLevel)
		return atom
	case "warn":
		atom.SetLevel(zapcore.WarnLevel)
		return atom
	case "info":
		atom.SetLevel(zapcore.InfoLevel)
		return atom
	case "debug":
		atom.SetLevel(zapcore.DebugLevel)
		return atom
	default:
		atom.SetLevel(zapcore.InfoLevel)
		return atom
	}
}
