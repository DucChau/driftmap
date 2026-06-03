#!/usr/bin/env bash
# Runs a quick demo of driftmap's core features.
set -euo pipefail

BIN="./driftmap"

if [ ! -f "$BIN" ]; then
  echo "Building driftmap first..."
  go build -o driftmap .
fi

echo "============================="
echo "  driftmap demo"
echo "============================="
echo ""

echo "--- 1. Capture baseline snapshot ---"
$BIN capture demo-baseline
echo ""

echo "--- 2. Set a new env var and capture again ---"
export DRIFTMAP_DEMO_VAR="hello-from-demo"
$BIN capture demo-after
echo ""

echo "--- 3. List all snapshots ---"
$BIN list
echo ""

echo "--- 4. Diff baseline vs latest ---"
$BIN diff demo-baseline demo-after
echo ""

echo "--- 5. Replay baseline into 'env' subprocess ---"
echo "(showing only DRIFTMAP_DEMO_VAR; should be empty in baseline)"
$BIN replay demo-baseline -- sh -c 'echo DRIFTMAP_DEMO_VAR=${DRIFTMAP_DEMO_VAR:-<not set>}'
echo ""

echo "Demo complete!"
