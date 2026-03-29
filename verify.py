#!/usr/bin/env python3
"""Verifier for Maximum Induced Matching solutions."""
import sys


def verify(input_file, output_file):
    with open(input_file) as f:
        n = int(f.readline())
        adj = []
        for _ in range(n):
            row = list(map(int, f.readline().split()))
            adj.append(row)

    with open(output_file) as f:
        c = int(f.readline())
        pairs = []
        for _ in range(c):
            u, v = map(int, f.readline().split())
            pairs.append((u - 1, v - 1))  # convert to 0-indexed

    if c != len(pairs):
        print(f"INVALID: declared {c} pairs but found {len(pairs)}")
        return False

    # Check 1: Each pair is an edge
    for u, v in pairs:
        if u < 0 or u >= n or v < 0 or v >= n:
            print(f"INVALID: vertex out of range in pair ({u+1}, {v+1})")
            return False
        if adj[u][v] != 1:
            print(f"INVALID: ({u+1}, {v+1}) is not an edge")
            return False

    # Check 2: Valid matching (no shared vertices)
    used = set()
    for u, v in pairs:
        if u in used or v in used:
            print(f"INVALID: vertex reused in pair ({u+1}, {v+1})")
            return False
        used.add(u)
        used.add(v)

    # Check 3: Induced property (no adjacency between endpoints of different pairs)
    for i in range(len(pairs)):
        for j in range(i + 1, len(pairs)):
            u1, v1 = pairs[i]
            u2, v2 = pairs[j]
            for a in [u1, v1]:
                for b in [u2, v2]:
                    if adj[a][b] == 1:
                        print(
                            f"INVALID: endpoint {a+1} of pair ({u1+1},{v1+1}) "
                            f"adjacent to endpoint {b+1} of pair ({u2+1},{v2+1})"
                        )
                        return False

    print(f"VALID: induced matching of size {c}")
    return True


if __name__ == "__main__":
    if len(sys.argv) != 3:
        print(f"Usage: {sys.argv[0]} input.txt output.txt")
        sys.exit(1)
    ok = verify(sys.argv[1], sys.argv[2])
    sys.exit(0 if ok else 1)
