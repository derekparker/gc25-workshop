package main

import (
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

	wg.Add(2)
	go sender(ch, &wg)
	go receiver(ch, &wg)

	wg.Wait()
}

func sender(ch chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := range 5 {
		resp, err := http.Get("https://httpbin.org/delay/0.01")
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		ch <- fmt.Sprintf("message %d: %s", i, resp.Body)
	}
	close(ch)
}

func receiver(ch chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for msg := range ch {
		fmt.Printf("Received: size %d\n", len(msg))
	}
}
