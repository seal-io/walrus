package log

import (
	"strings"

	"github.com/go-logr/logr"
)

// AsLogrSink converts the Logger to Logr sink.
func AsLogrSink(logger Logger, filters ...Filter) logr.LogSink {
	return logrSinker{logger: logger, filters: filters}
}

// AsLogr converts the Logger to Logr logger.
func AsLogr(logger Logger, filters ...Filter) logr.Logger {
	return logr.New(AsLogrSink(logger, filters...))
}

type logrSinker struct {
	logger  Logger
	filters []Filter
}

func (l logrSinker) Init(info logr.RuntimeInfo) {}

func (l logrSinker) Enabled(level int) bool {
	return l.logger.V(uint64(level)).Enabled()
}

func (l logrSinker) Info(v int, msg string, keysAndValues ...any) {
	if uint64(v) <= GetVerbosity() {
		msg = strings.TrimSuffix(msg, "\n")

		for i := range l.filters {
			if l.filters[i] == nil {
				continue
			}

			var ok bool

			msg, keysAndValues, ok = l.filters[i](msg, keysAndValues)
			if !ok {
				return
			}
		}

		l.logger.InfoS(msg, keysAndValues...)
	}
}

func (l logrSinker) Error(err error, msg string, keysAndValues ...any) {
	msg = strings.TrimSuffix(msg, "\n")

	for i := range l.filters {
		if l.filters[i] == nil {
			continue
		}

		var ok bool

		msg, keysAndValues, ok = l.filters[i](msg, keysAndValues)
		if !ok {
			return
		}
	}

	l.logger.ErrorS(err, msg, keysAndValues...)
}

func (l logrSinker) WithValues(keysAndValues ...any) logr.LogSink {
	return logrSinker{
		logger: l.logger.WithValues(keysAndValues...),
	}
}

func (l logrSinker) WithName(name string) logr.LogSink {
	return logrSinker{
		logger: l.logger.WithName(name),
	}
}
