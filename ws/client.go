package ws

import (
	"log"

	"github.com/bwmarrin/snowflake"
	"github.com/gorilla/websocket"
)

type ClientId snowflake.ID

type Client struct {
	Conn  *websocket.Conn
	Id    ClientId
	Relay *Relay
	Send  chan []byte
}

func NewClient(conn *websocket.Conn, relay *Relay, id ClientId) *Client {
	return &Client{
		Conn:  conn,
		Id:    id,
		Relay: relay,
		Send:  make(chan []byte),
	}
}

func (c *Client) Read() {
	defer c.Conn.Close()

	for {
		_, data, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("read error: ", err)
			}
			break
		}

		c.Relay.Broadcast <- data
	}
}

func (c *Client) Write() {
	defer c.Conn.Close()

	for {
		msg := <-c.Send

		err := c.Conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("write error: ", err)
			break
		}
	}
}
