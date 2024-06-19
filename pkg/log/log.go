package log

import (
	"net/http"
	"os"

	"go.elastic.co/apm/module/apmzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	encodingJSON    = "json"
	encodingConsole = "console"

	stdout = "stdout"
	stderr = "stderr"

	keyTime       = "time"
	keyLevel      = "severity"
	keyName       = "name"
	keyCaller     = "caller"
	keyMsg        = "logMsg"
	keyStacktrace = "stacktrace"
)

type Logger struct {
	*zap.SugaredLogger
	zapLogger *zap.Logger

	mode          Mode
	sensitiveData bool // sensitive data in log, currently, support for gorm
	enableAPM     bool
	enableDebug   bool
	traceIdKey    string

	testingMode bool // using for testing

	encoderConfig *zapcore.EncoderConfig
}

var l *Logger

func init() {
	if os.Getenv("LOGGER_MODE") == "prod" {
		l = New(
			WithModeFromString("prod"),
			WithDebug(false),
		)
	} else {
		l = New(
			WithModeFromString("dev"),
			WithDebug(true),
		)
	}
}

func Log() *zap.SugaredLogger {
	return l.SugaredLogger
}

func Zap() *zap.Logger {
	return l.zapLogger
}

// WithZapOption clone logger with additional zap options
func WithZapOption(opts ...zap.Option) *Logger {
	clone := *l
	zapLogger := clone.newZapLogger(opts...)
	clone.SugaredLogger = zapLogger.Sugar()
	return &clone
}

func Close() {
	_ = l.Sync()
}

// Init initialize an logger for package usage
func Init(logger *Logger) {
	l = logger
}

// Env return env from current package logger
func Env() Mode {
	return l.mode
}

// Sensitive return if logging data needs sensitive or not
func Sensitive() bool {
	return l.sensitiveData
}

func TraceIdKey() string {
	return l.traceIdKey
}

// TestingMode setup testing mode
// This is useful for case of log.Fatal which call os.Exit as normal (could NOT be handled)
func TestingMode(testing bool) {
	l.testingMode = testing
}

// New initialize new logger instance
func New(options ...Option) *Logger {
	logger := &Logger{}
	for _, option := range options {
		option(logger)
	}
	zapLogger := logger.newZapLogger()
	logger.zapLogger = zapLogger
	logger.SugaredLogger = zapLogger.Sugar()
	return logger
}

func (logger *Logger) newZapLogger(opts ...zap.Option) *zap.Logger {
	zapConfig := logger.newZapConfig()
	zapOptions := []zap.Option{zap.AddCallerSkip(1)}
	if logger.enableAPM {
		zapOptions = append(zapOptions, zap.WrapCore((&apmzap.Core{}).WrapCore))
	}
	zapOptions = append(zapOptions, opts...)
	// build zap logger
	zapLogger, _ := zapConfig.Build(zapOptions...)
	return zapLogger
}

func (logger *Logger) newZapConfig() zap.Config {
	var (
		logLevel      zap.AtomicLevel
		isDevelopment bool
		sampling      *zap.SamplingConfig
		encodingType  string
	)

	if logger.enableDebug {
		logLevel = zap.NewAtomicLevelAt(zap.DebugLevel)
	} else {
		logLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	switch logger.mode {
	case ModeDevelopment:
		isDevelopment = true
		encodingType = encodingConsole
	case ModeProduction:
		isDevelopment = false
		sampling = &zap.SamplingConfig{
			Initial:    http.StatusOK,
			Thereafter: 50,
		}
		encodingType = encodingJSON
	}

	return zap.Config{
		Level:             logLevel,
		Development:       isDevelopment,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          sampling,
		Encoding:          encodingType,
		EncoderConfig:     logger.newZapEncoderConfig(),
		OutputPaths:       []string{stdout},
		ErrorOutputPaths:  []string{stderr},
		InitialFields:     nil,
	}
}

func (logger *Logger) newZapEncoderConfig() zapcore.EncoderConfig {
	if logger.encoderConfig != nil {
		return *logger.encoderConfig
	}
	var (
		encodeLevel    zapcore.LevelEncoder
		encodeDuration zapcore.DurationEncoder
	)

	switch logger.mode {
	case ModeDevelopment:
		encodeLevel = zapcore.LowercaseColorLevelEncoder
		encodeDuration = zapcore.StringDurationEncoder
	case ModeProduction:
		encodeLevel = encodeGcpLevel()
		encodeDuration = zapcore.SecondsDurationEncoder
	}

	return zapcore.EncoderConfig{
		MessageKey:     keyMsg,
		LevelKey:       keyLevel,
		TimeKey:        keyTime,
		NameKey:        keyName,
		CallerKey:      keyCaller,
		StacktraceKey:  keyStacktrace,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    encodeLevel,
		EncodeTime:     zapcore.RFC3339TimeEncoder,
		EncodeDuration: encodeDuration,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
}

func encodeGcpLevel() zapcore.LevelEncoder {
	return func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		switch l {
		case zapcore.DebugLevel:
			enc.AppendString("DEBUG")
		case zapcore.InfoLevel:
			enc.AppendString("INFO")
		case zapcore.WarnLevel:
			enc.AppendString("WARNING")
		case zapcore.ErrorLevel:
			enc.AppendString("ERROR")
		case zapcore.DPanicLevel:
			enc.AppendString("CRITICAL")
		case zapcore.PanicLevel:
			enc.AppendString("ALERT")
		case zapcore.FatalLevel:
			enc.AppendString("EMERGENCY")
		}
	}
}
