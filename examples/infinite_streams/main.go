package main

import (
	"context"
	"github.com/4thel00z/fakenews/pkg/v1/fakenews"
	"log"
	"time"
)

type Message struct {
	Content string
}

func main() {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	in := make(chan Message)
	nameRule := fakenews.RandomPick[Message, string]("Content", "Hello", "World")
	message := Message{}
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case val := <-in:
				println("Received:", val.Content)
			}
		}
	}()

	log.Fatalln(fakenews.InfiniteStream[Message](ctx, message, in, nameRule))
}
