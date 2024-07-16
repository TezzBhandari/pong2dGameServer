package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

func main() {
	// var sigChan chan bool
	// fmt.Println(sigChan) // output: nil

	// sigChan := make(chan bool)
	// go anotherThread(sigChan)
	// fmt.Println(<-sigChan)
	// fmt.Println(sigChan) // output: 0xc000094060
	// ch1, ch2 := make(chan bool), make(chan bool)

	// go thread(ch1, ch2)

	// select {
	// case msg := <-ch1:
	// 	fmt.Println("channel 1: ", msg)
	// case msg := <-ch2:
	// 	fmt.Println("channel 2:", msg)
	// }

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	sleepAndTalk(ctx, 5*time.Second, "hello")

}

func sleepAndTalk(ctx context.Context, d time.Duration, msg string) {
	select {
	case <-time.After(d):
		fmt.Println(msg)
	case <-ctx.Done():
		log.Fatal(ctx.Err())
	}
}

// func anotherThread(channel chan bool) {

// 	time.Sleep(1 * time.Second)
// 	channel <- true

// }

func thread(ch1 chan bool, ch2 chan bool) {
	ch1 <- true
	ch2 <- false

}
