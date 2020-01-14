package iris

import (
	"github.com/laconiz/eros/network/websocket"
	"github.com/laconiz/eros/oceanus"
)

var Encoder = websocket.NameEncoder

const (
	NodeUser    oceanus.NodeType = "user"
	NodeGateway oceanus.NodeType = "gateway"
)

const (
	UserProxyMessage oceanus.MsgType = iota
)
