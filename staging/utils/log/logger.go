package log

import (
	"io"
	"strconv"

	"go.uber.org/atomic"
)

type LoggingLevel uint8

const (
	minLevel LoggingLevel = iota
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
	maxLevel
)

var verbosity = atomic.NewUint64(^uint64(0) >> 1)

func (l LoggingLevel) String() string {
	if minLevel < l && l < maxLevel {
		switch l {
		case DebugLevel:
			return "debug"
		case InfoLevel:
			return "info"
		case WarnLevel:
			return "warn"
		case ErrorLevel:
			return "error"
		case FatalLevel:
			return "fatal"
		}
	}

	return "unknown level " + strconv.FormatUint(uint64(l), 10)
}

type RecoverLogger interface {
	Recovering()
	Recover(i any)
}

type ValueLogger interface {
	Debug(args ...any)
	Info(args ...any)
	Warn(args ...any)
	Error(args ...any)
	Fatal(args ...any)
}

type FormatLogger interface {
	Debugf(format string, args ...any)
	Infof(format string, args ...any)
	Warnf(format string, args ...any)
	Errorf(format string, args ...any)
	Fatalf(format string, args ...any)
}

type StructLogger interface {
	DebugS(msg string, keysAndValues ...any)
	InfoS(msg string, keysAndValues ...any)
	WarnS(msg string, keysAndValues ...any)
	ErrorS(err error, msg string, keysAndValues ...any)
	FatalS(msg string, keysAndValues ...any)
}

type PrinterLogger interface {
	Print(args ...any)
	Printf(format string, args ...any)
	PrintS(msg string, keysAndValues ...any)
	Println(args ...any)
}

type VerbosityLogger interface {
	io.Writer

	Enabled() bool
	Info(args ...any)
	Infof(format string, args ...any)
	InfoS(msg string, keysAndValues ...any)
	Error(args ...any)
	ErrorS(err error, msg string, keysAndValues ...any)
}

type Logger interface {
	io.Writer
	RecoverLogger
	ValueLogger
	FormatLogger
	StructLogger
	PrinterLogger

	Enabled(v LoggingLevel) bool
	SetLevel(v LoggingLevel)
	GetLevel() LoggingLevel
	V(v uint64) VerbosityLogger
	WithName(name string) Logger
	WithValues(keysAndValues ...any) Logger
}

type LegacyLogger interface {
	Warning(args ...any)
	Warningf(format string, args ...any)
	WarningS(msg string, keysAndValues ...any)
}

// logger holds the global logger.
var logger = DelegatedLogger{
	Delegate: NewDevelopmentWrappedZapperAsLogger(),
}

// Write exposes the io.Writer implementation of the global logger.
func Write(p []byte) (n int, err error) {
	return logger.Write(p)
}

// Recovering exposes the RecoverLogger implementation of the global logger.
func Recovering() {
	logger.Recovering()
}

// Recover exposes the RecoverLogger implementation of the global logger.
func Recover(i any) {
	logger.Recover(i)
}

// Debug exposes the ValueLogger implementation of the global logger.
func Debug(args ...any) {
	logger.Debug(args...)
}

// Info exposes the ValueLogger implementation of the global logger.
func Info(args ...any) {
	logger.Info(args...)
}

// Warn exposes the ValueLogger implementation of the global logger.
func Warn(args ...any) {
	logger.Warn(args...)
}

// Error exposes the ValueLogger implementation of the global logger.
func Error(args ...any) {
	logger.Error(args...)
}

// Fatal exposes the ValueLogger implementation of the global logger.
func Fatal(args ...any) {
	logger.Fatal(args...)
}

// Debugf exposes the FormatLogger implementation of the global logger.
func Debugf(format string, args ...any) {
	logger.Debugf(format, args...)
}

// Infof exposes the FormatLogger implementation of the global logger.
func Infof(format string, args ...any) {
	logger.Infof(format, args...)
}

// Warnf exposes the FormatLogger implementation of the global logger.
func Warnf(format string, args ...any) {
	logger.Warnf(format, args...)
}

// Errorf exposes the FormatLogger implementation of the global logger.
func Errorf(format string, args ...any) {
	logger.Errorf(format, args...)
}

// Fatalf exposes the FormatLogger implementation of the global logger.
func Fatalf(format string, args ...any) {
	logger.Fatalf(format, args...)
}

// DebugS exposes the StructLogger implementation of the global logger.
func DebugS(msg string, keysAndValues ...any) {
	logger.DebugS(msg, keysAndValues...)
}

// InfoS exposes the StructLogger implementation of the global logger.
func InfoS(msg string, keysAndValues ...any) {
	logger.InfoS(msg, keysAndValues...)
}

// WarnS exposes the StructLogger implementation of the global logger.
func WarnS(msg string, keysAndValues ...any) {
	logger.WarnS(msg, keysAndValues...)
}

// ErrorS exposes the StructLogger implementation of the global logger.
func ErrorS(err error, msg string, keysAndValues ...any) {
	logger.ErrorS(err, msg, keysAndValues...)
}

// FatalS exposes the StructLogger implementation of the global logger.
func FatalS(msg string, keysAndValues ...any) {
	logger.FatalS(msg, keysAndValues...)
}

// Print exposes the PrinterLogger implementation of the global logger.
func Print(args ...any) {
	logger.Print(args...)
}

// Printf exposes the PrinterLogger implementation of the global logger.
func Printf(format string, args ...any) {
	logger.Printf(format, args...)
}

// PrintS exposes the PrinterLogger implementation of the global logger.
func PrintS(msg string, keysAndValues ...any) {
	logger.PrintS(msg, keysAndValues...)
}

// Enabled exposes the Logger implementation of the global logger.
func Enabled(v LoggingLevel) bool {
	return logger.Enabled(v)
}

// SetLevel set the Logger level of the global logger.
func SetLevel(v LoggingLevel) {
	logger.SetLevel(v)
}

// GetLevel exposes the Logger level of the global logger.
func GetLevel() LoggingLevel {
	return logger.GetLevel()
}

// V exposes the Logger implementation of the global logger.
func V(v uint64) VerbosityLogger {
	return logger.V(v)
}

// WithName exposes the Logger implementation of the global logger.
func WithName(name string) Logger {
	return logger.WithName(name)
}

// WithValues exposes the Logger implementation of the global logger.
func WithValues(keysAndValues ...any) Logger {
	return logger.WithValues(keysAndValues)
}

// Warning exposes the LegacyLogger implementation of the global logger.
func Warning(args ...any) {
	logger.Warn(args...)
}

// Warningf exposes the LegacyLogger implementation of the global logger.
func Warningf(format string, args ...any) {
	logger.Warnf(format, args...)
}

// WarningS exposes the LegacyLogger implementation of the global logger.
func WarningS(msg string, keysAndValues ...any) {
	logger.WarnS(msg, keysAndValues...)
}

// GetLogger returns the global Logger implement.
func GetLogger() Logger {
	return logger
}

// SetLogger configures the global Logger.
func SetLogger(delegate Logger) {
	logger.Delegate = delegate
}

// GetVerbosity returns the verbosity of the global Logger.
func GetVerbosity() uint64 {
	return verbosity.Load()
}

// SetVerbosity configures the verbosity of the global Logger.
func SetVerbosity(level uint64) {
	verbosity.Store(level)
}

// WrapAsVerbosityLogger wraps a given Logger as VerbosityLogger.
func WrapAsVerbosityLogger(v uint64, r Logger) VerbosityLogger {
	return verboseLogger{v: v, Delegate: r}
}
