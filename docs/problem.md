# Problem Definition

## Overview

The problem "Konstrukcija maksimalne ortogonalne CC-množice" (Construction of Maximum Orthogonal CC-Set) asks for a **Maximum Induced Matching** in an undirected graph.

Two vertices are *neighbors* if they are connected by an edge. The goal is to find the largest subset of edges such that the endpoints of any selected edge are not adjacent to the endpoints of any other selected edge.

## Formal Definition

Given graph G = (V, E) with vertex set V = {v1, v2, ..., vn} and edge set E:

Find the largest set S of edges such that for every two pairs (u1, u2), (v1, v2) in S:

```
(u1, u2), (u1, v2), (v1, u2), (v1, v2) are NOT in E
```

Note that (u1, u2) being "not in E" is automatically satisfied since S consists of edges — the condition really states that no endpoint of one selected edge may be connected to any endpoint of another selected edge.

### Equivalent Formulation

S is an **induced matching**: a matching (no shared vertices) where the subgraph induced by the matched vertices contains exactly the matching edges and no others.

## Complexity

The Maximum Induced Matching problem is NP-hard in general graphs. It remains NP-hard even on planar bipartite graphs of maximum degree 3. However, polynomial-time algorithms exist for special graph classes (trees, chordal graphs, interval graphs).

## Input Format

Input is a text file containing:

1. First line: integer `n` (number of vertices)
2. Next `n` lines: each containing `n` space-separated values (0 or 1)

The value at row `j`, column `i` is 1 if vertices `j` and `i` are connected, 0 otherwise. The matrix is symmetric (undirected graph) with zeros on the diagonal (no self-loops).

**Vertices are 1-indexed** in the input/output (row 1 = vertex 1).

### Example Input

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

This defines a graph with 7 vertices and edges:
- (1, 3), (2, 6), (3, 7), (4, 5), (4, 7)

## Output Format

Output is a text file containing:

1. First line: integer `c` (number of selected edge pairs)
2. Next `c` lines: each containing two space-separated integers `u v` (1-indexed vertex pair)

### Example Output

```
3
1 3
2 6
4 5
```

This selects 3 edges: (1,3), (2,6), (4,5).

### Verification

For the example, the induced matching property holds:
- Vertex 1's neighbors: {3} — not among {2, 6, 4, 5}
- Vertex 3's neighbors: {1, 7} — not among {2, 6, 4, 5}
- Vertex 2's neighbors: {6} — not among {1, 3, 4, 5}
- Vertex 6's neighbors: {2} — not among {1, 3, 4, 5}
- Vertex 4's neighbors: {5, 7} — not among {1, 3, 2, 6}
- Vertex 5's neighbors: {4} — not among {1, 3, 2, 6}

No cross-pair adjacencies exist, so S is a valid induced matching of size 3.

## Known Results for Special Graphs

| Graph | n | MIM Size | Reasoning |
|-------|---|----------|-----------|
| Complete Kn | n | 1 | Any two edges share neighbors |
| Path Pn | n | floor(n/3) | Every 3rd edge can be selected |
| Cycle Cn | n | floor(n/3) | Same as path for large n |
| Star K(1,k) | k+1 | 1 | All edges share the center |
| k disjoint edges | 2k | k | No adjacencies between components |
| Petersen graph | 10 | 3 | Verified by solver |
| Empty graph | n | 0 | No edges to select |
