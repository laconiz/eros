package oceanus

// 远程通道
type Course struct {
	channel *Channel
	burl    *Burl
	router  *Router
}

func (c *Course) Info() *Channel {
	return c.channel
}

func (c *Course) Node() Mesh {
	return c.burl
}

func (c *Course) Push(message *Message) error {
	return c.burl.Push(message)
}
