# Getting Started Guide for AI Agents

## Quick Start (60 seconds)

1. **Start the application:**
   ```bash
   docker-compose up --build -d
   ```

2. **Verify it's running:**
   ```bash
   curl http://localhost:8080/health
   # Expected: {"status":"ok","service":"hunter-seeker"}
   ```

3. **Access the dashboard:**
   ```bash
   open http://localhost:8080
   # Or visit in browser
   ```

4. **Add test data:**
   ```bash
   go run cmd/debug/main.go add-test-data
   ```

5. **Run tests:**
   ```bash
   ./scripts/test-application.sh
   ```

## First-Time Setup

### Prerequisites Check
```bash
# Check Docker
docker --version
docker-compose --version

# Check Go (for development)
go version

# Check if port 8080 is free
lsof -i :8080
```

### Start from Scratch
```bash
# 1. Start the application
docker-compose up --build -d

# 2. Wait for it to be ready (should take 10-30 seconds)
docker-compose logs -f hunter-seeker

# 3. Test basic functionality
curl http://localhost:8080/health
curl -s http://localhost:8080/ | grep "Hunter-Seeker"

# 4. Add some test data
go run cmd/debug/main.go add-test-data

# 5. Verify data was added
go run cmd/debug/main.go
```

## Core Functionality Tests

### Test the Main Features
```bash
# 1. Create a job application
curl -X POST http://localhost:8080/create \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "date_applied=2025-08-27&job_title=Software Engineer&company=TechCorp&status=Applied&job_url=https://example.com&notes=Test application"

# 2. Check it appears on dashboard
curl -s http://localhost:8080/ | grep "TechCorp"

# 3. Test delete with non-existent ID (should redirect with error)
curl -X POST http://localhost:8080/delete/999 -I

# 4. Test delete success message
curl -s "http://localhost:8080/?error=notfound&id=999" | grep "not found"

# 5. Test filter (known to have issues but should partially work)
curl -s http://localhost:8080/filter?status=Applied | grep "Filter by status"
```

### Verify Database Operations
```bash
# Check database exists and has data
docker exec -i hunter-seeker-app ls -la data/jobs.db

# View all jobs in database
go run cmd/debug/main.go

# Check database size
docker exec -i hunter-seeker-app du -h data/jobs.db
```

## Development Workflow

### Making Changes
```bash
# 1. Edit Go files or templates
# 2. Rebuild and restart
docker-compose up --build -d

# 3. Check logs for errors
docker-compose logs --tail=20 hunter-seeker

# 4. Test changes
curl http://localhost:8080/health
```

### Debugging Issues
```bash
# View real-time logs
docker-compose logs -f hunter-seeker

# Check container status
docker-compose ps

# Access container shell
docker exec -it hunter-seeker-app /bin/sh

# Run tests
go test ./...
./scripts/test-application.sh
```

## Common Issues and Solutions

### Container Won't Start
```bash
# Check for port conflicts
lsof -i :8080

# Check Docker status
docker info

# Force rebuild
docker-compose down
docker-compose up --build --force-recreate
```

### Database Issues
```bash
# Reset database (WARNING: deletes all data)
docker-compose down
rm -f data/jobs.db
docker-compose up -d

# Add fresh test data
go run cmd/debug/main.go add-test-data
```

### Template/Handler Errors
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

## Key Endpoints to Test

### Web Interface
- `GET /` - Main dashboard
- `GET /add` - Add job form
- `POST /create` - Create job (redirects to /)
- `GET /edit/{id}` - Edit job form
- `POST /update/{id}` - Update job
- `POST /delete/{id}` - Delete job (redirects with status)
- `GET /filter?status=Applied` - Filter jobs

### API Endpoints
- `GET /health` - Health check
- `GET /api/stats` - Job statistics JSON

## Expected Behavior

### Working Features ✅
- Health check endpoint
- Main dashboard loads and displays jobs
- Add new job applications
- Edit existing job applications  
- Delete functionality with proper error messages
- Database persistence
- Statistics display

### Known Issues ⚠️
- Filter pages show "Internal server error" at bottom
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
   # Create
   curl -X POST http://localhost:8080/create -d "date_applied=2025-08-27&job_title=Test&company=Test&status=Applied"
   
   # View
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

## Useful Commands Reference

```bash
# Essential commands
docker-compose up --build -d          # Start application
docker-compose down                    # Stop application
docker-compose logs -f hunter-seeker   # View logs
curl http://localhost:8080/health      # Check health
go run cmd/debug/main.go              # View database
./scripts/test-application.sh         # Run tests

# Development
go test ./...                         # Run Go tests
go build ./cmd/server                 # Check compilation
docker-compose restart hunter-seeker  # Quick restart
docker exec -it hunter-seeker-app /bin/sh  # Container shell

# Troubleshooting
docker-compose ps                     # Check status
docker stats hunter-seeker-app       # Resource usage
lsof -i :8080                        # Check port conflicts
docker system prune                  # Clean up Docker
```

This guide should get any AI agent up and running with the Hunter-Seeker application in under 5 minutes.