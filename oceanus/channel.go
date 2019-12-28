package oceanus

type ChannelInfo struct {
	UUID string
	Type string
	Name string
	Peer string
}

type Channel struct {
	ChannelInfo
}
