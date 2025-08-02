#!/bin/bash

# =============================================================================
# PostEaze Development Environment Deployment Script
# =============================================================================

set -e  # Exit on any error

# Configuration
PROJECT_DIR="/root/PostEaze"
BACKUP_DIR="/root/backups"
LOG_FILE="/root/deploy.log"
GIT_REPO="origin"
GIT_BRANCH="dev"
COMPOSE_FILE="docker-compose.yml"

# Load environment variables from .env
if [[ -f "$PROJECT_DIR/.env" ]]; then
    set -o allexport
    source "$PROJECT_DIR/.env"
    set +o allexport
fi

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

log() { echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')]${NC} $1" | tee -a "$LOG_FILE"; }
error() { echo -e "${RED}[ERROR $(date +'%Y-%m-%d %H:%M:%S')]${NC} $1" | tee -a "$LOG_FILE"; }
warning() { echo -e "${YELLOW}[WARNING $(date +'%Y-%m-%d %H:%M:%S')]${NC} $1" | tee -a "$LOG_FILE"; }
info() { echo -e "${BLUE}[INFO $(date +'%Y-%m-%d %H:%M:%S')]${NC} $1" | tee -a "$LOG_FILE"; }

# check_user() {
#     if [[ $EUID -eq 0 ]]; then
#         error "This script should not be run as root for security reasons"
#         exit 1
#     fi
# }

check_prerequisites() {
    log "Checking prerequisites..."
    command -v docker &>/dev/null || { error "Docker is not installed"; exit 1; }
    docker info &>/dev/null || { error "Docker is not running or permission denied"; exit 1; }
    docker compose version &>/dev/null || { error "Docker Compose is not available"; exit 1; }
    command -v git &>/dev/null || { error "Git is not installed"; exit 1; }
    log "Prerequisites check passed ✓"
}

setup_directories() {
    log "Setting up directories..."
    mkdir -p "$BACKUP_DIR"
    mkdir -p "$(dirname "$LOG_FILE")"
    log "Directories setup complete ✓"
}

check_project_directory() {
    if [[ ! -d "$PROJECT_DIR" ]]; then
        error "Project directory $PROJECT_DIR does not exist"
        exit 1
    fi
    cd "$PROJECT_DIR"
}

get_current_commit() {
    CURRENT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
    log "Current commit: $CURRENT_COMMIT"
}

backup_database() {
    log "Creating database backup..."
    local timestamp=$(date +"%Y%m%d_%H%M%S")
    local backup_file="$BACKUP_DIR/posteaze_dev_${timestamp}_${CURRENT_COMMIT}.sql"

    if docker compose ps postgres | grep -q "Up"; then
        docker compose exec -T postgres pg_dump -U "$POSTGRES_USER" -d "$POSTGRES_DB" > "$backup_file" 2>/dev/null
        if [[ $? -eq 0 && -s "$backup_file" ]]; then
            log "Database backup created: $backup_file ✓"
            gzip "$backup_file"
            log "Backup compressed: ${backup_file}.gz ✓"
            ls -t "$BACKUP_DIR"/posteaze_dev_*.sql.gz | tail -n +11 | xargs -r rm
        else
            warning "Database backup failed or empty"
            rm -f "$backup_file"
        fi
    else
        warning "PostgreSQL container is not running, skipping DB backup"
    fi
}

backup_app_files() {
    log "Creating application files backup..."
    local timestamp=$(date +"%Y%m%d_%H%M%S")
    local backup_file="$BACKUP_DIR/posteaze_app_${timestamp}_${CURRENT_COMMIT}.tar.gz"

    tar -czf "$backup_file" \
        --exclude='node_modules' \
        --exclude='.git' \
        --exclude='dist' \
        --exclude='build' \
        --exclude='*.log' \
        --exclude='tmp' \
        --exclude='temp' \
        . 2>/dev/null

    if [[ $? -eq 0 ]]; then
        log "App backup created: $backup_file ✓"
        ls -t "$BACKUP_DIR"/posteaze_app_*.tar.gz | tail -n +6 | xargs -r rm
    else
        warning "App backup failed"
    fi
}

pull_changes() {
    log "Pulling latest changes from Git..."
    if ! git diff --quiet || ! git diff --staged --quiet; then
        warning "Stashing local changes..."
        git stash push -m "Auto-stash before deploy $(date)"
    fi

    git fetch "$GIT_REPO" "$GIT_BRANCH"
    local old_commit=$(git rev-parse --short HEAD)
    git pull "$GIT_REPO" "$GIT_BRANCH"
    local new_commit=$(git rev-parse --short HEAD)

    if [[ "$old_commit" != "$new_commit" ]]; then
        log "Updated from $old_commit to $new_commit ✓"
        info "Changes:"
        git log --oneline "$old_commit..$new_commit" | head -10
    else
        log "No new changes"
    fi
}

check_migrations() {
    log "Checking for migrations..."
    local dirs=("backend/migrations" "migrations" "db/migrations")
    MIGRATION_DIR=""

    for dir in "${dirs[@]}"; do
        if [[ -d "$dir" && -n "$(ls -A "$dir" 2>/dev/null)" ]]; then
            MIGRATION_DIR="$dir"
            log "Found migration directory: $MIGRATION_DIR ✓"
            return
        fi
    done

    warning "No migration directory found"
}

run_migrations() {
    if [[ -z "$MIGRATION_DIR" ]]; then
        log "No migrations to run"
        return
    fi

    log "Running database migrations..."
    DB_URL="postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@posteaze-postgres:5432/${POSTGRES_DB}?sslmode=disable"

    docker run --rm \
        --network posteaze \
        -v "$PROJECT_DIR/$MIGRATION_DIR":/migrations \
        migrate/migrate \
        -path=/migrations \
        -database "$DB_URL" \
        up

    if [[ $? -eq 0 ]]; then
        log "Migrations applied successfully ✓"
    else
        warning "Migrations failed. Please check your migration files and DB state."
    fi
}

deploy_containers() {
    log "Deploying containers..."
    docker compose down --timeout 30
    docker compose pull
    docker compose up -d --build --remove-orphans
    log "Containers started ✓"
    sleep 30
}

health_check() {
    log "Running health checks..."
    local failed_services=()

    if ! docker compose exec backend curl -f http://localhost:8080/health &>/dev/null; then
        failed_services+=("backend")
    fi

    if ! docker compose exec postgres pg_isready -U "$POSTGRES_USER" -d "$POSTGRES_DB" &>/dev/null; then
        failed_services+=("postgres")
    fi

    if ! docker compose exec frontend nginx -t &>/dev/null; then
        failed_services+=("frontend")
    fi

    if [[ ${#failed_services[@]} -eq 0 ]]; then
        log "All health checks passed ✓"
    else
        warning "Health check failed for: ${failed_services[*]}"
        info "Use 'docker compose logs <service>' to inspect"
    fi
}

cleanup_docker() {
    log "Cleaning up Docker..."
    docker image prune -f &>/dev/null || true
    docker container prune -f &>/dev/null || true
    log "Docker cleanup done ✓"
}

rollback() {
    error "Deployment failed. Rolling back..."
    local previous_commit=$(git log --oneline -2 | tail -1 | cut -d' ' -f1)
    if [[ -n "$previous_commit" ]]; then
        warning "Rolling back to: $previous_commit"
        git checkout "$previous_commit"
        docker compose down
        docker compose up -d --build
        warning "Rollback complete"
    else
        error "Unable to determine previous commit"
    fi
}

send_notification() {
    local status=$1
    local msg="PostEaze Dev Deployment $status at $(date)"
    # Optional: Add webhook integration here
    log "Deployment $status"
}

main() {
    log "Starting PostEaze Dev Deployment"
    log "============================================="

    trap 'rollback' ERR

    # check_user
    check_prerequisites
    setup_directories
    check_project_directory
    get_current_commit
    backup_database
    backup_app_files
    pull_changes
    check_migrations
    deploy_containers
    run_migrations
    health_check
    cleanup_docker

    trap - ERR

    log "============================================="
    log "Deployment completed successfully! 🎉"
    log "URL: https://dev.posteaze.in"
    log "Log: $LOG_FILE"
    send_notification "COMPLETED"
}

main "$@"
