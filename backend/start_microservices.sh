#!/bin/bash

# Terminate all child processes on exit
cleanup() {
    echo ""
    echo "Shutting down all services..."
    kill "$AUTH_PID" "$DOC_PID" "$GATEWAY_PID" 2>/dev/null
    exit 0
}

trap cleanup SIGINT SIGTERM EXIT

# Start Auth & User Service
echo "Starting Auth & User Service on port 8081..."
go run services/auth/main.go &
AUTH_PID=$!

# Start Document & Workflow Service
echo "Starting Document & Workflow Service on port 8082..."
go run services/document/main.go &
DOC_PID=$!

# Give the microservices a second to initialize before launching the gateway
sleep 2

# Start API Gateway
echo "Starting API Gateway on port 8080..."
go run services/gateway/main.go &
GATEWAY_PID=$!

# Wait for background jobs to finish
wait
