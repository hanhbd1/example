package log

import "go.uber.org/zap/zapcore"

type Option func(*Logger)

// WithModeFromString get logger mode from string
func WithModeFromString(value string) Option {
	return func(logger *Logger) {
		m := GetModeFromString(value)
		if m == ModeUnknown {
			panic("unknown logger mode")
		}
		logger.mode = m
	}
}

// WithAPM add APM agent
func WithAPM(enable bool) Option {
	return func(logger *Logger) {
		logger.enableAPM = enable
	}
}

// WithDebug enable debug level of log
func WithDebug(enable bool) Option {
	return func(logger *Logger) {
		logger.enableDebug = enable
	}
}

// WithEncoderConfig custom zap encoder config
func WithEncoderConfig(cfg *zapcore.EncoderConfig) Option {
	return func(logger *Logger) {
		logger.encoderConfig = cfg
	}
}

// EnableSensitive decide to enable sensitive log data or not
func EnableSensitive(enable bool) Option {
	return func(logger *Logger) {
		logger.sensitiveData = enable
	}
}

func WithTraceIdKey(traceIdKey string) Option {
	return func(logger *Logger) {
		logger.traceIdKey = traceIdKey
	}
}
