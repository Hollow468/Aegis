# API Gateway Phase 1 - Benchmark Results

## Test Environment
- OS: Linux
- Go: 1.26.1

## How to Run

```bash
cd benchmark
./run_benchmarks.sh
```

## Test Scenarios

| Scenario | Route Type | Path | Concurrency | Requests |
|----------|-----------|------|-------------|----------|
| Test 1 | Exact | /api/users | 100 | 10,000 |
| Test 2 | Regex | /api/users/{id} | 100 | 10,000 |
| Test 3 | Prefix | /api/files/* | 100 | 10,000 |

## Results

*Run `./run_benchmarks.sh` to generate results.*
