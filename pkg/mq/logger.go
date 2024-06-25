package mq

import (
	"example/pkg/log"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/samber/lo"
)

type Logger struct {
	ignoreMsg map[watermill.LogLevel][]string
	level     watermill.LogLevel
	fields    watermill.LogFields
}

func NewLogger() watermill.LoggerAdapter {
	level := watermill.InfoLogLevel
	ignoreMsg := map[watermill.LogLevel][]string{
		watermill.ErrorLogLevel: {
			"Handler returned error",
		},
		watermill.InfoLogLevel:  {},
		watermill.DebugLogLevel: {},
		watermill.TraceLogLevel: {},
	}
	return &Logger{level: level, ignoreMsg: ignoreMsg}
}

func (l *Logger) Error(msg string, err error, fields watermill.LogFields) {
	if l.level <= watermill.ErrorLogLevel {
		fields = l.fields.Add(fields)
		messages := []interface{}{}
		for k, v := range fields {
			messages = append(messages, k, v)
		}
		messages = append(messages, "error", err)
		if !lo.Contains(l.ignoreMsg[watermill.ErrorLogLevel], msg) {
			log.Errorw(msg, messages...)
		}
	}
}

func (l *Logger) Info(msg string, fields watermill.LogFields) {
	if l.level <= watermill.InfoLogLevel {
		fields = l.fields.Add(fields)
		messages := []interface{}{}
		for k, v := range fields {
			messages = append(messages, k, v)
		}
		if !lo.Contains(l.ignoreMsg[watermill.InfoLogLevel], msg) {
			log.Infow(msg, messages...)
		}
	}
}

func (l *Logger) Debug(msg string, fields watermill.LogFields) {
	if l.level <= watermill.DebugLogLevel {
		fields = l.fields.Add(fields)
		messages := []interface{}{}
		for k, v := range fields {
			messages = append(messages, k, v)
		}
		if !lo.Contains(l.ignoreMsg[watermill.DebugLogLevel], msg) {
			log.Debugw(msg, messages...)
		}
	}
}

func (l *Logger) Trace(msg string, fields watermill.LogFields) {
	if l.level <= watermill.TraceLogLevel {
		fields = l.fields.Add(fields)
		messages := []interface{}{}
		for k, v := range fields {
			messages = append(messages, k, v)
		}
		if !lo.Contains(l.ignoreMsg[watermill.TraceLogLevel], msg) {
			log.Debugw(msg, messages...)
		}
	}
}

func (l *Logger) With(fields watermill.LogFields) watermill.LoggerAdapter {
	return &Logger{
		level:  l.level,
		fields: l.fields.Add(fields),
	}
}
