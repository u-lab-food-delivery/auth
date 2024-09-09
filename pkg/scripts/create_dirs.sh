#!/bin/bash

# Function to print usage
usage() {
  echo "Usage: $0 <path>"
  exit 1
}

# Check if the path is provided
if [ -z "$1" ]; then
  usage
fi

# Get the path from the first argument
TARGET_PATH=$1

# Check if the target path exists and is a directory
if [ ! -d "$TARGET_PATH" ]; then
  echo "The specified path does not exist or is not a directory."
  exit 1
fi

# Create default directories and files
echo "Creating default directory structure under $TARGET_PATH..."

# Create api directory with subdirectories and files
API_DIR="$TARGET_PATH/api"
mkdir -p "$API_DIR/handlers"
mkdir -p "$API_DIR/middleware"
mkdir -p "$API_DIR/docs"
touch "$API_DIR/routers.go"
touch "$API_DIR/utils.go"
echo "package handlers" > "$API_DIR/handlers/main_handler.go"
echo "package middleware" > "$API_DIR/middleware/auth_middleware.go"
echo "package api" > "$API_DIR/routers.go"
echo "package api" > "$API_DIR/utils.go"

# Create cmd directory with main.go file
CMD_DIR="$TARGET_PATH/cmd"
mkdir -p "$CMD_DIR"
touch "$CMD_DIR/main.go"
echo "package main" > "$CMD_DIR/main.go"

# Create storage directory with subdirectories
STORAGE_DIR="$TARGET_PATH/storage"
mkdir -p "$STORAGE_DIR/postgres"
mkdir -p "$STORAGE_DIR/redis"
touch "$STORAGE_DIR/postgres/postgres.go"
touch "$STORAGE_DIR/redis/redis.go"
echo "package postgres" > "$STORAGE_DIR/postgres/postgres.go"
echo "package redis" > "$STORAGE_DIR/redis/redis.go"

# Create migrations directory
MIGRATION_DIR="$TARGET_PATH/migrations"
mkdir -p "$MIGRATION_DIR"

# Create tests directory with subdirectories and files
TESTS_DIR="$TARGET_PATH/tests"
mkdir -p "$TESTS_DIR/handlers"
mkdir -p "$TESTS_DIR/middleware"
mkdir -p "$TESTS_DIR/storage/postgres"
mkdir -p "$TESTS_DIR/storage/redis"

LOGGER_DIR="$TARGET_PATH/logger"
mkdir -p "$LOGGER_DIR"
touch "$LOGGER_DIR/logger.go"
echo "package logger" > "$LOGGER_DIR/logger.go"

# Create test files with package declarations
touch "$TESTS_DIR/handlers/main_handler_test.go"
touch "$TESTS_DIR/middleware/auth_middleware_test.go"
touch "$TESTS_DIR/storage/postgres/postgres_test.go"
touch "$TESTS_DIR/storage/redis/redis_test.go"

# Write package declarations and simple test boilerplate to the test files
echo "package handlers" > "$TESTS_DIR/handlers/main_handler_test.go"
echo "import \"testing\"" >> "$TESTS_DIR/handlers/main_handler_test.go"
echo "" >> "$TESTS_DIR/handlers/main_handler_test.go"
echo "func TestExample(t *testing.T) {}" >> "$TESTS_DIR/handlers/main_handler_test.go"

echo "package middleware" > "$TESTS_DIR/middleware/auth_middleware_test.go"
echo "import \"testing\"" >> "$TESTS_DIR/middleware/auth_middleware_test.go"
echo "" >> "$TESTS_DIR/middleware/auth_middleware_test.go"
echo "func TestExample(t *testing.T) {}" >> "$TESTS_DIR/middleware/auth_middleware_test.go"

echo "package postgres" > "$TESTS_DIR/storage/postgres/postgres_test.go"
echo "import \"testing\"" >> "$TESTS_DIR/storage/postgres/postgres_test.go"
echo "" >> "$TESTS_DIR/storage/postgres/postgres_test.go"
echo "func TestExample(t *testing.T) {}" >> "$TESTS_DIR/storage/postgres/postgres_test.go"

echo "package redis" > "$TESTS_DIR/storage/redis/redis_test.go"
echo "import \"testing\"" >> "$TESTS_DIR/storage/redis/redis_test.go"
echo "" >> "$TESTS_DIR/storage/redis/redis_test.go"
echo "func TestExample(t *testing.T) {}" >> "$TESTS_DIR/storage/redis/redis_test.go"

echo "Default directory structure created successfully."

# create .env
touch "config.env"