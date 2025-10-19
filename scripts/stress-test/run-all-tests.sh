#!/bin/bash

# Run all stress tests sequentially
# This script runs both 10 QPS and 100 QPS tests

set -e

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

echo "=================================================="
echo "Transport Analytics API - Complete Stress Test"
echo "=================================================="
echo ""

# Check if vegeta is installed
if ! command -v vegeta &> /dev/null; then
    echo "Error: vegeta is not installed"
    echo "Install with: brew install vegeta (macOS) or go install github.com/tsenart/vegeta@latest"
    exit 1
fi

# Check if the API is running
if ! curl -s http://localhost:9091/health > /dev/null 2>&1; then
    echo "Error: API is not running at http://localhost:9091"
    echo "Please start the application first with: ./run.sh or docker-compose up"
    exit 1
fi

cd "$SCRIPT_DIR"

echo "API is running. Starting test suite..."
echo ""

# Run 10 QPS test
echo "=================================================="
echo "Test 1/2: Running 10 QPS test..."
echo "=================================================="
./run-10qps.sh

echo ""
echo "Waiting 10 seconds before next test..."
sleep 10
echo ""

# Run 100 QPS test
echo "=================================================="
echo "Test 2/2: Running 100 QPS test..."
echo "=================================================="
./run-100qps.sh

echo ""
echo "=================================================="
echo "All Tests Complete!"
echo "=================================================="
echo ""
echo "Results are saved in the results/ directory"
echo ""
echo "To compare results:"
echo "  ls -lht results/"
echo ""
echo "To view latest 10 QPS results:"
echo "  vegeta report \$(ls -t results/results_10qps_*.bin | head -1)"
echo ""
echo "To view latest 100 QPS results:"
echo "  vegeta report \$(ls -t results/results_100qps_*.bin | head -1)"
echo ""
echo "=================================================="
