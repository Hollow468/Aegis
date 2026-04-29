#!/bin/bash
set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

BENCH_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(dirname "$BENCH_DIR")"
RESULTS_FILE="$BENCH_DIR/RESULTS.md"

echo -e "${GREEN}=== API Gateway Phase 1 Benchmark Suite ===${NC}"
echo ""

# Check for hey
if ! command -v hey &> /dev/null; then
    echo -e "${YELLOW}Warning: 'hey' not found. Install with: go install github.com/rakyll/hey@latest${NC}"
    echo "Falling back to basic curl benchmark..."
    USE_CURL=1
else
    USE_CURL=0
fi

# Build gateway
echo "Building API Gateway..."
cd "$PROJECT_DIR"
go build -o /tmp/apigateway ./cmd/apigateway/

# Start upstream servers
echo "Starting upstream servers..."
go run "$BENCH_DIR/upstream_server.go" 9001 &
UPSTREAM_PID_1=$!
go run "$BENCH_DIR/upstream_server.go" 9002 &
UPSTREAM_PID_2=$!
go run "$BENCH_DIR/upstream_server.go" 9003 &
UPSTREAM_PID_3=$!

sleep 1

# Start gateway
echo "Starting API Gateway..."
/tmp/apigateway &
GW_PID=$!
sleep 2

# Cleanup function
cleanup() {
    echo ""
    echo "Cleaning up..."
    kill $GW_PID 2>/dev/null || true
    kill $UPSTREAM_PID_1 2>/dev/null || true
    kill $UPSTREAM_PID_2 2>/dev/null || true
    kill $UPSTREAM_PID_3 2>/dev/null || true
    rm -f /tmp/apigateway
}
trap cleanup EXIT

# Initialize results file
cat > "$RESULTS_FILE" << 'HEADER'
# API Gateway Phase 1 - Benchmark Results

Generated: $(date)

## Test Environment
- OS: $(uname -sr)
- Go: $(go version)
- CPU: $(grep 'model name' /proc/cpuinfo | head -1 | cut -d: -f2 | xargs)

## Results

HEADER

# Run benchmarks
run_hey_bench() {
    local name=$1
    local url=$2
    local concurrency=$3
    local total=$4

    echo -e "${GREEN}Running: $name${NC} (c=$concurrency, n=$total)"
    echo "### $name" >> "$RESULTS_FILE"
    echo '```' >> "$RESULTS_FILE"

    if [ "$USE_CURL" -eq 0 ]; then
        hey -c "$concurrency" -n "$total" "$url" 2>&1 | tee -a "$RESULTS_FILE"
    else
        # Basic curl benchmark
        start=$(date +%s%N)
        for i in $(seq 1 "$total"); do
            curl -s "$url" > /dev/null &
            if [ $((i % concurrency)) -eq 0 ]; then
                wait
            fi
        done
        wait
        end=$(date +%s%N)
        elapsed=$(( (end - start) / 1000000 ))
        qps=$(echo "scale=2; $total * 1000 / $elapsed" | bc)
        avg_latency=$(echo "scale=2; $elapsed / $total" | bc)
        echo "Requests: $total" | tee -a "$RESULTS_FILE"
        echo "Concurrency: $concurrency" | tee -a "$RESULTS_FILE"
        echo "Total time: ${elapsed}ms" | tee -a "$RESULTS_FILE"
        echo "QPS: $qps" | tee -a "$RESULTS_FILE"
        echo "Avg latency: ${avg_latency}ms" | tee -a "$RESULTS_FILE"
    fi

    echo '```' >> "$RESULTS_FILE"
    echo "" >> "$RESULTS_FILE"
}

# Test 1: Exact match routing
run_hey_bench "Exact Match: GET /api/users" "http://localhost:8080/api/users" 100 10000

# Test 2: Regex match routing
run_hey_bench "Regex Match: GET /api/users/123" "http://localhost:8080/api/users/123" 100 10000

# Test 3: Prefix match routing
run_hey_bench "Prefix Match: GET /api/files/test.txt" "http://localhost:8080/api/files/test.txt" 100 10000

echo -e "${GREEN}=== Benchmark complete! ===${NC}"
echo "Results saved to: $RESULTS_FILE"
