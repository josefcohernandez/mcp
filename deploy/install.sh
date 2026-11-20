#!/usr/bin/env bash
# Deploy the DevFest MCP server (workshop step 05) on a Linux VM with Docker.
#
# Usage (as a user with sudo):
#   curl -fsSL https://raw.githubusercontent.com/josefcohernandez/mcp/main/deploy/install.sh -o install.sh
#   less install.sh                      # read it before running it
#   bash install.sh                      # behind your existing HTTPS reverse proxy
#   bash install.sh --domain mcp.example.com   # bundled Caddy gets a Let's Encrypt certificate
#
# Options:
#   --domain HOST   Public hostname. Starts Caddy on ports 80/443 with automatic HTTPS.
#                   Omit it if Nginx/Caddy/Traefik already serves HTTPS on this VM.
#   --dir PATH      Install directory (default: /opt/devfest-mcp)
#   --branch NAME   Git branch to deploy (default: main)
#   --port N        Local port for the MCP container (default: 8080, bound to 127.0.0.1)
#   --repo URL      Git repository (default: https://github.com/josefcohernandez/mcp.git)
#   --yes           Do not ask before installing Docker
#   --uninstall     Stop and remove the containers (keeps the checkout)
set -euo pipefail

REPO="https://github.com/josefcohernandez/mcp.git"
BRANCH="main"
DIR="/opt/devfest-mcp"
PORT="8080"
DOMAIN=""
ASSUME_YES=0
UNINSTALL=0

log()  { printf '\033[1;34m==>\033[0m %s\n' "$*"; }
warn() { printf '\033[1;33m!!\033[0m  %s\n' "$*" >&2; }
die()  { printf '\033[1;31mxx\033[0m  %s\n' "$*" >&2; exit 1; }

while [ $# -gt 0 ]; do
  case "$1" in
    --domain)    DOMAIN="${2:?--domain needs a value}"; shift 2 ;;
    --dir)       DIR="${2:?--dir needs a value}"; shift 2 ;;
    --branch)    BRANCH="${2:?--branch needs a value}"; shift 2 ;;
    --port)      PORT="${2:?--port needs a value}"; shift 2 ;;
    --repo)      REPO="${2:?--repo needs a value}"; shift 2 ;;
    --yes|-y)    ASSUME_YES=1; shift ;;
    --uninstall) UNINSTALL=1; shift ;;
    -h|--help)   sed -n '2,22p' "$0"; exit 0 ;;
    *) die "unknown option: $1 (see --help)" ;;
  esac
done

[ "$(uname -s)" = "Linux" ] || die "this script is for Linux"
case "$PORT" in ''|*[!0-9]*) die "--port must be a number" ;; esac

SUDO=""
if [ "$(id -u)" -ne 0 ]; then
  command -v sudo >/dev/null || die "run as root or install sudo"
  SUDO="sudo"
fi

# docker may need sudo until the user is in the docker group.
DOCKER="docker"
if ! docker info >/dev/null 2>&1; then DOCKER="$SUDO docker"; fi
compose() {
  (cd "$DIR" && MCP_PORT="$PORT" DOMAIN="$DOMAIN" $DOCKER compose -f deploy/compose.yaml "$@")
}

if [ "$UNINSTALL" -eq 1 ]; then
  [ -d "$DIR" ] || die "$DIR not found"
  log "Stopping containers"
  compose --profile https down
  log "Done. Checkout kept in $DIR (remove it with: $SUDO rm -rf $DIR)"
  exit 0
fi

# ---------------------------------------------------------------- prerequisites
install_pkgs() {
  if command -v apt-get >/dev/null; then
    $SUDO apt-get update -qq && $SUDO apt-get install -y -qq "$@"
  elif command -v dnf >/dev/null; then
    $SUDO dnf install -y -q "$@"
  elif command -v yum >/dev/null; then
    $SUDO yum install -y -q "$@"
  else
    die "install manually: $*"
  fi
}

command -v git  >/dev/null || { log "Installing git and curl"; install_pkgs git curl; }
command -v curl >/dev/null || install_pkgs curl

if ! command -v docker >/dev/null; then
  warn "Docker is not installed. It will be installed with Docker's official script (https://get.docker.com)."
  if [ "$ASSUME_YES" -ne 1 ]; then
    read -r -p "Continue? [y/N] " ans
    case "$ans" in y|Y|s|S) ;; *) die "aborted; install Docker and run again" ;; esac
  fi
  curl -fsSL https://get.docker.com -o /tmp/get-docker.sh
  $SUDO sh /tmp/get-docker.sh
  $SUDO systemctl enable --now docker 2>/dev/null || true
  DOCKER="$SUDO docker"
fi
$DOCKER compose version >/dev/null 2>&1 || die "Docker Compose v2 plugin missing (package docker-compose-plugin)"

# ------------------------------------------------------------- port pre-checks
port_busy() { (command -v ss >/dev/null && ss -ltnH "sport = :$1" | grep -q .) ; }
if port_busy "$PORT" && ! $DOCKER ps --format '{{.Names}}' | grep -qx devfest-mcp; then
  die "port $PORT is already in use; choose another with --port"
fi
if [ -n "$DOMAIN" ]; then
  for p in 80 443; do
    if port_busy "$p" && ! $DOCKER ps --format '{{.Names}}' | grep -qx devfest-caddy; then
      die "port $p is in use (an existing web server?). Run without --domain and add a proxy rule instead (see deploy/README.md)"
    fi
  done
fi

# ------------------------------------------------------------------- checkout
if [ -d "$DIR/.git" ]; then
  log "Updating $DIR ($BRANCH)"
  $SUDO git -C "$DIR" fetch -q origin "$BRANCH"
  $SUDO git -C "$DIR" checkout -q -B "$BRANCH" FETCH_HEAD
else
  log "Cloning $REPO ($BRANCH) into $DIR"
  $SUDO git clone -q --depth 1 --branch "$BRANCH" "$REPO" "$DIR"
fi
[ -f "$DIR/deploy/compose.yaml" ] || die "deploy/compose.yaml not found on branch $BRANCH"

# --------------------------------------------------------------------- start
if [ -n "$DOMAIN" ]; then
  log "Building and starting devfest-mcp + Caddy (HTTPS for $DOMAIN)"
  compose --profile https up -d --build
else
  log "Building and starting devfest-mcp on 127.0.0.1:$PORT"
  compose up -d --build devfest-mcp
fi

log "Waiting for /health"
for _ in $(seq 1 30); do
  if curl -fsS "http://127.0.0.1:$PORT/health" >/dev/null 2>&1; then break; fi
  sleep 1
done
curl -fsS "http://127.0.0.1:$PORT/health" >/dev/null || { compose logs --tail 50; die "server did not become healthy"; }

TOOLS=$(curl -fsS "http://127.0.0.1:$PORT/mcp" \
  -H 'Content-Type: application/json' -H 'Accept: application/json, text/event-stream' \
  -d '{"jsonrpc":"2.0","id":1,"method":"tools/list"}' | grep -o '"name":"[a-z_]*"' | tr '\n' ' ')
log "Local MCP OK: $TOOLS"

if [ -n "$DOMAIN" ]; then
  log "Checking https://$DOMAIN/health (certificate issuance can take ~30 s)"
  ok=0
  for _ in $(seq 1 30); do
    if curl -fsS "https://$DOMAIN/health" >/dev/null 2>&1; then ok=1; break; fi
    sleep 2
  done
  if [ "$ok" -eq 1 ]; then
    log "Public HTTPS OK"
  else
    warn "https://$DOMAIN is not answering yet. Check that DNS points here and ports 80/443 are open:"
    warn "  $DOCKER logs devfest-caddy"
  fi
  URL="https://$DOMAIN/mcp"
else
  URL="https://<your-domain>/mcp"
  warn "No --domain given: add a rule in your reverse proxy to 127.0.0.1:$PORT (Nginx/Caddy/Traefik examples in $DIR/deploy/README.md)."
fi

cat <<EOF

Done. Add it to Claude (web, desktop and then the Android app):
  claude.ai -> Settings -> Connectors -> Add custom connector
  URL: $URL

Update later:   bash $DIR/deploy/install.sh ${DOMAIN:+--domain $DOMAIN }--branch $BRANCH
Logs:           $DOCKER logs -f devfest-mcp
Uninstall:      bash $DIR/deploy/install.sh --uninstall
EOF
