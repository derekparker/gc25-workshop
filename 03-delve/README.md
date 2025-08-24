# 03-delve: Debugging Go Concurrency with Delve

## Overview
This section teaches you how to use Delve, the Go debugger, to diagnose and fix complex concurrency issues. You'll learn to inspect goroutines, analyze channel states, trace deadlocks, and find goroutine leaks.

## Why Delve for Concurrency?
While the race detector finds data races and the execution tracer shows timing issues, Delve lets you:
- Inspect the state of all goroutines at once
- See exactly what channels are blocking
- Execute debugging commands in specific goroutine contexts
- Step through concurrent code execution
- Find the root cause of deadlocks

## What You'll Learn
- List and examine all goroutines in your program
- Identify goroutines blocked on channel operations
- Trace through deadlock chains to find root causes
- Detect and fix goroutine leaks
- Use advanced Delve features for concurrent debugging

## Prerequisites

### Install Delve
```bash
go install github.com/go-delve/delve/cmd/dlv@latest
```

### Verify Installation
```bash
dlv version
```

### Build with Debug Symbols
For optimal debugging, always build with:
```bash
go build -gcflags="all=-N -l" -o program program.go
```
- `-N`: Disable optimizations
- `-l`: Disable inlining

## Additional Resources

- [Delve Documentation](https://github.com/go-delve/delve/tree/master/Documentation)
- [Delve Command Reference](https://github.com/go-delve/delve/blob/master/Documentation/cli/README.md)
- [Debugging Go Code with Delve](https://blog.golang.org/debugging-what-you-deploy)

Remember: Delve is your window into the runtime behavior of concurrent Go programs. Master it, and you'll never fear debugging code again!
