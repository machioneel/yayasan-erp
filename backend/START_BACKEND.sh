#!/bin/bash

# Backend Startup Script for Linux/Mac
# Run this to start the backend server

echo "🚀 Starting Yayasan ERP Backend..."
echo ""

# Check if Go is installed
echo "🔍 Checking Go installation..."
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed!"
    echo "   Please install Go 1.21+ from https://go.dev/dl/"
    exit 1
fi
echo "✅ Go installed: $(go version)"
echo ""

# Check if .env exists
echo "🔍 Checking .env file..."
if [ ! -f ".env" ]; then
    echo "⚠️  .env file not found. Creating from template..."
    cp .env.example .env 2>/dev/null
    if [ -f ".env" ]; then
        echo "✅ .env created from template"
        echo "   Please update .env with your database credentials"
    else
        echo "❌ Could not create .env file"
        echo "   Please create .env manually"
    fi
    echo ""
fi

# Download dependencies
echo "📦 Downloading dependencies..."
go mod tidy
if [ $? -ne 0 ]; then
    echo "❌ Failed to download dependencies!"
    echo "   Try running: go mod tidy"
    exit 1
fi
echo "✅ Dependencies downloaded"
echo ""

# Check if go.sum exists
if [ ! -f "go.sum" ]; then
    echo "⚠️  go.sum not found. Running go mod download..."
    go mod download
    echo ""
fi

# Start the server
echo "🚀 Starting server..."
echo "   Server will be available at: http://localhost:8080"
echo "   API documentation: http://localhost:8080/api/v1/health"
echo ""
echo "   Press Ctrl+C to stop the server"
echo ""

go run cmd/api/main.go
