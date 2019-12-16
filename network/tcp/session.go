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
	queue   queue.Queue       // 发送队列
	data    sync.Map          // 附加数据
	logger  *log.Logger       // 日志接口
}

func (ses *Session) ID() network.SessionID {
	return ses.id
}

func (ses *Session) Addr() string {
	return ses.conn.RemoteAddr().String()
}

func (ses *Session) Send(msg interface{}) {

}

func (ses *Session) SendStream(stream []byte) {

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
			ses.logger.Info(err)
			return
		}

		ses.invoke(event)
	}
}

func (ses *Session) write() {

	// 写线程退出时关闭连接
	defer ses.conn.Close()

	conf := ses.config

	for {

		// 设置写超时
		ses.conn.SetWriteDeadline(time.Now().Add(conf.WriteTimeout))

		// 获取消息列表
		events, closed := ses.queue.Pick()

		for _, e := range events {

			event := e.(*network.Event)

			// 序列化消息
			if event.Stream == nil {
				var err error
				if event, err = ses.encoder.Decode(event.Msg); err != nil {
					ses.logger.Info(err)
					return
				}
			}

			ses.logger.Infof("write: %s", ses.encoder.String(event))

			// 发送消息
			n, err := ses.conn.Write(event.Stream)
			if err != nil {
				ses.logger.Errorf("write stream error: %v", err)
				return
			}
			if n != len(event.Stream) {
				ses.logger.Errorf("write stream error: has %d bytes, %d wrote", len(event.Stream), n)
				return
			}
		}

		if closed {
			return
		}
	}
}

func (ses *Session) run() {

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
) *Session {

	name := fmt.Sprintf("%s.session.%d", peer, id)
	logger := log.Std(name)
	logger.SetLevel()

	return &Session{
		id:     id,
		conn:   nil,
		queue:  queue.NewQueue(),
		logger: logger,
	}
}
