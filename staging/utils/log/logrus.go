package log

import (
	"strings"

	"github.com/sirupsen/logrus"
)

func AsLogrusFormatter(logger Logger, filters ...Filter) logrus.Formatter {
	return logrusFormatter{logger: logger, filters: filters}
}

type logrusFormatter struct {
	logger  Logger
	filters []Filter
}

func (l logrusFormatter) Format(entry *logrus.Entry) (bs []byte, err error) {
	lvl := DebugLevel

	switch entry.Level {
	case logrus.InfoLevel:
		lvl = InfoLevel
	case logrus.WarnLevel:
		lvl = WarnLevel
	case logrus.ErrorLevel:
		lvl = ErrorLevel
	case logrus.FatalLevel, logrus.PanicLevel:
		lvl = FatalLevel
	}

	if l.logger.Enabled(lvl) {
		var (
			msg           = entry.Message
			keysAndValues = make([]any, 0, 2*len(entry.Data))
		)

		for k := range entry.Data {
			keysAndValues = append(keysAndValues, k, entry.Data[k])
		}

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

		switch lvl {
		case DebugLevel:
			l.logger.DebugS(msg, keysAndValues...)
		case InfoLevel:
			l.logger.InfoS(msg, keysAndValues...)
		case WarnLevel:
			l.logger.WarnS(msg, keysAndValues...)
		case ErrorLevel:
			l.logger.ErrorS(nil, msg, keysAndValues...)
		case FatalLevel:
			l.logger.FatalS(msg, keysAndValues...)
		}
	}

	return
}
