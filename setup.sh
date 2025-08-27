#!/bin/bash

# Hunter-Seeker Job Tracker Setup Script
# This script helps you get started with the job application tracker

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
    echo "🎯 Hunter-Seeker Job Application Tracker Setup"
    echo "=============================================="
    echo ""
}

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check dependencies
check_dependencies() {
    print_status "Checking dependencies..."

    local missing_deps=()

    if ! command_exists docker; then
        missing_deps+=("docker")
    fi

    if ! command_exists docker-compose; then
        missing_deps+=("docker-compose")
    fi

    if [ ${#missing_deps[@]} -gt 0 ]; then
        print_error "Missing dependencies: ${missing_deps[*]}"
        echo ""
        echo "Please install the following:"
        for dep in "${missing_deps[@]}"; do
            case $dep in
                docker)
                    echo "  - Docker: https://docs.docker.com/get-docker/"
                    ;;
                docker-compose)
                    echo "  - Docker Compose: https://docs.docker.com/compose/install/"
                    ;;
            esac
        done
        echo ""
        exit 1
    fi

    print_success "All dependencies are installed"
}

# Setup function for Docker
setup_docker() {
    print_status "Setting up with Docker..."

    # Create data directory
    mkdir -p ./data
    print_status "Created data directory"

    # Build the application
    print_status "Building Docker image (this may take a few minutes)..."
    docker-compose build
    print_success "Docker image built successfully"

    # Start the application
    print_status "Starting the application..."
    docker-compose up -d
    print_success "Application started"

    # Wait a moment for the application to start
    sleep 3

    # Check if the application is running
    if curl -s http://localhost:8080/health > /dev/null 2>&1; then
        print_success "Application is running and healthy"
    else
        print_warning "Application may still be starting up"
    fi
}

# Setup function for local Go development
setup_local() {
    print_status "Setting up for local development..."

    # Check if Go is installed
    if ! command_exists go; then
        print_error "Go is not installed. Please install Go 1.21 or later from https://golang.org/dl/"
        exit 1
    fi

    # Check Go version
    go_version=$(go version | awk '{print $3}' | sed 's/go//')
    required_version="1.21"

    if [ "$(printf '%s\n' "$required_version" "$go_version" | sort -V | head -n1)" != "$required_version" ]; then
        print_error "Go version $go_version is too old. Please install Go $required_version or later"
        exit 1
    fi

    print_success "Go version $go_version is compatible"

    # Install dependencies
    print_status "Installing Go dependencies..."
    go mod download
    go mod tidy
    print_success "Dependencies installed"

    # Create data directory
    mkdir -p ./data
    print_status "Created data directory"

    # Build the application
    print_status "Building application..."
    mkdir -p ./bin
    CGO_ENABLED=0 go build -o ./bin/hunter-seeker ./cmd/server
    print_success "Application built successfully"

    print_status "You can now run the application with: ./bin/hunter-seeker"
    print_status "Or run in development mode with: make dev (if you have make installed)"
}

# Add sample data
add_sample_data() {
    print_status "Adding sample job applications..."

    if [ "$1" = "docker" ]; then
        # For Docker setup, we'll add data directly since the volume is mounted
        go run cmd/sample-data/main.go 2>/dev/null || {
            print_warning "Could not add sample data via Go. You can add it manually later."
            return
        }
    else
        # For local setup
        go run cmd/sample-data/main.go
    fi

    print_success "Sample data added successfully"
}

# Display final instructions
show_final_instructions() {
    local setup_type=$1

    echo ""
    echo "🎉 Setup Complete!"
    echo "=================="
    echo ""

    if [ "$setup_type" = "docker" ]; then
        echo "Your Hunter-Seeker job tracker is now running via Docker!"
        echo ""
        echo "📍 Access your application:"
        echo "   Web Interface: http://localhost:8080"
        echo "   Health Check:  http://localhost:8080/health"
        echo "   API Stats:     http://localhost:8080/api/stats"
        echo ""
        echo "🔧 Useful Docker commands:"
        echo "   Stop:    docker-compose down"
        echo "   Start:   docker-compose up -d"
        echo "   Logs:    docker-compose logs -f"
        echo "   Rebuild: docker-compose build"
        echo ""
    else
        echo "Your Hunter-Seeker job tracker is set up for local development!"
        echo ""
        echo "🚀 To start the application:"
        echo "   ./bin/hunter-seeker"
        echo ""
        echo "📍 Once started, access your application:"
        echo "   Web Interface: http://localhost:8080"
        echo ""
        echo "🔧 Development commands:"
        echo "   make dev     - Run with live reload"
        echo "   make build   - Build the application"
        echo "   make test    - Run tests"
        echo "   make help    - Show all available commands"
        echo ""
    fi

    echo "📊 Your application comes with sample data to get you started."
    echo "📁 Database location: ./data/jobs.db"
    echo ""
    echo "Happy job hunting! 🎯"
    echo ""
}

# Main setup logic
main() {
    print_header

    # Parse command line arguments
    SETUP_TYPE=""
    ADD_SAMPLE="yes"

    while [[ $# -gt 0 ]]; do
        case $1 in
            --docker)
                SETUP_TYPE="docker"
                shift
                ;;
            --local)
                SETUP_TYPE="local"
                shift
                ;;
            --no-sample)
                ADD_SAMPLE="no"
                shift
                ;;
            --help|-h)
                echo "Usage: $0 [options]"
                echo ""
                echo "Options:"
                echo "  --docker     Setup using Docker (recommended)"
                echo "  --local      Setup for local Go development"
                echo "  --no-sample  Don't add sample data"
                echo "  --help       Show this help message"
                echo ""
                echo "If no option is specified, you'll be prompted to choose."
                exit 0
                ;;
            *)
                print_error "Unknown option: $1"
                echo "Use --help for usage information"
                exit 1
                ;;
        esac
    done

    # If no setup type specified, prompt user
    if [ -z "$SETUP_TYPE" ]; then
        echo "How would you like to set up Hunter-Seeker?"
        echo ""
        echo "1) Docker (recommended) - Easy setup, no local Go required"
        echo "2) Local Development - For Go developers who want to modify the code"
        echo ""
        read -p "Enter your choice (1 or 2): " choice

        case $choice in
            1)
                SETUP_TYPE="docker"
                ;;
            2)
                SETUP_TYPE="local"
                ;;
            *)
                print_error "Invalid choice. Please run the script again and choose 1 or 2."
                exit 1
                ;;
        esac
    fi

    # Check dependencies based on setup type
    if [ "$SETUP_TYPE" = "docker" ]; then
        check_dependencies
        setup_docker
    else
        setup_local
    fi

    # Add sample data if requested
    if [ "$ADD_SAMPLE" = "yes" ]; then
        echo ""
        read -p "Would you like to add sample job applications to get started? (y/N): " add_sample_choice
        if [[ $add_sample_choice =~ ^[Yy]$ ]]; then
            add_sample_data "$SETUP_TYPE"
        fi
    fi

    # Show final instructions
    show_final_instructions "$SETUP_TYPE"
}

# Run main function with all arguments
main "$@"
