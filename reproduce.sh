#!/bin/bash
set -e

# DCC Tetragon Bridge: Automated Reproduction Script
# Validates the causal enforcement logic in a standalone environment.

echo "--- [1/3] Building Tetragon Causal Bridge ---"
mkdir -p bin
go build -o bin/tetragon-dcc src/main.go

echo "--- [2/3] Verifying Causal Logic (Demo Mode) ---"

echo "Test A: Verified PID (1234)"
RESULT_A=$(./bin/tetragon-dcc --demo --pid 1234 --intent 256)
echo "$RESULT_A"
if [[ "$RESULT_A" == *"ALLOW"* ]]; then
    echo "[PASS] Verified process allowed."
else
    echo "[FAIL] Verified process rejected."
    exit 1
fi

echo -e "\nTest B: Orphaned PID (9999)"
# Expect failure (Exit code 1)
if ./bin/tetragon-dcc --demo --pid 9999 --intent 256 > /dev/null 2>&1; then
    echo "[FAIL] Orphaned process allowed (Security Breach)."
    exit 1
else
    echo "[PASS] Orphaned process blocked (Fail-Closed)."
fi

echo -e "\n--- [3/3] Deployment Model ---"
echo "In a production BioOS environment, this bridge performs direct"
echo "eBPF map lookups via: LoadPinnedMap(\"$DCCMapPath\")"

echo -e "\nSUCCESS: Tetragon Causal Enforcement logic verified."
