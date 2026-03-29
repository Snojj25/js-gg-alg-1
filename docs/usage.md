# Usage Guide

## Prerequisites

- **Go 1.21+** for building the solver
- **Python 3** for test generation and verification

## Building

```bash
go build -o solver solver.go
```

This produces a single static binary `solver` with no external dependencies.

## Running the Solver

```bash
./solver <input_file> <output_file> [time_limit_seconds]
```

| Argument | Required | Description |
|----------|----------|-------------|
| `input_file` | Yes | Path to input graph (adjacency matrix format) |
| `output_file` | Yes | Path where solution will be written |
| `time_limit_seconds` | No | B&B time limit in seconds (default: 300) |

### Examples

```bash
# Basic usage
./solver tests/example.in tests/example.out

# With 60-second time limit
./solver tests/random_200_sparse.in output.txt 60

# Quick heuristic-only (1-second limit skips most B&B)
./solver tests/random_500_sparse.in output.txt 1
```

### Diagnostic Output

The solver prints a status line to stderr:

```
Solution size: 3, nodes: 11, time: 252.541us, timed_out: false
```

- **Solution size**: number of edges in the induced matching
- **nodes**: number of B&B nodes explored
- **time**: wall-clock time
- **timed_out**: whether the B&B was terminated early

## Input Format

Text file with:
1. First line: `n` (number of vertices, integer)
2. Next `n` lines: `n` space-separated values (0 or 1) forming the adjacency matrix

```
7
0 0 1 0 0 0 0
0 0 0 0 0 1 0
1 0 0 0 0 0 1
0 0 0 0 1 0 1
0 0 0 1 0 0 0
0 1 0 0 0 0 0
0 0 1 1 0 0 0
```

Vertices are numbered 1 through n. Row j, column i = 1 means vertices j and i are connected.

## Output Format

Text file with:
1. First line: `c` (number of edge pairs selected)
2. Next `c` lines: `u v` (1-indexed vertex pair per line)

```
3
1 3
2 6
4 5
```

## Generating Test Cases

```bash
# Generate all predefined tests into tests/ directory
python3 gen.py

# Generate into a custom directory
python3 gen.py my_tests
```

This creates 13 test files covering structured graphs (paths, cycles, stars, complete, Petersen) and random graphs of varying sizes and densities.

## Verifying Solutions

```bash
python3 verify.py <input_file> <output_file>
```

The verifier checks:
1. Every pair (u,v) in the output is an actual edge in the graph
2. No vertex appears in more than one pair (valid matching)
3. No endpoint of one pair is adjacent to any endpoint of another pair (induced property)

Output is either `VALID: induced matching of size X` or a specific error message.

### Batch Verification

```bash
for f in tests/*.in; do
    out="${f%.in}.out"
    echo -n "$(basename $f): "
    ./solver "$f" "$out" 60
    python3 verify.py "$f" "$out"
done
```

## Performance Tuning

- **Small graphs (n < 50)**: Default settings solve exactly in milliseconds.
- **Medium graphs (n = 50-100)**: A 60-second limit is usually sufficient for exact solutions on sparse graphs.
- **Large graphs (n > 100)**: The B&B will likely time out. Use a short limit (1-5 seconds) to get the heuristic result quickly, or a long limit (300+ seconds) to let the B&B improve it.
- **Dense graphs**: Even small dense graphs (n = 30, p = 0.5) may be hard for B&B. The heuristic alone often gives a near-optimal result.
