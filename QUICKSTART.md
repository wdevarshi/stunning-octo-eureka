# Quick Start Guide

## One-Command Setup

The easiest way to get started:

```bash
./run.sh start
```

This single command will:
- Check all prerequisites (Docker, Go, Node.js)
- Build and start all services (database, backend, frontend, nginx)
- Run database migrations automatically
- Wait for health checks
- Optionally populate sample data
- Show you all access points

## Prerequisites

**Required:**
- Docker (with Docker Compose)
- Docker must be running

**Optional (for local development):**
- Go 1.21+
- Node.js 20+
- Make

## Common Commands

```bash
# Start everything (first time setup)
./run.sh start

# Check status of all services
./run.sh status

# View logs (all services)
./run.sh logs

# View logs (specific service)
./run.sh logs app      # Backend
./run.sh logs db       # Database
./run.sh logs frontend # Frontend

# Populate sample data
./run.sh populate

# Stop all services
./run.sh stop

# Restart all services
./run.sh restart

# Clean everything (including database)
./run.sh clean

# Run tests
./run.sh test

# Run stress tests
./run.sh stress

# Show all available commands
./run.sh help
```

## Different Run Modes

### 1. Full Docker Mode (Recommended)
Everything runs in Docker - easiest and most consistent:
```bash
./run.sh docker-bg    # Background mode
./run.sh docker       # Attached mode (see logs)
```

### 2. Local Development Mode
Database in Docker, app runs locally (requires Go and Make):
```bash
./run.sh local
```

### 3. Database Only
Start just the database for development:
```bash
./run.sh db
```

## Access Points

Once running, access your application at:

- **Frontend Dashboard**: http://localhost:3000
- **Swagger UI (API docs)**: http://localhost:9091/swagger/
- **Backend HTTP API**: http://localhost:9091
- **Backend gRPC API**: http://localhost:9090
- **Nginx Proxy**: http://localhost:8080
- **Database**: localhost:5433

## Database Connection

```
Host:     localhost
Port:     5433
Database: transport_reliability
User:     lta_user
Password: lta_pass_2025
```

## Database Migrations

Migrations run **automatically** when the database container starts!

The `database/init.sql` script creates:
- All tables (lines, stations, incidents)
- Indexes for performance
- Initial seed data (6 MRT lines and 150+ stations)

## Troubleshooting

### Docker not running
```bash
# macOS
open -a Docker

# Then wait for Docker to start and try again
./run.sh start
```

### Check what's running
```bash
./run.sh status
```

### View logs if something fails
```bash
./run.sh logs
```

### Clean slate (reset everything)
```bash
./run.sh clean
./run.sh start
```

### Port already in use
If you get a port conflict, stop other services using those ports or change the ports in `docker-compose.yml`:
- 5433 (database)
- 9091 (backend HTTP)
- 9090 (backend gRPC)
- 3000 (frontend)
- 8080 (nginx)

## Development Workflow

### First Time Setup
```bash
# Install dependencies
./run.sh install

# Start everything
./run.sh start

# Populate sample data
./run.sh populate
```

### Daily Development
```bash
# Start services
./run.sh docker-bg

# Check status
./run.sh status

# View logs while developing
./run.sh logs app

# Make code changes...

# Restart to see changes
./run.sh restart

# Run tests
./run.sh test

# Stop when done
./run.sh stop
```

### Running Stress Tests
```bash
# Start application
./run.sh start

# Run stress tests (requires vegeta)
./run.sh stress

# Or use make
make stress-test
```

## Sample Data

The populate script creates:
- **5 MRT Lines**: North-South, East-West, Circle, Downtown, Thomson-East Coast
- **80+ Stations** across all lines
- **15 Sample Incidents** with different types and timings

```bash
# Populate after starting the app
./run.sh populate

# Or use make
make populate-data
```

## Next Steps

1. Visit http://localhost:3000 to see the frontend dashboard
2. Visit http://localhost:9091/swagger/ to explore the API
3. Try creating incidents, viewing analytics, and exploring the data
4. Check out the API documentation in Swagger

## Getting Help

```bash
# Show all available commands
./run.sh help

# Check the README for detailed information
cat README.md
```

Enjoy building with the Transport Reliability Analytics System!
