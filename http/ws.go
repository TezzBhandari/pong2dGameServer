package http

import (
	"fmt"
	"net/http"

	"github.com/TezzBhandari/pong/ws"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func (s *Server) handleWebsocketConnection() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(rw, r, nil)
		if err != nil {
			fmt.Println("failed to establish websocket connection: ", err)
			return
		}

		clientId := s.snowflake.Generate()
		client := ws.NewClient(conn, s.Relay, ws.ClientId(clientId))
		client.Relay.Register <- client
		fmt.Printf("[CONNECTION] client %s connected\n", r.RemoteAddr)
		go client.Read()
		go client.Write()

	})
}
