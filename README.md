# Go Concurrency: Debugging Goroutines and Channels
## Gophercon 2025 Workshop

A comprehensive 4-hour workshop covering Go's concurrency debugging tools: Race Detector, Execution Tracer, and Delve Debugger.

## Quick Start

```bash
# Setup environment
./setup.sh
```

## Workshop Structure

- **Part I**: Race Detector
- **Part II**: Execution Tracer
- **Part III**: Delve Debugger

## Prerequisites

- Go 1.24+
- Basic Go experience (6+ months)
- Familiarity with goroutines and channels

Also note that if WIFI becomes an issue for some examples which make HTTP requests,
this should be solvable by running a local httpbin container:

  docker run -p 80:80 kennethreitz/httpbin

## Repository Structure

```
race-detector/      # Race detection examples and exercises
execution-tracer/   # Execution tracing examples and microservice
delve/             # Delve debugging examples
scripts/           # Demo and utility scripts
```

## Getting Help

- Run `./scripts/exercise_checker.sh` to validate solutions
- Use `./cleanup.sh` to reset environment

Happy debugging! 🐛➡️✅
