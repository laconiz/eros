package logis

import (
	"fmt"
	"github.com/laconiz/eros/utils/mathe"
)

type Level int8

func (level Level) Valid() bool {
	switch level {
	case DEBUG, INFO, WARN, ERROR, FATAL:
		return true
	}
	return false
}

func (level Level) Grade() Grade {
	switch level {
	case DEBUG:
		return GradeDebug
	case INFO:
		return GradeInfo
	case WARN:
		return GradeWarn
	case ERROR:
		return GradeError
	case FATAL:
		return GradeFatal
	default:
		return Grade(fmt.Sprintf("unknown[%d]", level))
	}
}

type Grade string

func (grade Grade) Level() Level {
	switch grade {
	case GradeDebug:
		return DEBUG
	case GradeInfo:
		return INFO
	case GradeWarn:
		return WARN
	case GradeError:
		return ERROR
	case GradeFatal:
		return FATAL
	default:
		return INVALID
	}
}

const (
	GradeDebug = "debug"
	GradeInfo  = "info"
	GradeWarn  = "warn"
	GradeError = "error"
	GradeFatal = "fatal"
)

const (
	DEBUG Level = 1 << iota
	INFO
	WARN
	ERROR
	FATAL
	INVALID
)

func MinLevel(a, b Level) Level {
	return Level(mathe.MinInt8(int8(a), int8(b)))
}
