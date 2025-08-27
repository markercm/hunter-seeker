# Database Management Guide

This guide covers database management for the Hunter-Seeker job tracking application.

## Database Overview

Hunter-Seeker uses SQLite as its database engine with a single database file:
- **Location**: `./data/jobs.db`
- **Type**: SQLite 3.x database
- **Purpose**: Stores all job application data locally

## Database Configuration

### Application Settings
All components use consistent database paths:

| Component | Database Path | Configuration |
|-----------|---------------|---------------|
| Main Server | `./data/jobs.db` | Environment variable `DB_PATH` (default: `./data/jobs.db`) |
| Debug Tool | `./data/jobs.db` | Hardcoded in `cmd/debug/main.go` |
| Sample Data Tool | `./data/jobs.db` | Hardcoded in `cmd/sample-data/main.go` |
| Docker Container | `./data/jobs.db` | Set via `DB_PATH=./data/jobs.db` in docker-compose.yml |

### Environment Variables
- `DB_PATH`: Override default database path (optional)
- Default behavior: Creates `./data/jobs.db` if not specified

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

## Common Operations

### View Database Contents
```bash
# Using debug tool
go run cmd/debug/main.go

# View specific database file
go run cmd/debug/main.go ./path/to/database.db

# Using SQLite CLI (if installed)
sqlite3 ./data/jobs.db "SELECT * FROM job_applications;"
```

### Add Sample Data
```bash
# Add comprehensive sample data (12 applications)
go run cmd/sample-data/main.go

# Add simple test data (3 applications)
go run cmd/debug/main.go add-test-data

# Using make command
make sample-data
```

### Backup Database
```bash
# Manual backup
cp ./data/jobs.db ./data/backup-$(date +%Y%m%d-%H%M%S).db

# Using make command
make backup

# Docker backup
docker cp hunter-seeker-app:/root/data/jobs.db ./backup.db
```

### Restore Database
```bash
# Stop application first
docker-compose down

# Restore from backup
cp ./data/backup-20250827-120000.db ./data/jobs.db

# Restart application
docker-compose up -d
```

### Clear All Data
```bash
# Interactive clear with backup
make clear-data

# Force clear without confirmation
./scripts/clear_data.sh --force

# Manual clear
rm -f ./data/jobs.db
```

## Troubleshooting

### Database File Not Found
```bash
# Check if data directory exists
ls -la ./data/

# Create data directory if missing
mkdir -p ./data

# Restart application to create database
docker-compose restart
```

### Database Permissions Issues
```bash
# Check permissions
ls -la ./data/jobs.db

# Fix permissions (if needed)
chmod 644 ./data/jobs.db

# For Docker
docker exec -it hunter-seeker-app ls -la data/
```

### Multiple Database Files
If you find multiple database files in `./data/`:
```bash
# Check which one is active
docker-compose logs hunter-seeker | grep "Database file"

# Check file sizes and modification times
ls -la ./data/*.db

# Remove unused files (after backing up)
rm ./data/unused-database.db
```

### Database Corruption
```bash
# Check database integrity
sqlite3 ./data/jobs.db "PRAGMA integrity_check;"

# Dump and recreate if corrupted
sqlite3 ./data/jobs.db ".dump" > backup.sql
rm ./data/jobs.db
sqlite3 ./data/jobs.db < backup.sql
```

## Best Practices

### Development
1. **Never commit database files** to version control
2. **Use sample data** for testing: `make sample-data`
3. **Regular backups** before major changes: `make backup`
4. **Use debug tool** to inspect data: `go run cmd/debug/main.go`

### Production
1. **Regular backups** using automated scripts
2. **Monitor disk space** - SQLite files can grow large
3. **Backup before updates** to the application
4. **Use environment variables** for custom database paths

### Data Safety
1. **Always backup** before clearing data
2. **Test restore procedures** regularly
3. **Keep multiple backup generations**
4. **Monitor database file integrity**

## Data Migration

### Exporting Data
```bash
# Export to SQL
sqlite3 ./data/jobs.db ".dump" > export.sql

# Export to CSV
sqlite3 -header -csv ./data/jobs.db "SELECT * FROM job_applications;" > export.csv

# Export specific data
sqlite3 -header -csv ./data/jobs.db "SELECT * FROM job_applications WHERE status='Applied';" > applied_jobs.csv
```

### Importing Data
```bash
# Import from SQL dump
sqlite3 ./data/jobs.db < import.sql

# Note: CSV import requires custom scripting or SQL INSERT statements
```

## Monitoring

### Database Size
```bash
# Check file size
du -h ./data/jobs.db

# In Docker
docker exec hunter-seeker-app du -h data/jobs.db
```

### Application Count
```bash
# Total applications
go run cmd/debug/main.go | head -1

# Using SQLite
sqlite3 ./data/jobs.db "SELECT COUNT(*) FROM job_applications;"
```

### Status Distribution
```bash
# Using API endpoint
curl -s http://localhost:8080/api/stats | jq

# Using SQLite
sqlite3 ./data/jobs.db "SELECT status, COUNT(*) FROM job_applications GROUP BY status;"
```

## Security Considerations

1. **Local Storage**: Database contains personal job search information
2. **No Encryption**: SQLite file is not encrypted by default
3. **File Permissions**: Ensure appropriate file system permissions
4. **Backup Security**: Secure backup files appropriately
5. **No Network Access**: Database is not accessible remotely (by design)

## File Location Reference

```
hunter-seeker/
├── data/
│   ├── jobs.db              # Main database file
│   ├── backup-*.db          # Backup files (created manually)
│   └── (temp files)         # Test databases (auto-cleaned)
├── cmd/debug/main.go        # Database inspection tool
├── cmd/sample-data/main.go  # Sample data generator
└── scripts/clear_data.sh    # Data clearing script
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `DB_PATH` | `./data/jobs.db` | Database file path |
| `PORT` | `8080` | Server port (doesn't affect database) |

## Quick Commands Reference

```bash
# View data
go run cmd/debug/main.go

# Add sample data
make sample-data

# Backup database
make backup

# Clear data safely
make clear-data

# Check application health
curl http://localhost:8080/health

# Get statistics
curl http://localhost:8080/api/stats
```

This guide covers the essential aspects of database management for Hunter-Seeker. For application-specific issues, refer to the main documentation in the `docs/` directory.