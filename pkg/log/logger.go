package log

import (
	"context"
	"github.com/donghc/goutils/pkg/netutil"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"io"
	"os"
)

// Logger is a logger that supports log levels, context and structured logging.
type Logger interface {
	// With returns a logger based off the root logger and decorates it with
	// the given context and arguments.
	With(ctx context.Context, args ...interface{}) Logger

	// Debug uses fmt.Sprint to construct and log a message at DEBUG level
	Debug(args ...interface{})
	// Info uses fmt.Sprint to construct and log a message at INFO level
	Info(args ...interface{})
	// Error uses fmt.Sprint to construct and log a message at ERROR level
	Error(args ...interface{})
	// Warn uses fmt.Sprint to construct and log a message at WARN level
	Warn(args ...interface{})
	// Panic uses fmt.Sprint to construct and log a message at Panic level
	Panic(args ...interface{})

	// Debugf uses fmt.Sprintf to construct and log a message at DEBUG level
	Debugf(format string, args ...interface{})
	// Infof uses fmt.Sprintf to construct and log a message at INFO level
	Infof(format string, args ...interface{})
	// Warnf uses fmt.Sprintf to construct and log a message at WARN level
	Warnf(format string, args ...interface{})
	// Errorf uses fmt.Sprintf to construct and log a message at ERROR level
	Errorf(format string, args ...interface{})
	// Panicf uses fmt.Sprintf to construct and log a message at Panic level
	Panicf(format string, args ...interface{})
}

type logger struct {
	*zap.SugaredLogger
}

type contextKey int

const (
	requestIDKey contextKey = iota
	correlationIDKey
)

var (
	hostname, _ = os.Hostname()
	addrs, _    = netutil.ExternalIP()
)

// New creates a new logger using the default configuration.
func New() Logger {
	l, _ := zap.NewProduction()
	return NewWithZap(l)
}

func NewCustomLogger(level string, ws []io.Writer, opts []zap.Option) Logger {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		FunctionKey:    zapcore.OmitKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // ???????????????
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC ????????????
		EncodeDuration: zapcore.SecondsDurationEncoder, //	?????????
		EncodeCaller:   zapcore.ShortCallerEncoder,     // ??????????????????
		EncodeName:     zapcore.FullNameEncoder,
	}
	var zws []zapcore.WriteSyncer
	for i := range ws {
		zws = append(zws, zapcore.AddSync(ws[i]))
	}
	loggingLevel := getLoggingLevel(level)
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig), // ???????????????
		zapcore.NewMultiWriteSyncer(zws...),   // ???????????????????????????
		loggingLevel,                          // ????????????
	)

	return NewWithZap(zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.PanicLevel),
		zap.Fields(zap.String("hostname", hostname)), zap.Fields(zap.String("ip", addrs.String()))), opts...)

}

// NewCustom creates a new logger using a custom configuration
// given a log level.
func NewCustom(level string, w io.Writer, opts []zap.Option) Logger {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "line",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder, // ???????????????
		EncodeTime:     zapcore.ISO8601TimeEncoder,    // ISO8601 UTC ????????????
		EncodeDuration: zapcore.MillisDurationEncoder, //	?????????
		EncodeCaller:   zapcore.FullCallerEncoder,     // ??????????????????
		EncodeName:     zapcore.FullNameEncoder,
	}
	loggingLevel := getLoggingLevel(level)
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),                                       // ???????????????
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(w)), // ???????????????????????????
		loggingLevel, // ????????????
	)

	return NewWithZap(zap.New(core, zap.Fields(zap.String("hostname", hostname)), zap.Fields(zap.String("ip", addrs.String()))), opts...)
}

// NewWithZap creates a new logger using the preconfigured zap logger.
func NewWithZap(l *zap.Logger, opts ...zap.Option) Logger {
	l = l.WithOptions(opts...)
	return &logger{l.Sugar()}
}

// NewForTest returns a new logger and the corresponding observed logs which
// can be used in unit tests to verify log entries.
func NewForTest() (Logger, *observer.ObservedLogs) {
	core, recorded := observer.New(zapcore.InfoLevel)
	return NewWithZap(zap.New(core)), recorded
}

// With returns a logger based off the root logger and decorates it with
// the given context and arguments.
//
// If the context contains request ID and/or correlation ID information
//(recorded via WithRequestID() and WithCorrelationID()), they will be
// added to every log message generated by the new logger.
//
// The arguments should be specified as a sequence of name, value pairs
// with names being strings.
// The arguments will also be added to every log message generated by the logger.
func (l *logger) With(ctx context.Context, args ...interface{}) Logger {
	if ctx != nil {
		if id, ok := ctx.Value(requestIDKey).(string); ok {
			args = append(args, zap.String("request_id", id))
		}
		if id, ok := ctx.Value(correlationIDKey).(string); ok {
			args = append(args, zap.String("correlation_id", id))
		}
	}
	if len(args) > 0 {
		return &logger{l.SugaredLogger.With(args...)}
	}
	return l
}
