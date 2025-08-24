package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/trace"
	"sync"
	"time"
)

type APIResult struct {
	URL      string
	Duration time.Duration
	Error    error
}

func main() {
	// Add trace instrumentation here.
	// NOTE: Pretend this is a long running process and consider using
	// a flight recorder instead of a program wide trace:
	// https://pkg.go.dev/golang.org/x/exp/trace#NewFlightRecorder

	f, err := os.Create("io.trace")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	trace.Start(f)
	defer trace.Stop()

	ctx, task := trace.NewTask(context.Background(), "main")
	defer task.End()

	trace.WithRegion(ctx, "simulateNetworkLoad", func() {
		simulateNetworkLoad()
	})
}

func simulateNetworkLoad() {
	urls := []string{
		"https://httpbin.org/delay/1",
		"https://httpbin.org/delay/2",
		"https://httpbin.org/delay/3",
	}

	// Create channels for each API call result
	resultChannels := make([]chan APIResult, len(urls))
	for i := range resultChannels {
		resultChannels[i] = make(chan APIResult, 1)
	}

	// Create a channel for aggregated results
	aggregatedResults := make(chan []APIResult, 1)

	var wg sync.WaitGroup

	// Start API call goroutines
	for i, url := range urls {
		wg.Add(1)
		go func(url string, resultChan chan APIResult) {
			defer wg.Done()

			// TODO: Add trace instrumentation here
			ctx, task := trace.NewTask(context.Background(), "apiCall")
			defer task.End()

			trace.WithRegion(ctx, "callAPI:"+url, func() {
				start := time.Now()
				err := callAPI(ctx, url)
				// Can we use this to determine when to output trace results?
				duration := time.Since(start)

				// Send result to channel
				resultChan <- APIResult{
					URL:      url,
					Duration: duration,
					Error:    err,
				}
			})
		}(url, resultChannels[i])
	}

	// Start goroutines that wait on individual results (these will be blocked)
	for i, ch := range resultChannels {
		wg.Add(1)
		go func(idx int, resultChan chan APIResult) {
			defer wg.Done()

			// TODO: Add trace instrumentation here
			ctx, task := trace.NewTask(context.Background(), "resultWaiter")
			defer task.End()

			// This goroutine will be blocked waiting for the result
			var result APIResult
			trace.WithRegion(ctx, fmt.Sprintf("waitForResult:%d", idx), func() {
				result = <-resultChan
			})

			if result.Error != nil {
				fmt.Printf("Result %d: %s failed: %v\n", idx, result.URL, result.Error)
			} else {
				fmt.Printf("Result %d: %s completed in %v\n", idx, result.URL, result.Duration)
			}
		}(i, ch)
	}

	// Start a goroutine that aggregates all results (will be blocked until all are ready)
	wg.Add(1)
	go func() {
		defer wg.Done()

		// TODO: Add trace instrumentation here
		ctx, task := trace.NewTask(context.Background(), "aggregator")
		defer task.End()

		results := make([]APIResult, 0, len(resultChannels))

		// This will block on each channel sequentially
		for i, ch := range resultChannels {
			// Create a new channel to re-send the result
			// since the individual waiter goroutines also need it
			resCopy := make(chan APIResult, 1)
			go func(idx int, originalCh chan APIResult) {
				trace.WithRegion(ctx, fmt.Sprintf("readOriginalCh:%d", i), func() {
					r := <-originalCh
					resCopy <- r
					originalCh <- r // Send it back for the individual waiter
				})
			}(i, ch)
			var result APIResult
			trace.WithRegion(ctx, fmt.Sprintf("aggregateWait:%d", i), func() {
				result = <-resCopy
			})
			results = append(results, result)
		}

		aggregatedResults <- results
	}()

	// Start a final goroutine that waits for all aggregated results
	wg.Add(1)
	go func() {
		defer wg.Done()

		// TODO: Add trace instrumentation here
		ctx, task := trace.NewTask(context.Background(), "summarizer")
		defer task.End()

		// This will block until aggregation is complete
		var allResults []APIResult
		trace.WithRegion(ctx, "waitForAggregated", func() {
			allResults = <-aggregatedResults
		})

		var totalDuration time.Duration
		var successCount int

		for _, result := range allResults {
			if result.Error == nil {
				totalDuration += result.Duration
				successCount++
			}
		}

		if successCount > 0 {
			avgDuration := totalDuration / time.Duration(successCount)
			fmt.Printf("\n=== Summary ===\n")
			fmt.Printf("Successful calls: %d/%d\n", successCount, len(allResults))
			fmt.Printf("Average duration: %v\n", avgDuration)
			fmt.Printf("Total duration: %v\n", totalDuration)
		}
	}()

	// Create additional monitoring goroutines that will show blocking behavior
	done := make(chan struct{})

	// Monitor goroutine that periodically checks status
	go func() {
		ctx, task := trace.NewTask(context.Background(), "monitor")
		defer task.End()

		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				// TODO: Add trace instrumentation here
				trace.WithRegion(ctx, "statusCheck", func() {
					fmt.Printf(".")
				})
			case <-done:
				return
			}
		}
	}()

	wg.Wait()
	close(done)
	fmt.Println("\nAll operations completed")
}

func callAPI(ctx context.Context, url string) error {
	// Add timeout to context
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// TODO: Add trace instrumentation here
	trace.WithRegion(ctx, "preprocessing", func() {
		// Simulate some processing
		time.Sleep(10 * time.Millisecond)
	})

	r, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	client := &http.Client{}

	// TODO: Add trace instrumentation here
	var resp *http.Response
	trace.WithRegion(ctx, "httpRequest", func() {
		resp, err = client.Do(r)
	})

	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	return nil
}
