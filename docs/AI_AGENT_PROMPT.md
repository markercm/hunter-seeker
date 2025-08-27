# AI Agent Prompt for Hunter-Seeker Project

## Project Overview

Hunter-Seeker is a local web application for tracking job applications built with Go, SQLite, and a clean web interface. It allows users to track job applications without relying on external services.

## Key Technologies
- **Backend**: Go 1.21+
- **Database**: SQLite
- **Frontend**: HTML templates with vanilla JavaScript
- **Deployment**: Docker & Docker Compose
- **Router**: Gorilla Mux

## Project Structure

```
hunter-seeker/
├── cmd/
│   ├── server/          # Main application entry point
│   └── debug/           # Debug utilities
├── internal/
│   ├── database/        # Database operations and models
│   ├── handlers/        # HTTP request handlers
│   └── models/          # Data structures
├── web/
│   ├── templates/       # HTML templates
│   └── static/          # CSS, JS, images (if any)
├── data/               # SQLite database (auto-created)
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
└── README.md
```

## How to Run the Application

### Using Docker (Recommended)

1. **Start the application**:
   ```bash
   docker-compose up --build -d
   ```

2. **Check if running**:
   ```bash
   docker-compose ps
   ```

3. **View logs**:
   ```bash
   docker-compose logs -f hunter-seeker
   ```

4. **Stop the application**:
   ```bash
   docker-compose down
   ```

5. **Access the application**:
   - Main dashboard: http://localhost:8080
   - Health check: http://localhost:8080/health

### Using Go Directly (Development)

1. **Install dependencies**:
   ```bash
   go mod download
   ```

2. **Run the server**:
   ```bash
   go run cmd/server/main.go
   ```

3. **Build binary**:
   ```bash
   go build -o hunter-seeker ./cmd/server
   ```

## Testing and Debugging

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests for specific package
go test ./internal/database
go test ./internal/handlers

# Run tests with coverage
go test -cover ./...
```

### Debug Tools

1. **Debug endpoint for database content**:
   ```bash
   go run cmd/debug/main.go
   ```

2. **Add test data**:
   ```bash
   go run cmd/debug/main.go add-test-data
   ```

3. **Check database with specific path**:
   ```bash
   go run cmd/debug/main.go ./data/jobs.db
   ```

### Health Checks

1. **Application health**:
   ```bash
   curl http://localhost:8080/health
   ```
   Expected response: `{"status":"ok","service":"hunter-seeker"}`

2. **Database connectivity**:
   ```bash
   # Check if database file exists
   ls -la data/jobs.db
   
   # Or with Docker
   docker exec -i hunter-seeker-app ls -la data/
   ```

## API Endpoints

### Web Routes
- `GET /` - Main dashboard (HomeHandler)
- `GET /add` - Add job form (AddJobHandler)
- `POST /create` - Create new job (CreateJobHandler)
- `GET /edit/{id}` - Edit job form (EditJobHandler)
- `POST /update/{id}` - Update job (UpdateJobHandler)
- `POST /delete/{id}` - Delete job (DeleteJobHandler)
- `GET /filter?status={status}` - Filter jobs by status (FilterHandler)

### API Routes
- `GET /api/stats` - Get job statistics as JSON
- `GET /health` - Health check endpoint

## Database Schema

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

## Common Tasks

### Adding Test Data

Use the debug tool or create via web interface:

```bash
# Via debug tool
go run cmd/debug/main.go add-test-data

# Via curl (create new job)
curl -X POST http://localhost:8080/create \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "date_applied=2025-08-27&job_title=Software Engineer&company=TechCorp&status=Applied&job_url=https://example.com&notes=Applied via website"
```

### Testing Delete Functionality

```bash
# Delete existing job (should redirect with success message)
curl -X POST http://localhost:8080/delete/1 -I

# Delete non-existent job (should redirect with error message)
curl -X POST http://localhost:8080/delete/999 -I

# Check the response page
curl -s http://localhost:8080/?success=deleted
curl -s http://localhost:8080/?error=notfound&id=999
```

### Testing Filter Functionality

```bash
# Filter by status
curl -s http://localhost:8080/filter?status=Applied
curl -s http://localhost:8080/filter?status=Interview

# Check if filters work in browser
open http://localhost:8080/filter?status=Applied
```

## Environment Variables

- `PORT` - Server port (default: 8080)
- `DB_PATH` - Database file path (default: ./data/jobs.db)

## Data Persistence

- Database file: `./data/jobs.db`
- Docker volume: `./data` directory mounted to `/root/data` in container
- Data persists between container restarts

## Troubleshooting

### Container Issues

```bash
# Check container status
docker-compose ps

# View detailed logs
docker-compose logs hunter-seeker

# Restart container
docker-compose restart hunter-seeker

# Rebuild from scratch
docker-compose down
docker-compose up --build
```

### Database Issues

```bash
# Check if database file exists and has data
go run cmd/debug/main.go

# Check database permissions
ls -la data/

# Reset database (WARNING: deletes all data)
rm data/jobs.db
docker-compose restart
```

### Template/Handler Issues

```bash
# Check for compilation errors
go build ./cmd/server

# Run with verbose logging
go run cmd/server/main.go

# Check specific handler functionality
curl -v http://localhost:8080/[endpoint]
```

## Development Workflow

1. **Make changes** to Go files or templates
2. **Test locally**:
   ```bash
   go run cmd/server/main.go
   ```
3. **Run tests**:
   ```bash
   go test ./...
   ```
4. **Test with Docker**:
   ```bash
   docker-compose up --build -d
   ```
5. **Verify functionality** via browser or curl
6. **Check logs** for errors:
   ```bash
   docker-compose logs hunter-seeker
   ```

## Known Issues

### Filter Functionality
- **Issue**: Status filter pages show "Internal server error" at the bottom
- **Workaround**: Main dashboard works fine, filtering can be tested via direct URLs
- **Debug**: Check `docker-compose logs` for template execution errors

### Template Functions
- Templates use custom functions: `replace`, `lower`, `formatDate`
- These are defined in `internal/handlers/handlers.go`
- Template execution errors often indicate missing or malformed data

## Security Notes

- **Local use only** - No authentication required
- **No HTTPS** - HTTP only for local development
- **SQLite** - No access controls, file-based database
- **DO NOT** expose to internet without proper security measures

## Useful Commands Reference

```bash
# Quick setup
docker-compose up -d

# Add test data
go run cmd/debug/main.go add-test-data

# Check health
curl http://localhost:8080/health

# View all jobs
curl -s http://localhost:8080/ | grep -A5 "Job Applications"

# Test delete
curl -X POST http://localhost:8080/delete/1 -I

# Clean restart
docker-compose down && docker-compose up --build -d

# Follow logs
docker-compose logs -f hunter-seeker
```

## Template Structure

Templates are standalone HTML files (not using a base template system):
- `index.html` - Main dashboard with job list and filters
- `add_job.html` - Form to add new job application
- `edit_job.html` - Form to edit existing job application

Template data structures include:
- `Jobs` - Array of job applications
- `StatusCounts` - Map of status to counts
- `Statuses` - Array of available statuses
- `CurrentFilter` - Current filter status
- `StatusMessage` - Success/error messages
- `StatusType` - "success" or "error"

This prompt should help AI agents understand how to work with, test, and debug the Hunter-Seeker application effectively.