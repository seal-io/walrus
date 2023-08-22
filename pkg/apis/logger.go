package apis

import (
	stdlog "log"
	"strings"

	"github.com/seal-io/walrus/utils/log"
)

func newStdErrorLogger(delegate log.Logger) *stdlog.Logger {
	return stdlog.New(logWriter{logger: delegate}, "", 0)
}

type logWriter struct {
	logger log.Logger
}

func (l logWriter) Write(p []byte) (int, error) {
	// Trim the trailing newline.
	s := strings.TrimSuffix(string(p), "\n")

	ok := true

	switch {
	case strings.HasPrefix(s, "http: TLS handshake error from"):
		switch {
		case strings.HasSuffix(s, "tls: unknown certificate"):
			// Ignore self-generated certificate errors from client.
			ok = false
		case strings.HasSuffix(s, "connection reset by peer"):
			// Reset TLS handshake errors from client.
			ok = false
		case strings.HasSuffix(s, "EOF"):
			// Terminate TLS handshake errors by client.
			ok = false
		}
	case strings.Contains(s, "broken pipe"):
		// Ignore the underlying error of broken pipe.
		ok = false
	}

	if ok {
		l.logger.Warn(s)
	}

	return len(p), nil
}
