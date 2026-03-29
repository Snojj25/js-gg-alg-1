# Maximum Induced Matching Solver

An exact solver for the **Maximum Induced Matching (MIM)** problem, implemented for the "Konstrukcija maksimalne ortogonalne CC-množice" assignment.

## Problem

Given an undirected graph G = (V, E), find the largest subset S of edges such that for any two edges (u1, u2) and (v1, v2) in S, none of the pairs (u1, v1), (u1, v2), (u2, v1), (u2, v2) are edges in G. In other words, no endpoint of a selected edge may be adjacent to any endpoint of another selected edge.

This is NP-hard in general, so the solver combines a fast heuristic with an exact branch-and-bound algorithm.

## Approach

The solver runs in two phases:

1. **Greedy + Local Search** — builds a good initial solution in milliseconds by greedily selecting low-conflict edges, then improving via (1,2)-swaps.
2. **Branch-and-Bound** — searches for the optimal solution using min-degree branching, bitset-accelerated graph operations, and pruning against the heuristic lower bound. Falls back to the heuristic result if the time limit is reached.

All graph operations use a custom bitset implementation over `[]uint64` arrays with hardware `popcount`, making neighborhood unions, intersections, and vertex removal run in O(n/64) time.

## Quick Start

```bash
# Build
go build -o solver solver.go

# Run on an input file
./solver input.txt output.txt

# With custom time limit (seconds)
./solver input.txt output.txt 120
```

## Testing

```bash
# Generate all test cases
python3 gen.py tests

# Run solver on all tests and verify
for f in tests/*.in; do
    out="${f%.in}.out"
    ./solver "$f" "$out" 60
    python3 verify.py "$f" "$out"
done
```

## File Structure

| File | Description |
|------|-------------|
| `solver.go` | Main solver (greedy + local search + branch-and-bound) |
| `go.mod` | Go module file |
| `gen.py` | Test case generator (paths, cycles, stars, random graphs, etc.) |
| `verify.py` | Solution verifier (checks edge validity, matching, induced property) |
| `tests/` | Generated test input/output files |
| `docs/` | Detailed documentation |

## Results

| Test | n | Type | |S| | Time | Exact? |
|------|---|------|-----|------|--------|
| example | 7 | PDF example | 3 | <1ms | Yes |
| complete_6 | 6 | Complete K6 | 1 | <1ms | Yes |
| path_10 | 10 | Path | 3 | <1ms | Yes |
| path_12 | 12 | Path | 4 | <1ms | Yes |
| cycle_9 | 9 | Cycle | 3 | <1ms | Yes |
| star_6 | 6 | Star | 1 | <1ms | Yes |
| disjoint_10 | 10 | 5 disjoint edges | 5 | <1ms | Yes |
| petersen | 10 | Petersen graph | 3 | <1ms | Yes |
| random_50 | 50 | Sparse (p=0.05) | 13 | 16ms | Yes |
| random_100 | 100 | Medium (p=0.10) | 17 | 60s | Heuristic |
| random_200 | 200 | Sparse (p=0.03) | 41 | 60s | Heuristic |
| random_500 | 500 | Sparse (p=0.02) | 80 | 60s | Heuristic |

## Documentation

- [Problem Definition](docs/problem.md) — formal definition, input/output format, examples
- [Algorithm Design](docs/algorithm.md) — greedy heuristic, local search, branch-and-bound details
- [Usage Guide](docs/usage.md) — building, running, command-line options
- [Testing](docs/testing.md) — test case descriptions, generator usage, verification
