# Transport Reliability Analytics

A full-stack analytics system for tracking and analyzing transport incidents in Singapore's MRT and bus network. Built with Go, gRPC, PostgreSQL, Next.js, and Docker.

## Overview

Production-ready system providing real-time incident tracking and analytics for the Land Transport Authority (LTA). Features a modern web dashboard with automated data refresh, comprehensive REST/gRPC APIs, and high-performance analytics.

### Key Features

- **Real-time Dashboard** - Next.js frontend with interactive charts and metrics (auto-refreshes every 30s)
- **Dual API** - gRPC-first design with automatic HTTP/REST gateway
- **Advanced Analytics** - MTBF calculations, top breakdowns, recent disruptions
- **High Performance** - Optimized PostgreSQL queries with covering indexes
- **Production Ready** - Dockerized deployment with health checks and monitoring
- **Stress Tested** - Verified performance under 10-100 QPS loads

## Quick Start

### Automated Setup (Recommended)

The setup script installs all dependencies and builds the application:

```bash
./setup.sh
```

This will:
1. Detect your OS (macOS/Linux)
2. Install Docker, Go, Node.js, and other required tools
3. Build the application
4. Setup the database with seed data
5. Optionally run stress tests

### Manual Setup

If you have Docker already installed:

```bash
docker-compose up
```

This starts:
- PostgreSQL database with seed data
- Go backend API (gRPC + HTTP)
- Next.js frontend dashboard
- nginx reverse proxy

### Access Points

| Service | URL | Description |
|---------|-----|-------------|
| **Dashboard** | http://localhost:3000 | Interactive analytics dashboard |
| **API (nginx)** | http://localhost:8080 | REST/HTTP API endpoint |
| **Backend Direct** | http://localhost:9091 | Direct backend access |
| **gRPC API** | localhost:9090 | gRPC endpoint |
| **API Docs** | http://localhost:9091/swagger/ | Swagger UI documentation |
| **Database** | localhost:5433 | PostgreSQL (user: `lta_user`, db: `transport_reliability`) |

### Health Checks

```bash
# Application health
curl http://localhost:8080/health

# Database connectivity
curl http://localhost:8080/ready

# Frontend accessibility
curl http://localhost:3000
```

## Architecture

### System Design

```
┌─────────────┐      ┌────────┐      ┌─────────────┐      ┌──────────────┐
│   Browser   │─────▶│ nginx  │─────▶│  Go Backend │─────▶│  PostgreSQL  │
│             │      │ (8080) │      │  (9091)     │      │  (5433)      │
└─────────────┘      └────────┘      └─────────────┘      └──────────────┘
       │                                     │
       │                                     │
       ▼                                     ▼
┌─────────────┐                      ┌─────────────┐
│   Next.js   │                      │    gRPC     │
│   (3000)    │                      │    (9090)   │
└─────────────┘                      └─────────────┘
```

### Technology Stack

**Backend:**
- Go 1.24
- go-coldbrew (gRPC framework)
- gRPC-Gateway (automatic REST mapping)
- Protocol Buffers
- PostgreSQL 15

**Frontend:**
- Next.js 15
- React 18
- TypeScript
- Tailwind CSS
- Recharts

**Infrastructure:**
- Docker & Docker Compose
- nginx (reverse proxy & CORS)
- PostgreSQL with optimized indexes

## API Documentation

All API endpoints are accessible through nginx at `http://localhost:8080` or directly at `http://localhost:9091`.

### 1. Create Incident

Submit a new transport incident.

**Endpoint:** `POST /incidents`

```bash
curl -X POST http://localhost:8080/incidents \
  -H "Content-Type: application/json" \
  -d '{
    "line": "North South Line",
    "station": "Orchard",
    "timestamp": "2025-10-16T08:32:00Z",
    "duration_minutes": 45,
    "incident_type": "signal"
  }'
```

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "line": "North South Line",
  "station": "Orchard",
  "timestamp": "2025-10-16T08:32:00Z",
  "duration_minutes": 45,
  "incident_type": "signal",
  "line_id": "123e4567-e89b-12d3-a456-426614174000",
  "station_id": "789e0123-e89b-12d3-a456-426614174000",
  "status": "open"
}
```

**Validation:**
- `line`: 1-100 characters
- `station`: 1-100 characters
- `timestamp`: ISO 8601, cannot be in future
- `duration_minutes`: 0-1440 (24 hours max)
- `incident_type`: `mechanical`, `power`, `signal`, `weather`, `other`

### 2. Top Breakdowns

Get top N lines or stations by incident count.

**Endpoint:** `GET /analytics/top_breakdowns`

**Parameters:**
- `scope`: `line` or `station` (required)
- `limit`: number of results (default: 5, max: 100)

```bash
# Top 5 stations
curl "http://localhost:8080/analytics/top_breakdowns?scope=station&limit=5"

# Top 10 lines
curl "http://localhost:8080/analytics/top_breakdowns?scope=line&limit=10"
```

**Response:**
```json
{
  "scope": "station",
  "items": [
    {"name": "Jurong East", "count": 156},
    {"name": "City Hall", "count": 148},
    {"name": "Tampines", "count": 45}
  ]
}
```

### 3. Mean Time Between Failures (MTBF)

Calculate MTBF for all lines.

**Endpoint:** `GET /analytics/mean_time_between_failures`

```bash
curl http://localhost:8080/analytics/mean_time_between_failures
```

**Response:**
```json
{
  "lines": [
    {"name": "North South Line", "mtbf_minutes": 741.2},
    {"name": "East West Line", "mtbf_minutes": 940.5},
    {"name": "Circle Line", "mtbf_minutes": 1135.2}
  ]
}
```

### 4. Recent Disruptions

Get recent incidents with optional filtering.

**Endpoint:** `GET /analytics/recent_disruptions`

**Parameters:**
- `line`: filter by line name (optional)
- `station`: filter by station name (optional)
- `limit`: number of results (default: 20, max: 100)

```bash
# Last 20 incidents
curl "http://localhost:8080/analytics/recent_disruptions?limit=20"

# Recent incidents on North South Line
curl "http://localhost:8080/analytics/recent_disruptions?line=North%20South%20Line"
```

**Response:**
```json
{
  "items": [
    {
      "line": "East West Line",
      "station": "Tampines",
      "timestamp": "2025-10-15T14:23:00Z",
      "duration_minutes": 32,
      "incident_type": "mechanical",
      "status": "investigating"
    }
  ]
}
```

## Database Schema

### Tables

**lines**
```sql
CREATE TABLE lines (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT UNIQUE NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

**stations**
```sql
CREATE TABLE stations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    line_id UUID REFERENCES lines(id),
    status TEXT DEFAULT 'active',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(name, line_id)
);
```

**incidents**
```sql
CREATE TABLE incidents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    station_id UUID REFERENCES stations(id),
    line_id UUID REFERENCES lines(id),
    ts TIMESTAMPTZ NOT NULL,
    duration_minutes INT CHECK (duration_minutes BETWEEN 0 AND 1440),
    incident_type TEXT NOT NULL,
    status TEXT DEFAULT 'open',
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

### Indexes

Optimized for analytics queries:

```sql
CREATE INDEX idx_incidents_ts ON incidents(ts DESC);
CREATE INDEX idx_incidents_line_id ON incidents(line_id);
CREATE INDEX idx_incidents_station_id ON incidents(station_id);
CREATE INDEX idx_incidents_line_ts ON incidents(line_id, ts DESC);
CREATE INDEX idx_incidents_station_ts ON incidents(station_id, ts DESC);
CREATE INDEX idx_incidents_status ON incidents(status);
CREATE INDEX idx_stations_status ON stations(status);
```

### Seed Data

Pre-populated with realistic data:
- 6 MRT lines (NSL, EWL, CCL, DTL, TEL, NEL)
- 150+ stations across all lines
- 400+ incidents over the last 90 days
- Clustered incidents at major interchanges (Jurong East, City Hall)

## Stress Testing

Comprehensive performance testing included.

### Quick Run

```bash
# Run all stress tests
make stress-test

# Individual tests
make stress-test-10qps   # Light load: 10 req/sec
make stress-test-100qps  # Heavy load: 100 req/sec
```

### Manual Execution

```bash
cd scripts/stress-test

# Run tests
./run-10qps.sh      # 600 requests over 60s
./run-100qps.sh     # 6,000 requests over 60s
./run-all-tests.sh  # All tests
```

### Performance Baselines

**10 QPS (Light Load):**
- Success Rate: ~100%
- p50 Latency: <50ms
- p95 Latency: <100ms

**100 QPS (Heavy Load):**
- Success Rate: >95%
- p50 Latency: <100ms
- p95 Latency: <500ms

Results saved in `scripts/stress-test/results/` with:
- Binary results (`.bin`)
- JSON reports (`.json`)
- Latency histograms
- Success rates and throughput metrics

## Local Development

### Prerequisites

- Go 1.24+
- Node.js 20+
- Docker & Docker Compose
- Make (optional)

### Backend Development

```bash
# Start database only
docker-compose up db -d

# Set environment
export DATABASE_URL="postgres://lta_user:lta_pass_2025@localhost:5432/transport_reliability?sslmode=disable"

# Install dependencies
go mod download

# Generate proto files
make generate
# or: buf generate

# Run backend
make run
# or: go run main.go
```

### Frontend Development

```bash
cd frontend

# Install dependencies
npm install

# Set API URL
export NEXT_PUBLIC_API_URL="http://localhost:8080"

# Run development server
npm run dev

# Build for production
npm run build
```

### Running Tests

```bash
# Backend tests
go test ./...

# Backend tests with coverage
go test -v -count=1 ./...

# Frontend tests (if added)
cd frontend
npm test
```

## Project Structure

```
.
├── backend/               # Go backend services
│   ├── models.go         # Data models
│   ├── repository.go     # Database layer
│   ├── service.go        # Business logic
│   └── service_test.go   # Unit tests
├── config/               # Configuration management
│   └── config.go
├── database/             # Database setup
│   ├── Dockerfile
│   └── init/
│       ├── 01-schema.sql      # Schema definitions
│       └── 02-seed-data.sql   # Initial data
├── frontend/             # Next.js dashboard
│   ├── src/
│   │   ├── app/         # Next.js app router
│   │   ├── components/  # React components
│   │   ├── lib/         # API client
│   │   └── types/       # TypeScript types
│   ├── Dockerfile
│   └── package.json
├── nginx/                # Reverse proxy config
│   └── nginx.conf
├── proto/                # Protocol buffer definitions
│   ├── transport.proto
│   └── *.pb.go          # Generated Go code
├── scripts/
│   └── stress-test/     # Performance testing
│       ├── run-10qps.sh
│       └── run-100qps.sh
├── third_party/          # Third-party integrations
│   └── OpenAPI/         # Swagger UI
├── docker-compose.yml    # Multi-container setup
├── Dockerfile           # Backend container
├── Makefile             # Build automation
├── setup.sh             # Automated setup
├── run.sh               # Run script
└── README.md            # This file
```

## Configuration

Environment variables:

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `DATABASE_URL` | PostgreSQL connection string | - | Yes |
| `ENVIRONMENT` | Environment name | `dev` | No |
| `LOG_LEVEL` | Log level (DEBUG, INFO, WARN, ERROR) | `INFO` | No |
| `HTTP_PORT` | HTTP server port | `9091` | No |
| `GRPC_PORT` | gRPC server port | `9090` | No |
| `NEXT_PUBLIC_API_URL` | Frontend API URL (build-time) | `http://localhost:8080` | No |
| `API_URL` | Server-side API URL | `http://nginx:8080` | No |

## Production Considerations

### Scaling

- **Horizontal Scaling:** Run multiple app instances behind load balancer
- **Database:** Use read replicas for analytics queries
- **Connection Pooling:** Configured for 25 max connections, 5 idle
- **Frontend:** Deploy to CDN with serverless functions

### Monitoring

- Health checks: `/health` and `/ready`
- Structured logging with request IDs
- Prometheus metrics endpoint
- OpenTelemetry tracing support

### Security

- Use strong passwords (change defaults!)
- Enable SSL/TLS for all connections
- Implement API authentication
- Add rate limiting middleware
- Input validation at protocol level
- CORS configured via nginx

### Performance

- All analytics use covering indexes
- MTBF optimized with window functions
- Connection pooling prevents exhaustion
- Prepared statements prevent SQL injection
- Static asset caching via nginx

## Troubleshooting

### Docker Issues

```bash
# Check container status
docker-compose ps

# View logs
docker-compose logs app
docker-compose logs frontend
docker-compose logs nginx
docker-compose logs db

# Restart services
docker-compose restart

# Clean rebuild
docker-compose down
docker-compose up --build
```

### Database Connection

```bash
# Check database
docker-compose ps db

# Connect to database
docker-compose exec db psql -U lta_user -d transport_reliability

# View database logs
docker-compose logs db
```

### Port Conflicts

If default ports are in use, edit `docker-compose.yml`:

```yaml
services:
  nginx:
    ports:
      - "8081:8080"  # Change nginx port
  frontend:
    ports:
      - "3001:3000"  # Change frontend port
```

### Frontend API Errors

If you see "Network error: Unable to connect to the API":

1. Verify nginx is running: `docker-compose ps nginx`
2. Check API is accessible: `curl http://localhost:8080/health`
3. Hard refresh browser: Cmd+Shift+R (Mac) or Ctrl+Shift+F5 (Windows)
4. Clear browser cache for localhost:3000

## Future Enhancements

### Features
- Real-time alerts when MTBF drops below threshold
- Predictive maintenance using ML
- Historical trend analysis (YoY comparisons)
- Passenger impact metrics
- CSV/PDF export functionality
- Multi-language support

### Technical
- Redis caching for frequent queries
- Kafka for high-volume incident ingestion
- GraphQL API alternative
- Multi-tenancy support
- Kubernetes deployment
- CI/CD pipeline

## License

Built for LTA assessment. All rights reserved.

---

**Tech Stack:** Go · gRPC · PostgreSQL · Next.js · TypeScript · Docker · nginx

**Last Updated:** 2025-10-19
