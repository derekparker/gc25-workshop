# Demo: Order Processing Pipeline Deadlock

## Overview
This demo showcases a multi-stage order processing pipeline that has several concurrency bugs leading to deadlocks. The pipeline has four stages:

1. **Receiver** - Accepts incoming orders
2. **Validator** - Validates orders have items
3. **Processor** - Processes valid orders (2 workers)
4. **Shipper** - Ships processed orders

## The Problem
The pipeline uses unbuffered channels for communication between stages, which creates multiple potential deadlock scenarios:

1. **Channel blocking** - Unbuffered channels block senders until receivers are ready
2. **Capacity mismatch** - Shipper only handles 10 orders but we send 15
3. **Synchronous sending** - Main goroutine blocks when sending orders
4. **Chain reaction** - One blocked stage can freeze the entire pipeline

## Using Delve to Debug

### Initial Setup
```bash
# Build with debugging symbols
go build -gcflags="all=-N -l" -o pipeline pipeline.go

# Start Delve
dlv exec ./pipeline
```

### Key Delve Commands for This Demo

#### 1. Examining Goroutines
```
(dlv) goroutines
# Shows all goroutines and their current state
# Look for goroutines in "chan send" or "chan receive" state

(dlv) goroutine 5
# Switch to a specific goroutine

(dlv) bt
# Show stack trace for current goroutine
```

#### 2. Finding Blocked Channels
```
(dlv) goroutines -group chan
# Groups goroutines by what they're blocked on
# Perfect for seeing channel bottlenecks

(dlv) goroutine 3
(dlv) frame 0
(dlv) locals
# See local variables including which channel is blocking
```

#### 3. Setting Strategic Breakpoints
```
# Break when orders are sent to validation
(dlv) break pipeline.go:71

# Break in the shipper after 9 orders
(dlv) break pipeline.go:137
(dlv) condition 2 shipped == 9

# Continue execution
(dlv) continue
```

#### 4. Examining Pipeline State
```
# Switch to the main goroutine
(dlv) goroutine 1
(dlv) print pipeline
# Shows all channel states

# Check individual channels
(dlv) print pipeline.validation
(dlv) print pipeline.processing
```

#### 5. Using Goroutine Expressions
```
# Execute code in a specific goroutine context
(dlv) goroutine 5 print order.ID
(dlv) goroutine 5 call fmt.Printf("Order %d stuck\n", order.ID)
```

## Debugging Walkthrough

### Step 1: Run Until Deadlock
```
(dlv) continue
# Wait for "Pipeline timeout - possible deadlock!" message
# Then press Ctrl+C to break
```

### Step 2: Analyze Goroutine States
```
(dlv) goroutines
# You'll see something like:
# Goroutine 1 - main.main (select)
# Goroutine 5 - main.(*Pipeline).receiver (chan send) 
# Goroutine 6 - main.(*Pipeline).validator (chan send)
# Goroutine 7 - main.(*Pipeline).processor (chan send)
# Goroutine 8 - main.(*Pipeline).processor (chan send)
# Goroutine 9 - main.(*Pipeline).shipper (finished)
```

### Step 3: Trace the Deadlock Chain
```
# Check what the receiver is blocked on
(dlv) goroutine 5
(dlv) bt
(dlv) frame 0
(dlv) locals
# Shows it's trying to send to p.validation

# Check what the validator is blocked on
(dlv) goroutine 6
(dlv) locals
# Shows it's trying to send to p.processing

# Check what processors are blocked on
(dlv) goroutine 7
(dlv) locals
# Shows it's trying to send to p.shipping
```

### Step 4: Find Root Cause
```
# Check if shipper is still running
(dlv) goroutines | grep shipper
# If not listed, it has exited

# Or check shipper's last state
(dlv) goroutine 9
# May show it exited after 10 orders
```

## The Bugs

1. **Unbuffered channels** - All channels should have buffers to prevent blocking
2. **Shipper capacity** - Shipper stops at 10 orders but we send 15
3. **Missing channel closure** - Some channels aren't properly closed
4. **Synchronous order sending** - Main blocks when sending orders

## Solution Hints

To fix the deadlock:
1. Add buffers to channels: `make(chan Order, 10)`
2. Fix shipper to handle all orders or use a done signal
3. Ensure proper channel closure in the pipeline
4. Consider using select with timeout for sends

## Advanced Delve Features to Try

1. **Watchpoints** (if supported):
   ```
   (dlv) watch pipeline.shipping
   ```

2. **Tracepoints**:
   ```
   (dlv) trace main.(*Pipeline).SendOrder
   ```

3. **Display expressions**:
   ```
   (dlv) display -a order.ID
   ```

4. **Goroutine filters**:
   ```
   (dlv) goroutines -s running
   (dlv) goroutines -s 'chan send'
   ```

## Learning Objectives

After this demo, you should understand how to use Delve to:
- List and examine all goroutines in your program
- Identify goroutines blocked on channel operations
- Trace through a deadlock to find the root cause
- Execute debugging commands in specific goroutine contexts
- Set conditional breakpoints for complex scenarios