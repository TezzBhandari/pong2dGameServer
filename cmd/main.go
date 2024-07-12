package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/TezzBhandari/pong/http"
)

type Cofig struct {
	Addr string
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	go func() { <-sigChan; cancel() }()

	m := NewMain()
	if err := m.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		if err = m.Close(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(1)
	}

	// wait for ctrl c
	<-ctx.Done()

	// cleanup the program
	if err := m.Close(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type Main struct {
	HTTPServer *http.Server
}

func NewMain() *Main {
	return &Main{
		HTTPServer: http.NewHttpServer(":9090"),
	}
}

func (m *Main) Close() error {
	if m.HTTPServer != nil {
		if err := m.HTTPServer.Close(); err != nil {
			return err
		}
	}

	return nil
}

func (m *Main) Run() error {
	if err := m.HTTPServer.Open(); err != nil {
		return err
	}

	log.Printf("running: url=%q", m.HTTPServer.URL())
	return nil
}
