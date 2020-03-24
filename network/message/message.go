package message

type Message struct {
	Meta   Meta        // 消息元数据
	Msg    interface{} // 消息体
	Raw    []byte      // 消息体流
	Stream []byte      // 消息流
}
