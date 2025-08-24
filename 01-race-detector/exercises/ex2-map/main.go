package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Service simulates a service with both a cache and metrics tracking
type Service struct {
	// cache stores expensive computation results - read frequently after initial computation
	// Keys are stable (user IDs, computation keys)
	cache map[string]*CachedResult

	// metrics tracks real-time metrics - written frequently, read occasionally
	metrics map[string]int
}

// CachedResult represents an expensive computation result
type CachedResult struct {
	Value      string
	ComputedAt time.Time
	HitCount   int
}

// NewService creates a new service instance
func NewService() *Service {
	return &Service{
		cache:   make(map[string]*CachedResult),
		metrics: make(map[string]int),
	}
}

// GetFromCache retrieves a cached value or computes and caches it
func (s *Service) GetFromCache(key string) (*CachedResult, bool) {
	// Check if value exists in cache
	if result, exists := s.cache[key]; exists {
		// Increment hit count
		result.HitCount++
		s.RecordMetric("cache.hits")
		return result, true
	}

	s.RecordMetric("cache.misses")
	return nil, false
}

// StoreInCache stores a computed value in the cache
func (s *Service) StoreInCache(key string, value string) {
	s.cache[key] = &CachedResult{
		Value:      value,
		ComputedAt: time.Now(),
		HitCount:   0,
	}
}

// InvalidateCacheEntry removes a cache entry (rare operation)
func (s *Service) InvalidateCacheEntry(key string) {
	delete(s.cache, key)
	s.RecordMetric("cache.invalidations")
}

// GetCacheStats returns cache statistics
func (s *Service) GetCacheStats() (entries int, totalHits int) {
	for _, result := range s.cache {
		entries++
		totalHits += result.HitCount
	}
	return entries, totalHits
}

// RecordMetric increments a metric counter (called very frequently)
func (s *Service) RecordMetric(metric string) {
	s.metrics[metric]++
}

// GetMetric reads a metric value (called occasionally)
func (s *Service) GetMetric(metric string) int {
	return s.metrics[metric]
}

// GetAllMetrics returns a snapshot of all metrics (called rarely)
func (s *Service) GetAllMetrics() map[string]int {
	result := make(map[string]int)
	for k, v := range s.metrics {
		result[k] = v
	}
	return result
}

// expensiveComputation simulates an expensive operation like DB lookup or complex calculation
func expensiveComputation(key string) string {
	// Simulate expensive work
	time.Sleep(time.Millisecond * 2)
	return fmt.Sprintf("computed_value_for_%s_%d", key, rand.Intn(1000))
}

// cacheWorker simulates workers that heavily use the cache
func cacheWorker(id int, svc *Service, wg *sync.WaitGroup) {
	defer wg.Done()

	readCount := 0
	writeCount := 0

	// Stable set of keys that this worker uses
	myKeys := make([]string, 20)
	for i := 0; i < 20; i++ {
		myKeys[i] = fmt.Sprintf("user_%d_%d", id, i)
	}

	// Simulate realistic cache usage pattern
	for i := 0; i < 1000; i++ {
		// Pick a key from our stable set
		key := myKeys[rand.Intn(len(myKeys))]

		// Try to get from cache (95% of operations)
		if result, hit := svc.GetFromCache(key); hit {
			// Use cached value
			_ = result.Value
			readCount++
		} else {
			// Cache miss - compute and store
			value := expensiveComputation(key)
			svc.StoreInCache(key, value)
			writeCount++
		}

		// Occasionally invalidate a cache entry
		if rand.Float32() < 0.005 {
			keyToInvalidate := myKeys[rand.Intn(len(myKeys))]
			svc.InvalidateCacheEntry(keyToInvalidate)
			writeCount++
		}

		// Small delay between operations
		if i%100 == 0 {
			time.Sleep(time.Microsecond * 10)
		}
	}

	cacheHits := svc.GetMetric("cache.hits")
	cacheMisses := svc.GetMetric("cache.misses")
	hitRate := float64(cacheHits) * 100 / float64(cacheHits+cacheMisses)

	fmt.Printf("Cache Worker %d: %d reads (hits), %d writes (misses+invalidations), %.1f%% hit rate\n",
		id, readCount, writeCount, hitRate)
}

// metricsCollector simulates workers that primarily write metrics
func metricsCollector(id int, svc *Service, wg *sync.WaitGroup) {
	defer wg.Done()

	writeCount := 0
	readCount := 0

	for i := 0; i < 1000; i++ {
		// heavily write-biased
		if rand.Float32() < 0.90 {
			// Record various metrics
			metrics := []string{
				fmt.Sprintf("worker.%d.requests", id),
				fmt.Sprintf("worker.%d.errors", id),
				fmt.Sprintf("worker.%d.latency", id),
				"total.requests",
				"total.operations",
			}

			for _, metric := range metrics {
				svc.RecordMetric(metric)
				writeCount++
			}
		} else {
			// Occasionally read metrics for reporting
			metric := fmt.Sprintf("worker.%d.requests", id)
			value := svc.GetMetric(metric)
			_ = value
			readCount++
		}

		// Small delay to simulate work
		if i%100 == 0 {
			time.Sleep(time.Microsecond * 10)
		}
	}

	fmt.Printf("Metrics Collector %d: %d writes, %d reads (%.1f%% writes)\n",
		id, writeCount, readCount, float64(writeCount)*100/float64(writeCount+readCount))
}

// monitor periodically reports on both cache and metrics
func monitor(svc *Service, done chan bool) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Get cache statistics
			entries, totalHits := svc.GetCacheStats()

			// Get metrics
			metrics := svc.GetAllMetrics()
			cacheHits := metrics["cache.hits"]
			cacheMisses := metrics["cache.misses"]

			hitRate := float64(0)
			if cacheHits+cacheMisses > 0 {
				hitRate = float64(cacheHits) * 100 / float64(cacheHits+cacheMisses)
			}

			fmt.Printf("[Monitor] Cache: %d entries, %d total hits, %.1f%% hit rate | Metrics: %d unique counters\n",
				entries, totalHits, hitRate, len(metrics))

		case <-done:
			return
		}
	}
}

// verifier checks cache consistency
func verifier(svc *Service, wg *sync.WaitGroup) {
	defer wg.Done()

	// Let other workers warm up the cache
	time.Sleep(100 * time.Millisecond)

	errors := 0

	// Create some known keys
	testKeys := []string{"verify_1", "verify_2", "verify_3"}

	for _, key := range testKeys {
		// Store a value
		expectedValue := fmt.Sprintf("verified_%s", key)
		svc.StoreInCache(key, expectedValue)
	}

	// Verify we can read them back consistently
	for i := 0; i < 30; i++ {
		for _, key := range testKeys {
			result, exists := svc.GetFromCache(key)
			if !exists {
				fmt.Printf("[Verifier] Error: Key %s disappeared from cache!\n", key)
				errors++
			} else if result.Value != fmt.Sprintf("verified_%s", key) {
				fmt.Printf("[Verifier] Error: Key %s has wrong value: %s!\n", key, result.Value)
				errors++
			}
		}
		time.Sleep(10 * time.Millisecond)
	}

	// Verify hit counts are incrementing
	key := testKeys[0]
	result1, _ := svc.GetFromCache(key)
	hitCount1 := result1.HitCount

	// Access it a few more times
	for i := 0; i < 5; i++ {
		svc.GetFromCache(key)
	}

	result2, _ := svc.GetFromCache(key)
	hitCount2 := result2.HitCount

	if hitCount2 <= hitCount1 {
		fmt.Printf("[Verifier] Error: Hit count not incrementing properly for %s (was %d, now %d)\n",
			key, hitCount1, hitCount2)
		errors++
	}

	if errors == 0 {
		fmt.Println("[Verifier] No consistency errors detected (but races may still exist!)")
	} else {
		fmt.Printf("[Verifier] Found %d consistency errors\n", errors)
	}
}

func main() {
	svc := NewService()
	var wg sync.WaitGroup

	// Start monitor
	done := make(chan bool)
	go monitor(svc, done)

	// Start cache workers (cache-heavy operations)
	fmt.Println("Starting cache workers...")
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go cacheWorker(i, svc, &wg)
	}

	// Start metrics collectors (write-heavy operations)
	fmt.Println("Starting metrics collectors...")
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go metricsCollector(i, svc, &wg)
	}

	// Start verifier
	fmt.Println("Starting verifier...")
	wg.Add(1)
	go verifier(svc, &wg)

	fmt.Println()

	// Wait for workers
	wg.Wait()

	// Stop monitor
	done <- true
	time.Sleep(10 * time.Millisecond)

	// Final statistics
	fmt.Println("\n=== Final Statistics ===")

	// Cache statistics
	entries, totalHits := svc.GetCacheStats()
	fmt.Printf("Cache entries: %d\n", entries)
	fmt.Printf("Total cache hits across all entries: %d\n", totalHits)

	// Metrics
	allMetrics := svc.GetAllMetrics()
	fmt.Printf("Unique metrics tracked: %d\n", len(allMetrics))

	cacheHits := allMetrics["cache.hits"]
	cacheMisses := allMetrics["cache.misses"]
	invalidations := allMetrics["cache.invalidations"]

	if cacheHits+cacheMisses > 0 {
		hitRate := float64(cacheHits) * 100 / float64(cacheHits+cacheMisses)
		fmt.Printf("Overall cache hit rate: %.1f%% (%d hits, %d misses, %d invalidations)\n",
			hitRate, cacheHits, cacheMisses, invalidations)
	}

	totalOps := 0
	for k, v := range allMetrics {
		if k == "total.operations" {
			totalOps = v
			break
		}
	}
	fmt.Printf("Total operations recorded: %d\n", totalOps)
}
