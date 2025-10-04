package main

import "fmt"

// generateDockerfile creates Dockerfile content
func generateDockerfile(moduleName string) string {
	return fmt.Sprintf(`FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o %s cmd/%s/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/%s .

EXPOSE 8080

CMD ["./%s"]
`, moduleName, moduleName, moduleName, moduleName)
}

// generateDockerCompose creates docker-compose.yml content
func generateDockerCompose(moduleName string) string {
	return fmt.Sprintf(`services:
  %s:
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=db
      - DB_NAME=%s
      - DB_USER=postgres
      - DB_PWD=postgres
      - DB_PORT=5432
      - SVC_PORT=8080
    depends_on:
      db:
        condition: service_healthy

  db:
    image: postgres:15
    environment:
      POSTGRES_DB: %s
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./migrations:/docker-entrypoint-initdb.d/
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres -d %s"]
      interval: 10s
      timeout: 5s
      retries: 5
      start_period: 10s

volumes:
  postgres_data:
`, moduleName, moduleName, moduleName, moduleName)
}

// generateBuildScript creates build.sh script for Linux/macOS
func generateBuildScript(moduleName string) string {
	return fmt.Sprintf(`#!/bin/bash

# Build script for %s (Linux/macOS)

echo "Building %s..."

# Create build directory
mkdir -p build

# Navigate to service directory
cd cmd/%s

# Build the application
go build -o ../../build/%s

echo "Build complete! Binary available at build/%s"
`, moduleName, moduleName, moduleName, moduleName, moduleName)
}

// generateBuildScriptWindows creates build.bat script for Windows
func generateBuildScriptWindows(moduleName string) string {
	return fmt.Sprintf(`@echo off
REM Build script for %s (Windows)

echo Building %s...

REM Create build directory
if not exist "build" mkdir build

REM Navigate to service directory
cd cmd\%s

REM Build the application
go build -o ..\..\build\%s.exe

echo Build complete! Binary available at build\%s.exe
pause
`, moduleName, moduleName, moduleName, moduleName, moduleName)
}

// generateBuildScriptPowerShell creates build.ps1 script for Windows PowerShell
func generateBuildScriptPowerShell(moduleName string) string {
	return fmt.Sprintf(`# PowerShell build script for %s (Windows)

Write-Host "Building %s..." -ForegroundColor Green

# Create build directory
if (-not (Test-Path "build")) {
    New-Item -ItemType Directory -Path "build"
}

# Navigate to service directory
Set-Location "cmd\%s"

# Build the application
go build -o "..\..\build\%s.exe"

if ($LASTEXITCODE -eq 0) {
    Write-Host "Build complete! Binary available at build\%s.exe" -ForegroundColor Green
} else {
    Write-Host "Build failed!" -ForegroundColor Red
    exit 1
}

# Return to root directory
Set-Location "..\.."

Write-Host "Press any key to continue..." -ForegroundColor Yellow
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
`, moduleName, moduleName, moduleName, moduleName, moduleName)
}

// generateCrossPlatformBuildScript creates build script for all platforms
func generateCrossPlatformBuildScript(moduleName string) string {
	return fmt.Sprintf(`#!/bin/bash

# Cross-platform build script for %s

echo "Building %s for multiple platforms..."

# Create build directory
mkdir -p build

# Build for Windows (AMD64)
echo "Building for Windows AMD64..."
GOOS=windows GOARCH=amd64 go build -o build/%s-windows-amd64.exe ./cmd/%s

# Build for Linux (AMD64)
echo "Building for Linux AMD64..."
GOOS=linux GOARCH=amd64 go build -o build/%s-linux-amd64 ./cmd/%s

# Build for macOS (AMD64)
echo "Building for macOS AMD64..."
GOOS=darwin GOARCH=amd64 go build -o build/%s-darwin-amd64 ./cmd/%s

# Build for macOS (ARM64 - Apple Silicon)
echo "Building for macOS ARM64..."
GOOS=darwin GOARCH=arm64 go build -o build/%s-darwin-arm64 ./cmd/%s

# Build for Linux (ARM64)
echo "Building for Linux ARM64..."
GOOS=linux GOARCH=arm64 go build -o build/%s-linux-arm64 ./cmd/%s

echo "Cross-platform build complete!"
echo "Binaries available in build/ directory:"
ls -la build/
`, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName)
}

// generateMakefile creates Makefile content
func generateMakefile(moduleName string) string {
	return fmt.Sprintf(`# Build the application
build:
	cd cmd/%s && go build -o ../../build/%s

# Run the application
run:
	cd cmd/%s && go run main.go

# Run tests
test:
	go test ./...

# Run tests with coverage
test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	go vet ./...
	gofmt -l .
	golint ./...
	golangci-lint run

# Clean build artifacts
clean:
	rm -rf build/
	rm -f coverage.out

# Install dependencies
deps:
	go mod download
	go mod tidy

# Docker build
docker-build:
	docker build -t %s .

# Docker run
docker-run:
	docker-compose up

.PHONY: build run test test-coverage fmt lint clean deps docker-build docker-run
`, moduleName, moduleName, moduleName, moduleName)
}
