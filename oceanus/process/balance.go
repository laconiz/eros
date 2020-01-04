package process

import (
	"github.com/laconiz/eros/oceanus"
	"github.com/laconiz/eros/oceanus/channel"
	"math/rand"
	"time"
)

type group struct {
	typo     channel.Type
	dirty    bool
	channels map[channel.Key]channel.Channel
	balances []channel.Channel
	rand     *rand.Rand
}

func (g *group) Dirty() {
	g.dirty = true
}

func (g *group) Update(c channel.Channel) {
	g.dirty = true
	g.channels[c.Info().Key] = c
}

func (g *group) Get(key channel.Key) (channel.Channel, bool) {
	c, ok := g.channels[key]
	return c, ok
}

func (g *group) Rand(message *oceanus.Message) error {
	return nil
}

func (g *group) Balance(message *oceanus.Message) error {

	return nil
}

func newGroup(typo channel.Type) *group {
	return &group{
		typo:     typo,
		channels: map[channel.Key]channel.Channel{},
		rand:     rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}
