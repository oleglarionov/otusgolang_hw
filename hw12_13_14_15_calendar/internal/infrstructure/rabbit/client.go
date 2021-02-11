package rabbit

import (
	"context"
	"encoding/json"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/application/broker"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"log"
	"sync"
)

type Config struct {
	DSN   string
	Queue string
}

type Client struct {
	cfg         Config
	isConnected bool
	mux         sync.Mutex
	conn        *amqp.Connection
	ch          *amqp.Channel
	q           amqp.Queue
}

var _ broker.Pusher = (*Client)(nil)

var _ broker.Reader = (*Client)(nil)

func NewClient(cfg Config) *Client {
	return &Client{cfg: cfg}
}

func (c *Client) Push(_ context.Context, data interface{}) error {
	err := c.connect()
	if err != nil {
		return err
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		return errors.WithStack(err)
	}

	err = c.ch.Publish(
		"",
		c.q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        bytes,
		})
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (c *Client) connect() error {
	if c.isConnected {
		return nil
	}

	c.mux.Lock()
	defer c.mux.Unlock()
	if c.isConnected {
		return nil
	}

	conn, err := amqp.Dial(c.cfg.DSN)
	if err != nil {
		return errors.WithStack(err)
	}
	c.conn = conn
	log.Println("connection to rabbit established")

	ch, err := conn.Channel()
	if err != nil {
		return errors.WithStack(err)
	}
	c.ch = ch

	q, err := c.ch.QueueDeclare(
		c.cfg.Queue,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return errors.WithStack(err)
	}
	c.q = q

	c.isConnected = true

	return nil
}

func (c *Client) Read(ctx context.Context) (<-chan []byte, error) {
	err := c.connect()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	msgs, err := c.ch.Consume(
		c.cfg.Queue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	out := make(chan []byte, 1)
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-msgs:
				if !ok {
					return
				}
				l := len(msg.Body)
				data := make([]byte, l, l)
				copy(data, msg.Body)
				out <- data
			}
		}
	}()

	return out, nil
}

func (c *Client) Close() {
	if c.ch != nil {
		err := c.ch.Close()
		if err != nil {
			log.Println(err)
		}
	}

	if c.conn != nil {
		err := c.conn.Close()
		if err != nil {
			log.Println(err)
		}
	}
}
