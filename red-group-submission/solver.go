package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"strconv"
	"time"
)

// Bitset supports up to 1024 vertices.
const maxN = 1024
const maxWords = maxN / 64

var nWords int

type Bitset [maxWords]uint64

func (b *Bitset) Set(i int)      { b[i>>6] |= 1 << uint(i&63) }
func (b *Bitset) Clear(i int)    { b[i>>6] &^= 1 << uint(i&63) }
func (b *Bitset) Has(i int) bool { return b[i>>6]&(1<<uint(i&63)) != 0 }

func (b *Bitset) OrWith(c *Bitset) {
	for i := 0; i < nWords; i++ {
		b[i] |= c[i]
	}
}

func (b *Bitset) AndNot(c *Bitset) {
	for i := 0; i < nWords; i++ {
		b[i] &^= c[i]
	}
}

func (b *Bitset) And(c *Bitset) Bitset {
	var r Bitset
	for i := 0; i < nWords; i++ {
		r[i] = b[i] & c[i]
	}
	return r
}

func (b Bitset) PopCount() int {
	c := 0
	for i := 0; i < nWords; i++ {
		c += bits.OnesCount64(b[i])
	}
	return c
}

func (b Bitset) IsZero() bool {
	for i := 0; i < nWords; i++ {
		if b[i] != 0 {
			return false
		}
	}
	return true
}

func (b Bitset) FirstSet() int {
	for i := 0; i < nWords; i++ {
		if b[i] != 0 {
			return i*64 + bits.TrailingZeros64(b[i])
		}
	}
	return -1
}

type Edge struct{ u, v int }

var (
	n   int
	adj [maxN]Bitset

	bestSize  int
	bestEdges []Edge
	curEdges  []Edge

	startTime time.Time
	timeLimit = 300 * time.Second
	timedOut  bool
	nodes     int
)

// greedySolve picks edges by minimum |N(u) ∪ N(v)| to build an induced matching.
func greedySolve() []Edge {
	var active Bitset
	for i := 0; i < n; i++ {
		active.Set(i)
	}
	var result []Edge
	for {
		bestScore := n*2 + 1
		bestU, bestV := -1, -1
		for u := 0; u < n; u++ {
			if !active.Has(u) {
				continue
			}
			nbU := adj[u].And(&active)
			for wi := 0; wi < nWords; wi++ {
				w := nbU[wi]
				// Only consider v > u to avoid duplicates.
				if wi < u>>6 {
					w = 0
				} else if wi == u>>6 {
					w &= ^((uint64(1) << uint((u&63)+1)) - 1)
				}
				for w != 0 {
					bit := bits.TrailingZeros64(w)
					v := wi*64 + bit
					w &= w - 1
					var combined Bitset
					combined = adj[u]
					combined.OrWith(&adj[v])
					score := combined.And(&active).PopCount()
					if score < bestScore {
						bestScore = score
						bestU = u
						bestV = v
					}
				}
			}
		}
		if bestU == -1 {
			break
		}
		result = append(result, Edge{bestU, bestV})
		var remove Bitset
		remove = adj[bestU]
		remove.OrWith(&adj[bestV])
		active.AndNot(&remove)
	}
	return result
}

// localSearch tries (1,2)-swaps: remove 1 edge, try to insert 2.
func localSearch(edges []Edge) []Edge {
	improved := true
	for improved {
		improved = false
		for i := 0; i < len(edges); i++ {
			// Forbidden vertices from remaining edges.
			var forbidden Bitset
			for j := 0; j < len(edges); j++ {
				if j == i {
					continue
				}
				forbidden.OrWith(&adj[edges[j].u])
				forbidden.OrWith(&adj[edges[j].v])
			}
			// Available vertices.
			var avail Bitset
			for k := 0; k < n; k++ {
				avail.Set(k)
			}
			avail.AndNot(&forbidden)
			// Greedily find induced matching edges in available subgraph.
			localAvail := avail
			var newEdges []Edge
			for len(newEdges) < 2 {
				foundU, foundV := -1, -1
				fBest := n*2 + 1
				for u := 0; u < n; u++ {
					if !localAvail.Has(u) {
						continue
					}
					nbU := adj[u].And(&localAvail)
					for wi := 0; wi < nWords; wi++ {
						w := nbU[wi]
						if wi < u>>6 {
							w = 0
						} else if wi == u>>6 {
							w &= ^((uint64(1) << uint((u&63)+1)) - 1)
						}
						for w != 0 {
							bit := bits.TrailingZeros64(w)
							v := wi*64 + bit
							w &= w - 1
							var combined Bitset
							combined = adj[u]
							combined.OrWith(&adj[v])
							score := combined.And(&localAvail).PopCount()
							if score < fBest {
								fBest = score
								foundU = u
								foundV = v
							}
						}
					}
				}
				if foundU == -1 {
					break
				}
				newEdges = append(newEdges, Edge{foundU, foundV})
				var remove Bitset
				remove = adj[foundU]
				remove.OrWith(&adj[foundV])
				localAvail.AndNot(&remove)
			}
			if len(newEdges) >= 2 {
				newSol := make([]Edge, 0, len(edges)-1+len(newEdges))
				for j := 0; j < len(edges); j++ {
					if j != i {
						newSol = append(newSol, edges[j])
					}
				}
				newSol = append(newSol, newEdges...)
				edges = newSol
				improved = true
				break
			}
		}
	}
	return edges
}

func solve(active Bitset) {
	if timedOut {
		return
	}
	nodes++
	if nodes&0xFFF == 0 {
		if time.Since(startTime) > timeLimit {
			timedOut = true
			return
		}
	}

	savedLen := len(curEdges)
	defer func() { curEdges = curEdges[:savedLen] }()

	// === Reductions ===
	changed := true
	for changed {
		changed = false
		// Remove isolated vertices (degree 0).
		for u := 0; u < n; u++ {
			if !active.Has(u) {
				continue
			}
			if adj[u].And(&active).IsZero() {
				active.Clear(u)
				changed = true
			}
		}
		// Force isolated edges (both endpoints degree 1).
		for u := 0; u < n; u++ {
			if !active.Has(u) {
				continue
			}
			nb := adj[u].And(&active)
			if nb.PopCount() != 1 {
				continue
			}
			v := nb.FirstSet()
			if adj[v].And(&active).PopCount() != 1 {
				continue
			}
			curEdges = append(curEdges, Edge{u, v})
			active.Clear(u)
			active.Clear(v)
			changed = true
		}
	}

	// === Cheap upper bound: activeCount/2 ===
	activeCount := active.PopCount()
	if len(curEdges)+activeCount/2 <= bestSize {
		return
	}

	// === Find min-degree vertex + edge count ===
	minDeg := n + 1
	minV := -1
	totalDeg := 0
	for u := 0; u < n; u++ {
		if !active.Has(u) {
			continue
		}
		deg := adj[u].And(&active).PopCount()
		totalDeg += deg
		if deg > 0 && deg < minDeg {
			minDeg = deg
			minV = u
		}
	}
	edgeCount := totalDeg / 2

	// Tighter upper bound using edge count.
	ub := activeCount / 2
	if edgeCount < ub {
		ub = edgeCount
	}
	if len(curEdges)+ub <= bestSize {
		return
	}

	// No edges remain: check if current solution is best.
	if minV == -1 {
		if len(curEdges) > bestSize {
			bestSize = len(curEdges)
			bestEdges = make([]Edge, len(curEdges))
			copy(bestEdges, curEdges)
		}
		return
	}

	neighbors := adj[minV].And(&active)
	afterReduction := len(curEdges)

	// Branch: include edge (minV, nei) for each active neighbor.
	for wi := 0; wi < nWords; wi++ {
		w := neighbors[wi]
		for w != 0 {
			bit := bits.TrailingZeros64(w)
			nei := wi*64 + bit
			w &= w - 1
			if timedOut {
				return
			}
			newActive := active
			var remove Bitset
			remove = adj[minV]
			remove.OrWith(&adj[nei])
			newActive.AndNot(&remove)
			curEdges = append(curEdges[:afterReduction], Edge{minV, nei})
			solve(newActive)
		}
	}

	// Branch: exclude minV.
	if !timedOut {
		curEdges = curEdges[:afterReduction]
		newActive := active
		newActive.Clear(minV)
		solve(newActive)
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s input.txt output.txt [time_limit_secs]\n", os.Args[0])
		os.Exit(1)
	}

	startTime = time.Now()

	if len(os.Args) >= 4 {
		if secs, err := strconv.Atoi(os.Args[3]); err == nil {
			timeLimit = time.Duration(secs) * time.Second
		}
	}

	// Read input.
	fin, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	reader := bufio.NewReader(fin)
	fmt.Fscan(reader, &n)
	nWords = (n + 63) / 64
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			var val int
			fmt.Fscan(reader, &val)
			if val == 1 {
				adj[i].Set(j)
				adj[j].Set(i)
			}
		}
	}
	fin.Close()

	// Phase 1: Greedy + Local Search.
	heuristic := greedySolve()
	heuristic = localSearch(heuristic)
	bestSize = len(heuristic)
	bestEdges = make([]Edge, len(heuristic))
	copy(bestEdges, heuristic)

	// Phase 2: Branch and Bound.
	var active Bitset
	for i := 0; i < n; i++ {
		active.Set(i)
	}
	curEdges = make([]Edge, 0, n/2)
	solve(active)

	// Write output.
	fout, err := os.Create(os.Args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	writer := bufio.NewWriter(fout)
	fmt.Fprintln(writer, bestSize)
	for _, e := range bestEdges {
		u, v := e.u, e.v
		if u > v {
			u, v = v, u
		}
		fmt.Fprintf(writer, "%d %d\n", u+1, v+1)
	}
	writer.Flush()
	fout.Close()

	elapsed := time.Since(startTime)
	fmt.Fprintf(os.Stderr, "Solution size: %d, nodes: %d, time: %v, timed_out: %v\n",
		bestSize, nodes, elapsed, timedOut)
}
