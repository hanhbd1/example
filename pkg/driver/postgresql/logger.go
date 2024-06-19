package postgresql

import (
	"bytes"
	"context"
	"database/sql/driver"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"example/pkg/log"

	"go.uber.org/zap"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

var (
	sqlRegexp                = regexp.MustCompile(`\?`)
	numericPlaceHolderRegexp = regexp.MustCompile(`\$\d+`)
)

func isPrintable(s string) bool {
	for _, r := range s {
		if !unicode.IsPrint(r) {
			return false
		}
	}
	return true
}

func trimCallerPath(s string) string {
	// Find the last separator.
	//
	idx := strings.LastIndexByte(s, '/')
	if idx == -1 {
		return s
	}
	// Find the penultimate separator.
	idx = strings.LastIndexByte(s[:idx], '/')
	if idx == -1 {
		return s
	}
	buf := &bytes.Buffer{}
	// Keep everything after the penultimate separator.
	buf.WriteString(s[idx+1:])
	caller := buf.String()
	buf.Reset()
	return caller
}

func (l *gormLogger) logFormatter(values ...interface{}) (source string, messages []interface{}) {
	if len(values) > 1 {
		var (
			sql             string
			formattedValues []string
			level           = values[0]
		)

		source = trimCallerPath(fmt.Sprintf("%v", values[1]))

		if level == "sql" {
			// duration
			messages = append(messages,
				"duration", fmt.Sprintf("%.2fms", float64(values[2].(time.Duration).Nanoseconds()/1e4)/100.0))
			// sql
			for _, value := range values[4].([]interface{}) {
				indirectValue := reflect.Indirect(reflect.ValueOf(value))
				if indirectValue.IsValid() {
					value = indirectValue.Interface()
					if t, ok := value.(time.Time); ok {
						if t.IsZero() {
							formattedValues = append(formattedValues, l.sensitive(fmt.Sprintf("'%v'", "0000-00-00 00:00:00")))
						} else {
							formattedValues = append(formattedValues, l.sensitive(fmt.Sprintf("'%v'", t.Format("2006-01-02 15:04:05"))))
						}
					} else if b, ok := value.([]byte); ok {
						if str := string(b); isPrintable(str) {
							formattedValues = append(formattedValues, l.sensitive(fmt.Sprintf("'%v'", str)))
						} else {
							formattedValues = append(formattedValues, l.sensitive("'<binary>'"))
						}
					} else if r, ok := value.(driver.Valuer); ok {
						if value, err := r.Value(); err == nil && value != nil {
							formattedValues = append(formattedValues, l.sensitive(fmt.Sprintf("'%v'", value)))
						} else {
							formattedValues = append(formattedValues, l.sensitive("NULL"))
						}
					} else {
						switch value.(type) {
						case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool:
							formattedValues = append(formattedValues, l.sensitive(fmt.Sprintf("%v", value)))
						default:
							formattedValues = append(formattedValues, l.sensitive(fmt.Sprintf("'%v'", value)))
						}
					}
				} else {
					formattedValues = append(formattedValues, l.sensitive("NULL"))
				}
			}

			// differentiate between $n placeholders or else treat like ?
			if numericPlaceHolderRegexp.MatchString(values[3].(string)) {
				sql = values[3].(string)
				for index, value := range formattedValues {
					placeholder := fmt.Sprintf(`\$%d([^\d]|$)`, index+1)
					sql = regexp.MustCompile(placeholder).ReplaceAllString(sql, value+"$1")
				}
			} else {
				formattedValuesLength := len(formattedValues)
				for index, value := range sqlRegexp.Split(values[3].(string), -1) {
					sql += value
					if formattedValues != nil && index < formattedValuesLength {
						sql += formattedValues[index]
					}
				}
			}
			messages = append(messages, "sql", sql)
			messages = append(messages, "result",
				fmt.Sprintf("%v", strconv.FormatInt(values[5].(int64), 10)+" rows affected or returned"))
		} else {
			messages = append(messages, "other", values[2:])
		}
	}

	return
}

func (l *gormLogger) sensitive(s string) string {
	if l.enableSensitive {
		return "***"
	}

	return s
}

type zapOutputEncoder func(caller string, messages ...interface{}) (string, []interface{})

// due to zap dont allow passing callerEntry from external
// workaround by omit caller
// - add caller key-value for case of json encode
// - concat string with message field in case of console encode
var consoleEncode = func(caller string, messages ...interface{}) (string, []interface{}) {
	return caller + "\tgorm query", messages
}

var jsonEncode = func(caller string, messages ...interface{}) (string, []interface{}) {
	var msgs = []interface{}{"caller", caller}
	msgs = append(msgs, messages...)
	return "gorm query", msgs
}

type gormLogger struct {
	log             *log.Logger
	enc             zapOutputEncoder
	enableSensitive bool
	LogLevel        logger.LogLevel
}

// implement log v2 but not sure it work :))
func (l *gormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

func (l *gormLogger) Info(ctx context.Context, s string, v ...interface{}) {
	if l.LogLevel >= logger.Info {
		caller := trimCallerPath(utils.FileWithLineNum())
		messages := fmt.Sprintf(s, v...)
		msg, kvs := l.enc(caller, "other", messages)
		l.log.Infow(msg, kvs...)
	}
}

func (l *gormLogger) Warn(ctx context.Context, s string, v ...interface{}) {
	if l.LogLevel >= logger.Warn {
		caller := trimCallerPath(utils.FileWithLineNum())
		messages := fmt.Sprintf(s, v...)
		msg, kvs := l.enc(caller, "other", messages)
		l.log.Warnw(msg, kvs...)
	}
}

func (l *gormLogger) Error(ctx context.Context, s string, v ...interface{}) {
	if l.LogLevel >= logger.Error {
		caller := trimCallerPath(utils.FileWithLineNum())
		messages := fmt.Sprintf(s, v...)
		msg, kvs := l.enc(caller, "other", messages)
		l.log.Errorw(msg, kvs...)
	}
}

func (l *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel >= logger.Info {
		elapsed := time.Since(begin)
		sql, rows := fc()
		caller := trimCallerPath(utils.FileWithLineNum())
		var messages []interface{}
		messages = append(messages,
			"duration", fmt.Sprintf("%.2fms", float64(elapsed.Nanoseconds()/1e4)/100.0))
		sql = l.sensitive(sql)
		messages = append(messages, "sql", sql)
		messages = append(messages, "result",
			fmt.Sprintf("%v", strconv.FormatInt(rows, 10)+" rows affected or returned"))
		if err != nil {
			messages = append(messages, "result", err)
		}
		msg, kvs := l.enc(caller, messages...)
		l.log.Infow(msg, kvs...)
	}
}

// newLogger is a gorm logger to convert gorm log format into zap's log encoder format, must be initialized
// after initializing log package instance
func newLogger(env log.Mode, sensitive bool) *gormLogger {
	l := log.WithZapOption(zap.WithCaller(false))
	var enc zapOutputEncoder = consoleEncode
	logLevel := logger.Info
	if env == log.ModeProduction {
		// format console
		enc = jsonEncode
		logLevel = logger.Error
	}

	return &gormLogger{
		log:             l,
		enc:             enc,
		enableSensitive: sensitive,
		LogLevel:        logLevel,
	}
}

func (l *gormLogger) Print(v ...interface{}) {
	caller, messages := l.logFormatter(v...)
	msg, kvs := l.enc(caller, messages...)
	l.log.Infow(msg, kvs...)
}
