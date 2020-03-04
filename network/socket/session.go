package socket

import (
	"github.com/laconiz/eros/logis"
	"github.com/laconiz/eros/network"
	"github.com/laconiz/eros/network/cipher"
	"github.com/laconiz/eros/network/message"
	"github.com/laconiz/eros/network/queue"
	"github.com/laconiz/eros/network/session"
	"github.com/laconiz/eros/network/socket/packer"
	"net"
	"sync"
	"time"
)

// 生成一个session
func newSession(conn net.Conn, option *SessionOption, logger logis.Logger) *Session {
	id := session.Increment()
	return &Session{
		id:      id,
		conn:    conn,
		option:  option,
		queue:   queue.New(option.QueueLen),
		logger:  logger.Field(network.FieldSession, id),
		cipher:  option.Cipher.New(),
		encoder: option.Encoder.New(),
		packer:  option.Packer.New(),
	}
}

type Session struct {
	id      session.ID      // ID
	conn    net.Conn        // 连接
	option  *SessionOption  // 配置信息
	queue   *queue.Queue    // 发送队列
	data    sync.Map        // 附加数据
	logger  logis.Logger    // 日志接口
	encoder message.Encoder // 编码器
	cipher  cipher.Cipher   // 加密器
	packer  packer.Packer   // 包装器
}

func (session *Session) ID() session.ID {
	return session.id
}

func (session *Session) Addr() string {
	return session.conn.RemoteAddr().String()
}

func (session *Session) Send(msg interface{}) error {

	message, err := session.encoder.Encode(msg)
	if err != nil {
		return err
	}

	return session.queue.Add(message.Stream)
}

func (session *Session) SendRaw(raw []byte) error {
	return session.queue.Add(raw)
}

func (session *Session) Close() {
	session.queue.Close()
}

func (session *Session) Set(key, value interface{}) {
	session.data.Store(key, value)
}

func (session *Session) Get(key interface{}) (interface{}, bool) {
	return session.data.Load(key)
}

func (session *Session) read() {

	option := session.option

	for {

		session.conn.SetReadDeadline(time.Now().Add(option.Timeout))

		stream, err := session.packer.Decode(session.conn)
		if err != nil {
			session.logger.Err(err).Info("read stream error")
			return
		}

		raw, err := session.cipher.Decode(stream)
		if err != nil {
			session.logger.Err(err).Warn("cipher decode error")
			break
		}

		message, err := session.encoder.Decode(raw)
		if err != nil {
			session.logger.Err(err).Warn("encoder decode error")
			break
		}

		session.logger.Data(string(raw)).Debug("read message")
		session.invoke(&network.Event{Meta: message.Meta, Msg: message.Msg, Ses: session})
	}
}

func (session *Session) write() {

	option := session.option

	for {

		session.conn.SetWriteDeadline(time.Now().Add(option.Timeout))

		raws, exit := session.queue.Pick()
		for _, raw := range raws {

			stream, err := session.cipher.Encode(raw.([]byte))
			if err != nil {
				session.logger.Err(err).Warn("cipher encode error")
				goto BREAK
			}

			if err := session.packer.Encode(session.conn, stream); err != nil {
				session.logger.Err(err).Warn("write stream error")
				goto BREAK
			}

			session.logger.Data(string(raw.([]byte))).Debug("write message")
		}

		if exit {
			goto BREAK
		}
	}

BREAK:
}

func (session *Session) run(closeFunc func(*Session)) {

	session.logger.Info("connected")

	go func() {
		session.write()
		session.conn.Close()
	}()

	session.invoke(network.NewConnectedEvent(session))
	session.read()
	session.queue.Close()

	session.logger.Info("disconnected")
	closeFunc(session)
	session.invoke(network.NewDisconnectedEvent(session))
}

func (session *Session) invoke(event *network.Event) {
	defer func() {
		if err := recover(); err != nil {
			session.logger.Data(err).Error("invoke panic")
		}
	}()
	session.option.Invoker.Invoke(event)
}
