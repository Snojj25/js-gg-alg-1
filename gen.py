#!/usr/bin/env python3
"""Test case generator for Maximum Induced Matching."""
import random
import sys
import os


def write_graph(filename, n, edges):
    adj = [[0] * n for _ in range(n)]
    for u, v in edges:
        adj[u][v] = 1
        adj[v][u] = 1
    with open(filename, "w") as f:
        f.write(f"{n}\n")
        for row in adj:
            f.write(" ".join(map(str, row)) + "\n")


def path_graph(n):
    return [(i, i + 1) for i in range(n - 1)]


def cycle_graph(n):
    return [(i, (i + 1) % n) for i in range(n)]


def complete_graph(n):
    return [(i, j) for i in range(n) for j in range(i + 1, n)]


def star_graph(n):
    return [(0, i) for i in range(1, n)]


def disjoint_edges(k):
    return [(2 * i, 2 * i + 1) for i in range(k)]


def petersen_graph():
    outer = [(i, (i + 1) % 5) for i in range(5)]
    inner = [(5 + i, 5 + (i + 2) % 5) for i in range(5)]
    spokes = [(i, i + 5) for i in range(5)]
    return outer + inner + spokes


def random_graph(n, p, seed=None):
    if seed is not None:
        random.seed(seed)
    edges = []
    for i in range(n):
        for j in range(i + 1, n):
            if random.random() < p:
                edges.append((i, j))
    return edges


def generate_all(directory):
    os.makedirs(directory, exist_ok=True)

    tests = [
        ("example", 7, [(0, 2), (1, 5), (2, 6), (3, 4), (3, 6)]),
        ("empty_5", 5, []),
        ("complete_6", 6, complete_graph(6)),
        ("path_10", 10, path_graph(10)),
        ("path_12", 12, path_graph(12)),
        ("cycle_9", 9, cycle_graph(9)),
        ("star_6", 6, star_graph(6)),
        ("disjoint_10", 10, disjoint_edges(5)),
        ("petersen", 10, petersen_graph()),
    ]

    for name, n, edges in tests:
        write_graph(f"{directory}/{name}.in", n, edges)
        print(f"  {name}: n={n}, m={len(edges)}")

    random_tests = [
        ("random_50_sparse", 50, 0.05, 42),
        ("random_100_medium", 100, 0.10, 42),
        ("random_200_sparse", 200, 0.03, 42),
        ("random_500_sparse", 500, 0.02, 42),
    ]

    for name, n, p, seed in random_tests:
        edges = random_graph(n, p, seed)
        write_graph(f"{directory}/{name}.in", n, edges)
        print(f"  {name}: n={n}, m={len(edges)}, p={p}")


if __name__ == "__main__":
    directory = sys.argv[1] if len(sys.argv) > 1 else "tests"
    print(f"Generating test cases in {directory}/")
    generate_all(directory)
    print("Done.")
