// session

package websocket

import (
	"github.com/gorilla/websocket"
	"github.com/laconiz/eros/logis"
	"github.com/laconiz/eros/network"
	"github.com/laconiz/eros/network/cipher"
	"github.com/laconiz/eros/network/message"
	"github.com/laconiz/eros/network/queue"
	"github.com/laconiz/eros/network/session"
	"sync"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------

func newSession(conn *websocket.Conn, addr string, option *SessionOption, logger logis.Logger) *Session {
	id := session.Increment()
	return &Session{
		id:      id,
		addr:    addr,
		conn:    conn,
		option:  option,
		queue:   queue.New(option.QueueLen),
		logger:  logger.Field(network.FieldSession, id),
		encoder: option.Encoder.New(),
		cipher:  option.Cipher.New(),
	}
}

// ---------------------------------------------------------------------------------------------------------------------

type Session struct {
	id       session.ID      // session ID
	addr     string          // 连接地址
	conn     *websocket.Conn // websocket连接
	option   *SessionOption  // 配置
	queue    *queue.Queue    // 写入队列
	logger   logis.Logger    // 日志
	encoder  message.Encoder // 编码器
	cipher   cipher.Cipher   // 加密器
	sync.Map                 // 附加信息
}

// ---------------------------------------------------------------------------------------------------------------------

func (session *Session) ID() session.ID {
	return session.id
}

func (session *Session) Addr() string {
	return session.addr
}

func (session *Session) Close() {
	session.queue.Close()
}

// ---------------------------------------------------------------------------------------------------------------------

func (session *Session) Send(msg interface{}) error {

	message, err := session.encoder.Encode(msg)
	if err != nil {
		session.logger.Data(msg).Err(err).Error("encoder encode error")
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

	session.conn.SetReadLimit(option.ReadLimit)

	for {

		session.conn.SetReadDeadline(time.Now().Add(option.Timeout))

		_, stream, err := session.conn.ReadMessage()
		if err != nil {
			logger.Err(err).Info("read stream error")
			break
		}

		raw, err := session.cipher.Decode(stream)
		if err != nil {
			logger.Data(raw).Err(err).Warn("cipher decode error")
			break
		}

		message, err := session.encoder.Decode(raw)
		if err != nil {
			logger.Data(raw).Err(err).Warn("encoder decode error")
			break
		}

		logger.Data(string(raw)).Debug("read message")
		session.invoke(&network.Event{Meta: message.Meta, Msg: message.Msg, Ses: session})
	}
}

func (session *Session) write() {

	logger := session.logger
	option := session.option

	for {

		session.conn.SetWriteDeadline(time.Now().Add(option.Timeout))

		raws, closed := session.queue.Pick()

		for _, raw := range raws {

			stream, err := session.cipher.Encode(raw.([]byte))
			if err != nil {
				logger.Data(raw).Err(err).Warn("cipher encode error")
				return
			}

			if err := session.conn.WriteMessage(websocket.BinaryMessage, stream); err != nil {
				logger.Err(err).Warn("write stream error")
				return
			}
		}

		if closed {
			return
		}
	}
}

func (session *Session) run(callback func(*Session)) {

	session.logger.Info("connected")

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
