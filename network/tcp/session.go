package tcp

import (
	"fmt"
	"github.com/laconiz/eros/log"
	"github.com/laconiz/eros/network"
	"github.com/laconiz/eros/queue"
	"net"
	"sync"
	"time"
)

type Session struct {
	id      network.SessionID // session id
	conn    net.Conn          // 连接
	config  *SessionConfig    // 配置信息
	encoder Encoder           // 编码器
	queue   *queue.Queue      // 发送队列
	data    sync.Map          // 附加数据
	logger  *log.Logger       // 日志接口
}

func (ses *Session) ID() network.SessionID {
	return ses.id
}

func (ses *Session) Addr() string {
	return ses.conn.RemoteAddr().String()
}

func (ses *Session) Send(msg interface{}) error {
	return ses.queue.Add(&network.Event{Msg: msg})
}

func (ses *Session) SendStream(stream []byte) error {
	return ses.queue.Add(&network.Event{Stream: stream})
}

func (ses *Session) Close() {
	ses.queue.Close()
}

func (ses *Session) Set(key, value interface{}) {
	ses.data.Store(key, value)
}

func (ses *Session) Get(key interface{}) interface{} {
	if value, ok := ses.data.Load(key); ok {
		return value
	}
	return nil
}

func (ses *Session) read() {

	conf := ses.config

	for {

		ses.conn.SetReadDeadline(time.Now().Add(conf.ReadTimeout))

		event, err := ses.encoder.Encode(ses.conn)
		if err != nil {
			ses.logger.Warn(err)
			return
		}

		if ses.logger.LogLevel(log.Info) {
			ses.logger.Infof("read: %s", ses.encoder.String(event))
		}

		event.Session = ses
		ses.invoke(event)
	}
}

func (ses *Session) write() {

	conf := ses.config

	for {

		// 设置写超时
		ses.conn.SetWriteDeadline(time.Now().Add(conf.WriteTimeout))

		// 获取消息列表
		events, closed := ses.queue.Pick()

		for _, e := range events {

			event := e.(*network.Event)

			stream, err := ses.encoder.Decode(event)
			if err != nil {
				ses.logger.Error(err)
				return
			}

			if ses.logger.LogLevel(log.Info) {
				ses.logger.Infof("write: %s", ses.encoder.String(event))
			}

			// 发送消息
			n, err := ses.conn.Write(stream)
			if err != nil {
				ses.logger.Errorf("write stream error: %v", err)
				return
			}
			if n != len(stream) {
				ses.logger.Errorf("write stream error: has %d bytes, %d wrote", len(stream), n)
				return
			}
		}

		if closed {
			return
		}
	}
}

func (ses *Session) run(closeFunc func(*Session)) {

	ses.logger.Info("connected")

	// 连接成功回调
	ses.invoke(&network.Event{
		Meta:    network.MetaConnected,
		Msg:     &network.Connected{},
		Session: ses,
	})

	// 启动写线程
	go func() {
		ses.write()
		// 关闭读线程
		ses.conn.Close()
	}()

	// 启动读线程
	ses.read()
	// 关闭写线程
	ses.queue.Close()

	ses.logger.Info("disconnected")

	// 关闭回调
	closeFunc(ses)

	// 连接断开回调
	ses.invoke(&network.Event{
		Meta:    network.MetaDisconnected,
		Msg:     &network.Connected{},
		Session: ses,
	})
}

func (ses *Session) invoke(event *network.Event) {

	defer func() {
		if err := recover(); err != nil {
			ses.logger.Errorf("invoke panic: %v", err)
		}
	}()

	ses.config.Invoker.Invoke(event)
}

func newSession(
	peer string,
	id network.SessionID,
	conn net.Conn,
	conf *SessionConfig,
) *Session {

	// 日志
	name := fmt.Sprintf("%s.ses.%d", peer, id)
	logger := log.Std(name)
	logger.SetLevel(conf.LogLevel)

	return &Session{
		id:      id,
		conn:    conn,
		config:  conf,
		encoder: conf.EncoderMaker.New(),
		queue:   queue.New(conf.QueueLen),
		logger:  logger,
	}
}
