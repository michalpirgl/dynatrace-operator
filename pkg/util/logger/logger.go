package logger

import (
	"os"

	"github.com/go-logr/logr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	ctrlzap "sigs.k8s.io/controller-runtime/pkg/log/zap"
)

type logSink struct {
	infoLogger  logr.Logger
	errorLogger logr.Logger
}

func newLogger() logr.Logger {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder

	level, err := zapcore.ParseLevel(os.Getenv("LOG_LEVEL")) // Quick and dirty way to make it configurable
	if err != nil {
		level = zapcore.InfoLevel
	}

	// Removed our weird logger-sink implementation here in favor of just a simple zap one for this POC
	// Our weird logger is only useful to print out errors in a nicer way, but it totally breaks the log.V(123) functionality
	// The reason for this is that the Enabled(level) method of the sink should actually consider the level, which it doesn't. (we would have to store it and all)
	// It wouldn't be that hard, BUT zap handles verbosity/level in an inverted way, which we could/would actually just override here if wanted :D
	// IMO, there must be a way to print the errors in a nice way, so we don't have to do a home-grown solution, and I would like to avoid coming up with a similarly home-grown solution for managing log verbosity/level
	return ctrlzap.New(ctrlzap.Encoder(zapcore.NewJSONEncoder(config)), ctrlzap.Level(level))
}

func (dtl logSink) Init(logr.RuntimeInfo) {}

func (dtl logSink) Info(_ int, msg string, keysAndValues ...any) {
	dtl.infoLogger.Info(msg, keysAndValues...)
}

func (dtl logSink) Enabled(int) bool {
	return dtl.infoLogger.Enabled()
}

func (dtl logSink) Error(err error, msg string, keysAndValues ...any) {
	dtl.errorLogger.Error(err, msg, keysAndValues...)
}

func (dtl logSink) WithValues(keysAndValues ...any) logr.LogSink {
	return logSink{
		infoLogger:  dtl.infoLogger.WithValues(keysAndValues...),
		errorLogger: dtl.errorLogger.WithValues(keysAndValues...),
	}
}

func (dtl logSink) WithName(name string) logr.LogSink {
	return logSink{
		infoLogger:  dtl.infoLogger.WithName(name),
		errorLogger: dtl.errorLogger.WithName(name),
	}
}
