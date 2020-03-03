package logis

import (
	"github.com/laconiz/eros/logis/context"
	"time"
)

type Log struct {
	Level   Level            // 等级
	Message string           // 文本内容
	Time    time.Time        // 时间
	Value   interface{}      // 即时数据
	Context *context.Context // 上下文
}
