package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"strconv"
	"time"
)

// Bitset is a dynamically-sized bit array. nWords is fixed globally after reading n.
type Bitset []uint64

var (
	n      int
	nWords int
	adj    []Bitset
)

func newBitset() Bitset { return make(Bitset, nWords) }

func (b Bitset) Clone() Bitset {
	r := make(Bitset, len(b))
	copy(r, b)
	return r
}

func (b Bitset) Set(i int)      { b[i>>6] |= 1 << uint(i&63) }
func (b Bitset) Clear(i int)    { b[i>>6] &^= 1 << uint(i&63) }
func (b Bitset) Has(i int) bool { return b[i>>6]&(1<<uint(i&63)) != 0 }

func (b Bitset) OrWith(c Bitset) {
	for i := range b {
		b[i] |= c[i]
	}
}

func (b Bitset) AndNot(c Bitset) {
	for i := range b {
		b[i] &^= c[i]
	}
}

func (b Bitset) And(c Bitset) Bitset {
	r := newBitset()
	for i := range b {
		r[i] = b[i] & c[i]
	}
	return r
}

func (b Bitset) PopCount() int {
	c := 0
	for i := range b {
		c += bits.OnesCount64(b[i])
	}
	return c
}

func (b Bitset) IsZero() bool {
	for i := range b {
		if b[i] != 0 {
			return false
		}
	}
	return true
}

func (b Bitset) FirstSet() int {
	for i := range b {
		if b[i] != 0 {
			return i*64 + bits.TrailingZeros64(b[i])
		}
	}
	return -1
}

type Edge struct{ u, v int }

var (
	bestSize  int
	bestEdges []Edge
	curEdges  []Edge

	startTime time.Time
	timeLimit = 300 * time.Second
	timedOut  bool
	nodes     int
)

// findComponents returns connected components of the subgraph induced by active.
func findComponents(active Bitset) []Bitset {
	remaining := active.Clone()
	var comps []Bitset
	for !remaining.IsZero() {
		start := remaining.FirstSet()
		comp := newBitset()
		stack := []int{start}
		comp.Set(start)
		remaining.Clear(start)
		for len(stack) > 0 {
			u := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			nb := adj[u].And(remaining)
			for wi := 0; wi < nWords; wi++ {
				w := nb[wi]
				for w != 0 {
					bit := bits.TrailingZeros64(w)
					v := wi*64 + bit
					w &= w - 1
					comp.Set(v)
					remaining.Clear(v)
					stack = append(stack, v)
				}
			}
		}
		comps = append(comps, comp)
	}
	return comps
}

// edgeCount returns the number of edges in the subgraph induced by active.
func edgeCount(active Bitset) int {
	sumDeg := 0
	for u := 0; u < n; u++ {
		if !active.Has(u) {
			continue
		}
		sumDeg += adj[u].And(active).PopCount()
	}
	return sumDeg / 2
}

// isClique returns true if active induces a complete graph (so MIM = 1).
func isClique(active Bitset) bool {
	k := active.PopCount()
	if k < 2 {
		return false
	}
	expectedEdges := k * (k - 1) / 2
	return edgeCount(active) == expectedEdges
}

// upperBound combines several UBs for MIM on the subgraph induced by active.
//
// We use:
//   - |V_C| / 2   (each MIM edge uses 2 distinct vertices)
//   - edge_count  (each MIM edge consumes one edge)
//   - clique detection: K_t has MIM=1
//
// All three are easy and provably correct. For sparse / decomposed inputs
// these are tight at the component level (e.g. K_3 has |V|/2 = 1 = MIM).
// For dense / large connected components they can still be loose, but at
// least they never under-count.
func upperBound(active Bitset) int {
	k := active.PopCount()
	if k < 2 {
		return 0
	}
	if isClique(active) {
		return 1
	}
	ub := k / 2
	if ec := edgeCount(active); ec < ub {
		ub = ec
	}
	return ub
}

// greedySolve picks edges by minimum |N(u) ∪ N(v)| to build an induced matching
// on the subgraph induced by active.
func greedySolve(active Bitset) []Edge {
	work := active.Clone()
	var result []Edge
	for {
		bestScore := n*2 + 1
		bestU, bestV := -1, -1
		for u := 0; u < n; u++ {
			if !work.Has(u) {
				continue
			}
			nbU := adj[u].And(work)
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
					combined := adj[u].Clone()
					combined.OrWith(adj[v])
					score := combined.And(work).PopCount()
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
		remove := adj[bestU].Clone()
		remove.OrWith(adj[bestV])
		work.AndNot(remove)
	}
	return result
}

// localSearch tries (1,2)-swaps within the subgraph induced by active.
func localSearch(edges []Edge, active Bitset) []Edge {
	improved := true
	for improved {
		improved = false
		for i := 0; i < len(edges); i++ {
			forbidden := newBitset()
			for j := 0; j < len(edges); j++ {
				if j == i {
					continue
				}
				forbidden.OrWith(adj[edges[j].u])
				forbidden.OrWith(adj[edges[j].v])
			}
			avail := active.Clone()
			avail.AndNot(forbidden)
			localAvail := avail
			var newEdges []Edge
			for len(newEdges) < 2 {
				foundU, foundV := -1, -1
				fBest := n*2 + 1
				for u := 0; u < n; u++ {
					if !localAvail.Has(u) {
						continue
					}
					nbU := adj[u].And(localAvail)
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
							combined := adj[u].Clone()
							combined.OrWith(adj[v])
							score := combined.And(localAvail).PopCount()
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
				remove := adj[foundU].Clone()
				remove.OrWith(adj[foundV])
				localAvail.AndNot(remove)
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

// solveBB is the branch-and-bound routine. It mutates bestSize/bestEdges/curEdges
// to track the best found within the *current* component (caller manages saving
// these across components).
func solveBB(active Bitset) {
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

	// Reductions
	changed := true
	for changed {
		changed = false
		for u := 0; u < n; u++ {
			if !active.Has(u) {
				continue
			}
			if adj[u].And(active).IsZero() {
				active.Clear(u)
				changed = true
			}
		}
		for u := 0; u < n; u++ {
			if !active.Has(u) {
				continue
			}
			nb := adj[u].And(active)
			if nb.PopCount() != 1 {
				continue
			}
			v := nb.FirstSet()
			if adj[v].And(active).PopCount() != 1 {
				continue
			}
			curEdges = append(curEdges, Edge{u, v})
			active.Clear(u)
			active.Clear(v)
			changed = true
		}
	}

	activeCount := active.PopCount()
	if len(curEdges)+activeCount/2 <= bestSize {
		return
	}

	minDeg := n + 1
	minV := -1
	totalDeg := 0
	for u := 0; u < n; u++ {
		if !active.Has(u) {
			continue
		}
		deg := adj[u].And(active).PopCount()
		totalDeg += deg
		if deg > 0 && deg < minDeg {
			minDeg = deg
			minV = u
		}
	}
	edges := totalDeg / 2

	ub := activeCount / 2
	if edges < ub {
		ub = edges
	}
	if len(curEdges)+ub <= bestSize {
		return
	}

	if minV == -1 {
		if len(curEdges) > bestSize {
			bestSize = len(curEdges)
			bestEdges = make([]Edge, len(curEdges))
			copy(bestEdges, curEdges)
		}
		return
	}

	neighbors := adj[minV].And(active)
	afterReduction := len(curEdges)

	for wi := 0; wi < nWords; wi++ {
		w := neighbors[wi]
		for w != 0 {
			bit := bits.TrailingZeros64(w)
			nei := wi*64 + bit
			w &= w - 1
			if timedOut {
				return
			}
			newActive := active.Clone()
			remove := adj[minV].Clone()
			remove.OrWith(adj[nei])
			newActive.AndNot(remove)
			curEdges = append(curEdges[:afterReduction], Edge{minV, nei})
			solveBB(newActive)
		}
	}

	if !timedOut {
		curEdges = curEdges[:afterReduction]
		newActive := active.Clone()
		newActive.Clear(minV)
		solveBB(newActive)
	}
}

// solveComponent runs heuristic + UB checks + B&B on a single component.
// Returns the best edges found for this component.
func solveComponent(comp Bitset) []Edge {
	// Heuristic LB
	heur := greedySolve(comp)
	heur = localSearch(heur, comp)
	lb := len(heur)

	// Tight UB
	ub := upperBound(comp)

	if lb >= ub {
		// Heuristic already optimal (or component size 0/1).
		return heur
	}

	// B&B on this component, seeded with the heuristic.
	savedBest := bestSize
	savedEdges := bestEdges
	savedCurLen := len(curEdges)

	bestSize = lb
	bestEdges = make([]Edge, len(heur))
	copy(bestEdges, heur)
	curEdges = curEdges[:0]

	solveBB(comp.Clone())

	result := bestEdges

	bestSize = savedBest
	bestEdges = savedEdges
	curEdges = curEdges[:savedCurLen]

	return result
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

	fin, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	reader := bufio.NewReader(fin)
	fmt.Fscan(reader, &n)
	if n <= 0 {
		fmt.Fprintf(os.Stderr, "Invalid n=%d\n", n)
		os.Exit(1)
	}
	nWords = (n + 63) / 64
	adj = make([]Bitset, n)
	for i := 0; i < n; i++ {
		adj[i] = newBitset()
	}
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

	// Decompose the graph into connected components and solve each independently.
	full := newBitset()
	for i := 0; i < n; i++ {
		full.Set(i)
	}
	comps := findComponents(full)

	curEdges = make([]Edge, 0, n/2)
	var allEdges []Edge
	for _, comp := range comps {
		ce := solveComponent(comp)
		allEdges = append(allEdges, ce...)
	}

	finalSize := len(allEdges)

	fout, err := os.Create(os.Args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	writer := bufio.NewWriter(fout)
	fmt.Fprintln(writer, finalSize)
	for _, e := range allEdges {
		u, v := e.u, e.v
		if u > v {
			u, v = v, u
		}
		fmt.Fprintf(writer, "%d %d\n", u+1, v+1)
	}
	writer.Flush()
	fout.Close()

	elapsed := time.Since(startTime)
	fmt.Fprintf(os.Stderr, "Solution size: %d, components: %d, nodes: %d, time: %v, timed_out: %v\n",
		finalSize, len(comps), nodes, elapsed, timedOut)
}
