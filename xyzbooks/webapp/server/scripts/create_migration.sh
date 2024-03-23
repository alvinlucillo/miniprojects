#!/bin/bash

# Function to convert title to lowercase and replace spaces with underscores
convert_title() {
    local title=$1
    local lowercase_title=$(echo "$title" | tr '[:upper:]' '[:lower:]')
    local converted_title=$(echo "$lowercase_title" | tr ' ' '_')
    echo "$converted_title"
}

# Get the title from the first argument
title="$1"

# Check if the title is empty
if [ -z "$title" ]; then
    echo "Error: Migration title is required."
    exit 1
fi

# Generate timestamp
timestamp=$(date +"%Y%m%d%H%M%S")

# Convert title to lowercase and replace spaces with underscores
converted_title=$(convert_title "$title")

# Create migration files
up_migration_file="${timestamp}_${converted_title}.up.sql"
down_migration_file="${timestamp}_${converted_title}.down.sql"

# Create empty migration files
touch ../database/migrations/"$up_migration_file"
touch ../database/migrations/"$down_migration_file"

echo "Migration files created:"
echo "$up_migration_file"
echo "$down_migration_file"
