# Exercise 4: Banking System Challenge

## Problem
A complete banking system with multiple race conditions across accounts, balances, and operations.

## Your Task
Fix ALL race conditions in the banking system. Consider:
- Account creation and ID generation
- Balance access and modification
- Map operations
- Transfer operations (potential deadlocks!)

## Testing
```bash
# Detect all races
go run -race main.go

# Verify your solution
go run -race main.go  # Should complete without race warnings
```

## Hints
- Consider deadlock prevention in transfers
- Some operations might benefit from read-write locks
- Think about locking granularity

## Expected Output
Should complete without race warnings and maintain consistent total balance.
