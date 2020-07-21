package database

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"regexp"
	"runtime"
	"time"
	"unicode"

	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog"
)

// GormLogger is a custom logger for Gorm, making it use logrus.
type GormLogger struct {
	logger     zerolog.Logger
	WithColor  bool
	WithCaller bool
}

// Print handles log events from Gorm for the custom logger.
func (gl *GormLogger) Print(v ...interface{}) {
	messages := LogFormatter(v...)
	gormType := messages[0]

	var callers string
	pc := make([]uintptr, 10) // at least 1 entry needed
	n := runtime.Callers(8, pc)
	for i := 1; i < n; i++ {
		f := runtime.FuncForPC(pc[i])
		file, line := f.FileLine(pc[i])
		callers += fmt.Sprintf("\n%s:%d", file, line)
	}
	if gormType == "sql" {
		src := messages[1]
		latency := messages[3].(time.Duration)
		sql := messages[4]
		rowsAffected := messages[5]
		latencyStr := latency.String()
		srcStr := src
		sqlStr := sql
		rowsAffectedStr := rowsAffected
		if gl.WithColor {
			latencyStr = fmt.Sprintf("\033[33m[ %v ]\033[0m ", latency.String())
			srcStr = fmt.Sprintf("\033[35m[ %v ]\033[0m ", src)
			sqlStr = fmt.Sprintf("\033[34m[\n %v \n]\033[0m ", sql)
			rowsAffectedStr = fmt.Sprintf("\033[36;31m[ %v %s ]\033[0m ", rowsAffected, "rows affected or returned")
		}
		fields := map[string]interface{}{
			"gorm_type":     gormType,
			"src":           src,
			"latency":       int64(latency),
			"latency_human": latency.String(),
			"rows_affected": rowsAffected,
		}
		if gl.WithCaller {
			fields["caller"] = callers
		}
		gl.logger.Debug().Fields(fields).Msgf("\n%s\n%s\n%s\n%s\n", latencyStr, srcStr, sqlStr, rowsAffectedStr)
	} else {
		gl.logger.Debug().Fields(map[string]interface{}{
			"gorm_type": gormType,
		}).Msgf("%s", messages[1:]...)
	}
}

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

// LogFormatter for log formatter
var LogFormatter = func(values ...interface{}) (messages []interface{}) {
	if len(values) > 1 {
		var (
			sql             string
			formattedValues []string
			level           = values[0]
			currentTime     = "\n\033[33m[" + gorm.NowFunc().Format("2006-01-02 15:04:05") + "]\033[0m"
			source          = fmt.Sprintf(" %v ", values[1])
		)

		messages = []interface{}{level, source, currentTime}

		if level == "sql" {
			// duration
			messages = append(messages, values[2].(time.Duration))
			// sql

			for _, value := range values[4].([]interface{}) {
				indirectValue := reflect.Indirect(reflect.ValueOf(value))
				if indirectValue.IsValid() {
					value = indirectValue.Interface()
					if t, ok := value.(time.Time); ok {
						formattedValues = append(formattedValues, fmt.Sprintf("'%v'", t.Format("2006-01-02 15:04:05")))
					} else if b, ok := value.([]byte); ok {
						if str := string(b); isPrintable(str) {
							formattedValues = append(formattedValues, fmt.Sprintf("'%v'", str))
						} else {
							formattedValues = append(formattedValues, "'<binary>'")
						}
					} else if r, ok := value.(driver.Valuer); ok {
						if value, err := r.Value(); err == nil && value != nil {
							formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
						} else {
							formattedValues = append(formattedValues, "NULL")
						}
					} else {
						switch t := value.(type) {
						case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool:
							formattedValues = append(formattedValues, fmt.Sprintf("%v", value))
						case []byte:
							formattedValues = append(formattedValues, fmt.Sprintf("'%v'", string(t)))
						default:
							formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
						}
					}
				} else {
					formattedValues = append(formattedValues, "NULL")
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
					if index < formattedValuesLength {
						sql += formattedValues[index]
					}
				}
			}

			messages = append(messages, sql)
			messages = append(messages, values[5].(int64)) // rowsAffected
		} else {
			messages = append(messages, "\033[31;1m")
			messages = append(messages, values[2:]...)
			messages = append(messages, "\033[0m")
		}
	}

	return
}
