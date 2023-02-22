package log

import "github.com/go-logr/logr"

// AsLogrSink converts the Logger to Logr sink.
func AsLogrSink(logger Logger) logr.LogSink {
	return logrSinker{logger: logger}
}

// AsLogr converts the Logger to Logr logger.
func AsLogr(logger Logger) logr.Logger {
	return logr.New(AsLogrSink(logger))
}

type logrSinker struct {
	logger Logger
}

func (l logrSinker) Init(info logr.RuntimeInfo) {}

func (l logrSinker) Enabled(level int) bool {
	return l.logger.V(uint64(level)).Enabled()
}

func (l logrSinker) Info(level int, msg string, keysAndValues ...interface{}) {
	l.logger.V(uint64(level)).InfoS(msg, keysAndValues...)
}

func (l logrSinker) Error(err error, msg string, keysAndValues ...interface{}) {
	l.logger.ErrorS(err, msg, keysAndValues...)
}

func (l logrSinker) WithValues(keysAndValues ...interface{}) logr.LogSink {
	return logrSinker{
		logger: l.logger.WithValues(keysAndValues...),
	}
}

func (l logrSinker) WithName(name string) logr.LogSink {
	return logrSinker{
		logger: l.logger.WithName(name),
	}
}
