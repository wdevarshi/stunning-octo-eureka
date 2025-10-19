#!/bin/bash

# Stress Test - 100 QPS (Queries Per Second)
# This script runs a load test at 100 requests per second for 60 seconds

set -e

DURATION="60s"
RATE="100"
OUTPUT_DIR="results"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

# Create output directory if it doesn't exist
mkdir -p "$OUTPUT_DIR"

echo "=================================================="
echo "Starting Stress Test - 100 QPS"
echo "=================================================="
echo "Duration: $DURATION"
echo "Rate: $RATE requests/second"
echo "Total Expected Requests: $((100 * 60)) requests"
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

echo "Testing READ endpoints..."
vegeta attack \
    -targets=targets.txt \
    -rate=$RATE \
    -duration=$DURATION \
    -timeout=30s \
    -workers=10 \
    | tee "$OUTPUT_DIR/results_100qps_read_${TIMESTAMP}.bin" \
    | vegeta report -type=text

echo ""
echo "=================================================="
echo "Generating detailed reports..."
echo "=================================================="

# Generate JSON report
vegeta report -type=json "$OUTPUT_DIR/results_100qps_read_${TIMESTAMP}.bin" > "$OUTPUT_DIR/report_100qps_read_${TIMESTAMP}.json"

# Generate histogram
vegeta report -type='hist[0,2ms,4ms,6ms,8ms,10ms,20ms,50ms,100ms,200ms,500ms,1s,2s]' "$OUTPUT_DIR/results_100qps_read_${TIMESTAMP}.bin"

echo ""
echo "=================================================="
echo "Test Complete!"
echo "=================================================="
echo "Results saved to:"
echo "  - $OUTPUT_DIR/results_100qps_read_${TIMESTAMP}.bin"
echo "  - $OUTPUT_DIR/report_100qps_read_${TIMESTAMP}.json"
echo ""
echo "To view results later, run:"
echo "  vegeta report $OUTPUT_DIR/results_100qps_read_${TIMESTAMP}.bin"
echo ""
echo "To generate HTML plot:"
echo "  cat $OUTPUT_DIR/results_100qps_read_${TIMESTAMP}.bin | vegeta plot > $OUTPUT_DIR/plot_100qps_read_${TIMESTAMP}.html"
echo "=================================================="
