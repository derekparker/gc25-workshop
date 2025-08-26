package main

import (
	"fmt"
	"log"
	"os"
	"runtime/trace"
	"sync"
	"sync/atomic"
)

type incrementor struct {
	counter int64
}

func (inc *incrementor) increment(wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < 1000; i++ {
		atomic.AddInt64(&inc.counter, 1)
	}
}

func main() {
	f, err := os.Create("counter.trace")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	trace.Start(f)
	defer trace.Stop()

	var wg sync.WaitGroup
	inc := &incrementor{}

	wg.Add(2)
	go inc.increment(&wg)
	go inc.increment(&wg)

	wg.Wait()

	// Result will be unpredictable, but should be 2000
	fmt.Printf("Final counter: %d\n", inc.counter)
}
