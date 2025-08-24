package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type APIResult struct {
	URL      string
	Duration time.Duration
	Error    error
}

func main() {
	// TODO: Add flight recorder instrumentation here.
	// Consider using https://pkg.go.dev/golang.org/x/exp/trace#NewFlightRecorder
	// and outputting the trace if any request takes longer than a threshold.

	urls := []string{
		"https://httpbin.org/delay/1",
		"https://httpbin.org/delay/2",
		"https://httpbin.org/delay/3",
	}

	// Channel to communicate results
	results := make(chan APIResult, len(urls))

	// Goroutine that makes API calls
	go func() {
		for _, url := range urls {
			start := time.Now()
			err := callAPI(url)
			duration := time.Since(start)

			results <- APIResult{
				URL:      url,
				Duration: duration,
				Error:    err,
			}
		}
		close(results)
	}()

	// Goroutine that reads results
	go func() {
		for result := range results {
			if result.Error != nil {
				fmt.Printf("❌ %s failed: %v (took %v)\n", result.URL, result.Error, result.Duration)
			} else {
				fmt.Printf("✅ %s succeeded in %v\n", result.URL, result.Duration)
			}

			// TODO: Check if duration exceeds threshold (e.g., 2.5 seconds)
			// and if so, output the flight recorder trace to a file.
		}
	}()

	// Simple wait to ensure goroutines complete
	// In a real application, you'd use sync.WaitGroup or channels
	time.Sleep(15 * time.Second)
	fmt.Println("\nAll operations completed")
}

func callAPI(url string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	return nil
}
