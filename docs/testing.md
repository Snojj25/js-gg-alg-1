# Testing

## Test Suite

The test generator (`gen.py`) produces 13 test cases covering a range of graph structures.

### Structured Graphs

| Test File | n | Edges | Description | Expected MIM |
|-----------|---|-------|-------------|--------------|
| `example.in` | 7 | 5 | Example from the assignment PDF | 3 |
| `empty_5.in` | 5 | 0 | No edges | 0 |
| `complete_6.in` | 6 | 15 | Complete graph K6 | 1 |
| `path_10.in` | 10 | 9 | Path graph P10 | 3 |
| `path_12.in` | 12 | 11 | Path graph P12 | 4 |
| `cycle_9.in` | 9 | 9 | Cycle graph C9 | 3 |
| `star_6.in` | 6 | 5 | Star graph K(1,5) | 1 |
| `disjoint_10.in` | 10 | 5 | 5 disconnected edges | 5 |
| `petersen.in` | 10 | 15 | Petersen graph | 3 |

### Random Graphs

| Test File | n | Edges | Density | Solver MIM | Exact? |
|-----------|---|-------|---------|------------|--------|
| `random_50_sparse.in` | 50 | 55 | p=0.05 | 13 | Yes (16ms) |
| `random_100_medium.in` | 100 | 474 | p=0.10 | 17 | No (timed out) |
| `random_200_sparse.in` | 200 | 569 | p=0.03 | 41 | No (timed out) |
| `random_500_sparse.in` | 500 | 2437 | p=0.02 | 80 | No (timed out) |

Random graphs use seed 42 for reproducibility.

## Generating Tests

```bash
python3 gen.py [output_directory]
```

Default directory is `tests/`. The generator creates `.in` files for all tests listed above.

### Custom Graphs

To generate a custom random graph, use the module directly:

```python
from gen import random_graph, write_graph

edges = random_graph(n=150, p=0.08, seed=123)
write_graph("custom.in", 150, edges)
```

Available graph generators:
- `path_graph(n)` — path on n vertices
- `cycle_graph(n)` — cycle on n vertices
- `complete_graph(n)` — complete graph Kn
- `star_graph(n)` — star with n vertices (1 center + n-1 leaves)
- `disjoint_edges(k)` — k independent edges on 2k vertices
- `petersen_graph()` — Petersen graph (10 vertices)
- `random_graph(n, p, seed)` — Erdos-Renyi G(n, p)

## Verifier

```bash
python3 verify.py input.txt output.txt
```

### Checks Performed

1. **Edge validity**: Each output pair (u, v) must be an edge in the input graph.
2. **Matching property**: No vertex appears in more than one pair.
3. **Induced property**: For every two pairs (u1, u2) and (v1, v2), none of the four cross-pairs (u1,v1), (u1,v2), (u2,v1), (u2,v2) are edges.
4. **Count consistency**: Declared count `c` matches the actual number of pairs.

### Output

- Success: `VALID: induced matching of size N`
- Failure: specific error message (e.g., `INVALID: endpoint 3 of pair (3,7) adjacent to endpoint 4 of pair (4,5)`)

Exit code 0 on success, 1 on failure.

## Interpreting Results

### Exact Solutions

For small graphs (n <= 50 typically), the B&B completes within the time limit and the solution is provably optimal. The diagnostic output shows `timed_out: false`.

### Heuristic Solutions

For larger graphs, `timed_out: true` means the B&B did not finish. The reported solution is the best found — either from the greedy+local search heuristic or an improved solution discovered during partial B&B exploration. The solution is always valid but may not be optimal.

### Assessing Quality

For random graphs where the true optimum is unknown, you can:

1. Run with a longer time limit to see if the B&B can improve the result
2. Compare against the heuristic-only result (1-second time limit) to see how much the B&B improved it
3. For small instances, verify against brute force:

```python
# Brute-force for small graphs (n <= 15)
from itertools import combinations

def brute_force_mim(n, adj):
    edges = [(i, j) for i in range(n) for j in range(i+1, n) if adj[i][j]]
    best = 0
    for k in range(len(edges), 0, -1):
        for subset in combinations(edges, k):
            endpoints = set()
            valid = True
            for u, v in subset:
                if u in endpoints or v in endpoints:
                    valid = False
                    break
                endpoints.add(u)
                endpoints.add(v)
            if not valid:
                continue
            # Check induced property
            for i in range(len(subset)):
                for j in range(i+1, len(subset)):
                    for a in subset[i]:
                        for b in subset[j]:
                            if adj[a][b]:
                                valid = False
                                break
                        if not valid:
                            break
                    if not valid:
                        break
                if not valid:
                    break
            if valid:
                return k
    return 0
```
