# Algorithm Design

The solver uses a two-phase approach: a fast heuristic for a good initial solution, followed by an exact branch-and-bound that provably finds the optimum (or returns the heuristic result if time runs out).

## Data Structures

### Bitset

All graph operations are built on a fixed-size bitset type:

```go
type Bitset [maxWords]uint64  // maxWords = 16, supports up to 1024 vertices
```

Key operations and their complexity (where w = ceil(n/64)):

| Operation | Description | Time |
|-----------|-------------|------|
| `Set(i)` | Set bit i | O(1) |
| `Clear(i)` | Clear bit i | O(1) |
| `Has(i)` | Test bit i | O(1) |
| `OrWith(c)` | Bitwise OR in-place | O(w) |
| `AndNot(c)` | Clear all bits in c | O(w) |
| `And(c)` | Bitwise AND, return new | O(w) |
| `PopCount()` | Count set bits | O(w) |
| `IsZero()` | Check if empty | O(w) |
| `FirstSet()` | Lowest set bit index | O(w) |

`PopCount` uses `math/bits.OnesCount64` which compiles to a single hardware `POPCNT` instruction on modern CPUs.

### Graph Representation

- **Adjacency bitsets**: `adj[v]` is a Bitset where bit `u` is set iff edge (v, u) exists
- **Active vertex set**: a Bitset tracking which vertices are still candidates during search
- Vertex operations like "remove v and all its neighbors" become single bitwise operations:
  ```go
  active.AndNot(&adj[v])  // removes v's entire neighborhood in O(w)
  ```

## Phase 1: Greedy Heuristic

The greedy builds an initial induced matching by repeatedly selecting the "cheapest" edge.

### Algorithm

```
1. active = all vertices
2. while edges exist in active subgraph:
   a. For each edge (u,v) in active, compute score = |N(u) ∪ N(v) ∩ active|
   b. Select edge (u,v) with minimum score
   c. Add (u,v) to matching S
   d. Remove N(u) ∪ N(v) from active   (all neighbors of both endpoints)
3. return S
```

### Intuition

The score |N(u) ∪ N(v)| measures how many vertices we "consume" by selecting edge (u,v). Minimizing this preserves the most vertices for future edges.

### Complexity

O(m * k) where m = number of edges and k = size of resulting matching. Each iteration scans all remaining edges. With bitset operations, the inner scoring is O(w) per edge.

## Phase 1b: Local Search

After the greedy, a (1,2)-swap local search attempts to improve the solution.

### Algorithm

```
repeat:
  for each edge e in S:
    1. Temporarily remove e from S
    2. Compute forbidden = union of N(u) ∪ N(v) for all remaining edges
    3. available = all vertices not in forbidden
    4. Greedily find induced matching edges in available subgraph
    5. If 2+ edges found: accept swap (net gain >= +1), restart
  if no improving swap found: stop
```

### Why (1,2)-swaps

Removing 1 edge and inserting 2 gives a net gain of +1. The freed vertices (and their neighborhoods) may enable two edges where only one was before. This is a common improvement strategy for matching problems.

## Phase 2: Branch-and-Bound

The exact solver systematically explores all possible induced matchings, pruning branches that cannot improve on the best known solution.

### Reduction Rules

Before branching, the solver applies two safe reductions:

1. **Isolated vertex removal**: Vertices with degree 0 in the active subgraph cannot participate in any edge. Remove them.

2. **Isolated edge forcing**: If vertices u and v both have degree 1 and are each other's only neighbor, edge (u,v) must be included in any optimal solution (excluding either vertex wastes it). Force (u,v) into the matching and remove both vertices.

These reductions are applied in a loop until no more apply (each reduction may create new isolated vertices or edges).

**Important note**: The general pendant reduction (degree-1 vertex with a higher-degree neighbor) is NOT safe for MIM. Counterexample:

```
v - u - a - b
        |
    u - c - d
```

Forcing edge (v,u) gives matching {(v,u)}, size 1. But the optimal is {(a,b), (c,d)}, size 2. Only isolated edges (both endpoints degree 1) can be safely forced.

### Branching Strategy

The solver branches on the vertex v with **minimum degree** in the active subgraph:

- **Include branches**: For each neighbor w of v, include edge (v,w) in the matching. Remove N(v) ∪ N(w) from the active set (all neighbors of both endpoints).
- **Exclude branch**: Exclude v entirely (v will not be an endpoint of any matching edge). Remove only v from active.

This creates deg(v) + 1 subproblems. Choosing the minimum-degree vertex minimizes the branching factor.

Include branches are tried before the exclude branch because they increase the matching size, leading to stronger pruning.

### Pruning

Two upper bounds are computed at each node:

1. **Vertex bound**: `floor(|active| / 2)` — each induced matching edge needs 2 vertices.
2. **Edge bound**: `|edges in active subgraph|` — can't select more edges than exist. Computed as `sum(degrees) / 2` during the min-degree scan (no extra cost).

The tighter of the two is used:

```
ub = min(activeCount / 2, edgeCount)
if currentSize + ub <= bestSize: prune
```

The cheap vertex bound is checked first (single PopCount call) to skip obviously dead branches before the more expensive min-degree scan.

### Time Management

The solver checks `time.Since(startTime) > timeLimit` every 4096 nodes (controlled by `nodes & 0xFFF == 0`). When the time limit is reached, the `timedOut` flag propagates through the recursion and the solver returns the best solution found so far (at minimum, the heuristic result).

### Recursion Management

The solver uses a global `curEdges` slice with save/restore semantics:

```go
savedLen := len(curEdges)
defer func() { curEdges = curEdges[:savedLen] }()
```

Each recursive call saves the current edge count and restores it on return. Forced edges from reductions and branching decisions are appended/truncated efficiently without allocation.

## Performance Characteristics

| Graph Size | Heuristic | Branch-and-Bound |
|------------|-----------|------------------|
| n <= 20 | < 1ms | < 1ms (exact) |
| n = 50, sparse | < 1ms | ~16ms (exact) |
| n = 100, medium | < 1ms | Times out (heuristic used) |
| n = 200, sparse | < 1ms | Times out (heuristic used) |
| n = 500, sparse | < 10ms | Times out (heuristic used) |

The heuristic phase always completes in milliseconds regardless of graph size. The B&B phase solves small-to-medium instances exactly and gracefully degrades to the heuristic result for large instances.
