package logging

import (
	"sync/atomic"
)

type Logger struct {
	level Level
}

func (logger *Logger) SetLevel(level Level) {
	atomic.StoreInt32((*int32)(&logger.level), int32(level))
}

func (logger *Logger) GetLevel() Level {
	return Level(atomic.LoadInt32((*int32)(&logger.level)))
}

func (logger *Logger) IsLevelDisabled(level Level) bool {
	return logger.GetLevel() > level
}
