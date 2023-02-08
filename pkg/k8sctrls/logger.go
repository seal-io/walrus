package k8sctrls

import (
	"github.com/go-logr/logr"
	ctrlog "sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/seal-io/seal/utils/log"
)

func newLogrLogger(delegate log.Logger) logr.Logger {
	return logr.New(ctrlog.NewDelegatingLogSink(logSinker{logger: delegate}))
}

type logSinker struct {
	logger log.Logger
}

func (l logSinker) Init(info logr.RuntimeInfo) {}

func (l logSinker) Enabled(level int) bool {
	return l.logger.V(uint64(level)).Enabled()
}

func (l logSinker) Info(level int, msg string, keysAndValues ...interface{}) {
	l.logger.V(uint64(level)).InfoS(msg, keysAndValues...)
}

func (l logSinker) Error(err error, msg string, keysAndValues ...interface{}) {
	l.logger.ErrorS(err, msg, keysAndValues...)
}

func (l logSinker) WithValues(keysAndValues ...interface{}) logr.LogSink {
	return logSinker{
		logger: l.logger.WithValues(keysAndValues...),
	}
}

func (l logSinker) WithName(name string) logr.LogSink {
	return logSinker{
		logger: l.logger.WithName(name),
	}
}
