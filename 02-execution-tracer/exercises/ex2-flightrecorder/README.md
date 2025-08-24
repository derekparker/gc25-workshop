# Exercise 2: Mysterious long running network calls

## Problem
You have a program that makes multiple network calls. Some of them take longer than others...

## Your Task
1. Utilize the execution tracer in a way we haven't used it before. Look into
https://pkg.go.dev/golang.org/x/exp/trace#NewFlightRecorder and implement tracing based on
the flight recorder style, which is more suitable for long running production programs.

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
None, really... the network blocking is out of our control. This is more just to build the muscle memory
around creating a flight recorder style execution tracer within your program.

You should see a single trace output emitted conditionally based on your configured latency threshold.
