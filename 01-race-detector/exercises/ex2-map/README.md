# Exercise 2: Map Race Conditions - Cache vs Metrics

## Problem
A service uses two maps with fundamentally different access patterns:

1. **Cache Map**: Stores expensive computation results
   - **Stable keys**: User IDs, computation keys that don't change
   - **Cache pattern**: Once computed, read many times
   - **High hit rate**: After warm-up, mostly reads with occasional misses
   - **Invalidations**: Rare cache invalidations

2. **Metrics Map**: Tracks real-time performance metrics
   - **Growing keys**: New metric names as features are used
   - **Counter pattern**: Constant increments (writes)
   - **Occasional reads**: Periodic reporting and monitoring
   - **No deletion**: Metrics accumulate over time

Both maps have race conditions, but they need different solutions for optimal performance!

## Your Task
Fix the race conditions in both maps using the appropriate synchronization:

## Testing
```bash
# Detect the race conditions
go run -race main.go

# You'll see races on both maps:
# - cache map: concurrent read/write on CachedResult
# - metrics map: concurrent increments
# - GetCacheStats: iteration during modification
# - GetAllMetrics: range over map during writes

# Verify your solution
go run -race main.go  # Should show no race warnings
```

## Expected Output
- Cache workers should show high hit rates (>80% after warm-up)
- Metrics collectors should show ~90% write operations
- Verifier should find no consistency errors
- Monitor should show steady progress
- No race condition warnings

## Real-World Cache Patterns

This exercise models real caching scenarios:

### Database Query Cache
```go
// Expensive database lookup cached by user ID
userCache.Load("user_123") // Hit - no DB query needed
userCache.Store("user_456", expensiveDBLookup()) // Miss - cache result
```

### Computation Cache
```go
// ML model predictions cached by input hash
predictionCache.Load("input_hash_abc") // Reuse previous computation
```

### CDN/HTTP Cache
```go
// Rendered templates cached by route+params
pageCache.Load("/products/123?sort=price")
```

Remember: The goal isn't just to fix races, but to choose the optimal synchronization strategy for each use case!
