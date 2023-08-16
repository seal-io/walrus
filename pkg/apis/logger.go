package apis

import (
	stdlog "log"
	"strings"

	"github.com/seal-io/seal/utils/log"
)

func newStdErrorLogger(delegate log.Logger) *stdlog.Logger {
	return stdlog.New(logWriter{logger: delegate}, "", 0)
}

type logWriter struct {
	logger log.Logger
}

func (l logWriter) Write(p []byte) (int, error) {
	s := string(p)

	ok := true

	switch {
	case strings.HasPrefix(s, "http: TLS handshake error from"):
		switch {
		case strings.HasSuffix(s, "tls: unknown certificate\n"):
			// Ignore self-generated certificate errors from client.
			ok = false
		case strings.HasSuffix(s, "connection reset by peer\n"):
			// Reset TLS handshake errors from client.
			ok = false
		}
	case strings.Contains(s, "broken pipe"):
		// Terminate by client.
		ok = false
	}

	if ok {
		l.logger.Warn(s)
	}

	return len(p), nil
}
