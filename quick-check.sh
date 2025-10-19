#!/bin/bash

echo "Checking dependencies..."
echo ""

echo -n "Docker Compose: "
if docker compose version >/dev/null 2>&1 || command -v docker-compose >/dev/null 2>&1; then
    echo "✓ installed"
else
    echo "✗ missing"
fi

echo -n "Go: "
if command -v go >/dev/null 2>&1; then
    echo "✓ installed ($(go version))"
else
    echo "✗ missing"
fi

echo -n "Make: "
if command -v make >/dev/null 2>&1; then
    echo "✓ installed"
else
    echo "✗ missing"
fi

echo -n "Node.js: "
if command -v node >/dev/null 2>&1; then
    echo "✓ installed ($(node --version))"
else
    echo "✗ missing"
fi

echo -n "Vegeta: "
if command -v vegeta >/dev/null 2>&1; then
    echo "✓ installed"
else
    echo "✗ missing"
fi

echo -n "jq: "
if command -v jq >/dev/null 2>&1; then
    echo "✓ installed"
else
    echo "✗ missing (optional)"
fi
