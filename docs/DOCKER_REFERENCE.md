# Docker Quick Reference for Hunter-Seeker

## Basic Commands

### Start Application
```bash
# Start in background
docker-compose up -d

# Start with rebuild
docker-compose up --build -d

# Start in foreground (see logs)
docker-compose up --build
```

### Stop Application
```bash
# Stop containers
docker-compose down

# Stop and remove volumes
docker-compose down -v

# Stop specific service
docker-compose stop hunter-seeker
```

### View Status
```bash
# Check running containers
docker-compose ps

# Check container health
docker-compose ps hunter-seeker

# View resource usage
docker stats hunter-seeker-app
```

## Logs and Debugging

### View Logs
```bash
# Follow all logs
docker-compose logs -f

# Follow specific service logs
docker-compose logs -f hunter-seeker

# View last 20 lines
docker-compose logs --tail=20 hunter-seeker

# View logs since timestamp
docker-compose logs --since="2025-08-27T10:00:00" hunter-seeker
```

### Execute Commands in Container
```bash
# Open shell in running container
docker exec -it hunter-seeker-app /bin/sh

# Check files in container
docker exec -i hunter-seeker-app ls -la data/

# Check environment variables
docker exec -i hunter-seeker-app env

# Check processes
docker exec -i hunter-seeker-app ps aux
```

## Data Management

### Database Operations
```bash
# Check database file
docker exec -i hunter-seeker-app ls -la data/jobs.db

# Copy database out of container
docker cp hunter-seeker-app:/root/data/jobs.db ./backup.db

# Copy database into container
docker cp ./backup.db hunter-seeker-app:/root/data/jobs.db

# Check database size
docker exec -i hunter-seeker-app du -h data/jobs.db
```

### Volume Management
```bash
# List volumes
docker volume ls

# Inspect volume
docker volume inspect hunter-seeker_data

# Remove volume (WARNING: deletes data)
docker volume rm hunter-seeker_data
```

## Troubleshooting

### Container Won't Start
```bash
# Check build logs
docker-compose up --build

# Force rebuild without cache
docker-compose build --no-cache hunter-seeker

# Check Docker daemon
docker info

# Check port conflicts
lsof -i :8080
```

### Application Issues
```bash
# Check if app is responding
curl -I http://localhost:8080/health

# Check container logs for errors
docker-compose logs hunter-seeker | grep -i error

# Restart just the app container
docker-compose restart hunter-seeker

# Check container resource limits
docker exec -i hunter-seeker-app cat /proc/meminfo
```

### Network Issues
```bash
# Check exposed ports
docker port hunter-seeker-app

# Check network configuration
docker network ls
docker network inspect hunter-seeker_default

# Test connectivity from host
curl -v http://localhost:8080
```

## Development Workflow

### Code Changes
```bash
# Rebuild after code changes
docker-compose up --build -d

# Quick restart for template changes
docker-compose restart hunter-seeker

# Force recreation
docker-compose up --force-recreate -d
```

### Testing
```bash
# Run tests in container
docker exec -i hunter-seeker-app go test ./...

# Copy test results out
docker cp hunter-seeker-app:/root/test-results.xml ./

# Check Go version in container
docker exec -i hunter-seeker-app go version
```

## Performance Monitoring

### Resource Usage
```bash
# Monitor real-time stats
docker stats hunter-seeker-app

# Check container processes
docker exec -i hunter-seeker-app top

# Check disk usage
docker exec -i hunter-seeker-app df -h

# Check memory usage
docker exec -i hunter-seeker-app free -h
```

### Application Metrics
```bash
# Check HTTP response times
curl -w "@curl-format.txt" -o /dev/null -s http://localhost:8080/

# Monitor database size growth
watch "docker exec -i hunter-seeker-app du -h data/jobs.db"

# Check application logs for performance
docker-compose logs hunter-seeker | grep -E "(took|ms|seconds)"
```

## Cleanup Commands

### Clean Up Containers
```bash
# Remove stopped containers
docker container prune

# Remove unused images
docker image prune

# Remove everything unused
docker system prune

# Remove everything including volumes
docker system prune -a --volumes
```

### Reset Application
```bash
# Complete reset (keeps data)
docker-compose down
docker-compose up --build -d

# Complete reset (deletes data)
docker-compose down -v
docker-compose up --build -d
```

## Docker Compose Shortcuts

### Service-Specific Commands
```bash
# Build specific service
docker-compose build hunter-seeker

# View service configuration
docker-compose config

# Scale service (not applicable for this app)
docker-compose up --scale hunter-seeker=2 -d
```

### Environment Management
```bash
# Use different compose file
docker-compose -f docker-compose.dev.yml up

# Override environment variables
PORT=9090 docker-compose up -d

# Use environment file
docker-compose --env-file .env.local up -d
```

## Health Checks

### Application Health
```bash
# HTTP health check
curl http://localhost:8080/health

# Container health status
docker inspect hunter-seeker-app | grep -A5 Health

# Health check logs
docker-compose logs hunter-seeker | grep -i health
```

### System Health
```bash
# Check Docker daemon health
docker system df

# Check available disk space
df -h

# Check if ports are available
netstat -tlnp | grep :8080
```

## Backup and Restore

### Backup
```bash
# Backup database
docker cp hunter-seeker-app:/root/data/jobs.db ./backup-$(date +%Y%m%d).db

# Backup entire data directory
docker cp hunter-seeker-app:/root/data ./backup-data-$(date +%Y%m%d)

# Export container
docker export hunter-seeker-app > hunter-seeker-backup.tar
```

### Restore
```bash
# Stop application
docker-compose down

# Restore database
docker cp ./backup-20250827.db hunter-seeker-app:/root/data/jobs.db

# Start application
docker-compose up -d

# Verify restore
curl http://localhost:8080/health
```

## Common Issues and Solutions

### Port Already in Use
```bash
# Find process using port
lsof -i :8080

# Kill process
kill -9 $(lsof -t -i:8080)

# Use different port
PORT=8081 docker-compose up -d
```

### Database Permission Issues
```bash
# Check file permissions
docker exec -i hunter-seeker-app ls -la data/

# Fix permissions
docker exec -i hunter-seeker-app chmod 644 data/jobs.db

# Check ownership
docker exec -i hunter-seeker-app id
```

### Out of Disk Space
```bash
# Check disk usage
docker system df

# Clean up unused resources
docker system prune -a

# Remove old images
docker image prune -a --filter "until=24h"
```

This reference covers the most common Docker operations you'll need when working with the Hunter-Seeker application.