package main

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
}

func NewClient(c *websocket.Conn) *Client {
	n := &Client{}
	n.conn = c
	return n
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) ReadLoop() {
	defer c.Close()
	for {
		var msg ClientMessage
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			log.Fatal(err)
			break
		}
		fmt.Println("HANDLE MESSAGE, type", msg.Type, "Params", msg.Params)
	}
}

func (c *Client) WriteLoop() {

}
