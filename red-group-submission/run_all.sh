#!/bin/bash
# Požene solver na vseh testih in preveri pravilnost rešitev.

SOLVER="./solver"
VERIFY="verify.py"
TEST_DIR="tests"
TIME_LIMIT=60

if [ ! -f "$SOLVER" ]; then
    echo "Solver ni preveden. Prevajam..."
    go build -o solver solver.go || exit 1
fi

pass=0
fail=0

for input in "$TEST_DIR"/*.in; do
    name=$(basename "$input" .in)
    output="$TEST_DIR/$name.out"

    echo -n "$name: "
    "$SOLVER" "$input" "$output" "$TIME_LIMIT" 2>/dev/null
    result=$(python3 "$VERIFY" "$input" "$output" 2>&1)
    echo "$result"

    if echo "$result" | grep -q "^VALID"; then
        pass=$((pass + 1))
    else
        fail=$((fail + 1))
    fi
done

echo ""
echo "Rezultat: $pass uspešnih, $fail neuspešnih"
