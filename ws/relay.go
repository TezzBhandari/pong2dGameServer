package ws

import (
	"log"
	"sync"
)

type Relay struct {
	sync.RWMutex
	Clients    map[ClientId]*Client
	Register   chan *Client
	UnRegister chan *Client
	Broadcast  chan []byte
}

func NewRelay() *Relay {
	return &Relay{
		Clients:   make(map[ClientId]*Client),
		Register:  make(chan *Client),
		Broadcast: make(chan []byte),
	}
}

func (r *Relay) Run() {
	for {
		select {
		case client := <-r.Register:
			r.Lock()
			r.Clients[client.Id] = client
			r.Unlock()

		case client := <-r.UnRegister:
			r.Lock()
			delete(r.Clients, client.Id)
			r.Unlock()

		case data := <-r.Broadcast:
			for clientId := range r.Clients {
				client := r.Clients[clientId]

				select {
				case client.Send <- data:

				default:
					close(client.Send)
					delete(r.Clients, clientId)
				}
			}
			log.Println("i don't know")

		}
	}

}
