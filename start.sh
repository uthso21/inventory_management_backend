#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ML_DIR="$SCRIPT_DIR/Ml models"

# PIDs
FASTAPI_PID=""
GO_PID=""

# Cleanup function
cleanup() {
    echo -e "\n${YELLOW}Shutting down services...${NC}"
    [[ -n "$FASTAPI_PID" ]] && kill $FASTAPI_PID 2>/dev/null
    [[ -n "$GO_PID" ]] && kill $GO_PID 2>/dev/null
    echo -e "${GREEN}All services stopped.${NC}"
    exit 0
}

trap cleanup SIGINT SIGTERM

echo -e "${BLUE}╔════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║    ${GREEN}Inventory Management System${BLUE}                     ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════╝${NC}\n"

# 1. Check PostgreSQL
echo -e "${YELLOW}[1/3] Checking PostgreSQL...${NC}"
if command -v pg_isready &> /dev/null; then
    if pg_isready -q -p 5433; then
        echo -e "${GREEN}✓ PostgreSQL is running on port 5433${NC}"
    elif pg_isready -q -p 5432; then
        echo -e "${GREEN}✓ PostgreSQL is running on port 5432${NC}"
    else
        echo -e "${RED}PostgreSQL is not running. Starting...${NC}"
        sudo systemctl start postgresql 2>/dev/null || sudo service postgresql start 2>/dev/null
        sleep 2
        if pg_isready -q; then
            echo -e "${GREEN}✓ PostgreSQL started successfully${NC}"
        else
            echo -e "${RED}✗ Failed to start PostgreSQL. Please start it manually.${NC}"
            exit 1
        fi
    fi
else
    echo -e "${YELLOW}⚠ pg_isready not found. Assuming PostgreSQL is running.${NC}"
fi

# 2. Start FastAPI server
echo -e "\n${YELLOW}[2/3] Starting FastAPI ML service on port 8000...${NC}"
cd "$ML_DIR"
if [ -d ".venv" ]; then
    .venv/bin/uvicorn main:app --reload --host 0.0.0.0 --port 8000 > /dev/null 2>&1 &
else
    uvicorn main:app --reload --host 0.0.0.0 --port 8000 > /dev/null 2>&1 &
fi
FASTAPI_PID=$!
sleep 3

# Wait for FastAPI to be ready
for i in {1..10}; do
    if curl -s http://localhost:8000/health > /dev/null 2>&1; then
        echo -e "${GREEN}✓ FastAPI ML service started (PID: $FASTAPI_PID)${NC}"
        break
    fi
    if [ $i -eq 10 ]; then
        echo -e "${RED}✗ Failed to start FastAPI service${NC}"
        exit 1
    fi
    sleep 1
done

# 3. Start Go backend
echo -e "\n${YELLOW}[3/3] Starting Go backend on port 8080...${NC}"
cd "$SCRIPT_DIR"
go run cmd/app/main.go > /dev/null 2>&1 &
GO_PID=$!

# Wait for Go backend to be ready
for i in {1..15}; do
    if curl -s http://localhost:8080/ml/health > /dev/null 2>&1; then
        echo -e "${GREEN}✓ Go backend started (PID: $GO_PID)${NC}"
        break
    fi
    if [ $i -eq 15 ]; then
        echo -e "${RED}✗ Failed to start Go backend${NC}"
        cleanup
        exit 1
    fi
    sleep 1
done

echo -e "\n${BLUE}╔════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║  ${GREEN}All services are running!${BLUE}                        ║${NC}"
echo -e "${BLUE}╠════════════════════════════════════════════════════╣${NC}"
echo -e "${BLUE}║${NC}  PostgreSQL:      ${GREEN}localhost:5433${NC}                   ${BLUE}║${NC}"
echo -e "${BLUE}║${NC}  FastAPI ML:      ${GREEN}http://localhost:8000${NC}            ${BLUE}║${NC}"
echo -e "${BLUE}║${NC}  Go Backend API:  ${GREEN}http://localhost:8080${NC}            ${BLUE}║${NC}"
echo -e "${BLUE}╠════════════════════════════════════════════════════╣${NC}"
echo -e "${BLUE}║${NC}  ${YELLOW}Press Ctrl+C to stop all services${NC}                 ${BLUE}║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════╝${NC}\n"

# Wait for processes
wait
