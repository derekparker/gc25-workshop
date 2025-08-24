#!/bin/bash

echo "✅ Exercise Solution Checker"
echo "============================"

REPO_ROOT="$(dirname "$0")/.."

check_race_detector() {
    echo "🔍 Checking Race Detector Exercises..."
    
    cd "$REPO_ROOT/race-detector/exercises"
    
    for ex in ex1-counter ex2-loop ex3-map ex4-banking; do
        if [ -d "$ex" ]; then
            echo "Checking $ex..."
            cd "$ex"
            if go run -race main.go > /dev/null 2>&1; then
                echo "✅ $ex: No race detected"
            else
                echo "❌ $ex: Race still detected - solution needs work"
            fi
            cd ..
        fi
    done
}

check_execution_tracer() {
    echo "📊 Checking Execution Tracer Exercises..."
    
    cd "$REPO_ROOT/execution-tracer/exercises"
    
    for ex in ex1-scheduling ex2-io; do
        if [ -d "$ex" ]; then
            echo "Checking $ex..."
            cd "$ex"
            if go run main.go > /dev/null 2>&1; then
                echo "✅ $ex: Runs successfully"
                if [ -f "*.trace" ]; then
                    echo "✅ $ex: Trace file generated"
                fi
            else
                echo "❌ $ex: Runtime error"
            fi
            cd ..
        fi
    done
}

check_delve_skills() {
    echo "🐛 Delve skills cannot be automatically checked"
    echo "ℹ️  Practice these commands:"
    echo "   - goroutines"
    echo "   - break <function>"
    echo "   - print <variable>"
    echo "   - continue"
    echo "   - stack"
}

echo "Checking exercise solutions..."
echo ""

check_race_detector
echo ""
check_execution_tracer
echo ""
check_delve_skills

echo ""
echo "💡 Tip: Check solutions/ directory for reference implementations"
