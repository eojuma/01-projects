#!/bin/bash
# Script to build and start the social network application containers

echo "Building and starting Docker containers..."
docker compose up --build -d

echo ""
echo "Social Network is running!"
echo "- Frontend: http://localhost:3000"
echo "- Backend:  http://localhost:8080"
echo ""
echo "Fetching logs... (Press Ctrl+C to stop viewing logs, containers will keep running)"
docker compose logs -f
