package websocket

import (
	"fmt"
	"log"
	"net/url"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	URL       *url.URL
	conn      *websocket.Conn
	done      chan struct{}
	once      sync.Once
	onOpen    func(c *Client)
	onClose   func(c *Client)
	onMessage func(c *Client, payload []byte)
}

func NewClient(
	onOpen func(c *Client),
	onClose func(c *Client),
	onMessage func(c *Client, payload []byte),
) *Client {
	return &Client{
		done:      make(chan struct{}),
		once:      sync.Once{},
		onOpen:    onOpen,
		onClose:   onClose,
		onMessage: onMessage,
	}
}

func (c *Client) Connect(webSocketUrl string) error {
	var err error

	c.URL, err = url.Parse(webSocketUrl)
	if err != nil {
		return fmt.Errorf("failed to parse websocket url: %v", webSocketUrl)
	}

	c.conn, _, err = websocket.DefaultDialer.Dial(webSocketUrl, nil)
	if err != nil {
		return fmt.Errorf("websocket dial failed: %v", err)
	}

	if c.onOpen != nil {
		c.onOpen(c)
	}

	go c.messageHandler()

	return nil
}

func (c *Client) messageHandler() {
	defer c.Disconnect()

	// Receive message forever
	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket Client: read error", err)
			return
		}

		if c.onMessage != nil {
			c.onMessage(c, data)
		}
	}
}

func (c *Client) SendJSON(v any) error {
	return c.conn.WriteJSON(v)
}

func (c *Client) Disconnect() {
	// なんか二回呼び出されて close できずに panic で終了するから
	// once で囲う
	c.once.Do(func() {
		if c.conn != nil {
			log.Println("WebSocket Client: disconnecting...")
			c.conn.Close()

			if c.onClose != nil {
				c.onClose(c)
			}

			log.Println("WebSocket Client: successfully disconnected")
		}

		close(c.done)
	})
}

func (c *Client) Reconnect(url string) error {
	log.Println("WebSocket Client: reconnecting now...", url)
	c.Disconnect()
	return c.Connect(url)
}

func (c *Client) Wait() {
	<-c.done
	log.Println("WebSocket Client: connection closed")
}
