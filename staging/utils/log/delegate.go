package log

type silenceLogger struct{}

func (silenceLogger) Enabled() bool                        { return false }
func (silenceLogger) Write(p []byte) (int, error)          { return len(p), nil }
func (silenceLogger) Info(...interface{})                  {}
func (silenceLogger) Infof(string, ...interface{})         {}
func (silenceLogger) InfoS(string, ...interface{})         {}
func (silenceLogger) Error(...interface{})                 {}
func (silenceLogger) ErrorS(error, string, ...interface{}) {}

type verboseLogger struct {
	v             uint64
	Delegate      Logger
	KeysAndValues []interface{}
}

func (l verboseLogger) Enabled() bool {
	return l.v <= GetVerbosity()
}
func (l verboseLogger) Write(p []byte) (int, error) {
	if !l.Enabled() {
		return len(p), nil
	}
	return l.Delegate.Write(p)
}
func (l verboseLogger) Info(args ...interface{}) {
	if !l.Enabled() {
		return
	}
	l.Delegate.Info(args...)
}
func (l verboseLogger) Infof(format string, args ...interface{}) {
	if !l.Enabled() {
		return
	}
	l.Delegate.Infof(format, args...)
}
func (l verboseLogger) Error(args ...interface{}) {
	if !l.Enabled() {
		return
	}
	l.Delegate.Error(args...)
}
func (l verboseLogger) InfoS(msg string, keyAndValues ...interface{}) {
	if !l.Enabled() {
		return
	}
	l.Delegate.InfoS(msg, append(keyAndValues, l.KeysAndValues...)...)
}
func (l verboseLogger) ErrorS(err error, msg string, keysAndValues ...interface{}) {
	if !l.Enabled() {
		return
	}
	l.Delegate.ErrorS(err, msg, append(keysAndValues, l.KeysAndValues...)...)
}

type DelegatedLogger struct {
	Delegate Logger
}

func (l DelegatedLogger) Write(p []byte) (int, error) {
	if l.Delegate == nil {
		return len(p), nil
	}
	return l.Delegate.Write(p)
}

func (l DelegatedLogger) Recovering() {
	if l.Delegate == nil {
		return
	}
	l.Delegate.Recovering()
}

func (l DelegatedLogger) Recover(p interface{}) {
	if l.Delegate == nil {
		return
	}
	l.Delegate.Recover(p)
}

func (l DelegatedLogger) Debug(args ...interface{}) {
	if l.Delegate == nil {
		return
	}
	l.Delegate.Debug(args...)
}
func (l DelegatedLogger) Info(args ...interface{}) {
	if l.Delegate == nil {
		return
	}
	l.Delegate.Info(args...)
}
func (l DelegatedLogger) Warn(args ...interface{}) {
	if l.Delegate == nil {
		return
	}
	l.Delegate.Warn(args...)
}
func (l DelegatedLogger) Error(args ...interface{}) {
	if l.Delegate == nil {
		return
	}
	l.Delegate.Error(args...)
}
func (l DelegatedLogger) Fatal(args ...interface{}) {
	if l.Delegate == nil {
		return
	}
	l.Delegate.Fatal(args...)
}
func (l DelegatedLogger) Debugf(format string, args ...interface{}) {
	if l.Delegate == nil {
		return
	}
	l.Delegate.Debugf(format, args...)
}
func (l DelegatedLogger) Infof(format string, args ...interface{}) {
	if l.Delegate == nil {
		return
	}
	l.Delegate.Infof(format, args...)
}
func (l DelegatedLogger) Warnf(format string, args ...interface{}) {
	if l.Delegate == nil {
		return
	}
	l.Delegate.Warnf(format, args...)
}
func (l DelegatedLogger) Errorf(format string, args ...interface{}) {
	if l.Delegate == nil {
		return
	}
	l.Delegate.Errorf(format, args...)
}
func (l DelegatedLogger) Fatalf(format string, args ...interface{}) {
	if l.Delegate == nil {
		return
	}
	l.Delegate.Fatalf(format, args...)
}
func (l DelegatedLogger) DebugS(msg string, keysAndValues ...interface{}) {
	if l.Delegate == nil {
		return
	}
	l.Delegate.DebugS(msg, keysAndValues...)
}
func (l DelegatedLogger) InfoS(msg string, keysAndValues ...interface{}) {
	if l.Delegate == nil {
		return
	}
	l.Delegate.InfoS(msg, keysAndValues...)
}
func (l DelegatedLogger) WarnS(msg string, keysAndValues ...interface{}) {
	if l.Delegate == nil {
		return
	}
	l.Delegate.WarnS(msg, keysAndValues...)
}
func (l DelegatedLogger) ErrorS(err error, msg string, keysAndValues ...interface{}) {
	if l.Delegate == nil {
		return
	}
	l.Delegate.ErrorS(err, msg, keysAndValues...)
}
func (l DelegatedLogger) FatalS(msg string, keysAndValues ...interface{}) {
	if l.Delegate == nil {
		return
	}
	l.Delegate.FatalS(msg, keysAndValues...)
}
func (l DelegatedLogger) Print(args ...interface{}) {
	if l.Delegate == nil {
		return
	}
	l.Delegate.Print(args...)
}
func (l DelegatedLogger) Printf(format string, args ...interface{}) {
	if l.Delegate == nil {
		return
	}
	l.Delegate.Printf(format, args...)
}
func (l DelegatedLogger) PrintS(msg string, keysAndValues ...interface{}) {
	if l.Delegate == nil {
		return
	}
	l.Delegate.PrintS(msg, keysAndValues...)
}
func (l DelegatedLogger) Enabled(v LoggingLevel) bool {
	if l.Delegate == nil {
		return false
	}
	if v < minLevel || v > maxLevel {
		return false
	}
	return l.Delegate.Enabled(v)
}
func (l DelegatedLogger) SetLevel(v LoggingLevel) {
	if l.Delegate == nil {
		return
	}
	l.Delegate.SetLevel(v)
}
func (l DelegatedLogger) GetLevel() LoggingLevel {
	if l.Delegate == nil {
		return minLevel
	}
	return l.Delegate.GetLevel()
}
func (l DelegatedLogger) V(v uint64) VerbosityLogger {
	if l.Delegate == nil {
		return silenceLogger{}
	}
	return l.Delegate.V(v)
}
func (l DelegatedLogger) WithName(name string) Logger {
	if l.Delegate == nil {
		return DelegatedLogger{}
	}
	return l.Delegate.WithName(name)
}
func (l DelegatedLogger) WithValues(keysAndValues ...interface{}) Logger {
	if l.Delegate == nil {
		return DelegatedLogger{}
	}
	return l.Delegate.WithValues(keysAndValues...)
}
