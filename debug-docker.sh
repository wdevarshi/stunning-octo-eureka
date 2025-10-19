#!/bin/bash

echo "=== Docker Debug Report ==="
echo ""

echo "1. Checking Docker daemon..."
if pgrep -f "Docker" > /dev/null; then
    echo "✓ Docker Desktop is running"
else
    echo "✗ Docker Desktop is NOT running"
    echo "  → Start Docker Desktop and try again"
    exit 1
fi

echo ""
echo "2. Checking if Docker is responding..."
timeout 5 docker version > /dev/null 2>&1
if [ $? -eq 0 ]; then
    echo "✓ Docker is responding"
else
    echo "✗ Docker is NOT responding or too slow"
    echo "  → Docker might be stuck. Restart Docker Desktop:"
    echo "     osascript -e 'quit app \"Docker\"'"
    echo "     sleep 5"
    echo "     open -a Docker"
    exit 1
fi

echo ""
echo "3. Checking for transport containers..."
CONTAINERS=$(timeout 5 docker ps -a --filter "name=transport" --format "{{.Names}}\t{{.Status}}" 2>&1)
if [ $? -eq 0 ]; then
    if [ -n "$CONTAINERS" ]; then
        echo "Found containers:"
        echo "$CONTAINERS"
    else
        echo "No transport containers found"
    fi
else
    echo "✗ Could not list containers (timeout)"
fi

echo ""
echo "4. Checking ports..."
if lsof -i :5433 > /dev/null 2>&1; then
    echo "Port 5433 (database): IN USE"
    lsof -i :5433 | grep LISTEN
else
    echo "Port 5433 (database): AVAILABLE"
fi

if lsof -i :9091 > /dev/null 2>&1; then
    echo "Port 9091 (HTTP): IN USE"
    lsof -i :9091 | grep LISTEN
else
    echo "Port 9091 (HTTP): AVAILABLE"
fi

echo ""
echo "=== Recommendations ==="
echo ""
echo "Option 1: Restart Docker Desktop"
echo "  osascript -e 'quit app \"Docker\"'"
echo "  sleep 10"
echo "  open -a Docker"
echo "  sleep 20  # Wait for Docker to start"
echo ""
echo "Option 2: Clean up and start fresh"
echo "  docker-compose down -v"
echo "  docker system prune -f"
echo "  docker-compose up -d"
echo ""
echo "Option 3: Check Docker Desktop resources"
echo "  Docker Desktop → Settings → Resources"
echo "  Ensure: Memory >= 4GB, CPUs >= 2"
