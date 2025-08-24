# Exercise 1: Scheduling Problem

## Problem
This program creates too many goroutines for CPU-bound work, causing scheduler overhead.

## Your Task
1. Generate trace and analyze the scheduling problem
2. Implement worker pool solution  
3. Compare before/after performance using traces

## Testing
```bash
# Generate trace of broken version
go run main.go

# Analyze trace
go tool trace broken.trace

# Look for:
# - Many goroutines created
# - High scheduler latency
# - Poor processor utilization
```

## Expected Improvement
Should see significant performance improvement with proper worker pool pattern.
