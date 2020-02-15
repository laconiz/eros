// socket配置信息

package socket

import (
	"github.com/laconiz/eros/network/cipher"
	"github.com/laconiz/eros/network/encoder"
	"github.com/laconiz/eros/network/socket/packer"
	"time"

	"github.com/laconiz/eros/network"
)

const module = "socket"

// session配置
type SesOption struct {
	Timeout  time.Duration   // 超时
	QueueLen int             // 写缓冲区长度
	Invoker  network.Invoker // 回调接口
	Encoder  encoder.Maker   // 编码器
	Cipher   cipher.Maker    // 加密器
	Packer   packer.Maker    // 包装器
}

func (opt *SesOption) parse() {
	if opt.Timeout <= 0 {
		opt.Timeout = time.Second * 15
	}
	if opt.QueueLen <= 0 {
		opt.QueueLen = 32
	}
	if opt.Invoker == nil {
		opt.Invoker = network.NewStdInvoker()
	}
	if opt.Encoder == nil {
		opt.Encoder = encoder.NewIDMaker()
	}
	if opt.Cipher == nil {
		opt.Cipher = cipher.NewEmptyMaker()
	}
	if opt.Packer == nil {
		opt.Packer = packer.NewSizeMaker()
	}
}

// acceptor配置
type AccOption struct {
	Name    string
	Addr    string
	Session SesOption
}

func (opt *AccOption) parse() {
	if opt.Name == "" {
		opt.Name = "acceptor"
	}
	opt.Session.parse()
}

// connector配置
type ConnOption struct {
	Name      string          // 名称
	Addr      string          // 连接地址
	Reconnect bool            // 是否重连
	Delays    []time.Duration // 重连延迟
	Session   SesOption       // session配置
}

func (opt *ConnOption) parse() {
	if opt.Name == "" {
		opt.Name = "connector"
	}
	if len(opt.Delays) == 0 {
		opt.Delays = []time.Duration{
			time.Millisecond * 10,
			time.Millisecond * 200,
			time.Millisecond * 800,
			time.Millisecond * 1200,
			time.Millisecond * 2000,
			time.Millisecond * 3600,
			time.Millisecond * 6000,
			time.Millisecond * 9000,
			time.Millisecond * 15000,
		}
	}
	opt.Session.parse()
}
