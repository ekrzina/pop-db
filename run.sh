#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND_DIR="$SCRIPT_DIR/backend"
BACKEND="$BACKEND_DIR/popdb"
FRONTEND="$SCRIPT_DIR/frontend"

chmod +x "$BACKEND"

# Start backend from its folder
echo "[backend] Starting..."
cd "$BACKEND_DIR"
./popdb &
backend_pid=$!

# Start frontend
echo "[frontend] Starting (port 3000)..."
cd "$FRONTEND"
npm run start &
frontend_pid=$!

echo "[info] Backend PID: $backend_pid, Frontend PID: $frontend_pid"
echo "Press Ctrl-C to stop both."

wait -n $backend_pid $frontend_pid
kill -TERM $backend_pid $frontend_pid 2>/dev/null || true