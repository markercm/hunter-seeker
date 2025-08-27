# AI Agent Documentation

This directory contains documentation and tools specifically designed for AI agents working with the Hunter-Seeker job tracking application.

## Files Overview

### 📋 Main Documentation
- **`docs/AI_AGENT_PROMPT.md`** - Comprehensive guide for AI agents with project overview, commands, and troubleshooting
- **`docs/AI_GETTING_STARTED.md`** - Quick start guide to get up and running in 60 seconds
- **`docs/DOCKER_REFERENCE.md`** - Complete Docker commands reference for the project

### 🧪 Testing Tools
- **`scripts/test-application.sh`** - Automated test suite to verify application functionality
- **`cmd/debug/main.go`** - Debug tool to inspect database contents and add test data
- **`cmd/sample-data/main.go`** - Sample data generation tool for testing

## Quick Start for AI Agents

```bash
# 1. Start the application
docker-compose up --build -d

# 2. Verify it works
curl http://localhost:8080/health

# 3. Add test data
go run cmd/debug/main.go add-test-data

# 4. Run full test suite
./scripts/test-application.sh

# 5. Access the application
open http://localhost:8080
```

## What This Application Does

Hunter-Seeker is a local job application tracking system that allows users to:
- Track job applications with status, company, dates, and notes
- Filter applications by status (Applied, Interview, Offer, etc.)
- View statistics and manage job search progress
- All data stored locally in SQLite database

## Key Technical Details

- **Language**: Go 1.21+
- **Database**: SQLite (file-based, no server required)
- **Web Framework**: Gorilla Mux router with HTML templates
- **Deployment**: Docker Compose for easy setup
- **Architecture**: Simple MVC pattern with handlers, database layer, and templates

## Working Features ✅

- ✅ Health check endpoint (`/health`)
- ✅ Main dashboard with job listings
- ✅ Add/edit/delete job applications
- ✅ Delete functionality with proper error messages:
  - Shows "Job application with ID X not found" for non-existent jobs
  - Shows "Job application deleted successfully" for successful deletions
- ✅ Database persistence between restarts
- ✅ Statistics display
- ✅ Responsive web interface

## Known Issues ⚠️

- ⚠️ Filter pages (`/filter?status=Applied`) show "Internal server error" at the bottom
- ⚠️ Template rendering issues in FilterHandler (under investigation)
- ⚠️ No authentication (by design - local use only)

## Common Tasks

### Add Test Data
```bash
# Using debug tool
go run cmd/debug/main.go add-test-data

# Using web interface
curl -X POST http://localhost:8080/create \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "date_applied=2025-08-27&job_title=Software Engineer&company=TechCorp&status=Applied&job_url=https://example.com&notes=Test job"
```

### Test Delete Functionality
```bash
# Test deleting non-existent job (should show error message)
curl -X POST http://localhost:8080/delete/999 -I
curl -s "http://localhost:8080/?error=notfound&id=999" | grep "not found"

# Test deleting existing job (should show success message)
curl -X POST http://localhost:8080/delete/1 -I
curl -s "http://localhost:8080/?success=deleted" | grep "deleted successfully"
```

### View Database Contents
```bash
# See all job applications
go run cmd/debug/main.go

# Check database file
docker exec -i hunter-seeker-app ls -la data/jobs.db
```

### Debug Issues
```bash
# View real-time logs
docker-compose logs -f hunter-seeker

# Run comprehensive tests
./scripts/test-application.sh

# Check container health
docker-compose ps
```

## File Purposes

### docs/AI_AGENT_PROMPT.md
Complete reference guide including:
- Project structure and architecture
- All commands for running, testing, debugging
- API endpoints and database schema
- Troubleshooting procedures
- Development workflow

### docs/AI_GETTING_STARTED.md
Focused quick-start guide for:
- 60-second setup process
- First-time configuration
- Essential functionality tests
- Common issues and solutions

### docs/DOCKER_REFERENCE.md
Docker-specific commands for:
- Container management
- Log viewing and debugging
- Data backup/restore
- Performance monitoring
- Cleanup operations

### scripts/test-application.sh
Automated test suite that verifies:
- Health check endpoint
- Main dashboard functionality
- Job creation and deletion
- Database operations
- Container health
- API endpoints

## Success Criteria

The application is working correctly when:

1. Health check returns `{"status":"ok","service":"hunter-seeker"}`
2. Main dashboard loads and shows job applications
3. Can create new job applications via web form or API
4. Delete functionality shows appropriate success/error messages
5. Data persists between container restarts
6. All tests in `test-application.sh` pass

## Need Help?

1. **Start here**: `docs/AI_GETTING_STARTED.md`
2. **For detailed reference**: `docs/AI_AGENT_PROMPT.md`
3. **For Docker issues**: `docs/DOCKER_REFERENCE.md`
4. **To verify everything works**: Run `./scripts/test-application.sh`

This documentation is specifically designed to help AI agents quickly understand, run, test, and debug the Hunter-Seeker application.