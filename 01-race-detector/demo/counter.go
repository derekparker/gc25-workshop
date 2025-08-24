package main

import (
	"fmt"
	"sync"
)

type incrementor struct {
	counter int
}

func (inc *incrementor) increment(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 1000; i++ {
		inc.counter++
	}
}

func main() {
	var wg sync.WaitGroup
	inc := &incrementor{}

	wg.Add(2)
	go inc.increment(&wg)
	go inc.increment(&wg)

	wg.Wait()

	// Result will be unpredictable, but should be 2000
	fmt.Printf("Final counter: %d\n", inc.counter)
}
