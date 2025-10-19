#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Configuration
DB_PORT=5433
APP_HTTP_PORT=9091
APP_GRPC_PORT=9090
FRONTEND_PORT=3000
NGINX_PORT=8080
DB_CONTAINER="transport_db"
APP_CONTAINER="transport_app"
FRONTEND_CONTAINER="transport_frontend"
NGINX_CONTAINER="transport_nginx"

# Print colored message
print_msg() {
    color=$1
    shift
    echo -e "${color}$@${NC}"
}

print_success() {
    print_msg "$GREEN" "✓ $@"
}

print_error() {
    print_msg "$RED" "✗ $@"
}

print_info() {
    print_msg "$BLUE" "ℹ $@"
}

print_warning() {
    print_msg "$YELLOW" "⚠ $@"
}

print_step() {
    print_msg "$CYAN" "===> $@"
}

# Print banner
print_banner() {
    echo ""
    print_info "================================================"
    print_info "  Transport Reliability Analytics System"
    print_info "================================================"
    echo ""
}

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check prerequisites
check_prerequisites() {
    print_step "Checking prerequisites..."
    local missing=0

    # Check Docker
    if ! command_exists docker; then
        print_error "Docker is not installed"
        print_info "Install Docker from: https://docs.docker.com/get-docker/"
        missing=1
    else
        # Check if Docker daemon is running
        if ! docker ps >/dev/null 2>&1; then
            print_error "Docker is installed but not running"
            print_info "Please start Docker Desktop"
            if [[ "$OSTYPE" == "darwin"* ]]; then
                print_info "Run: open -a Docker"
            fi
            exit 1
        fi
        print_success "Docker is installed and running"
    fi

    # Check Docker Compose
    if ! command_exists docker-compose && ! docker compose version >/dev/null 2>&1; then
        print_error "Docker Compose is not installed"
        missing=1
    else
        print_success "Docker Compose is installed"
    fi

    # Check Go (for local mode)
    if command_exists go; then
        print_success "Go is installed ($(go version | awk '{print $3}'))"
    else
        print_warning "Go is not installed (only needed for local development)"
    fi

    # Check Make (for local mode)
    if command_exists make; then
        print_success "Make is installed"
    else
        print_warning "Make is not installed (only needed for local development)"
    fi

    # Check Node.js (for frontend development)
    if command_exists node; then
        print_success "Node.js is installed ($(node --version))"
    else
        print_warning "Node.js is not installed (only needed for frontend development)"
    fi

    if [ $missing -eq 1 ]; then
        echo ""
        print_error "Some required dependencies are missing. Please install them first."
        exit 1
    fi

    echo ""
}

# Check if container is running
check_container_running() {
    docker ps --format '{{.Names}}' | grep -q "^$1$" && return 0 || return 1
}

# Check if database is running
check_db_running() {
    check_container_running "$DB_CONTAINER"
}

# Wait for database to be ready
wait_for_db() {
    print_info "Waiting for database to be ready..."
    local max_attempts=30
    local attempt=0

    while [ $attempt -lt $max_attempts ]; do
        if docker exec "$DB_CONTAINER" pg_isready -U lta_user -d transport_reliability > /dev/null 2>&1; then
            print_success "Database is ready!"
            return 0
        fi
        attempt=$((attempt + 1))
        echo -n "."
        sleep 1
    done

    print_error "Database failed to start within 30 seconds"
    print_info "Check logs with: docker-compose logs db"
    return 1
}

# Wait for backend API to be ready
wait_for_api() {
    print_info "Waiting for backend API to be ready..."
    local max_attempts=60
    local attempt=0

    while [ $attempt -lt $max_attempts ]; do
        if curl -s http://localhost:$APP_HTTP_PORT/health > /dev/null 2>&1; then
            print_success "Backend API is ready!"
            return 0
        fi
        attempt=$((attempt + 1))
        echo -n "."
        sleep 1
    done

    print_error "Backend API failed to start within 60 seconds"
    print_info "Check logs with: docker-compose logs app"
    return 1
}

# Wait for frontend to be ready
wait_for_frontend() {
    print_info "Waiting for frontend to be ready..."
    local max_attempts=60
    local attempt=0

    while [ $attempt -lt $max_attempts ]; do
        if curl -s http://localhost:$FRONTEND_PORT > /dev/null 2>&1; then
            print_success "Frontend is ready!"
            return 0
        fi
        attempt=$((attempt + 1))
        echo -n "."
        sleep 1
    done

    print_warning "Frontend failed to start within 60 seconds (non-critical)"
    print_info "Check logs with: docker-compose logs frontend"
    return 1
}

# Start database
start_database() {
    print_step "Starting database..."

    if check_db_running; then
        print_success "Database is already running"
    else
        docker-compose up -d db
        wait_for_db || exit 1
    fi
}

# Build application (local)
build_app() {
    print_step "Building application..."

    # Check if Go is installed
    if ! command_exists go; then
        print_error "Go is not installed. Cannot build locally."
        print_info "Use 'docker' or 'docker-bg' mode instead, or install Go first."
        exit 1
    fi

    # Check if Make is installed
    if ! command_exists make; then
        print_error "Make is not installed. Cannot build."
        exit 1
    fi

    make build
    print_success "Application built successfully"
}

# Populate sample data
populate_data() {
    print_step "Populating database with sample data..."

    # Check if API is running
    if ! curl -s http://localhost:$APP_HTTP_PORT/health > /dev/null 2>&1; then
        print_error "API is not running. Cannot populate data."
        print_info "Start the application first with: ./run.sh start"
        return 1
    fi

    # Run populate script
    if [ -f "scripts/populate-data/main.go" ]; then
        cd scripts/populate-data && go run main.go && cd ../..
        print_success "Sample data populated successfully"
    else
        print_error "Populate script not found at scripts/populate-data/main.go"
        return 1
    fi
}

# Show status of all services
show_status() {
    print_banner
    print_step "Service Status:"
    echo ""

    local all_running=true

    # Check database
    if check_container_running "$DB_CONTAINER"; then
        print_success "Database:     Running"
    else
        print_error "Database:     Not running"
        all_running=false
    fi

    # Check app
    if check_container_running "$APP_CONTAINER"; then
        print_success "Backend App:  Running"
    else
        print_error "Backend App:  Not running"
        all_running=false
    fi

    # Check frontend
    if check_container_running "$FRONTEND_CONTAINER"; then
        print_success "Frontend:     Running"
    else
        print_warning "Frontend:     Not running"
    fi

    # Check nginx
    if check_container_running "$NGINX_CONTAINER"; then
        print_success "Nginx:        Running"
    else
        print_warning "Nginx:        Not running"
    fi

    echo ""

    if [ "$all_running" = true ]; then
        print_step "Access Points:"
        echo ""
        print_info "📱 Frontend Dashboard:  http://localhost:$FRONTEND_PORT"
        print_info "📚 Swagger UI:          http://localhost:$APP_HTTP_PORT/swagger/"
        print_info "📡 HTTP API:            http://localhost:$APP_HTTP_PORT"
        print_info "📡 gRPC API:            http://localhost:$APP_GRPC_PORT"
        print_info "🔀 Nginx Proxy:         http://localhost:$NGINX_PORT"
        print_info "🗄️  Database:           localhost:$DB_PORT"
        echo ""
    else
        print_warning "Some services are not running. Use './run.sh start' to start all services."
        echo ""
    fi
}

# Start everything (comprehensive)
start_all() {
    print_banner

    # Check prerequisites
    check_prerequisites

    print_step "Starting all services..."
    echo ""

    # Build and start with Docker Compose
    print_info "Building and starting Docker containers..."
    docker-compose up -d --build

    echo ""

    # Wait for database
    wait_for_db || exit 1

    echo ""

    # Wait for backend API
    wait_for_api || exit 1

    echo ""

    # Wait for frontend (optional)
    wait_for_frontend || true

    echo ""

    # Show success message
    print_success "=========================================="
    print_success "  All Services Started Successfully!"
    print_success "=========================================="
    echo ""

    print_step "Access Points:"
    echo ""
    print_info "📱 Frontend Dashboard:  http://localhost:$FRONTEND_PORT"
    print_info "📚 Swagger UI:          http://localhost:$APP_HTTP_PORT/swagger/"
    print_info "📡 HTTP API:            http://localhost:$APP_HTTP_PORT"
    print_info "📡 gRPC API:            http://localhost:$APP_GRPC_PORT"
    print_info "🔀 Nginx Proxy:         http://localhost:$NGINX_PORT"
    print_info "🗄️  Database:           localhost:$DB_PORT"
    echo ""

    print_step "Useful Commands:"
    echo ""
    print_info "View logs:        ./run.sh logs"
    print_info "View status:      ./run.sh status"
    print_info "Populate data:    ./run.sh populate"
    print_info "Run tests:        ./run.sh test"
    print_info "Stop all:         ./run.sh stop"
    echo ""

    # Ask if user wants to populate sample data
    read -p "$(echo -e ${CYAN}Would you like to populate the database with sample data? [Y/n]: ${NC})" -n 1 -r
    echo ""

    if [[ $REPLY =~ ^[Yy]$ ]] || [[ -z $REPLY ]]; then
        echo ""
        populate_data || print_warning "Failed to populate data. You can try again with: ./run.sh populate"
        echo ""
    fi

    print_success "Setup complete! Your application is ready to use."
    echo ""
}

# Run local mode
run_local() {
    print_banner

    # Check prerequisites for local mode
    if ! command_exists go || ! command_exists make; then
        print_error "Local mode requires Go and Make to be installed"
        print_info "Use './run.sh docker-bg' to run in Docker instead"
        exit 1
    fi

    # Start database
    start_database

    # Build application
    build_app

    # Show access URLs
    echo ""
    print_success "=========================================="
    print_success "  Application is starting..."
    print_success "=========================================="
    echo ""
    print_info "📡 API Endpoints:"
    print_info "   gRPC:    http://localhost:$APP_GRPC_PORT"
    print_info "   HTTP:    http://localhost:$APP_HTTP_PORT"
    echo ""
    print_info "📚 Swagger UI:"
    print_info "   http://localhost:$APP_HTTP_PORT/swagger/"
    echo ""
    print_info "🗄️  Database:"
    print_info "   Host:     localhost:$DB_PORT"
    print_info "   Database: transport_reliability"
    print_info "   User:     lta_user"
    echo ""
    print_warning "Press Ctrl+C to stop the application"
    echo ""

    # Run application
    set -a
    source local.env
    set +a
    ./bin/myapp
}

# Run docker mode
run_docker() {
    print_banner

    check_prerequisites

    print_info "Starting application in Docker mode (attached)..."
    docker-compose up --build
}

# Run docker detached mode
run_docker_detached() {
    print_banner

    check_prerequisites

    print_info "Starting application in Docker mode (detached)..."
    docker-compose up -d --build

    # Wait for services to be healthy
    echo ""
    wait_for_db || exit 1
    echo ""
    wait_for_api || exit 1
    echo ""

    # Show access URLs
    echo ""
    print_success "=========================================="
    print_success "  Application is running!"
    print_success "=========================================="
    echo ""
    print_info "📱 Frontend:    http://localhost:$FRONTEND_PORT"
    print_info "📚 Swagger UI:  http://localhost:$APP_HTTP_PORT/swagger/"
    print_info "📡 HTTP API:    http://localhost:$APP_HTTP_PORT"
    print_info "📡 gRPC API:    http://localhost:$APP_GRPC_PORT"
    print_info "🔀 Nginx:       http://localhost:$NGINX_PORT"
    print_info "🗄️  Database:   localhost:$DB_PORT"
    echo ""
    print_info "📋 Commands:"
    print_info "   View logs:    ./run.sh logs"
    print_info "   View status:  ./run.sh status"
    print_info "   Stop:         ./run.sh stop"
    echo ""
}

# Stop all services
stop_all() {
    print_step "Stopping all Docker containers..."
    docker-compose down
    print_success "All containers stopped"
}

# Restart all services
restart_all() {
    print_banner
    print_step "Restarting all services..."
    docker-compose restart
    echo ""
    wait_for_db || exit 1
    echo ""
    wait_for_api || exit 1
    echo ""
    print_success "All services restarted successfully"
}

# Show logs
show_logs() {
    if [ -z "$2" ]; then
        # Show all logs
        docker-compose logs -f
    else
        # Show specific service logs
        docker-compose logs -f "$2"
    fi
}

# Run tests
run_tests() {
    print_banner
    print_step "Running tests..."

    if ! command_exists go; then
        print_error "Go is not installed. Cannot run tests."
        exit 1
    fi

    echo ""
    make test
    echo ""
    print_success "Tests completed!"
}

# Run stress tests
run_stress_tests() {
    print_banner
    print_step "Running stress tests..."

    # Check if API is running
    if ! curl -s http://localhost:$APP_HTTP_PORT/health > /dev/null 2>&1; then
        print_error "API is not running. Please start the application first."
        print_info "Run: ./run.sh start"
        exit 1
    fi

    # Check if vegeta is installed
    if ! command_exists vegeta; then
        print_error "Vegeta is not installed."
        print_info "Install with: brew install vegeta (macOS) or go install github.com/tsenart/vegeta@latest"
        exit 1
    fi

    echo ""
    make stress-test
    echo ""
    print_success "Stress tests completed!"
}

# Clean everything
clean_all() {
    print_step "Cleaning up everything..."

    # Stop and remove containers, networks, volumes
    docker-compose down -v

    # Clean build artifacts
    if command_exists make; then
        make clean 2>/dev/null || true
    fi

    print_success "Cleanup complete"
}

# Install dependencies
install_deps() {
    print_banner
    print_step "Installing dependencies..."

    # Check if we're in the right directory
    if [ ! -f "go.mod" ]; then
        print_error "go.mod not found. Are you in the project root directory?"
        exit 1
    fi

    # Install Go dependencies
    if command_exists go; then
        print_info "Installing Go dependencies..."
        go mod download
        print_success "Go dependencies installed"
    fi

    # Install development tools
    if command_exists make; then
        print_info "Installing development tools..."
        make install 2>/dev/null || print_warning "Some dev tools failed to install"
    fi

    # Install frontend dependencies
    if [ -d "frontend" ] && [ -f "frontend/package.json" ]; then
        if command_exists npm; then
            print_info "Installing frontend dependencies..."
            cd frontend && npm install && cd ..
            print_success "Frontend dependencies installed"
        else
            print_warning "npm not found. Skipping frontend dependencies."
        fi
    fi

    echo ""
    print_success "Dependencies installed!"
}

# Show help
show_help() {
    cat << EOF
$(print_msg "$CYAN" "Transport Reliability Analytics System - Run Script")

$(print_msg "$YELLOW" "Usage:") ./run.sh [COMMAND] [OPTIONS]

$(print_msg "$YELLOW" "Main Commands:")
    start           🚀 Start everything (recommended for first time)
                       - Checks prerequisites
                       - Builds and starts all services
                       - Waits for health checks
                       - Optionally populates sample data

    stop            ⏹️  Stop all Docker containers

    restart         🔄 Restart all services

    status          📊 Show status of all services

$(print_msg "$YELLOW" "Run Modes:")
    local           💻 Run locally (database in Docker, app on host) [DEFAULT]
    docker          🐳 Run everything in Docker (attached mode, see logs)
    docker-bg       🐳 Run everything in Docker (detached/background mode)

$(print_msg "$YELLOW" "Development Commands:")
    build           🔨 Build the application only
    test            🧪 Run unit tests
    stress          ⚡ Run stress tests (requires running app)
    populate        📊 Populate database with sample data (requires running app)
    install         📦 Install all dependencies (Go, npm, dev tools)

$(print_msg "$YELLOW" "Utility Commands:")
    logs [service]  📝 Show Docker logs (all or specific service)
                       Examples: ./run.sh logs
                                ./run.sh logs app
                                ./run.sh logs db
    db              🗄️  Start only the database
    clean           🧹 Stop containers and clean up everything
    help            ❓ Show this help message

$(print_msg "$YELLOW" "Examples:")
    ./run.sh start              # First time setup - does everything
    ./run.sh                    # Run in local mode
    ./run.sh docker-bg          # Run in Docker background
    ./run.sh status             # Check what's running
    ./run.sh populate           # Add sample data
    ./run.sh logs app           # View backend logs
    ./run.sh test               # Run tests
    ./run.sh stop               # Stop everything

$(print_msg "$YELLOW" "Access Points (when running):")
    📱 Frontend:      http://localhost:$FRONTEND_PORT
    📚 Swagger UI:    http://localhost:$APP_HTTP_PORT/swagger/
    📡 HTTP API:      http://localhost:$APP_HTTP_PORT
    📡 gRPC API:      http://localhost:$APP_GRPC_PORT
    🔀 Nginx:         http://localhost:$NGINX_PORT
    🗄️  Database:     localhost:$DB_PORT

$(print_msg "$YELLOW" "Database Connection:")
    Host:     localhost
    Port:     $DB_PORT
    Database: transport_reliability
    User:     lta_user
    Password: lta_pass_2025

EOF
}

# Main script
main() {
    case "${1:-local}" in
        start|init|setup)
            start_all
            ;;
        local)
            run_local
            ;;
        docker)
            run_docker
            ;;
        docker-bg|detached|bg)
            run_docker_detached
            ;;
        stop)
            stop_all
            ;;
        restart)
            restart_all
            ;;
        status)
            show_status
            ;;
        logs)
            show_logs "$@"
            ;;
        build)
            print_banner
            build_app
            ;;
        test|tests)
            run_tests
            ;;
        stress|stress-test)
            run_stress_tests
            ;;
        populate|populate-data)
            print_banner
            populate_data
            ;;
        install|deps)
            install_deps
            ;;
        db|db-only)
            print_banner
            start_database
            print_success "Database is running on localhost:$DB_PORT"
            ;;
        clean)
            clean_all
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            print_error "Unknown command: $1"
            echo ""
            show_help
            exit 1
            ;;
    esac
}

# Handle Ctrl+C gracefully
trap 'echo ""; print_warning "Shutting down..."; exit 0' INT TERM

# Run main
main "$@"
