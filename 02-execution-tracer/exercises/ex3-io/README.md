# Exercise 2: Program that never terminates

## Problem
We have a program that runs forever without terminating.

## Your Task
1. Use the `go tool trace` visualizations and information to track down the bug and fix it.

## Testing
```bash
# Generate trace
go run main.go

# Analyze I/O patterns
go tool trace io.trace

# Look for:
# - Network blocking patterns
# - Missing timeout handling
# - Poor error visibility
```

## Expected Improvement
The program terminates after completing each network call.
