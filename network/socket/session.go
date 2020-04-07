package socket

import (
	"github.com/laconiz/eros/logis"
	"github.com/laconiz/eros/network"
	"github.com/laconiz/eros/network/cipher"
	"github.com/laconiz/eros/network/encoder"
	"github.com/laconiz/eros/network/queue"
	"github.com/laconiz/eros/network/session"
	"github.com/laconiz/eros/network/socket/reader"
	"net"
	"sync"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------

func newSession(conn net.Conn, option *SessionOption, logger logis.Logger) *Session {

	id := session.Increment()

	return &Session{
		id:      id,
		conn:    conn,
		option:  option,
		queue:   queue.New(option.Queue),
		logger:  logger.Field(network.FieldSession, id),
		cipher:  option.Cipher.New(),
		encoder: option.Encoder.New(),
		reader:  option.Reader.New(),
	}
}

// ---------------------------------------------------------------------------------------------------------------------

type Session struct {
	id       session.ID      // ID
	conn     net.Conn        // 连接
	option   *SessionOption  // 配置
	queue    *queue.Queue    // 发送队列
	logger   logis.Logger    // 日志接口
	reader   reader.Reader   // 包装器
	cipher   cipher.Cipher   // 加密器
	encoder  encoder.Encoder // 编码器
	sync.Map                 // 附加信息
}

// ---------------------------------------------------------------------------------------------------------------------

func (session *Session) ID() session.ID {
	return session.id
}

func (session *Session) Addr() string {
	return session.conn.RemoteAddr().String()
}

func (session *Session) Close() {
	session.queue.Close()
}

// ---------------------------------------------------------------------------------------------------------------------

func (session *Session) Send(msg interface{}) error {

	message, err := session.encoder.Marshal(msg)
	if err != nil {
		session.logger.Data(msg).Err(err).Error("marshal error")
		return err
	}

	return session.SendRaw(message.Stream)
}

func (session *Session) SendRaw(raw []byte) error {
	return session.queue.Add(raw)
}

// ---------------------------------------------------------------------------------------------------------------------

func (session *Session) read() {

	logger := session.logger
	option := session.option

	for {

		deadline := time.Now().Add(option.Timeout)
		session.conn.SetReadDeadline(deadline)

		stream, err := session.reader.Read(session.conn)
		if err != nil {
			logger.Err(err).Info("read error")
			return
		}

		raw, err := session.cipher.Decode(stream)
		if err != nil {
			logger.Err(err).Warn("decode error")
			break
		}

		message, err := session.encoder.Unmarshal(raw)
		if err != nil {
			logger.Err(err).Warn("unmarshal error")
			break
		}

		logger.Data(string(raw)).Debug("recv message")

		session.invoke(&network.Event{
			Meta: message.Meta,
			Msg:  message.Msg,
			Ses:  session,
		})
	}
}

func (session *Session) write() {

	logger := session.logger
	option := session.option

	for {

		deadline := time.Now().Add(option.Timeout)
		session.conn.SetWriteDeadline(deadline)

		events, closed := session.queue.Pick()

		for _, event := range events {

			raw := event.([]byte)

			stream, err := session.cipher.Encode(raw)
			if err != nil {
				logger.Data(raw).Err(err).Warn("encode error")
				return
			}

			if err := session.reader.Write(session.conn, stream); err != nil {
				logger.Data(stream).Err(err).Warn("write error")
				return
			}

			session.logger.Data(string(raw)).Debug("send message")
		}

		if closed {
			return
		}
	}
}

func (session *Session) run(callback func(*Session)) {

	session.logger.Data(session.Addr()).Info("connected")

	go func() {
		session.write()
		session.conn.Close()
	}()

	session.invoke(network.NewConnectedEvent(session))
	session.read()
	session.queue.Close()

	session.logger.Info("disconnected")
	callback(session)
	session.invoke(network.NewDisconnectedEvent(session))
}

// ---------------------------------------------------------------------------------------------------------------------

func (session *Session) invoke(event *network.Event) {

	defer func() {
		if err := recover(); err != nil {
			session.logger.Data(err).Error("invoke panic")
		}
	}()

	session.option.Invoker.Invoke(event)
}
