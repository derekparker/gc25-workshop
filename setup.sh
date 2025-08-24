#!/bin/bash

echo "🚀 Setting up Go Concurrency Workshop..."

# Check Go installation
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21 or later."
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo "✅ Go version: $GO_VERSION"

# Check Go version (minimum 1.22 for best tracer performance)
if [[ $(echo "$GO_VERSION 1.22" | tr ' ' '\n' | sort -V | head -n1) != "1.22" ]]; then
    echo "⚠️  Go 1.22+ recommended for best execution tracer performance"
fi

# Install Delve if not present
if ! command -v dlv &> /dev/null; then
    echo "📦 Installing Delve debugger..."
    go install github.com/go-delve/delve/cmd/dlv@latest
    echo "✅ Delve installed"
else
    echo "✅ Delve already installed"
fi

echo "✅ Workshop environment ready!"
echo ""
echo "📋 Quick verification:"
echo "   go version: $(go version)"
echo "   dlv version: $(dlv version 2>/dev/null | head -n1 || echo 'Not installed')"
echo ""
echo "🎯 You're ready for the workshop!"
echo ""
echo "Next steps:"
echo "  cd race-detector/exercises/      # Start with exercises"
