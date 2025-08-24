package main

import (
	"fmt"
	"log"
	"os"
	"runtime/trace"
	"sync"
	"time"
)

func main() {
	f, err := os.Create("scheduling.trace")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	trace.Start(f)
	defer trace.Stop()

	var wg sync.WaitGroup
	start := time.Now()

	for i := range 1000 {
		wg.Add(1)
		go expensiveComputation(i, &wg)
	}

	wg.Wait()
	fmt.Printf("Took: %v\n", time.Since(start))
}

func expensiveComputation(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	// Simulate CPU-intensive work with actual computation
	result := float64(id)
	iterations := 100000 + (id * 1000) // Vary iterations based on id

	for i := 0; i < iterations; i++ {
		// Perform various mathematical operations
		result = result * 1.000001
		result = result / 1.0000005
		if i%100 == 0 {
			result = result + float64(id)
			result = result - float64(id)/2.0
		}
		// Add some more complex operations
		if i%500 == 0 {
			for j := 0; j < 10; j++ {
				result = result * 1.01 / 1.005
			}
		}
	}

	// Use the result to prevent compiler optimization
	_ = result
}
