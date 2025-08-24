# Exercise 1: Statistics Tracker Race Conditions

## Problem
A statistics tracking system is being used by multiple workers to record their progress. The system includes:
- Multiple workers processing items in batches
- A progress monitor that reports statistics in real-time
- An audit logger that takes periodic snapshots
- Methods to record work and query current totals
- Time tracking to monitor when stats were last updated

However, the implementation has multiple race conditions that cause incorrect results and data races.

## Your Task
Fix all race conditions in the code. The system has several types of shared state:
1. **Simple counters** (`processed`, `workers`)
2. **Time tracking** (`lastUpdated`)
3. **Multiple readers** - progress monitor and audit logger read frequently, can we optimize for this use case?

## Testing
```bash
# Detect the race conditions
go run -race main.go

# You should see multiple race warnings for:
# - Stats.processed field (counter race)
# - Stats.workers field (counter race)
# - Stats.lastUpdated field (time.Time race - cannot be fixed with atomics!)

# Verify your solution
go run -race main.go  # Should show no race warnings
```

## Expected Output
- Workers used: 5
- Total items processed: 5000 (should always be exactly 5000)
- No race condition warnings
- Progress monitor should show smooth progression
- Audit logs should show consistent snapshots
- Last update time should be valid
