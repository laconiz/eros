package oceanus

import (
	"errors"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/laconiz/eros/consul"
	"github.com/laconiz/eros/json"
	"github.com/laconiz/eros/log"
	"github.com/laconiz/eros/network"
	"github.com/laconiz/eros/network/tcp"
	"math/big"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"
)

func powerForAddr(addr string) (uint64, error) {

	ap := strings.Split(addr, ":")
	if len(ap) != 2 {
		return 0, errors.New("invalid addr format")
	}

	ip := net.ParseIP(ap[0])
	if ip == nil {
		return 0, errors.New("invalid ip address")
	}

	port, err := strconv.ParseUint(ap[1], 10, 64)
	if err != nil || port > 65535 {
		return 0, errors.New("invalid port address")
	}

	power := big.NewInt(0).SetBytes(ip.To4()).Uint64()
	return power<<16 | port, nil
}

// 同步节点连接
func (p *Process) OnNodeListUpdated(nodes []*Node) {

	p.mutex.Lock()
	defer p.mutex.Unlock()

	sp, err := powerForAddr(p.Node.Addr)
	if err != nil {
		return
	}

	for _, node := range nodes {

		if node.ID == p.Node.ID {
			continue
		}

		if _, ok := p.connectors[node.Addr]; ok {
			continue
		}

		power, err := powerForAddr(node.Addr)
		if err != nil {
			continue
		}

		if (sp > power && (sp-power)%2 == 0) || (sp < power && (power-sp)%2 != 0) {
			continue
		}

		conf := tcp.ConnectorConfig{
			Name:      fmt.Sprintf("oceanus.connector.%s", node.Addr),
			Addr:      node.Addr,
			Reconnect: true,
			Session: tcp.SessionConfig{
				ReadTimeout:  time.Second * 6,
				WriteTimeout: time.Second * 6,
				LogLevel:     log.Warn,
				QueueLen:     64,
				Invoker:      p.NewCommonInvoker(),
			},
		}

		connector := tcp.NewConnector(conf)
		go connector.Run()
		p.connectors[node.Addr] = connector
	}
}

func (p *Process) NewCommonInvoker() *network.StdInvoker {

	key := "node"

	invoker := network.NewStdInvoker()

	invoker.Register(Message{}, func(event *network.Event) {
		p.Push(event.Msg.(*Message))
	})

	invoker.Register(NodeJoinMsg{}, func(event *network.Event) {
		if node := event.Msg.(*NodeJoinMsg).Node; node != nil {
			event.Session.Set(key, node)
			p.OnNodeJoin(node, event.Session)
		}
	})

	invoker.Register(NodeQuitMsg{}, func(event *network.Event) {
		if node := event.Msg.(*NodeQuitMsg).Node; node != nil {
			p.OnNodeQuit(node)
		}
	})

	invoker.Register(ChannelJoinMsg{}, func(event *network.Event) {
		p.OnRouteJoin(event.Msg.(*ChannelJoinMsg).Channels)
	})

	invoker.Register(ChannelQuitMsg{}, func(event *network.Event) {
		p.OnRouteQuit(event.Msg.(*ChannelQuitMsg).Channels)
	})

	invoker.Register(network.Disconnected{}, func(event *network.Event) {
		if node := event.Session.Get(key); node != nil {
			p.OnNodeDisconnected(node.(*Node))
		}
	})

	return invoker
}

func (p *Process) Run() {

	conf := tcp.AcceptorConfig{
		Name: "oceanus.acceptor",
		Addr: p.Node.Addr,
		Session: tcp.SessionConfig{
			ReadTimeout:  time.Second * 5,
			WriteTimeout: time.Second * 5,
			LogLevel:     log.Warn,
			QueueLen:     64,
			Invoker:      p.NewCommonInvoker(),
		},
	}

	p.acceptor = tcp.NewAcceptor(conf)
	p.acceptor.Run()

	// 注销监听器
	defer p.acceptor.Stop()

	consulKey := fmt.Sprintf("%s%s", kvPrefix, p.Node.ID)

	// 注册节点信息
	if err := consul.KV().Store(consulKey, p.Node); err != nil {
		logger.Errorf("store node error: %w", err)
		return
	}
	// 注销节点信息
	defer func() {
		if err := consul.KV().Delete(consulKey); err != nil {
			logger.Errorf("delete node error: %w", err)
		}
	}()

	// 监听节点列表
	watcher, err := consul.NewKeyPrefixWatcher("oceanus/", func(pairs api.KVPairs) {

		var nodes []*Node

		for _, pair := range pairs {

			node := &Node{}
			if err := json.Unmarshal(pair.Value, node); err != nil {
				continue
			}

			nodes = append(nodes, node)
		}

		p.OnNodeListUpdated(nodes)
	})
	if err != nil {
		logger.Errorf("new node watcher error: %w", err)
	}
	go watcher.Run()
	// 取消监听
	defer watcher.Stop()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, os.Kill)

	for {
		select {
		case signal := <-exit:
			logger.Infof("exit signal received: %v", signal)
			p.OnDestroy()
			logger.Info("process destroyed")
			return
		case <-time.After(time.Second * 5):
			p.NotifyState()
		}
	}
}

const (
	kvPrefix = "oceanus/"
)
