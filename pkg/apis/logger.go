package apis

import (
	"bufio"
	"bytes"
	stdlog "log"
	"strings"

	"github.com/seal-io/seal/utils/log"
)

func newStdLogger(delegate log.Logger) *stdlog.Logger {
	return stdlog.New(logWriter{logger: delegate}, "", stdlog.Lshortfile)
}

type logWriter struct {
	logger log.Logger
}

func (l logWriter) Write(p []byte) (int, error) {
	s := bufio.NewScanner(bytes.NewReader(p))
	for s.Scan() {
		if strings.HasSuffix(s.Text(), "tls: unknown certificate") {
			continue
		}
		l.logger.Info(s.Text())
	}
	return len(p), s.Err()
}
