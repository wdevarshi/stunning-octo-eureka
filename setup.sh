#!/bin/bash

# Transport Reliability Analytics System - Setup Script
# This script installs all dependencies and prepares the environment for stress testing

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Print colored messages
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
    print_info "=========================================================="
    print_info "  Transport Reliability Analytics - Setup Script"
    print_info "=========================================================="
    echo ""
}

# Detect OS
detect_os() {
    if [[ "$OSTYPE" == "darwin"* ]]; then
        OS="macos"
        print_info "Detected OS: macOS"
    elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
        OS="linux"
        print_info "Detected OS: Linux"
    else
        print_error "Unsupported OS: $OSTYPE"
        exit 1
    fi
}

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check and install Homebrew (macOS)
check_brew() {
    if [[ "$OS" != "macos" ]]; then
        return
    fi

    print_step "Checking Homebrew..."
    if command_exists brew; then
        print_success "Homebrew is installed"
    else
        print_warning "Homebrew not found. Installing..."
        /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
        print_success "Homebrew installed"
    fi
}

# Check and install Docker
check_docker() {
    print_step "Checking Docker..."
    if command_exists docker; then
        print_success "Docker is installed ($(docker --version))"

        # Check if Docker daemon is running
        if docker ps >/dev/null 2>&1; then
            print_success "Docker daemon is running"
        else
            print_warning "Docker is installed but not running"
            print_info "Please start Docker Desktop and run this script again"
            if [[ "$OS" == "macos" ]]; then
                print_info "Run: open -a Docker"
            fi
            exit 1
        fi
    else
        print_warning "Docker not found. Installing..."
        if [[ "$OS" == "macos" ]]; then
            brew install --cask docker
            print_success "Docker installed"
            print_warning "Please start Docker Desktop manually and run this script again"
            print_info "Run: open -a Docker"
            exit 1
        else
            print_error "Please install Docker manually from https://docs.docker.com/engine/install/"
            exit 1
        fi
    fi
}

# Check and install Docker Compose
check_docker_compose() {
    print_step "Checking Docker Compose..."
    if command_exists docker-compose || docker compose version >/dev/null 2>&1; then
        print_success "Docker Compose is installed"
    else
        print_warning "Docker Compose not found. Installing..."
        if [[ "$OS" == "macos" ]]; then
            # Usually comes with Docker Desktop
            print_error "Docker Compose should come with Docker Desktop. Please reinstall Docker Desktop."
            exit 1
        else
            sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
            sudo chmod +x /usr/local/bin/docker-compose
            print_success "Docker Compose installed"
        fi
    fi
}

# Check and install Go
check_go() {
    print_step "Checking Go..."
    if command_exists go; then
        GO_VERSION=$(go version | awk '{print $3}')
        print_success "Go is installed ($GO_VERSION)"
    else
        print_warning "Go not found. Installing..."
        if [[ "$OS" == "macos" ]]; then
            brew install go
            print_success "Go installed"
        else
            print_error "Please install Go manually from https://golang.org/dl/"
            print_info "Recommended: Go 1.21 or later"
            exit 1
        fi
    fi
}

# Check and install Make
check_make() {
    print_step "Checking Make..."
    if command_exists make; then
        print_success "Make is installed"
    else
        print_warning "Make not found. Installing..."
        if [[ "$OS" == "macos" ]]; then
            xcode-select --install 2>/dev/null || print_warning "Xcode Command Line Tools already installed or installation in progress"
            print_success "Make will be available after Xcode Command Line Tools installation"
        else
            print_error "Please install Make manually: sudo apt-get install build-essential"
            exit 1
        fi
    fi
}

# Check and install Vegeta (load testing tool)
check_vegeta() {
    print_step "Checking Vegeta (stress testing tool)..."
    if command_exists vegeta; then
        VEGETA_VERSION=$(vegeta -version 2>&1 | head -n 1 || echo "unknown")
        print_success "Vegeta is installed ($VEGETA_VERSION)"
    else
        print_warning "Vegeta not found. Installing..."
        if [[ "$OS" == "macos" ]]; then
            brew install vegeta
            print_success "Vegeta installed"
        else
            print_info "Installing Vegeta via Go..."
            go install github.com/tsenart/vegeta@latest

            # Check if GOPATH/bin is in PATH
            if [[ ":$PATH:" != *":$HOME/go/bin:"* ]]; then
                print_warning "Adding \$HOME/go/bin to PATH"
                export PATH="$PATH:$HOME/go/bin"
                echo 'export PATH="$PATH:$HOME/go/bin"' >> ~/.bashrc
                print_info "Added \$HOME/go/bin to PATH in ~/.bashrc"
            fi
            print_success "Vegeta installed"
        fi
    fi
}

# Check and install jq (JSON processor, optional but useful)
check_jq() {
    print_step "Checking jq (JSON processor)..."
    if command_exists jq; then
        print_success "jq is installed"
    else
        print_warning "jq not found. Installing (optional but recommended)..."
        if [[ "$OS" == "macos" ]]; then
            brew install jq
            print_success "jq installed"
        else
            print_info "You can install jq with: sudo apt-get install jq"
        fi
    fi
}

# Check and install curl
check_curl() {
    print_step "Checking curl..."
    if command_exists curl; then
        print_success "curl is installed"
    else
        print_warning "curl not found. Installing..."
        if [[ "$OS" == "macos" ]]; then
            brew install curl
        else
            print_error "Please install curl: sudo apt-get install curl"
            exit 1
        fi
    fi
}

# Check and install Node.js
check_node() {
    print_step "Checking Node.js..."
    if command_exists node; then
        NODE_VERSION=$(node --version)
        print_success "Node.js is installed ($NODE_VERSION)"

        # Check npm
        if command_exists npm; then
            NPM_VERSION=$(npm --version)
            print_success "npm is installed ($NPM_VERSION)"
        else
            print_error "npm not found but Node.js is installed. Please reinstall Node.js"
            exit 1
        fi
    else
        print_warning "Node.js not found. Installing..."
        if [[ "$OS" == "macos" ]]; then
            brew install node@20
            print_success "Node.js 20 installed"
        else
            print_error "Please install Node.js 20 manually from https://nodejs.org/"
            print_info "Recommended: Node.js 20 LTS"
            exit 1
        fi
    fi
}

# Install Go dependencies
install_go_deps() {
    print_step "Installing Go dependencies..."

    if [[ ! -f "go.mod" ]]; then
        print_error "go.mod not found. Are you in the project root directory?"
        exit 1
    fi

    go mod download
    print_success "Go dependencies downloaded"

    # Install development tools
    print_info "Installing development tools (buf, mockery, golangci-lint)..."
    make install 2>/dev/null || {
        print_warning "Could not install all dev tools via make install. Continuing..."
    }
}

# Install frontend dependencies
install_frontend_deps() {
    print_step "Installing frontend dependencies..."

    if [[ ! -d "frontend" ]]; then
        print_warning "frontend directory not found. Skipping frontend setup."
        return 0
    fi

    cd frontend

    if [[ ! -f "package.json" ]]; then
        print_error "package.json not found in frontend directory"
        cd ..
        exit 1
    fi

    print_info "Installing npm packages (this may take a few minutes)..."
    npm install

    print_success "Frontend dependencies installed"
    cd ..
}

# Build the application
build_application() {
    print_step "Building application..."

    # Generate proto files if needed
    if command_exists buf; then
        print_info "Generating protocol buffer files..."
        make generate 2>/dev/null || print_warning "Could not generate proto files, continuing..."
    fi

    # Build the binary
    make build
    print_success "Application built successfully"
}

# Setup database
setup_database() {
    print_step "Setting up database..."

    # Check if database is already running
    if docker ps | grep -q "transport_db"; then
        print_warning "Database container is already running"
        print_info "Stopping existing database container..."
        docker-compose down
    fi

    # Start database
    print_info "Starting database container..."
    docker-compose up -d db

    # Wait for database to be ready
    print_info "Waiting for database to be ready..."
    local max_attempts=30
    local attempt=0

    while [ $attempt -lt $max_attempts ]; do
        if docker exec transport_db pg_isready -U lta_user -d transport_reliability > /dev/null 2>&1; then
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

# Verify stress test files
verify_stress_tests() {
    print_step "Verifying stress test files..."

    if [[ ! -d "scripts/stress-test" ]]; then
        print_error "scripts/stress-test directory not found"
        exit 1
    fi

    local required_files=(
        "scripts/stress-test/run-10qps.sh"
        "scripts/stress-test/run-100qps.sh"
        "scripts/stress-test/run-all-tests.sh"
        "scripts/stress-test/targets.txt"
    )

    for file in "${required_files[@]}"; do
        if [[ ! -f "$file" ]]; then
            print_error "Required file not found: $file"
            exit 1
        fi
    done

    # Make scripts executable
    chmod +x scripts/stress-test/*.sh

    print_success "All stress test files are present and executable"
}

# Run application and stress tests
run_stress_tests() {
    print_step "Starting application for stress testing..."

    # Start application in Docker background mode
    print_info "Starting application in Docker (detached mode)..."
    docker-compose up -d

    # Wait for backend API to be ready
    print_info "Waiting for backend API to be ready..."
    local max_attempts=60
    local attempt=0

    while [ $attempt -lt $max_attempts ]; do
        if curl -s http://localhost:9091/health > /dev/null 2>&1; then
            print_success "Backend API is ready!"
            break
        fi
        attempt=$((attempt + 1))
        echo -n "."
        sleep 1
    done

    if [ $attempt -eq $max_attempts ]; then
        print_error "Backend API failed to start within 60 seconds"
        print_info "Check logs with: docker-compose logs app"
        return 1
    fi

    # Wait for frontend to be ready
    print_info "Waiting for frontend to be ready..."
    attempt=0

    while [ $attempt -lt $max_attempts ]; do
        if curl -s http://localhost:3000 > /dev/null 2>&1; then
            print_success "Frontend is ready!"
            break
        fi
        attempt=$((attempt + 1))
        echo -n "."
        sleep 1
    done

    if [ $attempt -eq $max_attempts ]; then
        print_warning "Frontend failed to start within 60 seconds (non-critical)"
        print_info "Check logs with: docker-compose logs frontend"
    fi

    echo ""
    print_step "Running stress tests..."
    echo ""

    # Run stress tests
    cd scripts/stress-test
    ./run-all-tests.sh
    cd ../..

    echo ""
    print_success "=========================================================="
    print_success "  Stress tests completed!"
    print_success "=========================================================="
    echo ""
    print_info "Application is still running. Access points:"
    print_info "  - Frontend Dashboard:  http://localhost:3000"
    print_info "  - Backend HTTP API:    http://localhost:9091"
    print_info "  - Backend gRPC API:    http://localhost:9090"
    print_info "  - Database:            localhost:5433"
    echo ""
    print_info "To view logs:"
    print_info "  docker-compose logs -f app       # Backend logs"
    print_info "  docker-compose logs -f frontend  # Frontend logs"
    print_info "  docker-compose logs -f db        # Database logs"
    echo ""
    print_info "To stop all services:"
    print_info "  docker-compose down"
    echo ""
}

# Main setup flow
main() {
    print_banner

    print_info "This script will:"
    print_info "  1. Check and install required dependencies"
    print_info "  2. Build the application"
    print_info "  3. Setup the database"
    print_info "  4. Run stress tests"
    echo ""

    # Detect OS
    detect_os
    echo ""

    # Check and install dependencies
    check_brew
    check_docker
    check_docker_compose
    check_go
    check_make
    check_node
    check_vegeta
    check_jq
    check_curl
    echo ""

    # Install Go dependencies
    install_go_deps
    echo ""

    # Install frontend dependencies
    install_frontend_deps
    echo ""

    # Build application
    build_application
    echo ""

    # Verify stress test files
    verify_stress_tests
    echo ""

    # Ask user if they want to run stress tests now
    print_info "All dependencies are installed and the application is built."
    echo ""
    read -p "$(echo -e ${CYAN}Do you want to start the application and run stress tests now? [Y/n]: ${NC})" -n 1 -r
    echo ""

    if [[ $REPLY =~ ^[Yy]$ ]] || [[ -z $REPLY ]]; then
        echo ""
        run_stress_tests
    else
        echo ""
        print_success "=========================================================="
        print_success "  Setup completed successfully!"
        print_success "=========================================================="
        echo ""
        print_info "To start the application and run stress tests later:"
        print_info "  1. Start the application:"
        print_info "     docker-compose up -d"
        print_info "     OR"
        print_info "     ./run.sh docker-bg"
        echo ""
        print_info "  2. Run stress tests:"
        print_info "     make stress-test"
        print_info "     OR"
        print_info "     cd scripts/stress-test && ./run-all-tests.sh"
        echo ""
        print_info "For more options, run: ./run.sh help"
        echo ""
    fi
}

# Handle Ctrl+C gracefully
trap 'echo ""; print_warning "Setup interrupted"; exit 1' INT TERM

# Run main
main "$@"
