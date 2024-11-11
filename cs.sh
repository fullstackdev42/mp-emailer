#!/bin/bash

# Function to create a directory if it doesn't exist
create_dir() {
    if [ ! -d "$1" ]; then
        mkdir -p "$1"
        echo "Created directory $1"
    fi
}

# Function to create a file and add a comment if the file doesn't exist
create_file() {
    if [ ! -f "$1" ]; then
        touch "$1"
        echo "$2" >> "$1"
        echo "Created file $1 with comment"
    fi
}

# Create directories
create_dir "database/core"
create_dir "database/decorators"
create_dir "database/config"
create_dir "database/migrations"

# Create files with comments
create_file "database/core/interface.go" "// Database interface definitions"
create_file "database/core/db.go" "// Core DB implementation"
create_file "database/core/result.go" "// Result interface and implementation"
create_file "database/decorators/logging.go" "// Logging decorator"
create_file "database/config/retry.go" "// Retry configuration"
create_file "database/config/connect.go" "// Connection logic"
create_file "database/migrations/migrator.go" "// Migration interface and implementation"
