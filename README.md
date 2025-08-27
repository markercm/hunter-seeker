# Hunter-Seeker 🎯

A local web application for tracking job applications built with Go, SQLite, and a clean web interface. Keep track of your job search progress without relying on external services or spreadsheets.

![screenshot of  hunter-seeker web app](https://github.com/markercm/hunter-seeker/blob/main/images/screenshot.png?raw=true)

## Features

- **Track Job Applications**: Record date applied, job title, company, status, job URL, and notes
- **CSV Import**: Bulk import job applications from CSV files with flexible date format support
- **Status Management**: Predefined statuses (Applied, In Review, Interview, etc.) with easy updates
- **Visual Dashboard**: Clean interface with status filtering and application statistics
- **Local-First**: Runs entirely on your machine with SQLite database
- **No Authentication Required**: Simple local-only access
- **Docker Support**: Easy deployment with Docker Compose

## Quick Start with Docker Compose

### Prerequisites

- Docker and Docker Compose installed on your machine

### Setup

1. **Start the application**

   ```bash
   docker-compose up -d
   ```

2. **Access the application**
   Open your browser and go to: <http://localhost:8080>

3. **Stop the application**

   ```bash
   docker-compose down
   ```

### Data Persistence

Your job application data is stored in a SQLite database that persists in the `./data` directory. This directory is automatically created and mounted as a volume in Docker, so your data will survive container restarts.

## Documentation

For detailed documentation and guides:

- **`docs/AI_CONTEXT.md`** - Complete AI agent context guide with commands, testing, and troubleshooting

## Security Notes

This application is designed for **local use only** and includes:

- No authentication system
- No HTTPS encryption
- SQLite database with no access controls

**Do not expose this application to the internet** without proper security measures.

## Development Notes

This application was largely built with the assistance of AI agents (Claude), demonstrating how AI can help create functional, well-structured applications with proper documentation and testing.
