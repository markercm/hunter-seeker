#!/bin/bash

# Hunter-Seeker Data Clearing Script
# This script clears all job application data from the database

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_header() {
    echo ""
    echo "🗑️  Hunter-Seeker Data Clearing Utility"
    echo "======================================="
    echo ""
}

# Check if database exists
check_database() {
    if [ -f "./data/jobs.db" ]; then
        return 0
    else
        return 1
    fi
}

# Backup database before clearing
backup_database() {
    if check_database; then
        local backup_name="backup-$(date +%Y%m%d-%H%M%S).db"
        print_status "Creating backup: $backup_name"
        cp ./data/jobs.db "./data/$backup_name"
        print_success "Backup created successfully"
        return 0
    else
        print_warning "No database file found to backup"
        return 1
    fi
}

# Clear all data
clear_data() {
    if check_database; then
        print_status "Removing database file..."
        rm -f ./data/jobs.db
        print_success "Database cleared successfully"
    else
        print_warning "No database file found to clear"
    fi
}

# Check if Docker container is running
check_docker_running() {
    if command -v docker-compose >/dev/null 2>&1; then
        if docker-compose ps | grep -q "hunter-seeker-app.*Up"; then
            return 0
        fi
    fi
    return 1
}

# Main function
main() {
    print_header

    # Parse command line arguments
    BACKUP="yes"
    FORCE="no"

    while [[ $# -gt 0 ]]; do
        case $1 in
            --no-backup)
                BACKUP="no"
                shift
                ;;
            --force)
                FORCE="yes"
                shift
                ;;
            --help|-h)
                echo "Usage: $0 [options]"
                echo ""
                echo "This script clears all job application data from Hunter-Seeker."
                echo ""
                echo "Options:"
                echo "  --no-backup  Don't create a backup before clearing"
                echo "  --force      Don't ask for confirmation"
                echo "  --help       Show this help message"
                echo ""
                echo "Examples:"
                echo "  $0                    # Clear data with backup and confirmation"
                echo "  $0 --no-backup       # Clear data without backup"
                echo "  $0 --force           # Clear data without confirmation"
                echo ""
                exit 0
                ;;
            *)
                print_error "Unknown option: $1"
                echo "Use --help for usage information"
                exit 1
                ;;
        esac
    done

    # Check if Docker container is running
    if check_docker_running; then
        print_warning "Hunter-Seeker Docker container is currently running."
        echo "You may want to stop it before clearing data to avoid conflicts."
        echo ""
        read -p "Continue anyway? (y/N): " continue_choice
        if [[ ! $continue_choice =~ ^[Yy]$ ]]; then
            echo "Aborted. Stop the container with: docker-compose down"
            exit 0
        fi
    fi

    # Check if database exists
    if ! check_database; then
        print_warning "No database file found. Nothing to clear."
        exit 0
    fi

    # Show current data count if possible
    if command -v sqlite3 >/dev/null 2>&1; then
        local count=$(sqlite3 ./data/jobs.db "SELECT COUNT(*) FROM job_applications;" 2>/dev/null || echo "unknown")
        if [ "$count" != "unknown" ]; then
            echo "Current database contains $count job application(s)."
        fi
    fi

    # Confirmation prompt
    if [ "$FORCE" != "yes" ]; then
        echo ""
        print_warning "This will permanently delete all your job application data!"
        if [ "$BACKUP" = "yes" ]; then
            echo "A backup will be created before clearing."
        else
            echo "No backup will be created."
        fi
        echo ""
        read -p "Are you sure you want to continue? (y/N): " confirm

        if [[ ! $confirm =~ ^[Yy]$ ]]; then
            echo "Operation cancelled."
            exit 0
        fi
    fi

    # Create backup if requested
    if [ "$BACKUP" = "yes" ]; then
        echo ""
        backup_database
    fi

    # Clear the data
    echo ""
    clear_data

    echo ""
    print_success "Data clearing completed!"

    if [ "$BACKUP" = "yes" ]; then
        echo ""
        echo "📁 Your data has been backed up in the ./data/ directory."
        echo "   To restore, copy a backup file over ./data/jobs.db"
    fi

    echo ""
    echo "🎯 You can now start fresh with your job application tracking!"
    echo "   Visit http://localhost:8080 to begin adding new applications."
    echo ""
}

# Run main function with all arguments
main "$@"
