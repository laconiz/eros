package oceanus

import "github.com/laconiz/eros/oceanus/proto"

func NewPackage() *Pkg {
	return &Pkg{}
}

type Pkg struct {
	IDs     []proto.NodeID
	Type    proto.NodeType
	Keys    []proto.NodeKey
	Headers proto.Headers
	Message interface{}
}

func (pack *Pkg) SetHeader(key string, value interface{}) *Pkg {

	if pack.Headers == nil {
		pack.Headers = proto.Headers{}
	}

	pack.Headers[key] = value
	return pack
}

func (pack *Pkg) SetHeaders(headers map[string]interface{}) *Pkg {

	for key, value := range headers {
		pack.SetHeader(key, value)
	}

	return pack
}

func (pack *Pkg) SetMessage(message interface{}) *Pkg {
	pack.Message = message
	return pack
}
