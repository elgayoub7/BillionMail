#!/bin/bash
# =============================================================================
# BillionMail Custom Build - Deploy Script
# Usage: ./deploy.sh [build|up|down|logs|update]
# =============================================================================

set -e

COMPOSE_FILE="docker-compose.custom.yml"
IMAGE_NAME="billionmail-custom"

case "${1:-help}" in
  build)
    echo "=== Building custom BillionMail image ==="
    echo "This will compile Go backend + build Vue frontend (may take 10-15 min)..."
    docker compose -f "$COMPOSE_FILE" build core-billionmail
    echo "=== Build complete ==="
    ;;

  up)
    echo "=== Starting BillionMail ==="
    docker compose -f "$COMPOSE_FILE" up -d
    echo ""
    echo "Waiting for services to initialize..."
    sleep 5
    echo ""
    echo "=== Services Status ==="
    docker compose -f "$COMPOSE_FILE" ps
    echo ""
    echo "Access the panel at: https://$(hostname -I 2>/dev/null | awk '{print $1}' || echo 'YOUR_SERVER_IP')"
    echo "Safe path: $(grep SafePath .env 2>/dev/null | cut -d= -f2 || echo 'billion')"
    echo ""
    echo "Admin username: $(grep ADMIN_USERNAME .env 2>/dev/null | cut -d= -f2 || echo 'Check .env file')"
    echo "Admin password: ****** (see .env)"
    ;;

  down)
    echo "=== Stopping BillionMail ==="
    docker compose -f "$COMPOSE_FILE" down
    ;;

  logs)
    SERVICE="${2:-core-billionmail}"
    docker compose -f "$COMPOSE_FILE" logs -f "$SERVICE"
    ;;

  update)
    echo "=== Pulling latest code ==="
    git pull origin main
    echo ""
    echo "=== Rebuilding custom image ==="
    docker compose -f "$COMPOSE_FILE" build core-billionmail
    echo ""
    echo "=== Restarting with new image ==="
    docker compose -f "$COMPOSE_FILE" up -d core-billionmail
    echo "=== Update complete ==="
    ;;

  init)
    if [ ! -f .env ]; then
      echo "=== Initializing .env from template ==="
      cp env_init .env
      echo "IMPORTANT: Edit .env and change ALL passwords before running 'deploy.sh up'"
      echo ""
      echo "Generate strong passwords:"
      echo "  DBPASS:       $(openssl rand -hex 16)"
      echo "  REDISPASS:    $(openssl rand -hex 16)"
      echo "  JWT_SECRET:   $(openssl rand -hex 32)"
      echo ""
      echo "Then set BILLIONMAIL_HOSTNAME to your mail server hostname (e.g. mail.yourdomain.com)"
    else
      echo ".env already exists. Edit it directly."
    fi
    ;;

  *)
    echo "BillionMail Custom Build Manager"
    echo ""
    echo "Usage: ./deploy.sh [command]"
    echo ""
    echo "Commands:"
    echo "  init     - Create .env from template (first time only)"
    echo "  build    - Build custom Docker image (Go + Vue)"
    echo "  up       - Start all services"
    echo "  down     - Stop all services"
    echo "  logs     - View logs (option: service name)"
    echo "  update   - Pull code + rebuild + restart"
    echo ""
    echo "First time setup:"
    echo "  1. ./deploy.sh init"
    echo "  2. nano .env          # Edit passwords + hostname"
    echo "  3. ./deploy.sh build  # Build image (~10-15 min)"
    echo "  4. ./deploy.sh up     # Start everything"
    ;;
esac
