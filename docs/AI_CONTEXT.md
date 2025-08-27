# AI Context Guide for Hunter-Seeker

This document provides everything an AI agent needs to understand, develop, test, and troubleshoot the Hunter-Seeker job tracking application.

## Project Overview

Hunter-Seeker is a local web application for tracking job applications built with Go, SQLite, and a clean web interface. It allows users to track job applications without relying on external services.

### Key Technologies
- **Backend**: Go 1.23+
- **Database**: SQLite
- **Frontend**: HTML templates with vanilla JavaScript
- **Router**: Gorilla Mux
- **Deployment**: Docker & Docker Compose

### Project Structure
```
hunter-seeker/
├── cmd/
│   ├── server/main.go       # Main application entry point
│   ├── debug/main.go        # Debug utilities and test data
│   └── sample-data/main.go  # Sample data generation
├── internal/
│   ├── database/            # Database operations and models
│   ├── handlers/            # HTTP request handlers
│   └── models/              # Data structures
├── web/
│   ├── templates/           # HTML templates
│   └── static/             # CSS, JS, images
├── data/                   # SQLite database (auto-created)
├── bin/                    # Build output directory
├── scripts/                # Utility scripts
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
└── README.md
```

## Quick Start Commands

### Essential Commands (60 seconds to running)
```bash
# 1. Start the application
docker-compose up --build -d

# 2. Verify it's running
curl http://localhost:8080/health
# Expected: {"status":"ok","service":"hunter-seeker"}

# 3. Add test data
go run cmd/debug/main.go add-test-data

# 4. Access dashboard
open http://localhost:8080  # or visit in browser

# 5. Run full test suite
./scripts/test-application.sh
```

## Development Commands

### Building and Running
```bash
# Build application
mkdir -p bin && go build -o bin/hunter-seeker ./cmd/server

# Run directly (without Docker)
go run cmd/server/main.go

# Run with live reload (if air is installed)
air

# Install air for development
go install github.com/cosmtrek/air@latest
```

### Testing
```bash
# Run Go tests
go test ./...

# Run integration tests
./scripts/test-application.sh

# Check code formatting
go fmt ./...

# Vet code
go vet ./...

# Build test (check compilation)
go build ./cmd/server
```

## Database Management

### Database Details
- **Location**: `./data/jobs.db`
- **Type**: SQLite 3.x database
- **Schema**:
```sql
CREATE TABLE job_applications (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date_applied DATE NOT NULL,
    job_title TEXT NOT NULL,
    company TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'Applied',
    job_url TEXT,
    notes TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### Database Operations
```bash
# View database contents
go run cmd/debug/main.go

# Add sample data (12 applications)
go run cmd/sample-data/main.go

# Add simple test data (3 applications)
go run cmd/debug/main.go add-test-data

# Direct SQLite access
sqlite3 ./data/jobs.db "SELECT * FROM job_applications;"

# Backup database
cp ./data/jobs.db ./data/backup-$(date +%Y%m%d-%H%M%S).db

# Clear all data safely
./scripts/clear_data.sh

# Force clear without confirmation
./scripts/clear_data.sh --force

# Remove database file
rm -f ./data/jobs.db
```

## Docker Operations

### Basic Docker Commands
```bash
# Start application
docker-compose up --build -d

# Stop application
docker-compose down

# View logs
docker-compose logs -f hunter-seeker

# Check status
docker-compose ps

# Restart service
docker-compose restart hunter-seeker

# Force rebuild
docker-compose up --build --force-recreate -d

# Clean shutdown with volume removal
docker-compose down -v
```

### Docker Debugging
```bash
# Access container shell
docker exec -it hunter-seeker-app /bin/sh

# Check container files
docker exec -i hunter-seeker-app ls -la data/

# Copy database out
docker cp hunter-seeker-app:/root/data/jobs.db ./backup.db

# Copy database in
docker cp ./backup.db hunter-seeker-app:/root/data/jobs.db

# Monitor resources
docker stats hunter-seeker-app

# Check container health
docker inspect hunter-seeker-app | grep -A5 Health
```

## API Endpoints

### Web Interface
- `GET /` - Main dashboard with job listings
- `GET /add` - Add new job application form
- `POST /create` - Create job application (redirects to /)
- `GET /edit/{id}` - Edit job application form
- `POST /update/{id}` - Update job application
- `POST /delete/{id}` - Delete job application
- `GET /filter?status=Applied` - Filter by status
- `GET /import` - CSV import page
- `POST /import` - Process CSV import

### API Endpoints
- `GET /health` - Health check (returns JSON)
- `GET /api/stats` - Job statistics JSON

### Testing Endpoints
```bash
# Health check
curl http://localhost:8080/health

# Main dashboard
curl -s http://localhost:8080/ | grep "Hunter-Seeker"

# Create job application
curl -X POST http://localhost:8080/create \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "date_applied=2025-08-27&job_title=Software Engineer&company=TechCorp&status=Applied&job_url=https://example.com&notes=Test application"

# Check statistics
curl http://localhost:8080/api/stats

# Test delete with non-existent ID
curl -X POST http://localhost:8080/delete/999 -I

# Check error handling
curl -s "http://localhost:8080/?error=notfound&id=999" | grep "not found"
```

## Application Features

### Job Application Statuses
- **Applied**: Initial application submitted
- **In Review**: Application being reviewed
- **Phone Screen**: Phone/video screening
- **Interview**: In-person or video interview
- **Technical Test**: Coding challenge/assessment
- **Offer**: Job offer received
- **Rejected**: Application rejected
- **Withdrawn**: Application withdrawn
- **No Response**: No response from company

### CSV Import Format
```csv
Date Applied,Job Title,Company,Status,Job URL,Notes
2024-01-15,Senior Software Engineer,TechCorp,Applied,https://techcorp.com/jobs/123,Applied through website
2024-01-20,Full Stack Developer,StartupCo,In Review,https://startupco.com/careers,Remote position
```

Supported date formats:
- ISO: `2024-01-15` (recommended)
- US: `01/15/2024` or `1/15/2024`
- European: `15/01/2024` or `15/1/2024`
- Month names: `Jan 15 2024`, `January 15, 2024`

## Troubleshooting

### Common Issues and Solutions

#### Container Won't Start
```bash
# Check for port conflicts
lsof -i :8080

# Check Docker status
docker info

# Force rebuild
docker-compose down
docker-compose up --build --force-recreate
```

#### Database Issues
```bash
# Reset database (WARNING: deletes all data)
docker-compose down
rm -f data/jobs.db
docker-compose up -d

# Add fresh test data
go run cmd/debug/main.go add-test-data
```

#### Template/Handler Errors
```bash
# Check for Go compilation errors
go build ./cmd/server

# Look for template errors in logs
docker-compose logs hunter-seeker | grep -i template

# Test specific endpoints
curl -v http://localhost:8080/
curl -v http://localhost:8080/add
curl -v http://localhost:8080/health
```

#### Build Issues
```bash
# Clean build artifacts
rm -rf bin/
rm -f server sample-data debug hunter-seeker main

# Clean Go module cache
go clean -modcache

# Update dependencies
go mod download
go mod tidy
```

### Working Features ✅
- Health check endpoint
- Main dashboard loads and displays jobs
- Add new job applications
- Edit existing job applications
- Delete functionality with proper error messages
- Database persistence
- Statistics display
- CSV import functionality
- Docker deployment

### Known Issues ⚠️
- Filter pages may show "Internal server error" at bottom
- Some template rendering issues in FilterHandler
- No authentication/security (by design for local use)

## Success Criteria

Your application is working correctly if:

1. **Health check returns 200:**
   ```bash
   curl -f http://localhost:8080/health
   ```

2. **Dashboard loads without errors:**
   ```bash
   curl -s http://localhost:8080/ | grep -q "Job Applications"
   ```

3. **Can create and view jobs:**
   ```bash
   curl -X POST http://localhost:8080/create -d "date_applied=2025-08-27&job_title=Test&company=Test&status=Applied"
   curl -s http://localhost:8080/ | grep -q "Test"
   ```

4. **Delete shows proper error messages:**
   ```bash
   curl -s "http://localhost:8080/?error=notfound&id=999" | grep -q "not found"
   ```

5. **Database persists data:**
   ```bash
   docker-compose restart
   curl -s http://localhost:8080/ | grep -q "Test"
   ```

## Development Workflow

### Making Changes
1. **Edit Go files or templates**
2. **Rebuild and restart:**
   ```bash
   docker-compose up --build -d
   ```
3. **Check logs for errors:**
   ```bash
   docker-compose logs --tail=20 hunter-seeker
   ```
4. **Test changes:**
   ```bash
   curl http://localhost:8080/health
   ```

### Debugging Process
1. **Check container status:**
   ```bash
   docker-compose ps
   ```
2. **View real-time logs:**
   ```bash
   docker-compose logs -f hunter-seeker
   ```
3. **Test compilation:**
   ```bash
   go build ./cmd/server
   ```
4. **Run tests:**
   ```bash
   go test ./...
   ```

## File Locations

### Key Files to Modify
- `cmd/server/main.go` - Main application entry point
- `internal/handlers/` - HTTP request handlers
- `internal/database/` - Database operations
- `web/templates/` - HTML templates
- `web/static/` - CSS and JavaScript
- `docker-compose.yml` - Docker configuration

### Generated/Temporary Files (Do Not Commit)
- `data/jobs.db` - SQLite database file
- `bin/` - Compiled binaries
- Root directory binaries (server, debug, etc.)

### Configuration Files
- `go.mod` and `go.sum` - Go dependencies
- `.air.toml` - Live reload configuration
- `.gitignore` - Git ignore rules
- `Dockerfile` - Container build instructions

## Environment Variables

### Available Variables
- `PORT` - HTTP server port (default: 8080)
- `DB_PATH` - Database file path (default: ./data/jobs.db)

### Docker Environment
Set in `docker-compose.yml`:
```yaml
environment:
  - PORT=8080
  - DB_PATH=./data/jobs.db
```

### Local Development
```bash
PORT=9090 go run cmd/server/main.go
DB_PATH=/custom/path/jobs.db go run cmd/server/main.go
```

## Quick Reference Commands

### Daily Development
```bash
docker-compose up --build -d          # Start with rebuild
docker-compose logs -f hunter-seeker  # View logs
curl http://localhost:8080/health      # Health check
go run cmd/debug/main.go              # View database
docker-compose down                   # Stop application
```

### Testing and Debugging
```bash
go test ./...                         # Run tests
go build ./cmd/server                 # Check compilation
./scripts/test-application.sh         # Integration tests
docker exec -it hunter-seeker-app /bin/sh  # Container shell
go run cmd/debug/main.go add-test-data # Add test data
```

### Troubleshooting
```bash
docker-compose ps                     # Check status
lsof -i :8080                        # Check port conflicts
docker system prune                  # Clean Docker resources
rm -f data/jobs.db                   # Reset database
docker-compose up --force-recreate -d # Force container recreate
```

This guide provides everything needed to work effectively with the Hunter-Seeker application. For any specific issues not covered here, check the application logs and use the debugging commands provided.