package log

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"runtime/debug"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewZapper creates a Zap logger.
func NewZapper(asJSON, inProduction, toStdout bool) (*zap.Logger, zap.AtomicLevel) {
	zapWriteSyncer := zapcore.AddSync(os.Stderr)
	if toStdout {
		zapWriteSyncer = zapcore.AddSync(os.Stdout)
	}
	zapOptions := []zap.Option{
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zap.ErrorLevel),
		zap.ErrorOutput(zapWriteSyncer),
	}

	zapLevel := zap.NewAtomicLevelAt(zap.DebugLevel)
	zapEncoderConfig := zap.NewDevelopmentEncoderConfig()
	if inProduction {
		zapLevel.SetLevel(zap.InfoLevel)
		zapEncoderConfig = zap.NewProductionEncoderConfig()
		zapOptions = append(zapOptions,
			zap.WrapCore(func(core zapcore.Core) zapcore.Core {
				return zapcore.NewSamplerWithOptions(core, time.Second, 100, 100)
			}),
		)
	}

	zapEncoderConfig.EncodeLevel = func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		s := "D"
		switch l {
		case zapcore.InfoLevel:
			s = "I"
		case zapcore.WarnLevel:
			s = "W"
		case zapcore.ErrorLevel:
			s = "E"
		case zapcore.DPanicLevel, zapcore.PanicLevel:
			s = "P"
		case zapcore.FatalLevel:
			s = "F"
		}
		enc.AppendString(s)
	}
	if asJSON {
		zapEncoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	}
	zapEncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	if asJSON {
		zapEncoderConfig.EncodeTime = zapcore.EpochTimeEncoder
	}

	zapEncoder := zapcore.NewConsoleEncoder(zapEncoderConfig)
	if asJSON {
		zapEncoder = zapcore.NewJSONEncoder(zapEncoderConfig)
	}

	return zap.New(zapcore.NewCore(zapEncoder, zapWriteSyncer, zapLevel), zapOptions...), zapLevel
}

// NewDevelopmentWrappedZapperAsLogger create a wrapped Zap logger as Logger with development config.
func NewDevelopmentWrappedZapperAsLogger() Logger {
	l, lv := NewZapper(false, false, false)
	return zapLogger{l: l, s: l.Sugar(), lv: lv}
}

// NewWrappedZapperAsLogger create a wrapped Zap logger as Logger.
func NewWrappedZapperAsLogger(asJSON, inProduction, toStdout bool) Logger {
	l, lv := NewZapper(asJSON, inProduction, toStdout)
	return zapLogger{l: l, s: l.Sugar(), lv: lv}
}

// WrapZapperAsLogger wraps a Zap logger as Logger.
func WrapZapperAsLogger(l *zap.Logger, lv zap.AtomicLevel) Logger {
	return zapLogger{l: l, s: l.Sugar(), lv: lv}
}

type zapLogger struct {
	l  *zap.Logger
	s  *zap.SugaredLogger
	lv zap.AtomicLevel
}

func (z zapLogger) Write(p []byte) (int, error) {
	s := bufio.NewScanner(bytes.NewReader(p))
	for s.Scan() {
		z.s.Info(s.Text())
	}
	return len(p), s.Err()
}

func (z zapLogger) Recovering() {
	if p := recover(); p != nil {
		z.Recover(p)
	}
}

func (z zapLogger) Recover(p interface{}) {
	z.s.Errorf("observing panic: %v, stack trace: %s", p, string(debug.Stack()))
}

func (z zapLogger) Debug(args ...interface{}) {
	z.s.Debug(args...)
}

func (z zapLogger) Info(args ...interface{}) {
	z.s.Info(args...)
}

func (z zapLogger) Warn(args ...interface{}) {
	z.s.Warn(args...)
}

func (z zapLogger) Error(args ...interface{}) {
	z.s.Error(args...)
}

func (z zapLogger) Fatal(args ...interface{}) {
	z.s.Fatal(args...)
}

func (z zapLogger) Debugf(format string, args ...interface{}) {
	z.s.Debugf(format, args...)
}

func (z zapLogger) Infof(format string, args ...interface{}) {
	z.s.Infof(format, args...)
}

func (z zapLogger) Warnf(format string, args ...interface{}) {
	z.s.Warnf(format, args...)
}

func (z zapLogger) Errorf(format string, args ...interface{}) {
	z.s.Errorf(format, args...)
}

func (z zapLogger) Fatalf(format string, args ...interface{}) {
	z.s.Fatalf(format, args...)
}

func (z zapLogger) DebugS(msg string, keysAndValues ...interface{}) {
	z.s.Debugw(msg, keysAndValues...)
}

func (z zapLogger) InfoS(msg string, keysAndValues ...interface{}) {
	z.s.Infow(msg, keysAndValues...)
}

func (z zapLogger) WarnS(msg string, keysAndValues ...interface{}) {
	z.s.Warnw(msg, keysAndValues...)
}

func (z zapLogger) ErrorS(err error, msg string, keysAndValues ...interface{}) {
	if err == nil {
		z.s.Errorw(msg, keysAndValues...)
		return
	}
	z.s.With(zap.Error(err)).Errorw(msg, keysAndValues...)
}

func (z zapLogger) FatalS(msg string, keysAndValues ...interface{}) {
	z.s.Fatalw(msg, keysAndValues...)
}

func (z zapLogger) Print(args ...interface{}) {
	z.s.Info(args...)
}

func (z zapLogger) Printf(format string, args ...interface{}) {
	z.s.Infof(format, args...)
}

func (z zapLogger) PrintS(msg string, keysAndValues ...interface{}) {
	z.s.Infow(msg, keysAndValues...)
}

func (z zapLogger) Enabled(v LoggingLevel) bool {
	lvl := toZapLevel(v)
	return z.l.Core().Enabled(lvl)
}

func (z zapLogger) SetLevel(v LoggingLevel) {
	lvl := toZapLevel(v)
	z.lv.SetLevel(lvl)
}

func (z zapLogger) GetLevel() LoggingLevel {
	switch z.l.Level() {
	case zapcore.DebugLevel:
		return DebugLevel
	case zapcore.InfoLevel:
		return InfoLevel
	case zapcore.WarnLevel:
		return WarnLevel
	case zapcore.ErrorLevel:
		return ErrorLevel
	case zapcore.FatalLevel:
		return FatalLevel
	default:
		return minLevel
	}
}

func (z zapLogger) V(v uint64) VerbosityLogger {
	return WrapAsVerbosityLogger(v, z)
}

func (z zapLogger) WithName(name string) Logger {
	return zapLogger{
		l: z.l.Named(name),
		s: z.s.Named(name),
	}
}

func (z zapLogger) WithValues(keysAndValues ...interface{}) Logger {
	return zapLogger{
		l: z.l.With(handleFields(keysAndValues...)...),
		s: z.s.With(keysAndValues...),
	}
}

func handleFields(args ...interface{}) (fields []zap.Field) {
	argSize := len(args)
	if argSize == 0 {
		return
	}
	for i := 0; i < argSize; {
		var field zap.Field
		arg := args[i]
		switch a := arg.(type) {
		case zap.Field:
			field = a
		case string:
			if i+1 < argSize {
				field = zap.Any(a, args[i+1])
				i++
			} else {
				field = zap.Any("#key$", a)
			}
		case error:
			field = zap.Any(fmt.Sprintf("#err%d", i+1), a)
		default:
			field = zap.Any(fmt.Sprintf("#key%d", i+1), a)
		}
		fields = append(fields, field)
		i++
	}
	return
}

func toZapLevel(l LoggingLevel) (lvl zapcore.Level) {
	switch l {
	case DebugLevel:
		lvl = zapcore.DebugLevel
	case InfoLevel:
		lvl = zapcore.InfoLevel
	case WarnLevel:
		lvl = zapcore.WarnLevel
	case ErrorLevel:
		lvl = zapcore.ErrorLevel
	case FatalLevel:
		lvl = zapcore.FatalLevel
	}
	return
}
