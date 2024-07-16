package http

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

type MessageRelay struct {
	Msg chan []byte
}

func NewMessageRealy() *MessageRelay {
	return &MessageRelay{
		Msg: make(chan []byte, 256),
	}
}

func (m *MessageRelay) Relay() {
	// for {
	// 	select {
	// 	case msg := <-m.Msg:
	// 		fmt.Println("recieved msg: ", string(msg))
	// 	default:
	// 		fmt.Println("i don't know")
	// 	}
	// }
	msg := <-m.Msg
	fmt.Println(msg)
}

func CreateWs(t time.Duration, m *MessageRelay) {
	conn, _, err := websocket.DefaultDialer.Dial("ws://192.168.1.100:9090/ws", nil)
	if err != nil {
		fmt.Println(err)
	}
	addr := conn.LocalAddr()
	fmt.Printf("connection %s established", addr)

	defer func() {
	}()

	count := 0

	go func() {
		defer func() {
			conn.Close()
			close(m.Msg)
			fmt.Printf("connection %s closed\n", addr)

		}()
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				fmt.Println(err)
				break
			}
			fmt.Println(string(msg))
			m.Msg <- msg

		}

	}()

	go func() {
		defer func() {
			conn.Close()
			fmt.Printf("connection %s closed\n", addr)
		}()

		for {
			err := conn.WriteMessage(websocket.BinaryMessage, []byte(fmt.Sprintf("client %s count %d\n", addr, count)))
			if err != nil {
				fmt.Println(err)
				break
			}
			count += 1
			time.Sleep(t)
		}
	}()

}
