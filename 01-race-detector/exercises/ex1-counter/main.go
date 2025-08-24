package main

import (
	"fmt"
	"sync"
	"time"
)

// Stats tracks processing statistics
type Stats struct {
	processed   int
	workers     int
	startTime   time.Time
	lastUpdated time.Time
}

// NewStats creates a statistics tracker
func NewStats() *Stats {
	now := time.Now()
	return &Stats{
		startTime:   now,
		lastUpdated: now,
	}
}

// RecordWork updates the statistics for completed work
func (s *Stats) RecordWork(items int) {
	// Simulate some processing time
	time.Sleep(time.Microsecond)
	s.processed += items
	s.lastUpdated = time.Now()
}

// RegisterWorker tracks active workers
func (s *Stats) RegisterWorker() {
	s.workers++
	s.lastUpdated = time.Now()
}

// GetTotal returns the total processed items
func (s *Stats) GetTotal() int {
	return s.processed
}

// GetWorkerCount returns the number of registered workers
func (s *Stats) GetWorkerCount() int {
	return s.workers
}

// GetElapsedTime returns time since stats creation
func (s *Stats) GetElapsedTime() time.Duration {
	return time.Since(s.startTime)
}

// GetLastUpdated returns the last update time
func (s *Stats) GetLastUpdated() time.Time {
	return s.lastUpdated
}

// GetTimeSinceUpdate returns duration since last update
func (s *Stats) GetTimeSinceUpdate() time.Duration {
	return time.Since(s.lastUpdated)
}

// IsStale checks if stats haven't been updated recently
func (s *Stats) IsStale() bool {
	return s.GetTimeSinceUpdate() > 100*time.Millisecond
}

// processItems simulates a worker processing items
func processItems(id int, stats *Stats, wg *sync.WaitGroup) {
	defer wg.Done()

	stats.RegisterWorker()

	// Process items in batches
	batchSize := 100
	for batch := 0; batch < 10; batch++ {
		// Process each item in the batch
		for item := 0; item < batchSize; item++ {
			stats.RecordWork(1)
		}

		// Small delay between batches
		if batch%3 == 0 {
			time.Sleep(time.Microsecond * 10)
		}

		// Occasionally check our own progress and staleness
		if batch == 5 {
			current := stats.GetTotal()
			lastUpdate := stats.GetLastUpdated()
			timeSince := time.Since(lastUpdate)
			fmt.Printf("Worker %d halfway: %d items (last update: %v ago)\n",
				id, current, timeSince.Round(time.Microsecond))
		}
	}

	// Check if we're the last to update
	updateTime := stats.GetLastUpdated()
	fmt.Printf("Worker %d completed at %v\n", id, updateTime.Format("15:04:05.000000"))
}

// monitorProgress periodically reports the current progress
func monitorProgress(stats *Stats, done chan bool) {
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			total := stats.GetTotal()
			workers := stats.GetWorkerCount()
			elapsed := stats.GetElapsedTime()
			sinceUpdate := stats.GetTimeSinceUpdate()

			// Calculate rate
			rate := float64(total) / elapsed.Seconds()

			// Check if processing has stalled
			status := "active"
			if stats.IsStale() {
				status = "STALE"
			}

			fmt.Printf("[Monitor] %s - %d items by %d workers (%.0f items/sec) - last update: %v ago\n",
				status, total, workers, rate, sinceUpdate.Round(time.Millisecond))

		case <-done:
			return
		}
	}
}

// auditLog periodically logs detailed statistics
func auditLog(stats *Stats, done chan bool) {
	ticker := time.NewTicker(75 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Read multiple fields for audit
			total := stats.GetTotal()
			workers := stats.GetWorkerCount()
			lastUpdate := stats.GetLastUpdated()

			if total > 0 {
				fmt.Printf("[Audit] Snapshot at %v: %d items, %d workers, last activity: %v\n",
					time.Now().Format("15:04:05.000"),
					total,
					workers,
					lastUpdate.Format("15:04:05.000000"))
			}

		case <-done:
			return
		}
	}
}

func main() {
	stats := NewStats()
	var wg sync.WaitGroup

	numWorkers := 5

	fmt.Printf("Starting %d workers at %v...\n", numWorkers, time.Now().Format("15:04:05.000000"))

	// Start progress monitor
	monitorDone := make(chan bool)
	go monitorProgress(stats, monitorDone)

	// Start audit logger
	auditDone := make(chan bool)
	go auditLog(stats, auditDone)

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go processItems(i, stats, &wg)
	}

	wg.Wait()

	// Stop monitors
	monitorDone <- true
	auditDone <- true
	time.Sleep(10 * time.Millisecond) // Let final output complete

	// Final statistics
	fmt.Printf("\nProcessing complete!\n")
	fmt.Printf("Workers used: %d\n", stats.GetWorkerCount())
	fmt.Printf("Total items processed: %d (expected: 5000)\n", stats.GetTotal())
	fmt.Printf("Total time: %v\n", stats.GetElapsedTime())
	fmt.Printf("Last update was: %v\n", stats.GetLastUpdated().Format("15:04:05.000000"))
	fmt.Printf("Time since last update: %v\n", stats.GetTimeSinceUpdate())
}
