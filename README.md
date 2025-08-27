# Hunter-Seeker 🎯

A local web application for tracking job applications built with Go, SQLite, and a clean web interface. Keep track of your job search progress without relying on external services or spreadsheets.

## Features

- **Track Job Applications**: Record date applied, job title, company, status, job URL, and notes
- **CSV Import**: Bulk import job applications from CSV files with flexible date format support
- **Status Management**: Predefined statuses (Applied, In Review, Interview, etc.) with easy updates
- **Visual Dashboard**: Clean interface with status filtering and application statistics
- **Local-First**: Runs entirely on your machine with SQLite database
- **No Authentication Required**: Simple local-only access
- **Docker Support**: Easy deployment with Docker Compose

## Screenshots

The application provides:
- A dashboard showing all your job applications with status badges
- Filtering by application status
- Statistics showing counts by status
- Forms to add and edit applications
- Clean, responsive design that works on desktop and mobile

## Quick Start with Docker Compose

### Prerequisites

- Docker and Docker Compose installed on your machine

### Setup

1. **Clone or download this repository**
   ```bash
   git clone <your-repo-url>
   cd hunter-seeker
   ```

2. **Start the application**
   ```bash
   docker-compose up -d
   ```

3. **Access the application**
   Open your browser and go to: http://localhost:8080

4. **Stop the application**
   ```bash
   docker-compose down
   ```

### Data Persistence

Your job application data is stored in a SQLite database that persists in the `./data` directory. This directory is automatically created and mounted as a volume in Docker, so your data will survive container restarts.

## Manual Setup (Without Docker)

### Prerequisites

- Go 1.21 or later
- SQLite3

### Setup

1. **Clone the repository**
   ```bash
   git clone <your-repo-url>
   cd hunter-seeker
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Run the application**
   ```bash
   go run cmd/server/main.go
   ```

4. **Access the application**
   Open your browser and go to: http://localhost:8080

## Usage Guide

### Adding a Job Application

1. Click "Add Application" from the dashboard
2. Fill in the required fields:
   - **Date Applied**: When you submitted the application
   - **Job Title**: The position you applied for
   - **Company**: The company name
   - **Status**: Current status (defaults to "Applied")
   - **Job URL**: Link to the job posting (optional)
   - **Notes**: Any additional information (optional)
3. Click "Add Application" to save

### Importing from CSV

1. Click "Import CSV" from the dashboard or navigation
2. Download the sample template if needed
3. Prepare your CSV file with the following columns (in order):
   - **Date Applied** (required): Date you applied for the job
   - **Job Title** (required): Position title
   - **Company** (required): Company name
   - **Status** (optional): Application status (defaults to "Applied")
   - **Job URL** (optional): Link to job posting
   - **Notes** (optional): Additional notes
4. Upload your CSV file and review the import results

#### Supported Date Formats
- ISO format: `2024-01-15` (recommended)
- US format: `01/15/2024` or `1/15/2024`
- European format: `15/01/2024` or `15/1/2024`
- Month names: `Jan 15 2024`, `January 15, 2024`, `Feb 20 2024`
- Various other formats are automatically detected

#### CSV Import Tips
- Include headers in the first row (they will be automatically detected)
- Empty rows are skipped automatically
- If status is empty, it defaults to "Applied"
- Maximum file size: 10MB
- Use quotes around fields containing commas

### Updating Application Status

1. From the dashboard, click "Edit" on any application
2. Update the status (e.g., from "Applied" to "Interview")
3. Add any new notes about progress
4. Click "Update Application"

### Filtering Applications

Use the filter buttons at the top of the dashboard to view applications by status:
- **All**: Shows all applications
- **Applied**: Recently submitted applications
- **In Review**: Applications being reviewed
- **Interview**: Applications in interview process
- **Offer**: Applications with job offers
- **Rejected**: Rejected applications
- **No Response**: Applications with no response

### Managing Applications

- **Edit**: Update any field of an existing application
- **Delete**: Remove an application (with confirmation)
- **View Statistics**: See counts of applications by status

## Application Statuses

The application comes with predefined statuses:

- **Applied**: Initial application submitted
- **In Review**: Application is being reviewed
- **Phone Screen**: Phone/video screening scheduled
- **Interview**: In-person or video interview
- **Technical Test**: Technical assessment or coding challenge
- **Offer**: Job offer received
- **Rejected**: Application rejected
- **Withdrawn**: You withdrew the application
- **No Response**: No response from the company

You can use any of these statuses or mix with custom ones.

## CSV Import Format

For bulk importing job applications, use a CSV file with the following structure:

```csv
Date Applied,Job Title,Company,Status,Job URL,Notes
2024-01-15,Senior Software Engineer,TechCorp,Applied,https://techcorp.com/jobs/123,Applied through website
2024-01-20,Full Stack Developer,StartupCo,In Review,https://startupco.com/careers,Remote position
01/25/2024,Backend Engineer,BigTech Inc,Interview,https://bigtech.com/jobs,Phone screen scheduled
```

Download a sample template from the Import CSV page in the application.

## Configuration

### Environment Variables

- `PORT`: Port to run the server on (default: 8080)
- `DB_PATH`: Path to SQLite database file (default: ./data/jobs.db)

### Docker Environment

The Docker setup automatically configures:
- Port 8080 exposed
- Data persistence in `./data` directory
- Automatic container restart

## Database Schema

The application uses a simple SQLite schema:

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

## Backup and Restore

### Backing Up Your Data

Your data is stored in `./data/jobs.db`. To backup:

```bash
# Copy the database file
cp ./data/jobs.db ./backup-$(date +%Y%m%d).db

# Or with Docker
docker-compose exec hunter-seeker cp ./data/jobs.db ./data/backup-$(date +%Y%m%d).db
```

### Restoring Data

```bash
# Stop the application
docker-compose down

# Replace the database file
cp your-backup.db ./data/jobs.db

# Start the application
docker-compose up -d
```

## Development

### Project Structure

```
hunter-seeker/
├── cmd/
│   ├── server/           # Main application entry point
│   ├── debug/           # Debug utilities and database inspection
│   └── sample-data/     # Sample data generation tool
├── internal/
│   ├── database/        # Database operations and models
│   ├── handlers/        # HTTP request handlers
│   └── models/          # Data structures and business logic
├── web/
│   ├── templates/       # HTML templates
│   └── static/         # CSS stylesheets and static assets
├── scripts/            # Utility scripts for development
├── docs/              # Documentation (AI guides, Docker reference)
├── data/              # SQLite database (auto-created)
├── bin/               # Compiled binaries (auto-created)
├── Dockerfile
├── docker-compose.yml
├── Makefile
└── README.md
```

### Running in Development

```bash
# Install dependencies
go mod download

# Run with live reload (if you have air installed)
air

# Or run directly
go run cmd/server/main.go

# Build for development (outputs to bin/ directory)
make build

# Or build manually (avoid this - creates binaries in root)
# go build ./cmd/server  # DON'T DO THIS - creates 'server' in root
```

### Building for Production

```bash
# Build binary (recommended - uses bin directory)
make build

# Or build manually into bin directory
mkdir -p bin
CGO_ENABLED=1 go build -o bin/hunter-seeker ./cmd/server

# Or use Docker
docker build -t hunter-seeker .
```

## Documentation

For detailed documentation and guides:

- **`docs/AI_README.md`** - Comprehensive guide for AI agents
- **`docs/AI_GETTING_STARTED.md`** - Quick 60-second setup guide
- **`docs/AI_AGENT_PROMPT.md`** - Complete reference with commands and troubleshooting
- **`docs/DOCKER_REFERENCE.md`** - Docker commands and container management

## Sample Data

To quickly populate your application with sample job applications:

```bash
# Using the sample data tool
go run cmd/sample-data/main.go

# Or using the debug tool
go run cmd/debug/main.go add-test-data
```

## Troubleshooting

### Common Issues

1. **Port already in use**
   - Change the port in docker-compose.yml or set PORT environment variable

2. **Database permission errors**
   - Ensure the `./data` directory is writable
   - Check Docker volume permissions

3. **Application won't start**
   - Check Docker logs: `docker-compose logs hunter-seeker`
   - Verify all files are in place

### Logs

View application logs:
```bash
# Docker logs
docker-compose logs -f hunter-seeker

# Or if running manually
# Logs are output to stdout
```

## Security Notes

This application is designed for **local use only** and includes:
- No authentication system
- No HTTPS encryption
- SQLite database with no access controls

**Do not expose this application to the internet** without proper security measures.

## Contributing

This is a personal productivity tool, but improvements are welcome:
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request

## License

[Add your preferred license here]

## Changelog

### v1.0.0
- Initial release
- Basic job application tracking
- Status management
- Docker support
- Responsive web interface