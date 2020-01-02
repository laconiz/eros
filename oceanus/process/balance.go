package process

import (
	"fmt"
	"github.com/laconiz/eros/oceanus"
	"github.com/laconiz/eros/oceanus/channel"
	"math/rand"
	"time"
)

type balance struct {
	typo channel.Type
	info []channel.Channel
	data map[channel.Key]channel.Channel
	rand *rand.Rand
}

func (b *balance) Replace(data map[channel.Key]channel.Channel) {
	b.data = data
}

func (b *balance) analyze() {

	b.info = nil

	for _, c := range b.data {
		b.info = append(b.info, c)
	}
}

func (b *balance) Balance(message *oceanus.Message) error {

	if b.data != nil {
		b.analyze()
		b.data = nil
	}

	length := len(b.info)
	if length == 0 {
		return fmt.Errorf("%w: type[%v]", ErrNotFound, b.typo)
	}

	index := b.rand.Intn(len(b.info))
	return b.info[index].Push(message)
}

func newBalance(typo channel.Type) *balance {
	return &balance{
		typo: typo,
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}
