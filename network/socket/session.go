package socket

import (
	"github.com/laconiz/eros/holder/message"
	"github.com/laconiz/eros/holder/queue"
	"github.com/laconiz/eros/logis"
	"github.com/laconiz/eros/network"
	"github.com/laconiz/eros/network/cipher"
	"github.com/laconiz/eros/network/session"
	"github.com/laconiz/eros/network/socket/packer"
	"net"
	"sync"
	"time"
)

// 生成一个session
func newSession(conn net.Conn, opt *SesOption, log logis.Logger) *Session {
	id := session.ID(session.Increment())
	return &Session{
		id:      id,
		conn:    conn,
		opt:     opt,
		queue:   queue.New(opt.QueueLen),
		log:     log.Field(network.FieldSession, id),
		cipher:  opt.Cipher.New(),
		encoder: opt.Encoder.New(),
		packer:  opt.Packer.New(),
	}
}

type Session struct {
	id      session.ID      // ID
	conn    net.Conn        // 连接
	opt     *SesOption      // 配置信息
	queue   *queue.Queue    // 发送队列
	data    sync.Map        // 附加数据
	log     logis.Logger    // 日志接口
	encoder message.Encoder // 编码器
	cipher  cipher.Cipher   // 加密器
	packer  packer.Packer   // 包装器
}

func (ses *Session) ID() session.ID {
	return ses.id
}

func (ses *Session) Addr() string {
	return ses.conn.RemoteAddr().String()
}

func (ses *Session) Send(msg interface{}) error {
	message, err := ses.encoder.Encode(msg)
	if err != nil {
		return err
	}
	return ses.queue.Add(message.Stream)
}

func (ses *Session) SendRaw(raw []byte) error {
	return ses.queue.Add(raw)
}

func (ses *Session) Close() {
	ses.queue.Close()
}

func (ses *Session) Set(key, value interface{}) {
	ses.data.Store(key, value)
}

func (ses *Session) Get(key interface{}) (interface{}, bool) {
	return ses.data.Load(key)
}

func (ses *Session) read() {

	opt := ses.opt

	for {

		ses.conn.SetReadDeadline(time.Now().Add(opt.Timeout))

		stream, err := ses.packer.Decode(ses.conn)
		if err != nil {
			ses.log.Infof("read stream error: %v", err)
			return
		}

		raw, err := ses.cipher.Decode(stream)
		if err != nil {
			ses.log.Warnf("cipher.decode error: %v", err)
			break
		}

		message, err := ses.encoder.Decode(raw)
		if err != nil {
			ses.log.Warnf("encoder.decode error: %v", err)
			break
		}

		ses.log.Debugf("read: %v", string(raw))
		ses.invoke(&network.Event{Meta: message.Meta, Msg: message.Msg, Ses: ses})
	}

	ses.log.Info("read loop break")
}

// 读取线程
func (ses *Session) write() {

	opt := ses.opt

	for {

		ses.conn.SetWriteDeadline(time.Now().Add(opt.Timeout))

		raws, ok := ses.queue.Pick()
		for _, raw := range raws {

			stream, err := ses.cipher.Encode(raw.([]byte))
			if err != nil {
				ses.log.Warnf("cipher.encode error: %v", err)
				goto BREAK
			}

			if err := ses.packer.Encode(ses.conn, stream); err != nil {
				ses.log.Warnf("write stream error: %v", err)
				goto BREAK
			}

			ses.log.Debugf("write: %v", string(raw.([]byte)))
		}

		if ok {
			goto BREAK
		}
	}

BREAK:
	ses.log.Info("write loop break")
}

func (ses *Session) run(closeFunc func(*Session)) {

	ses.log.Info("connected")

	go func() {
		ses.write()
		ses.conn.Close()
	}()

	ses.invoke(network.NewConnectedEvent(ses))
	ses.read()
	ses.queue.Close()

	ses.log.Info("disconnected")
	closeFunc(ses)
	ses.invoke(network.NewDisconnectedEvent(ses))
}

func (ses *Session) invoke(event *network.Event) {
	defer func() {
		if err := recover(); err != nil {
			ses.log.Errorf("invoke panic: %v", err)
		}
	}()
	ses.opt.Invoker.Invoke(event)
}
