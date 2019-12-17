package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/laconiz/eros/log"
	"github.com/laconiz/eros/network"
	queue2 "github.com/laconiz/eros/queue"
	"sync"
	"time"
)

type SessionMgr interface {
	Add(network.Session)
	Del(network.Session)
}

type Session struct {
	id     network.SessionID // session ID
	addr   string            // 连接地址
	conn   *websocket.Conn   // websocket连接
	queue  *queue2.Queue     // 写入队列
	config *SessionConfig    // 配置
	log    *log.Logger       // 日志
	data   sync.Map          // 携带信息
}

// session ID
func (ses *Session) ID() network.SessionID {
	return ses.id
}

// 连接地址
func (ses *Session) Addr() string {
	return ses.addr
}

// 发送消息
func (ses *Session) Send(msg interface{}) {

	event, err := ses.config.Encoder.Encode(msg)
	if err != nil {
		ses.log.Errorf("encode message[%+v] error: %v", msg, err)
		return
	}

	ses.queue.Add(event)
}

// 发送数据流
func (ses *Session) SendStream(stream []byte) {
	ses.queue.Add(&network.Event{Stream: stream})
}

// 关闭连接
func (ses *Session) Close() {
	ses.queue.Close()
}

// 设置附加数据
func (ses *Session) Set(key interface{}, value interface{}) {
	ses.data.Store(key, value)
}

// 获取附加数据
func (ses *Session) Get(key interface{}) interface{} {
	if value, ok := ses.data.Load(key); ok {
		return value
	}
	return nil
}

// 读取流程
func (ses *Session) read() {

	conf := ses.config

	// 设置消息最大字节数
	ses.conn.SetReadLimit(ses.config.ReadLimit)

	for {

		// 设置读取超时
		ses.conn.SetReadDeadline(time.Now().Add(conf.ReadTimeout))

		// 读取消息流
		_, stream, err := ses.conn.ReadMessage()
		if err != nil {
			ses.log.Infof("read break by %v", err)
			return
		}

		// 反序列化消息
		event, err := conf.Encoder.Decode(stream)
		if err != nil {
			ses.log.Warnf("decode error: %v", err)
			return
		}

		// 消息回调
		// ses.log.Infof("read: %s", conf.Encoder.String(event))
		event.Session = ses
		ses.invoke(event)
	}
}

// 写入流程
func (ses *Session) write() {

	conf := ses.config

	for {

		// 设置写入超时
		ses.conn.SetWriteDeadline(time.Now().Add(conf.WriteTimeout))

		// 读取队列
		events, exited := ses.queue.Pick()

		for _, e := range events {

			event := e.(*network.Event)

			// ses.log.Infof("write: %s", conf.Encoder.String(event))

			err := ses.conn.WriteMessage(websocket.BinaryMessage, event.Stream)
			if err != nil {
				ses.log.Errorf("write error: %v", err)
				return
			}
		}

		// 关闭连接
		if exited {
			return
		}
	}
}

func (ses *Session) run(closeFunc func(*Session)) {

	ses.log.Infof("connected")

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

	ses.log.Infof("disconnected")

	closeFunc(ses)

	// 连接断开回调
	ses.invoke(&network.Event{
		Meta:    network.MetaDisconnected,
		Msg:     &network.Connected{},
		Session: ses,
	})
}

func newSession(
	id network.SessionID,
	name string,
	addr string,
	conn *websocket.Conn,
	config *SessionConfig,
) *Session {

	logName := fmt.Sprintf("%s.ses.%d", name, id)

	return &Session{
		id:     id,
		addr:   addr,
		conn:   conn,
		queue:  queue2.New(config.WriteQueueLen),
		config: config,
		log:    log.Std(logName),
		data:   sync.Map{},
	}
}
