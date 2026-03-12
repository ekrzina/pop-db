#!/usr/bin/env bash
# build-and-run-linux.sh: build, package, or run PopDB backend + frontend (Linux only)

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# === parse arguments ===
PACKAGE=false
while [ "$#" -gt 0 ]; do
  case "$1" in
    package|dist) PACKAGE=true ;;
    *) echo "[warn] unknown argument: $1";;
  esac
  shift
done

echo "[info] Target OS: linux"

export CGO_ENABLED=1
mkdir -p "$SCRIPT_DIR/dbase" "$SCRIPT_DIR/backup" "$SCRIPT_DIR/bin"

# === build backend ===
BACKEND_BIN="$SCRIPT_DIR/bin/popdb"

if [ ! -x "$BACKEND_BIN" ] || [ "$SCRIPT_DIR/cmd" -nt "$BACKEND_BIN" ]; then
  echo "[backend] Building binary for Linux..."
  GOOS=linux GOARCH=amd64 go build -o "$BACKEND_BIN" ./cmd
else
  echo "[backend] Binary up to date."
fi

# === build frontend ===
cd "$SCRIPT_DIR/popdb-ui"
[ ! -d node_modules ] && npm install --production
echo "[ui] Building Next.js frontend..."
npm run build

# === package if requested ===
if [ "$PACKAGE" = true ]; then
  DIST="$SCRIPT_DIR/dist"
  echo "[package] Creating distribution in $DIST..."
  rm -rf "$DIST"
  mkdir -p "$DIST/backend" "$DIST/frontend"

  # backend
  cp -a "$BACKEND_BIN" "$DIST/backend/"
  cp -a "$SCRIPT_DIR/cmd/configs/config.yaml" "$DIST/backend/" || true
  mkdir -p "$DIST/backend/dbase" "$DIST/backend/backup"

  # frontend
  cp -a "$SCRIPT_DIR/popdb-ui/.next" "$DIST/frontend/"
  cp -a "$SCRIPT_DIR/popdb-ui/public" "$DIST/frontend/"
  cp -a "$SCRIPT_DIR/popdb-ui/package.json" "$DIST/frontend/"
  cp -a "$SCRIPT_DIR/popdb-ui/node_modules" "$DIST/frontend/"

  cp -a "$SCRIPT_DIR/run.sh" "$DIST/run.sh"
  chmod +x "$DIST/run.sh"

  echo "[package] Distribution ready!"
  exit 0
fi

# === run backend + frontend ===
cd "$SCRIPT_DIR"
echo "[backend] Starting..."
"$BACKEND_BIN" &
backend_pid=$!

echo "[ui] Starting (port 3000)..."
cd "$SCRIPT_DIR/popdb-ui"
npm run start &
ui_pid=$!

echo "[info] Backend PID: $backend_pid, UI PID: $ui_pid"
echo "Press Ctrl-C to stop both."

wait -n $backend_pid $ui_pid
kill -TERM $backend_pid $ui_pid 2>/dev/null || true