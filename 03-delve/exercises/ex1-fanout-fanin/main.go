package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

// DataItem represents input data to process
type DataItem struct {
	ID    int
	Value int
}

// ProcessedData represents the result of processing
type ProcessedData struct {
	ItemID   int
	Original int
	Result   int
	WorkerID int
}

// DataProcessor implements a fan-out/fan-in processing pipeline
type DataProcessor struct {
	numWorkers int
	input      chan DataItem
	output     chan ProcessedData
	errors     chan error
	wg         sync.WaitGroup
	mu         sync.Mutex
	processed  int
	errors_cnt int
}

// NewDataProcessor creates a new processor
func NewDataProcessor(numWorkers int) *DataProcessor {
	return &DataProcessor{
		numWorkers: numWorkers,
		input:      make(chan DataItem),      // BUG: unbuffered
		output:     make(chan ProcessedData), // BUG: unbuffered
		errors:     make(chan error, 1),      // BUG: buffer too small
	}
}

// Start begins the fan-out/fan-in processing
func (dp *DataProcessor) Start() {
	log.Printf("Starting processor with %d workers\n", dp.numWorkers)

	// Fan-out: start worker goroutines
	for i := 0; i < dp.numWorkers; i++ {
		dp.wg.Add(1)
		go dp.worker(i + 1)
	}

	// BUG: No goroutine to handle errors channel
	// BUG: No goroutine to collect results
}

// worker processes data items
func (dp *DataProcessor) worker(id int) {
	defer dp.wg.Done() // BUG: This might not always be called

	log.Printf("Worker %d started\n", id)

	for item := range dp.input {
		log.Printf("Worker %d processing item %d\n", id, item.ID)

		// Simulate processing with potential errors
		if item.Value < 0 {
			// BUG: Error channel might block if full
			dp.errors <- fmt.Errorf("worker %d: negative value %d for item %d",
				id, item.Value, item.ID)

			dp.mu.Lock()
			dp.errors_cnt++
			dp.mu.Unlock()
			continue
		}

		// Simulate processing delay
		processingTime := time.Duration(rand.Intn(500)+100) * time.Millisecond
		time.Sleep(processingTime)

		// Simulate occasional worker getting stuck
		if rand.Intn(20) == 0 {
			log.Printf("Worker %d stuck on item %d!\n", id, item.ID)
			time.Sleep(5 * time.Second)
			// BUG: Worker doesn't check for cancellation
		}

		result := ProcessedData{
			ItemID:   item.ID,
			Original: item.Value,
			Result:   item.Value * item.Value, // Square the value
			WorkerID: id,
		}

		// BUG: This will block if no one is reading
		dp.output <- result

		dp.mu.Lock()
		dp.processed++
		dp.mu.Unlock()

		log.Printf("Worker %d completed item %d\n", id, item.ID)
	}

	log.Printf("Worker %d shutting down\n", id)
}

// Process sends data items for processing
func (dp *DataProcessor) Process(items []DataItem) {
	log.Printf("Processing %d items\n", len(items))

	// Send items to workers
	go func() {
		for _, item := range items {
			log.Printf("Sending item %d for processing\n", item.ID)
			dp.input <- item // BUG: Can block if all workers are busy
		}
		// BUG: Should close input channel here
		log.Println("All items sent")
	}()
}

// CollectResults gathers processed data
func (dp *DataProcessor) CollectResults() []ProcessedData {
	var results []ProcessedData
	done := make(chan bool)

	go func() {
		// BUG: This will block forever if output is never closed
		for result := range dp.output {
			results = append(results, result)
			log.Printf("Collected result for item %d from worker %d\n",
				result.ItemID, result.WorkerID)
		}
		done <- true
	}()

	// Wait with timeout
	select {
	case <-done:
		log.Printf("Collected %d results\n", len(results))
	case <-time.After(10 * time.Second):
		log.Println("Timeout collecting results")
		// BUG: Goroutine leak - collector is still running
	}

	return results
}

// Wait waits for all workers to complete
func (dp *DataProcessor) Wait() {
	// BUG: This will hang if workers are blocked on sending
	dp.wg.Wait()
	// BUG: Should close output channel here
	log.Println("All workers completed")
}

// GetStats returns processing statistics
func (dp *DataProcessor) GetStats() (processed int, errors int) {
	dp.mu.Lock()
	defer dp.mu.Unlock()
	return dp.processed, dp.errors_cnt
}

// Merge combines multiple result channels (fan-in)
func Merge(channels ...<-chan ProcessedData) <-chan ProcessedData {
	out := make(chan ProcessedData) // BUG: unbuffered

	var wg sync.WaitGroup

	// Start a goroutine for each input channel
	for _, ch := range channels {
		wg.Add(1)
		go func(c <-chan ProcessedData) {
			defer wg.Done() // BUG: Won't be called if goroutine blocks
			for val := range c {
				out <- val // BUG: Can block if no reader
			}
		}(ch)
	}

	// BUG: This goroutine might close channel while senders are still active
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// SimulateMultiStage demonstrates a multi-stage pipeline with fan-out/fan-in
func SimulateMultiStage(items []DataItem) {
	// Stage 1: Split data into chunks for parallel processing
	numChunks := 3
	chunkSize := len(items) / numChunks

	var resultChannels []<-chan ProcessedData

	for i := 0; i < numChunks; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if i == numChunks-1 {
			end = len(items)
		}

		chunk := items[start:end]
		processor := NewDataProcessor(2) // 2 workers per chunk
		processor.Start()

		// Process chunk in parallel
		go func(p *DataProcessor, data []DataItem) {
			p.Process(data)
			// BUG: Never calls Wait() or closes channels
		}(processor, chunk)

		resultChannels = append(resultChannels, processor.output)
	}

	// Fan-in: merge all result channels
	merged := Merge(resultChannels...)

	// Collect all results
	var allResults []ProcessedData
	timeout := time.After(15 * time.Second)

	collecting := true
	for collecting {
		select {
		case result, ok := <-merged:
			if !ok {
				collecting = false
				break
			}
			allResults = append(allResults, result)
			log.Printf("Merged result for item %d\n", result.ItemID)
		case <-timeout:
			log.Println("Timeout in multi-stage processing")
			collecting = false
		}
	}

	log.Printf("Multi-stage processing collected %d results\n", len(allResults))
}

func generateData(count int) []DataItem {
	items := make([]DataItem, count)
	for i := 0; i < count; i++ {
		value := rand.Intn(100) - 10 // Some negative values to trigger errors
		items[i] = DataItem{
			ID:    i + 1,
			Value: value,
		}
	}
	return items
}

func main() {
	log.SetFlags(log.Lmicroseconds)
	rand.Seed(time.Now().UnixNano())

	log.Println("=== Starting Fan-out/Fan-in Processing Demo ===")

	// Test 1: Basic fan-out/fan-in
	log.Println("\n--- Test 1: Basic Processing ---")
	processor := NewDataProcessor(3)
	processor.Start()

	items := generateData(10)
	processor.Process(items)

	// Start result collection in background
	go func() {
		results := processor.CollectResults()
		log.Printf("Test 1 collected %d results\n", len(results))
	}()

	// Try to wait for completion
	done := make(chan bool)
	go func() {
		processor.Wait()
		done <- true
	}()

	select {
	case <-done:
		log.Println("Test 1 completed successfully")
	case <-time.After(20 * time.Second):
		log.Println("Test 1 timeout - possible deadlock")

		// Get stats to see what happened
		processed, errors := processor.GetStats()
		log.Printf("Stats: processed=%d, errors=%d\n", processed, errors)
	}

	// Test 2: Multi-stage pipeline
	log.Println("\n--- Test 2: Multi-Stage Pipeline ---")
	moreItems := generateData(15)
	SimulateMultiStage(moreItems)

	// Keep program alive for debugging
	log.Println("\nProgram finished but keeping alive for debugging...")
	log.Println("Use Delve to inspect goroutine state")
	select {} // Hang forever for debugging
}
