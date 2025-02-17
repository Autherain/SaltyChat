package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync/atomic"
	"time"

	resLogger "github.com/jirenius/go-res/logger"
)

type (
	Format   string
	LogLevel string
)

const (
	// Format constants
	JSONFormat Format = "json"
	TextFormat Format = "text"

	// Level constants
	TraceLevel LogLevel = "trace"
	DebugLevel LogLevel = "debug"
	InfoLevel  LogLevel = "info"
	WarnLevel  LogLevel = "warn"
	ErrorLevel LogLevel = "error"

	// Log attribute keys
	formattedMessageKey = "formattedMessage"
	messageKey          = "message"
	statusKey           = "status"
	timestampKey        = "timestamp"
)

// Interface definitions
type errorLogger interface {
	Error(msg string, attrs ...slog.Attr)
}

// Adapter implements the resLogger.Logger interface
type Adapter struct {
	logger *Logger
}

// Ensure Adapter implements resLogger.Logger interface
var _ resLogger.Logger = (*Adapter)(nil)

// Adapter methods for resgate compatibility
func (a *Adapter) Debugf(format string, v ...interface{})   { a.logger.Debugf(format, v...) }
func (a *Adapter) Errorf(format string, v ...interface{})   { a.logger.Errorf(format, v...) }
func (a *Adapter) Infof(format string, v ...interface{})    { a.logger.Infof(format, v...) }
func (a *Adapter) Tracef(format string, v ...interface{})   { a.logger.Tracef(format, v...) }
func (a *Adapter) Warningf(format string, v ...interface{}) { a.logger.Warnf(format, v...) }

func (f *Format) UnmarshalText(text []byte) error {
	switch string(text) {
	case string(JSONFormat), string(TextFormat):
		*f = Format(text)
		return nil
	default:
		*f = TextFormat
		return nil
	}
}

func (l *LogLevel) UnmarshalText(text []byte) error {
	switch string(text) {
	case string(TraceLevel), string(DebugLevel), string(InfoLevel), string(WarnLevel), string(ErrorLevel):
		*l = LogLevel(text)
		return nil
	default:
		*l = InfoLevel
		return nil
	}
}

type Config struct {
	Format    Format
	Level     LogLevel
	AddSource bool
}

// Logger is the main logger struct that handles both standard and RES logging
type Logger struct {
	slog    *slog.Logger
	Adapter *Adapter
}

var defaultLogger atomic.Pointer[Logger]

func init() {
	defaultLogger.Store(NewDefault())
}

// Default returns the default logger.
func Default() *Logger {
	return defaultLogger.Load()
}

// NewLogger creates a new configured logger that can be used for both standard and RES logging
func NewLogger(cfg Config) *Logger {
	level := convertLevel(cfg.Level)
	handler := createHandler(cfg, level)
	logger := &Logger{slog: slog.New(handler)}
	logger.Adapter = &Adapter{logger: logger}
	return logger
}

func NewDefault() *Logger {
	defaultConfig := Config{
		Format:    TextFormat,
		Level:     InfoLevel,
		AddSource: false,
	}
	return NewLogger(defaultConfig)
}

func convertLevel(level LogLevel) slog.Level {
	switch level {
	case TraceLevel, DebugLevel:
		return slog.LevelDebug
	case WarnLevel:
		return slog.LevelWarn
	case ErrorLevel:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func createHandler(cfg Config, level slog.Level) slog.Handler {
	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: cfg.AddSource,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.TimeKey:
				return slog.Attr{
					Key:   timestampKey,
					Value: slog.StringValue(time.Now().Format(time.RFC3339)),
				}
			case slog.LevelKey:
				return slog.Attr{
					Key:   statusKey,
					Value: slog.StringValue(strings.ToLower(a.Value.String())),
				}
			case slog.MessageKey:
				return slog.Attr{
					Key:   messageKey,
					Value: a.Value,
				}
			}
			return a
		},
	}

	if cfg.Format == JSONFormat {
		return slog.NewJSONHandler(os.Stdout, opts)
	}
	return slog.NewTextHandler(os.Stdout, opts)
}

// Helper function for formatting messages
func formatMessage(message string) string {
	return strings.ToLower(strings.ReplaceAll(message, " ", "-"))
}

// Standard logging methods with ...any args
func (l *Logger) Trace(msg string, args ...any)     { l.slog.Debug(msg, args...) }
func (l *Logger) Debug(msg string, args ...any)     { l.slog.Debug(msg, args...) }
func (l *Logger) Info(msg string, args ...any)      { l.slog.Info(msg, args...) }
func (l *Logger) Warn(msg string, args ...any)      { l.slog.Warn(msg, args...) }
func (l *Logger) ErrorArgs(msg string, args ...any) { l.slog.Error(msg, args...) }

// Interface-compatible Error method
func (l *Logger) Error(msg string, attrs ...slog.Attr) {
	l.slog.LogAttrs(context.TODO(), slog.LevelError, msg, attrs...)
}

// RES-compatible logging methods
func (l *Logger) Tracef(format string, v ...interface{}) { l.Trace(fmt.Sprintf(format, v...)) }
func (l *Logger) Debugf(format string, v ...interface{}) { l.Debug(fmt.Sprintf(format, v...)) }
func (l *Logger) Infof(format string, v ...interface{})  { l.Info(fmt.Sprintf(format, v...)) }
func (l *Logger) Warnf(format string, v ...interface{})  { l.Warn(fmt.Sprintf(format, v...)) }
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.ErrorArgs(fmt.Sprintf(format, v...))
}

// Attribute-based logging methods
func (l *Logger) TraceAttrs(msg string, attrs ...slog.Attr) {
	l.slog.LogAttrs(context.TODO(), slog.LevelDebug, msg, attrs...)
}

func (l *Logger) DebugAttrs(msg string, attrs ...slog.Attr) {
	l.slog.LogAttrs(context.TODO(), slog.LevelDebug, msg, attrs...)
}

func (l *Logger) InfoAttrs(msg string, attrs ...slog.Attr) {
	l.slog.LogAttrs(context.TODO(), slog.LevelInfo, msg, attrs...)
}

func (l *Logger) WarnAttrs(msg string, attrs ...slog.Attr) {
	l.slog.LogAttrs(context.TODO(), slog.LevelWarn, msg, attrs...)
}

func (l *Logger) ErrorAttrs(msg string, attrs ...slog.Attr) {
	l.Error(msg, attrs...)
}

// Formatted logging methods with Datadog integration
func (l *Logger) FormattedDebug(msg string, attrs ...slog.Attr) {
	l.DebugAttrs(msg, append(attrs,
		slog.String(formattedMessageKey, formatMessage(msg)))...)
}

func (l *Logger) FormattedInfo(msg string, attrs ...slog.Attr) {
	l.InfoAttrs(msg, append(attrs,
		slog.String(formattedMessageKey, formatMessage(msg)))...)
}

func (l *Logger) FormattedWarn(msg string, attrs ...slog.Attr) {
	l.WarnAttrs(msg, append(attrs,
		slog.String(formattedMessageKey, formatMessage(msg)))...)
}

func (l *Logger) FormattedError(msg string, attrs ...slog.Attr) {
	l.Error(msg, append(attrs,
		slog.String(formattedMessageKey, formatMessage(msg)))...)
}

// Chainable context methods
func (l *Logger) With(attrs ...slog.Attr) *Logger {
	newLogger := &Logger{
		slog: slog.New(l.slog.Handler().WithAttrs(attrs)),
	}
	newLogger.Adapter = &Adapter{logger: newLogger}
	return newLogger
}

func (l *Logger) WithGroup(name string) *Logger {
	newLogger := &Logger{
		slog: slog.New(l.slog.Handler().WithGroup(name)),
	}
	newLogger.Adapter = &Adapter{logger: newLogger}
	return newLogger
}

// SlogLogger returns the underlying slog.Logger if needed
func (l *Logger) SlogLogger() *slog.Logger {
	return l.slog
}
