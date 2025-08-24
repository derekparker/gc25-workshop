#!/bin/bash

echo "🧹 Cleaning up workshop files..."

# Remove generated files
find . -name "*.trace" -delete
find . -name "*.prof" -delete
find . -name "debug" -delete
find . -name "__debug_bin*" -delete

# Remove temporary go files
find . -name "*_broken.go" -delete
find . -name "*_debug.go" -delete
find . -name "hello_tracer.go" -delete
find . -name "channels_tracer.go" -delete

echo "✅ Cleanup complete"
