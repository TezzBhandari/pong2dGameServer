package main

import (
	"os"
	"os/signal"
	"time"

	"github.com/TezzBhandari/pong/http"
)

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	m := http.NewMessageRealy()
	go m.Relay()
	// for i := 0; i < 5; i++ {
	http.CreateWs(1*time.Second, m)
	// }

	// waits for ctrl-c
	<-interrupt
}
