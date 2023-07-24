package telemetry

import "github.com/seal-io/seal/utils/log"

var (
	logger     = log.GetLogger().WithName("telemetry")
	wrapLogger = phLogger{
		Logger: logger,
	}
)

type phLogger struct {
	log.Logger
}

func (l phLogger) Logf(format string, args ...interface{}) {
	l.Logger.Infof(format, args...)
}

func (l phLogger) Errorf(format string, args ...interface{}) {
	l.Logger.Errorf(format, args...)
}
