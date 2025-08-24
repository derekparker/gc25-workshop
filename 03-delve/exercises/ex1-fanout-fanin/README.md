# Exercise 2: Fan-out/Fan-in Pattern Debugging

## Problem
A data processing system uses the fan-out/fan-in pattern for parallel processing, but it's experiencing:
- Deadlocks during processing
- Goroutine leaks
- Lost data items
- Incomplete processing
- Channels blocking indefinitely

The system includes:
- Multiple workers processing data items in parallel (fan-out)
- Result collection and merging (fan-in)
- Multi-stage pipeline processing
- Error handling
- Statistics tracking

However, multiple synchronization bugs are causing the system to fail.

## Your Task
Use Delve to identify and fix all concurrency issues in the fan-out/fan-in implementation. Focus on:
1. Channel synchronization problems
2. Missing channel closures
3. Goroutine lifecycle management
4. Proper fan-in implementation
5. Deadlock resolution

## Debugging with Delve

### Setup
```bash
# Start Delve
dlv debug
```

Commands available during debug session:

  help

The above command will display all the commands, grouped logically, and explain what they do. Running `help` on a specific command will show more detailed information:

  help goroutines

Online documentation:

https://github.com/go-delve/delve/tree/master/Documentation/cli

## Bugs to Find

Use Delve to discover these issues:

1. **Unbuffered Channels**: Input and output channels cause blocking
2. **Missing Channel Closures**: Input channel is never closed, preventing worker shutdown
3. **No Result Collection**: No goroutine collecting results in some cases
4. **Error Channel Blocking**: Error channel buffer too small
5. **Goroutine Leaks**: Result collectors and merge goroutines may leak
6. **Incomplete Wait**: Workers blocked on send aren't properly waited for
7. **Race in Channel Closure**: Output channel might be closed while workers still sending

## Expected Output After Fixes
```
- Test 1: 10 items processed successfully
- Test 2: 15 items processed through multi-stage pipeline
- No timeout errors
- All goroutines properly terminated
- Consistent results without data loss
```
