package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/trace"
	"sync"
)

func main() {
	f, err := os.Create("channel.trace")
	if err != nil {
		log.Fatal("Error creating trace file:", err)
	}
	defer f.Close()

	trace.Start(f)
	defer trace.Stop()

	var wg sync.WaitGroup
	ch := make(chan string, 5)

	ctx := context.Background()
	ctx, task := trace.NewTask(ctx, "requestAndSend")
	defer task.End()

	wg.Add(2)
	go sender(ch, &wg, ctx)
	go receiver(ch, &wg, ctx)

	wg.Wait()
}

func sender(ch chan string, wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()
	for i := range 5 {
		trace.Log(ctx, "httpbin request", fmt.Sprintf("%d", i))

		var resp *http.Response
		trace.WithRegion(ctx, "HTTP request", func() {
			var err error
			resp, err = http.Get("https://httpbin.org/delay/0.01")
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()
		})

		trace.WithRegion(ctx, "channel send task", func() {
			ch <- fmt.Sprintf("message %d: %s", i, resp.Body)
		})
	}
	close(ch)
}

func receiver(ch chan string, wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()

	trace.WithRegion(ctx, "chan read loop", func() {
		for msg := range ch {
			fmt.Printf("Received: size %d\n", len(msg))
		}
	})
}
