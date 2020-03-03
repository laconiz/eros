package logis

import (
	"fmt"
)

type Level int8

func (level Level) Valid() bool {
	switch level {
	case DEBUG, INFO, WARN, ERROR, FATAL:
		return true
	}
	return false
}

func (level Level) Enable(other Level) bool {
	return level <= other
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
