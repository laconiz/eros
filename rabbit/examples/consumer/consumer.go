package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"strconv"
)

func main() {

}

func ConsumeUser(userID uint64) error {

	channel, err := connection.Channel()
	if err != nil {
		return err
	}

	channel.NotifyClose()
	queue := &Queue{Name: strconv.FormatUint(userID, 10)}

	if err := DeclareQueue(channel, queue); err != nil {
		return err
	}

	if err := channel.QueueBind(queue.Name, queue.Name, userProxy.Name, false, nil); err != nil {
		return err
	}

	delivery, err := channel.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	for message := range delivery {
		log.Printf("recv: %s", string(message.Body))
	}

	return nil
}

var (
	connection *amqp.Connection
)

func init() {

	var err error

	connection, err = amqp.Dial("amqp://thorin:Thorin-Rabbit@192.168.10.141:5672")
	if err != nil {
		panic(fmt.Errorf("connect to rabbit error: %w", err))
	}

	channel, err := connection.Channel()
	if err != nil {
		panic(fmt.Errorf("create channel error: %w", err))
	}

	if err := DeclareExchange(channel, userProxy); err != nil {
		panic(fmt.Errorf("declare user proxy exchange error: %w", err))
	}

	connection.NotifyBlocked()
}

type Exchange struct {
	Name       string
	Kind       string
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
}

var userProxy = &Exchange{
	Name:       "user.proxy",
	Kind:       amqp.ExchangeDirect,
	Durable:    true,
	AutoDelete: false,
	Internal:   false,
	NoWait:     false,
}

func DeclareExchange(c *amqp.Channel, e *Exchange) error {
	return c.ExchangeDeclare(e.Name, e.Kind, e.Durable, e.AutoDelete, e.Internal, e.NoWait, nil)
}

type Queue struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
}

func DeclareQueue(c *amqp.Channel, q *Queue) error {
	_, err := c.QueueDeclare(q.Name, q.Durable, q.AutoDelete, q.Exclusive, q.NoWait, nil)
	return err
}
