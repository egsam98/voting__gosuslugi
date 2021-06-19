package dbext

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
	sqldblogger "github.com/simukti/sqldb-logger"
)

const (
	LongDataFields = `query args`
	MaxDataLen     = 200
)

var _ sqldblogger.Logger = (*ZeroLogger)(nil)

// ZeroLogger is zerolog adapter for sqldblogger.Logger
type ZeroLogger struct {
	log zerolog.Logger
}

func NewZeroLogger(log zerolog.Logger) *ZeroLogger {
	return &ZeroLogger{log: log}
}

func (z *ZeroLogger) Log(_ context.Context, level sqldblogger.Level, msg string, data map[string]interface{}) {
	var lvl zerolog.Level
	switch level {
	case sqldblogger.LevelError:
		lvl = zerolog.ErrorLevel
	case sqldblogger.LevelInfo:
		lvl = zerolog.InfoLevel
	case sqldblogger.LevelDebug, sqldblogger.LevelTrace:
		fallthrough
	default:
		lvl = zerolog.DebugLevel
	}
	z.trimLongData(data)
	z.log.WithLevel(lvl).Fields(data).Msg(msg)
}

// trimLongData trims long data values with keys in "LongDataFields"
func (z *ZeroLogger) trimLongData(data map[string]interface{}) {
	for k, v := range data {
		if !strings.Contains(LongDataFields, k) {
			continue
		}

		str, ok := v.(string)
		if !ok {
			str = fmt.Sprintf("%v", v)
		}
		if length := len(str); length > MaxDataLen {
			data[k] = str[:MaxDataLen] + "... (" + strconv.Itoa(length) + " symbols)"
		}
	}
}
