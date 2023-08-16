package telemetry

import "github.com/seal-io/walrus/utils/log"

var (
	logger     = log.GetLogger().WithName("telemetry")
	wrapLogger = phLogger{
		Logger: logger,
	}
)

type phLogger struct {
	log.Logger
}

func (l phLogger) Logf(format string, args ...any) {
	l.Logger.Infof(format, args...)
}

func (l phLogger) Errorf(format string, args ...any) {
	l.Logger.Errorf(format, args...)
}
